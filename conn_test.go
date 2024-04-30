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
	if conn.IsValid() != false {
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
	conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
	conn.(*PolyphenyConn).isConnected.Store(statusDisconnected)
	conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != driver.ErrBadConn {
		t.Errorf("Expected to receive an ErrBadConn but got %s", err.Error())
	}
	conn.(*PolyphenyConn).isConnected.Store(statusServerConnected)
	conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err.Error() != "Ping: invalid connection to Polypheny server" {
		t.Errorf("Expected to receive a ClientError but got %s", err.Error())
	}
	conn.(*PolyphenyConn).isConnected.Store(statusPolyphenyConnected)
	conn.(*PolyphenyConn).pingInternal(errChan)
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
