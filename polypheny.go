package polypheny

type Connection struct {
	address string
	isOpen bool
	client *protoClient
}

func Connect(address string) *Connection {
	client := connect(address)
	conn := Connection {
		address: address,
		isOpen: true,
		client: client,
	}
	return &conn
}

func (conn *Connection) Execute(statement string, language string) {
	_ = conn.client.handleExecuteUnprepared(statement, language)
}

func (conn *Connection) Fetch() [][]interface{} {
	return conn.client.handleFetchResult()
}

func (conn *Connection) Commit() {
	conn.client.handleCommitRequest()
}

func (conn *Connection) Close() {
	conn.client.close()
}
