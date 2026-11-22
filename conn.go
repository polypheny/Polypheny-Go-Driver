package polypheny

import (
	"context"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
	"google.golang.org/protobuf/proto"
)

// The 8-byte length prefix can represent more than we can safely allocate.
const maxFrameSize = 256 << 20

// PolyphenyConn implements driver.Conn.
type PolyphenyConn struct {
	address     string
	username    string
	netConn     net.Conn
	isConnected atomic.Int32
	inTx        atomic.Bool

	mu            sync.Mutex
	nextRequestID uint64
}

// IsValid reports whether the connection can be reused by database/sql.
func (conn *PolyphenyConn) IsValid() bool {
	return conn.isConnected.Load() == statusPolyphenyConnected
}

// Ping checks that the server responds on this connection.
func (conn *PolyphenyConn) Ping(ctx context.Context) error {
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{ConnectionCheckRequest: &prism.ConnectionCheckRequest{}},
	})
	if err != nil {
		return err
	}
	if response.GetConnectionCheckResponse() == nil {
		return &ClientError{message: "Prism server did not acknowledge the connection check"}
	}
	return nil
}

// Prepare prepares a SQL statement on the server.
func (conn *PolyphenyConn) Prepare(query string) (driver.Stmt, error) {
	return conn.PrepareContext(context.Background(), query)
}

// PrepareContext prepares a SQL statement with positional parameters.
func (conn *PolyphenyConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return conn.prepare(ctx, query)
}

func (conn *PolyphenyConn) prepare(ctx context.Context, query string) (*PreparedStatement, error) {
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_PrepareIndexedStatementRequest{
			PrepareIndexedStatementRequest: &prism.PrepareStatementRequest{
				LanguageName: "sql",
				Statement:    query,
			},
		},
	})
	if err != nil {
		if ctx.Err() == nil {
			conn.finishStatement(0, true, false, &err)
		}
		return nil, err
	}
	signature := response.GetPreparedStatementSignature()
	if signature == nil {
		_ = conn.close()
		return nil, &ClientError{message: "Prism server did not return a prepared statement signature"}
	}
	if !conn.inTx.Load() {
		// Prism validates prepared SQL in a transaction even in auto-commit mode.
		// The prepared statement remains available after that transaction ends.
		if err := conn.endTransaction(ctx, true); err != nil {
			_ = conn.close()
			return nil, err
		}
	}
	args := make([]ArgSpec, 0, len(signature.GetParameterMetas()))
	for _, parameter := range signature.GetParameterMetas() {
		args = append(args, ArgSpec{
			precision:     parameter.GetPrecision(),
			scale:         parameter.GetScale(),
			typeName:      parameter.GetTypeName(),
			parameterName: parameter.GetParameterName(),
			name:          parameter.GetName(),
		})
	}
	return &PreparedStatement{conn: conn, id: signature.GetStatementId(), args: args}, nil
}

// Close disconnects from the server and releases the network connection.
func (conn *PolyphenyConn) Close() error {
	if !conn.IsValid() {
		return conn.close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_DisconnectRequest{DisconnectRequest: &prism.DisconnectRequest{}},
	})
	closeErr := conn.close()
	if err != nil {
		return err
	}
	return closeErr
}

// Begin starts a transaction.
func (conn *PolyphenyConn) Begin() (driver.Tx, error) {
	return conn.BeginTx(context.Background(), driver.TxOptions{})
}

// BeginTx starts a transaction with the server's default isolation level.
func (conn *PolyphenyConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if opts.Isolation != 0 || opts.ReadOnly {
		return nil, &ClientError{message: "Polypheny supports only the default transaction options"}
	}
	if !conn.inTx.CompareAndSwap(false, true) {
		return nil, &ClientError{message: "a transaction is already active"}
	}
	if err := conn.setAutoCommit(ctx, false); err != nil {
		conn.inTx.Store(false)
		return nil, err
	}
	return &PolyphenyTransaction{conn: conn}, nil
}

// Exec executes a SQL statement.
func (conn *PolyphenyConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return conn.ExecContext(context.Background(), query, namedValues(args))
}

