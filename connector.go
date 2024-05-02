package polypheny

import (
	bytes "bytes"
	driver "database/sql/driver"
	net "net"
	"sync/atomic"

	context "golang.org/x/net/context"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
)

// Connector implements the driver.Connector interface
type Connector struct {
	// TODO: add ConnectionProperties support
	address  string
	username string
	password string
}

// Connect connects to a Polypheny server and returns PolyphenyConn which implements driver.Conn
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	// Step 1, dial to the server
	// TODO: does polypheny also use other networks?
	var d net.Dialer
	netConn, err := d.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return nil, err
	}
	conn := PolyphenyConn{
		address:     c.address,
		username:    c.username,
		netConn:     netConn,
		isConnected: atomic.Int32{},
	}
	conn.isConnected.Store(statusServerConnected)
	// Step 2, exchange transport version
	recvVersion, err := conn.recv(1)
	if err != nil {
		return nil, err
	}
	sendVersion := []byte(transportVersion)
	err = conn.send(sendVersion, 1)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(recvVersion, sendVersion) {
		return nil, &ClientError{
			message: "The transport version is incompatible with server",
		}
	}
	// Step 3, send prism connection request
	request := prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &conn.username,
				Password:        &c.password,
			},
		},
	}
	_, err = conn.helperSendAndRecv(&request)
	if err != nil {
		return nil, err
	}
	conn.isConnected.Store(statusPolyphenyConnected)
	return &conn, nil
}

// Driver retuens a PolyphenyDriver
func (c *Connector) Driver() driver.Driver {
	return PolyphenyDriver{}
}
