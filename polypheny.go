package polypheny

import (
	protoclient "github.com/polypheny/Go-Driver/protoclient"
)

type Connection struct {
	address string
	isOpen bool
	client *protoclient.ProtoClient
}

func Connect(address string) *Connection {
	client := protoclient.Connect(address)
	conn := Connection {
		address: address,
		isOpen: true,
		client: client,
	}
	return &conn
}

func (conn *Connection) Execute(statement string, language string) {
	_ = conn.client.ExecuteUnprepared(statement, language)
}

func (conn *Connection) Fetch() [][]interface{} {
	return conn.client.FetchResult()
}

func (conn *Connection) Commit() {
	conn.client.Commit()
}

func (conn *Connection) Close() {
	conn.client.Close()
}
