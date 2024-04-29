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

// IsValid implements Validator interface
func (conn *PolyphenyConn) IsValid() bool {
	// don't think its necessary to do a connection check here
	return conn.isConnected.Load() == statusPolyphenyConnected
}

// ping pings a connection and sending the result via the error channel
func (conn *PolyphenyConn) pingInternal(err chan error) {
	status := conn.isConnected.Load()
	if status == statusDisconnected {
		err <- driver.ErrBadConn
	} else if status == statusServerConnected {
		err <- &ClientError{
			message: "Ping: invalid connection to Polypheny server",
		}
	}
	request := prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{
			ConnectionCheckRequest: &prism.ConnectionCheckRequest{},
		},
	}
	_, connErr := conn.helperSendAndRecv(&request)
	if connErr != nil {
		err <- driver.ErrBadConn
	}
	err <- nil
}

// Ping implmenets Pinger
func (conn *PolyphenyConn) Ping(ctx context.Context) error {
	errChan := make(chan error)
	go conn.pingInternal(errChan)
	var err error
	select {
	case <-ctx.Done():
		// context timeout or cancelled
		// TODO: or we return ErrBadConn?
		return &ClientError{
			message: "Context cancelled or timeout",
		}
	case err = <-errChan:
		return err
	}
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
// TODO: add support to ConnBeginTx
// Deprecated
func (conn *PolyphenyConn) Begin() (driver.Tx, error) {
	return &PolyphenyTranaction{
		conn: conn,
	}, nil
}

// BeginTx starts a transaction
func (conn *PolyphenyConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &PolyphenyTranaction{
			conn: conn,
		}, nil
	}
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

// execContextInternal is called by ExecContext, if the ctx is timeout or cancelled, ExecContext will return without waiting execContextInternal
//
// TODO add args support
func (conn *PolyphenyConn) execContextInternal(query string, resultChan chan *PolyphenyResult, errChan chan error) {
	queryLanguage, queryBody, err := parseQuery(query)
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
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
		resultChan <- nil
		errChan <- err
		return
	}
	requestID := response.GetStatementResponse().GetStatementId()
	buf, err := conn.recv(8) // this is the query result
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	err = proto.Unmarshal(buf, response)
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	// is this an error?
	if requestID != response.GetStatementResponse().GetStatementId() {
		resultChan <- nil
		errChan <- nil
		return
	}
	result, err := helperExtractResultFromStatementResult(response.GetStatementResponse().GetResult())
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	} else {
		resultChan <- result.(*PolyphenyResult)
		errChan <- err
		return
	}
}

// ExecContext executes a query that doesn't return data under Context
//
// TODO add fetch size, namespace and args support.
// TODO for args support, can we first prepare it and then exec?
func (conn *PolyphenyConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	errChan := make(chan error)
	resultChan := make(chan *PolyphenyResult)
	var err error
	var result *PolyphenyResult
	go conn.execContextInternal(query, resultChan, errChan)
	select {
	case <-ctx.Done():
		// context timeout or cancelled
		return nil, ctx.Err()
	case result = <-resultChan:
		err = <-errChan
		return result, err
	}
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

// queryContextInternal is called by ExecContext, if the ctx is timeout or cancelled, QueryContext will return without waiting queryContextInternal
//
// TODO add args support
func (conn *PolyphenyConn) queryContextInternal(query string, rowsChan chan *PolyphenyRows, errChan chan error) {
	queryLanguage, queryBody, err := parseQuery(query)
	if err != nil {
		rowsChan <- nil
		errChan <- err
		return
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
		rowsChan <- nil
		errChan <- err
		return
	}
	requestID := response.GetStatementResponse().GetStatementId()
	buf, err := conn.recv(8) // this is the query result
	if err != nil {
		rowsChan <- nil
		errChan <- err
		return
	}
	err = proto.Unmarshal(buf, response)
	if err != nil {
		rowsChan <- nil
		errChan <- err
		return
	}
	// is this an error?
	if requestID != response.GetStatementResponse().GetStatementId() {
		rowsChan <- nil
		errChan <- nil
		return
	}
	result, err := helperExtractRowsFromStatementResult(response.GetStatementResponse().GetResult())
	if err != nil {
		rowsChan <- nil
		errChan <- err
		return
	} else {
		rowsChan <- result.(*PolyphenyRows)
		errChan <- err
		return
	}
}

// QueryContext executes a query that returns data under Context
//
// TODO add fetch size, namespace and args support.
// TODO for args support, can we first prepare it and then exec?
func (conn *PolyphenyConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	errChan := make(chan error)
	rowsChan := make(chan *PolyphenyRows)
	var err error
	var result *PolyphenyRows
	go conn.queryContextInternal(query, rowsChan, errChan)
	select {
	case <-ctx.Done():
		// context timeout or cancelled
		return nil, ctx.Err()
	case result = <-rowsChan:
		err = <-errChan
		return result, err
	}
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
