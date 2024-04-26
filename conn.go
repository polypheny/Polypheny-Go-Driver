package polypheny

import (
	context "context"
	"database/sql/driver"
	binary "encoding/binary"
	math "math"
	net "net"
	atomic "sync/atomic"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"

	proto "google.golang.org/protobuf/proto"
)

// PolyphenyConn implements driver.Conn
type PolyphenyConn struct {
	address     string       // addr:port
	username    string       // username is stored, but password is not
	netConn     net.Conn     // the Conn struct returned by Dial
	isConnected atomic.Int32 // Connection status
}

// Ping a connection
// TODO: shall we add context cancel and timeout support to Ping too?
func (conn *PolyphenyConn) Ping(ctx context.Context) error {
	status := conn.isConnected.Load()
	if status == statusDisconnected {
		return driver.ErrBadConn
	} else if status == statusServerConnected {
		return &ClientError{
			message: "Ping: invalid connection to Polypheny server",
		}
	}
	request := prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{
			ConnectionCheckRequest: &prism.ConnectionCheckRequest{},
		},
	}
	_, err := conn.helperSendAndRecv(&request)
	if err != nil {
		// TODO: is there any other reasons causing the error?
		return driver.ErrBadConn
	}
	return nil
}

// Prepare a statement
func (conn *PolyphenyConn) Prepare(query string) (driver.Stmt, error) {
	queryLanguage, queryBody, err := parseQuery(query)
	if err != nil {
		return nil, err
	}
	request := prism.Request{
		Type: &prism.Request_PrepareIndexedStatementRequest{
			PrepareIndexedStatementRequest: &prism.PrepareStatementRequest{
				LanguageName:  queryLanguage,
				Statement:     queryBody,
				NamespaceName: nil,
			},
		},
	}
	response, err := conn.helperSendAndRecv(&request)
	if err != nil {
		return nil, err
	}
	var args []ArgSpec
	for _, parameter := range response.GetPreparedStatementSignature().GetParameterMetas() {
		args = append(args, ArgSpec{
			precision:     parameter.GetPrecision(),
			scale:         parameter.GetScale(),
			typeName:      parameter.GetTypeName(),
			parameterName: parameter.GetParameterName(),
			name:          parameter.GetName(),
		})
	}
	return &PreparedStatement{
		conn: conn,
		id:   response.GetPreparedStatementSignature().GetStatementId(),
		args: args,
	}, nil
}

// Close will close the network connection to Polypheny server
// TODO: add support to timeout
func (conn *PolyphenyConn) Close() error {
	request := prism.Request{
		Type: &prism.Request_DisconnectRequest{
			DisconnectRequest: &prism.DisconnectRequest{},
		},
	}
	_, err := conn.helperSendAndRecv(&request)
	if err != nil {
		return err
	}
	conn.isConnected.Store(statusServerConnected)
	err = conn.close()
	return err
}

// Begin starts a new transaction
// TODO: does prism interface currently have a BeginTransactionRequest?
// TODO: add support to ConnBeginTx
// Deprecated
func (conn *PolyphenyConn) Begin() (driver.Tx, error) {
	return nil, nil
}

// Exec executes a query that doesn't return data
// TODO: add fetch size, namespace and args support.
// TODO: for args support, can we first prepare it and then exec?
// Deprecated
func (conn *PolyphenyConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	queryLanguage, queryBody, err := parseQuery(query)
	if err != nil {
		return nil, err
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName:  queryLanguage,
				Statement:     queryBody,
				FetchSize:     nil,
				NamespaceName: nil,
			},
		},
	}
	response, err := conn.helperSendAndRecv(&request)
	if err != nil {
		return nil, err
	}
	requestID := response.GetStatementResponse().GetStatementId()
	buf, err := conn.recv(8) // this is the query result
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(buf, response)
	if err != nil {
		return nil, err
	}
	// is this an error?
	if requestID != response.GetStatementResponse().GetStatementId() {
		return nil, nil
	}
	return helperExtractResultFromStatementResult(response.GetStatementResponse().GetResult())
}

// Exec executes a query that doesn't return data
// TODO: add fetch size, namespace and args support.
// TODO: for args support, can we first prepare it and then exec?
// Deprecated
func (conn *PolyphenyConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	queryLanguage, queryBody, err := parseQuery(query)
	if err != nil {
		return nil, err
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName:  queryLanguage,
				Statement:     queryBody,
				FetchSize:     nil,
				NamespaceName: nil,
			},
		},
	}
	response, err := conn.helperSendAndRecv(&request)
	if err != nil {
		return nil, err
	}
	requestID := response.GetStatementResponse().GetStatementId()
	buf, err := conn.recv(8) // this is the query result
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(buf, response)
	if err != nil {
		return nil, err
	}
	// is this an error?
	if requestID != response.GetStatementResponse().GetStatementId() {
		return nil, nil
	}
	return helperExtractRowsFromStatementResult(response.GetStatementResponse().GetResult())
}

// Many requests to the server have a similar pattern
// Client first sends the length of the message, and then the message itself
// TODO: Currently only little endian is supported
func (conn *PolyphenyConn) send(serialized []byte, lengthSize int) error {
	if lengthSize > 8 {
		return &ClientError{
			message: "PolyphenyConn.send() expects the lengthSize parameter is not greater than 8",
		}
	}
	lengthBuf := make([]byte, 8)
	length := make([]byte, lengthSize)
	lenSerialized := len(serialized)
	if uint64(lenSerialized) >= uint64(math.Pow(2, float64(lengthSize*8-1))) {
		return &ClientError{
			message: "PolyphenyConn.send(): the size of the serialized message is too large to be put in a byte array of size lengthSize",
		}
	}
	binary.LittleEndian.PutUint64(lengthBuf, uint64(lenSerialized))
	for i := range length {
		length[i] = lengthBuf[i]
	}
	_, err := conn.netConn.Write(length)
	if err != nil {
		return nil
	}
	_, err = conn.netConn.Write(serialized)
	return err
}

// Many responses from the server also have a similar pattern
// Server first sents the length of a message and then the message itself
// TODO: Currently only little endian is supported
func (conn *PolyphenyConn) recv(lengthSize int) ([]byte, error) {
	if lengthSize > 8 {
		return nil, &ClientError{
			message: "PolyphenyConn.recv() expects the lengthSize parameter is not greater than 8",
		}
	}
	lengthBuf := make([]byte, 8)
	length := make([]byte, lengthSize)
	_, err := conn.netConn.Read(length)
	if err != nil {
		return nil, err
	}
	for i := range lengthBuf {
		if i < lengthSize {
			lengthBuf[i] = length[i]
		} else {
			lengthBuf[i] = 0
		}
	}
	recvLength := binary.LittleEndian.Uint64(lengthBuf)
	buf := make([]byte, recvLength)
	_, err = conn.netConn.Read(buf)
	return buf, err
}

func (conn *PolyphenyConn) close() error {
	err := conn.netConn.Close()
	if err != nil {
		return err
	}
	conn.isConnected.Store(statusDisconnected) //TODO: maybe add an error status
	return nil
}

// helperSendAndRecv is a helper function which serialize and send a request, then returns the responses
func (conn *PolyphenyConn) helperSendAndRecv(m proto.Message) (*prism.Response, error) {
	buf, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = conn.send(buf, 8)
	if err != nil {
		return nil, err
	}
	buf, err = conn.recv(8)
	if err != nil {
		return nil, err
	}
	var response prism.Response
	err = proto.Unmarshal(buf, &response)
	return &response, err
}
