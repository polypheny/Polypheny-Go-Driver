package polypheny

import (
	context "context"
	driver "database/sql/driver"
	atomic "sync/atomic"
	testing "testing"
)

func TestIsValid(t *testing.T) {
	conn := PolyphenyConn{
		address:     "",
		username:    "",
		netConn:     nil,
		isConnected: atomic.Int32{},
	}
	conn.isConnected.Store(statusDisconnected)
	if conn.IsValid() != false {
		t.Errorf("Failed to set connection status, the status is %d, but should be %d.", conn.isConnected.Load(), statusDisconnected)
	}
	conn.isConnected.Store(statusServerConnected)
	if conn.IsValid() != false {
		t.Errorf("Failed to set connection status, the status is %d, but should be %d.", conn.isConnected.Load(), statusServerConnected)
	}
	conn.isConnected.Store(statusPolyphenyConnected)
	if conn.IsValid() == false {
		t.Errorf("Failed to set connection status, the status is %d, but should be %d.", conn.isConnected.Load(), statusPolyphenyConnected)
	}
}

func TestPingInternal(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	errChan := make(chan error)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
	conn.(*PolyphenyConn).isConnected.Store(statusDisconnected)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != driver.ErrBadConn {
		t.Errorf("Expected to receive an ErrBadConn but got %s", err.Error())
	}
	conn.(*PolyphenyConn).isConnected.Store(statusServerConnected)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err.Error() != "Ping: invalid connection to Polypheny server" {
		t.Errorf("Expected to receive a ClientError but got %s", err.Error())
	}
	conn.(*PolyphenyConn).isConnected.Store(statusPolyphenyConnected)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).Ping(context.Background())
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = conn.(*PolyphenyConn).Ping(ctx)
	if err.Error() != "Context cancelled or timeout" {
		t.Errorf("Expected to receive a ClientError error but got %s", err.Error())
	}
}

func TestPrepare(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	stmt, err := conn.(*PolyphenyConn).Prepare("sql:SELECT * FROM emps WHERE name = ?")
	if err != nil {
		t.Error(err)
	}
	if len(stmt.(*PreparedStatement).args) != 1 {
		t.Errorf("Expected 1 arg in the prepared statement, but got %d", len(stmt.(*PreparedStatement).args))
	}
	stmt, err = conn.(*PolyphenyConn).Prepare("sql:SELECT * FROM emps WHERE name = ? AND salary = ?")
	if err != nil {
		t.Error(err)
	}
	if len(stmt.(*PreparedStatement).args) != 2 {
		t.Errorf("Expected 2 args in the prepared statement, but got %d", len(stmt.(*PreparedStatement).args))
	}
}

func TestClose(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).close()
	if err != nil {
		t.Error(err)
	}
	if conn.(*PolyphenyConn).isConnected.Load() != statusDisconnected {
		t.Errorf("The connection status should be %d, but got %d", statusDisconnected, conn.(*PolyphenyConn).isConnected.Load())
	}
}

func TestBegin(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	tx, err := conn.(*PolyphenyConn).Begin()
	if err != nil {
		t.Error(err)
	}
	if tx.(*PolyphenyTranaction).conn != conn {
		t.Error("Error when trys to start a transaction")
	}
}

func TestBeginTx(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	tx, err := conn.(*PolyphenyConn).BeginTx(context.Background(), driver.TxOptions{})
	if err != nil {
		t.Error(err)
	}
	if tx.(*PolyphenyTranaction).conn != conn {
		t.Error("Error when trys to start a transaction")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = conn.(*PolyphenyConn).BeginTx(ctx, driver.TxOptions{})
	if err != ctx.Err() {
		t.Error("context didn't work")
	}
}