// ExecContext executes a SQL statement with positional arguments.
func (conn *PolyphenyConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (result driver.Result, err error) {
	if len(args) != 0 {
		statement, prepareErr := conn.prepare(ctx, query)
		if prepareErr != nil {
			return nil, prepareErr
		}
		defer statement.recordCloseError(&err)
		return statement.ExecContext(ctx, args)
	}
	statementID, statementResult, err := conn.execute(ctx, "sql", query, nil)
	if statementID != 0 {
		defer conn.finishStatement(statementID, true, false, &err)
	}
	if err != nil {
		return nil, err
	}
	return helperExtractResultFromStatementResult(statementResult)
}

// Query executes a SQL query.
func (conn *PolyphenyConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return conn.QueryContext(context.Background(), query, namedValues(args))
}

// QueryContext executes a SQL query with positional arguments.
func (conn *PolyphenyConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	if len(args) != 0 {
		statement, prepareErr := conn.prepare(ctx, query)
		if prepareErr != nil {
			return nil, prepareErr
		}
		defer statement.recordCloseError(&err)
		return statement.QueryContext(ctx, args)
	}
	statementID, result, err := conn.execute(ctx, "sql", query, nil)
	paginated := result.GetFrame() != nil && !result.GetFrame().GetIsLast()
	if statementID != 0 {
		defer conn.finishStatement(statementID, true, paginated, &err)
	}
	if err != nil {
		return nil, err
	}
	return conn.readRows(ctx, statementID, result.GetFrame())
}

// execute receives the statement ID followed by the statement result in one
// exchange, so a second request cannot consume either message.
func (conn *PolyphenyConn) execute(ctx context.Context, language, query string, namespace *string) (int32, *prism.StatementResult, error) {
	responses, err := conn.exchange(ctx, &prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName:  language,
				Statement:     query,
				NamespaceName: namespace,
			},
		},
	})
	var statementID int32
	if len(responses) > 0 {
		statementID = responses[0].GetStatementResponse().GetStatementId()
	}
	if err != nil {
		if statementID == 0 && ctx.Err() == nil {
			conn.finishStatement(0, true, false, &err)
		}
		return statementID, nil, err
	}
	last := responses[len(responses)-1].GetStatementResponse()
	if last == nil || statementID <= 0 || last.GetStatementId() != statementID || last.GetResult() == nil {
		_ = conn.close()
		return statementID, nil, &ClientError{message: "Prism server returned an invalid statement response"}
	}
	return statementID, last.GetResult(), nil
}

func (conn *PolyphenyConn) fetch(ctx context.Context, statementID int32) (*prism.Frame, error) {
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_FetchRequest{
			FetchRequest: &prism.FetchRequest{StatementId: statementID},
		},
	})
	if err != nil {
		return nil, err
	}
	frame := response.GetFrame()
	if frame == nil {
		return nil, &ClientError{message: "Prism server did not return the next result frame"}
	}
	return frame, nil
}

func (conn *PolyphenyConn) closeStatement(ctx context.Context, statementID int32) error {
	_, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_CloseStatementRequest{
			CloseStatementRequest: &prism.CloseStatementRequest{StatementId: statementID},
		},
	})
	return err
}

func (conn *PolyphenyConn) finishStatement(statementID int32, temporary, paginated bool, operationErr *error) {
	if !conn.IsValid() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if statementID != 0 {
		if temporary {
			err = conn.closeStatement(ctx, statementID)
		} else {
			_, err = conn.call(ctx, &prism.Request{
				Type: &prism.Request_CloseResultRequest{
					CloseResultRequest: &prism.CloseResultRequest{StatementId: statementID},
				},
			})
		}
		if err != nil {
			if *operationErr == nil {
				*operationErr = err
			}
			_ = conn.close()
			return
		}
	}
	if conn.inTx.Load() || !conn.IsValid() {
		return
	}
	if *operationErr != nil {
		// Clear the implicit transaction left by a failed statement.
		if err := conn.endTransaction(ctx, false); err != nil {
			_ = conn.close()
		}
	} else if paginated {
		// Prism auto-commits an initial final frame, but not subsequent fetches.
		if err := conn.endTransaction(ctx, true); err != nil {
			*operationErr = err
			_ = conn.close()
		}
	}
}

