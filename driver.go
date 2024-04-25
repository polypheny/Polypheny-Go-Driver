package polypheny

import (
	sql "database/sql"
	driver "database/sql/driver"
	context "golang.org/x/net/context"
	strings "strings"
)

// PolyphenyDriver implements driver.Driver and DriverContext interface
type PolyphenyDriver struct{}

func init() {
	sql.Register("polypheny", PolyphenyDriver{})
}

func parseDSN(name string) (string, string, string, error) {
	connectionStrings := strings.Split(name, ",")
	if len(connectionStrings) != 2 {
		return "", "", "", &ClientError{
			message: "Connection string (DSN) should have a format like addr:port,user:password",
		}
	}
	addr := connectionStrings[0]
	user := strings.Split(connectionStrings[1], ":")
	if len(connectionStrings) != 2 {
		return "", "", "", &ClientError{
			message: "Connection string (DSN) should have a format like addr:port,user:password",
		}
	}
	username := user[0]
	password := user[1]
	return addr, username, password, nil
}

// Open will return a PolyphenyConn which implements driver.Conn and represents a Connection to Polypheny server
// The name, also called DSN, has the following format: addr:port,user:password
func (d PolyphenyDriver) Open(name string) (driver.Conn, error) {
	addr, username, password, err := parseDSN(name)
	if err != nil {
		return nil, err
	}
	connector := Connector{
		address:  addr,
		username: username,
		password: password,
	}
	return connector.Connect(context.Background())
}

// OpenConnector will return a Connector
func (d PolyphenyDriver) OpenConnector(name string) (driver.Connector, error) {
	addr, username, password, err := parseDSN(name)
	if err != nil {
		return nil, err
	}
	return &Connector{
		address:  addr,
		username: username,
		password: password,
	}, nil
}
