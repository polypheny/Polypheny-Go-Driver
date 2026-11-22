package polypheny

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestIntegrationStructuredErrors(t *testing.T) {
	db := integrationDatabase(t)

	t.Run("statement and rollback", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		table := fmt.Sprintf("go_error_%x", time.Now().UnixNano())
		if _, err := db.ExecContext(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY)"); err != nil {
			t.Fatalf("create test table: %v", err)
		}
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if _, err := db.ExecContext(ctx, "DROP TABLE "+table); err != nil {
				t.Errorf("drop test table: %v", err)
			}
		})

		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		if _, err := tx.ExecContext(ctx, "INSERT INTO "+table+" VALUES (1)"); err != nil {
			t.Fatalf("insert in transaction: %v", err)
		}
		var value int
		err = tx.QueryRowContext(ctx, "SELECT missing_column FROM "+table).Scan(&value)
		assertIntegrationSQLState(t, err, "42703")

		if err := tx.Rollback(); err != nil {
			t.Fatalf("rollback after statement error: %v", err)
		}
		if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&value); err != nil {
			t.Fatalf("query on the same connection after rollback: %v", err)
		}
		if value != 0 {
			t.Fatalf("found %d rows after rollback, want 0", value)
		}
	})

	t.Run("login", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		address, username, password, err := parseDSN(os.Getenv("POLYPHENY_TEST_DSN"))
		if err != nil {
			t.Fatal(err)
		}
		invalidDB := sql.OpenDB(&Connector{
			address:  address,
			username: username,
			password: password + "-invalid",
		})
		defer invalidDB.Close()
		assertIntegrationSQLState(t, invalidDB.PingContext(ctx), "28000")
	})
}

func assertIntegrationSQLState(t *testing.T, err error, want string) {
	t.Helper()
	var serverErr *ClientError
	if !errors.As(err, &serverErr) {
		t.Fatalf("got %T (%v), want *ClientError with SQLSTATE %s", err, err, want)
	}
	if got := serverErr.SQLState(); got != want {
		t.Fatalf("got SQLSTATE %q for %v, want %q", got, serverErr, want)
	}
}
