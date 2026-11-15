package polypheny

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
)

func integrationDatabase(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("POLYPHENY_TEST_DSN")
	if dsn == "" {
		t.Skip("set POLYPHENY_TEST_DSN to run tests against a Polypheny server")
	}
	db, err := sql.Open("polypheny", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	db.SetMaxOpenConns(1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping Polypheny: %v", err)
	}
	return db
}

func TestIntegrationSQLAndMQL(t *testing.T) {
	db := integrationDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	table := fmt.Sprintf("go_driver_test_%x", time.Now().UnixNano())
	if _, err := db.ExecContext(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY, name VARCHAR(30))"); err != nil {
		t.Fatalf("create test table: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := db.ExecContext(ctx, "DROP TABLE "+table); err != nil {
			t.Errorf("drop test table: %v", err)
		}
	})

	// Insert more rows than Prism's default frame size to exercise pagination.
	const rowCount = 110
	insert := "INSERT INTO " + table + " (id, name) VALUES (?, ?)"
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, insert)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	for id := 1; id <= rowCount; id++ {
		name := any(fmt.Sprintf("name-%d", id))
		if id == 2 {
			name = nil
		}
		result, err := stmt.ExecContext(ctx, id, name)
		if err != nil {
			t.Fatalf("insert row %d: %v", id, err)
		}
		if affected, err := result.RowsAffected(); err != nil || affected != 1 {
			t.Fatalf("insert row %d affected %d rows: %v", id, affected, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, "SELECT id FROM "+table+" ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		count++
		if id != count {
			t.Fatalf("row %d has id %d", count, id)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if count != rowCount {
		t.Fatalf("read %d rows, want %d", count, rowCount)
	}

	// database/sql should route positional arguments through a prepared statement.
	var name sql.NullString
	if err := db.QueryRowContext(ctx, "SELECT name FROM "+table+" WHERE id = ?", 2).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name.Valid {
		t.Fatalf("read %q, want SQL NULL", name.String)
	}

	tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, insert, rowCount+1, "rollback"); err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	var remaining int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&remaining); err != nil {
		t.Fatal(err)
	}
	if remaining != rowCount {
		t.Fatalf("found %d rows after rollback, want %d", remaining, rowCount)
	}

	if _, err := db.ExecContext(ctx, insert, rowCount+1, "auto-commit"); err != nil {
		t.Fatalf("insert outside a transaction: %v", err)
	}
	independentDB := integrationDatabase(t)
	if err := independentDB.QueryRowContext(ctx,
		"SELECT name FROM "+table+" WHERE id = ?", rowCount+1).Scan(&name); err != nil {
		t.Fatalf("read auto-committed row over a separate connection: %v", err)
	}
	if !name.Valid || name.String != "auto-commit" {
		t.Fatalf("found %v, want the auto-committed row", name)
	}

	if _, err := db.ExecContext(ctx, "SELECT missing_column FROM "+table); err == nil {
		t.Fatal("invalid SQL must return a server error")
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.Raw(func(raw any) error {
		documents, err := raw.(*PolyphenyConn).QueryMongoContext(ctx, "db."+table+".find({id: 1})")
		if err != nil {
			return err
		}
		if documents == nil || len(*documents) != 1 || (*documents)[0]["name"] != "name-1" {
			return fmt.Errorf("unexpected MQL documents: %v", documents)
		}
		return nil
	}); err != nil {
		t.Fatalf("query with MQL: %v", err)
	}
}
