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
	if err.Error() != ctx.Err().Error() {
		t.Error(err)
	}
}