func (conn *PolyphenyConn) send(serialized []byte, lengthSize int) error {
	if lengthSize < 1 || lengthSize > 8 {
		return &ClientError{message: "invalid Prism length field size"}
	}
	if len(serialized) > maxFrameSize || uint64(len(serialized)) >= 1<<uint(lengthSize*8-1) {
		return &ClientError{message: "Prism message is too large"}
	}
	var length [8]byte
	binary.LittleEndian.PutUint64(length[:], uint64(len(serialized)))
	if err := writeAll(conn.netConn, length[:lengthSize]); err != nil {
		return err
	}
	return writeAll(conn.netConn, serialized)
}

func writeAll(w io.Writer, data []byte) error {
	for len(data) != 0 {
		n, err := w.Write(data)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		data = data[n:]
	}
	return nil
}

func (conn *PolyphenyConn) recv(lengthSize int) ([]byte, error) {
	if lengthSize < 1 || lengthSize > 8 {
		return nil, &ClientError{message: "invalid Prism length field size"}
	}
	var length [8]byte
	if _, err := io.ReadFull(conn.netConn, length[:lengthSize]); err != nil {
		return nil, err
	}
	n := binary.LittleEndian.Uint64(length[:])
	if n > maxFrameSize || n >= 1<<uint(lengthSize*8-1) {
		return nil, &ClientError{message: "Prism message is too large"}
	}
	buf := make([]byte, n)
	_, err := io.ReadFull(conn.netConn, buf)
	return buf, err
}

func (conn *PolyphenyConn) close() error {
	if conn.isConnected.Swap(statusDisconnected) == statusDisconnected || conn.netConn == nil {
		return nil
	}
	return conn.netConn.Close()
}

// withDeadline interrupts an in-flight exchange when the context is canceled.
// A connection interrupted mid-frame cannot be safely returned to the pool.
func (conn *PolyphenyConn) withDeadline(ctx context.Context, operation func() error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if ctx.Done() == nil {
		return operation()
	}

	interruptDone := make(chan struct{})
	stop := context.AfterFunc(ctx, func() {
		_ = conn.netConn.SetDeadline(time.Now())
		close(interruptDone)
	})
	err := operation()
	if !stop() {
		<-interruptDone
		_ = conn.close()
		return ctx.Err()
	}
	return err
}

func (conn *PolyphenyConn) exchange(ctx context.Context, request *prism.Request) ([]*prism.Response, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.isConnected.Load() == statusDisconnected {
		return nil, driver.ErrBadConn
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	conn.nextRequestID++
	request.Id = conn.nextRequestID
	buf, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}
	var responses []*prism.Response
	err = conn.withDeadline(ctx, func() error {
		if err := conn.send(buf, 8); err != nil {
			_ = conn.close()
			return err
		}
		for {
			buf, err := conn.recv(8)
			if err != nil {
				_ = conn.close()
				return err
			}
			response := new(prism.Response)
			if err := proto.Unmarshal(buf, response); err != nil {
				_ = conn.close()
				return err
			}
			if response.GetId() != request.Id {
				_ = conn.close()
				return &ClientError{message: fmt.Sprintf("Prism response ID %d does not match request ID %d", response.GetId(), request.Id)}
			}
			if serverErr := response.GetErrorResponse(); serverErr != nil {
				clientErr := &ClientError{message: serverErr.GetMessage()}
				if !response.GetLast() {
					_ = conn.close()
					return clientErr
				}
				if details := serverErr.GetDetails(); details != nil {
					clientErr.sqlState = details.GetState()
					clientErr.errorCode = details.GetErrorCode()
					clientErr.hasErrorCode = details.ErrorCode != nil
				}
				return clientErr
			}
			responses = append(responses, response)
			if response.GetLast() {
				return nil
			}
		}
	})
	return responses, err
}

func (conn *PolyphenyConn) call(ctx context.Context, request *prism.Request) (*prism.Response, error) {
	responses, err := conn.exchange(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(responses) != 1 {
		return nil, &ClientError{message: "Prism server returned multiple responses for a single-response request"}
	}
	return responses[0], nil
}
