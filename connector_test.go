package polypheny

import (
	context "context"
	testing "testing"
)

func TestConnect(t *testing.T) {
	d := PolyphenyDriver{}
	connector, err := d.OpenConnector("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	conn, err := connector.(*Connector).Connect(context.Background())
	if err != nil || conn.(*PolyphenyConn).isConnected.Load() != statusPolyphenyConnected {
		t.Error("Failed to connect")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = connector.(*Connector).Connect(ctx)
	if err.Error() != "dial tcp: lookup localhost: operation was canceled" {
		t.Error(err)
	}
}

func TestDriver(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	connector.Driver()
}
