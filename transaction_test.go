package polypheny

import (
	context "context"
	driver "database/sql/driver"
	testing "testing"
)

func TestCommit(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer conn.(*PolyphenyConn).Close()
	tx, err := conn.(*PolyphenyConn).BeginTx(context.Background(), driver.TxOptions{})
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "DROP TABLE IF EXISTS mytable", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "CREATE TABLE mytable(id int not null, primary key(id))", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "insert into mytable values(1)", nil)
	if err != nil {
		t.Error(err)
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}

func TestRollback(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer conn.(*PolyphenyConn).Close()
	tx, err := conn.(*PolyphenyConn).BeginTx(context.Background(), driver.TxOptions{})
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "DROP TABLE IF EXISTS mytable", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "CREATE TABLE mytable(id int not null, primary key(id))", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "insert into mytable values(1)", nil)
	if err != nil {
		t.Error(err)
	}
	err = tx.Rollback()
	if err != nil {
		t.Error(err)
	}
}
