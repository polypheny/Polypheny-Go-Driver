package polypheny

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"net"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// Connector implements driver.Connector.
type Connector struct {
	address  string
	username string
	password string
}

// Connect connects to a Polypheny server using the Prism TCP transport.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	var dialer net.Dialer
	netConn, err := dialer.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return nil, err
	}

	conn := &PolyphenyConn{address: c.address, username: c.username, netConn: netConn}
	conn.isConnected.Store(statusServerConnected)
	defer func() {
		if !conn.IsValid() {
			_ = conn.close()
		}
	}()

	if err := conn.withDeadline(ctx, func() error {
		version, err := conn.recv(1)
		if err != nil {
			return err
		}
		expected := []byte(transportVersion)
		if !bytes.Equal(version, expected) {
			return &ClientError{message: "the Prism transport version is incompatible with the server"}
		}
		return conn.send(expected, 1)
	}); err != nil {
		return nil, err
	}

	autoCommit := true
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &conn.username,
				Password:        &c.password,
				ConnectionProperties: &prism.ConnectionProperties{
					IsAutoCommit: &autoCommit,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	version := response.GetConnectionResponse()
	if version == nil {
		return nil, &ClientError{message: "Prism server did not return connection properties"}
	}
	if !version.GetIsCompatible() {
		return nil, &ClientError{
			message: fmt.Sprintf("Prism API %d.%d is incompatible with server Prism API %d.%d",
				majorApiVersion, minorApiVersion, version.GetMajorApiVersion(), version.GetMinorApiVersion()),
		}
	}
	conn.isConnected.Store(statusPolyphenyConnected)
	return conn, nil
}

// Driver returns the underlying Polypheny driver.
func (c *Connector) Driver() driver.Driver {
	return PolyphenyDriver{}
}
