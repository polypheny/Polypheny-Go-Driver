package polypheny

import (
	context "context"
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
	args[0] = PolyphenyConn{}
	_, err = stmt.(*PreparedStatement).Exec(args)
	if err == nil {
		t.Error("Expecting error")
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
	args[0] = driver.NamedValue{
		Ordinal: 1,
		Value:   PolyphenyConn{},
	}
	go stmt.(*PreparedStatement).execContextInternal(args, resultChan, errChan)
	//_ = <-resultChan
	err = <-errChan
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestStmtExecContext(t *testing.T) {
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
	args := make([]driver.NamedValue, 1)
	args[0] = driver.NamedValue{
		Ordinal: 1,
		Value:   int32(1),
	}
	result, err := stmt.(*PreparedStatement).ExecContext(context.Background(), args)
	if err != nil {
		t.Error(err, result == nil)
	}
	if result.(*PolyphenyResult).rowsAffected != 1 {
		t.Errorf("The number of affected rows should be 1 but got %d", result.(*PolyphenyResult).rowsAffected)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = stmt.(*PreparedStatement).ExecContext(ctx, args)
	if err != ctx.Err() {
		t.Error(err)
	}
}

func TestStmtQuery(t *testing.T) {
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
	args := make([]driver.Value, 1)
	args[0] = "Bill"
	rows, err := stmt.(*PreparedStatement).Query(args)
	if err != nil {
		t.Error(err)
	}
	if len(rows.(*PolyphenyRows).Columns()) != 0 && rows.(*PolyphenyRows).Columns()[0] != "name" {
		t.Error("Error in Query")
	}
}

func TestStmtQueryInternal(t *testing.T) {
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
	args := make([]driver.NamedValue, 1)
	args[0] = driver.NamedValue{
		Ordinal: 1,
		Value:   "Bill",
	}
	errChan := make(chan error)
	rowsChan := make(chan *PolyphenyRows)
	var result *PolyphenyRows
	go stmt.(*PreparedStatement).queryContextInternal(args, rowsChan, errChan)
	result = <-rowsChan
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("Error in queryContextInternal, result should not be nil")
	} else if len(result.Columns()) != 0 && result.Columns()[0] != "name" {
		t.Error("Error in queryContextInternal")
	}
}

func TestStmtQueryContext(t *testing.T) {
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
	args := make([]driver.NamedValue, 1)
	args[0] = driver.NamedValue{
		Ordinal: 1,
		Value:   "Bill",
	}
	_, err = stmt.(*PreparedStatement).QueryContext(context.Background(), args)
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = stmt.(*PreparedStatement).QueryContext(ctx, args)
	if err != ctx.Err() {
		t.Error(err)
	}
}
