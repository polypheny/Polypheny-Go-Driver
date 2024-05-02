package polypheny

import (
	driver "database/sql/driver"
	testing "testing"
)

func TestStmtClose(t *testing.T) {
	d := PolyphenyDriver{}
	conn, err := d.Open("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	stmt, err := conn.Prepare("sql:SELECT name FROM emps WHERE name = ?")
	if err != nil {
		t.Error(err)
	}
	err = stmt.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestNumInput(t *testing.T) {
	d := PolyphenyDriver{}
	conn, err := d.Open("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	stmt, err := conn.Prepare("sql:SELECT name FROM emps WHERE name = ?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	if stmt.NumInput() != 1 {
		t.Error("Incorrect number of args for the prepared statement")
	}
}

func TestStmtExec(t *testing.T) {
	d := PolyphenyDriver{}
	conn, err := d.Open("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	_, err = conn.(*PolyphenyConn).Exec("sql:DROP TABLE IF EXISTS mytable", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).Exec("sql:CREATE TABLE mytable(id int not null, primary key(id))", nil)
	if err != nil {
		t.Error(err)
	}
	stmt, err := conn.Prepare("sql:insert into mytable values(?)")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	args := make([]driver.Value, 1)
	args[0] = int32(1)
	// here we didn't use the interace to call Exec so that it wont warn the usa of deprecated function
	result, err := stmt.(*PreparedStatement).Exec(args)
	if err != nil {
		t.Error(err)
	}
	if result.(*PolyphenyResult).rowsAffected != 1 {
		t.Errorf("The number of affected rows should be 1 but got %d", result.(*PolyphenyResult).rowsAffected)
	}
}

func TestStmtExecInternal(t *testing.T) {
	d := PolyphenyDriver{}
	conn, err := d.Open("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	_, err = conn.(*PolyphenyConn).Exec("sql:DROP TABLE IF EXISTS mytable", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).Exec("sql:CREATE TABLE mytable(id int not null, primary key(id))", nil)
	if err != nil {
		t.Error(err)
	}
	stmt, err := conn.Prepare("sql:insert into mytable values(?)")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	errChan := make(chan error)
	resultChan := make(chan *PolyphenyResult)
	var result *PolyphenyResult
	args := make([]driver.NamedValue, 1)
	args[0] = driver.NamedValue{
		Ordinal: 1,
		Value:   int32(1),
	}
	go stmt.(*PreparedStatement).execContextInternal(args, resultChan, errChan)
	result = <-resultChan
	err = <-errChan
	if err != nil {
		t.Error(err, result == nil)
	}
	if result.rowsAffected != 1 {
		t.Errorf("The number of affected rows should be 1 but got %d", result.rowsAffected)
	}
}
