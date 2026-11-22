package polypheny

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
	"google.golang.org/protobuf/proto"
)

func protocolPipe(t *testing.T) (*PolyphenyConn, *PolyphenyConn) {
	t.Helper()
	client, server := net.Pipe()
	for _, socket := range []net.Conn{client, server} {
		if err := socket.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = socket.Close() })
	}
	connection := &PolyphenyConn{netConn: client}
	connection.isConnected.Store(statusPolyphenyConnected)
	return connection, &PolyphenyConn{netConn: server}
}

func TestReceiveFragmentedPrismFrame(t *testing.T) {
	client, server := protocolPipe(t)
	payload := []byte("fragmented Prism message")
	frame := binary.LittleEndian.AppendUint64(nil, uint64(len(payload)))
	frame = append(frame, payload...)
	written := make(chan error, 1)
	go func() {
		for _, value := range frame {
			if _, err := server.netConn.Write([]byte{value}); err != nil {
				written <- err
				return
			}
		}
		written <- nil
	}()

	got, err := client.recv(8)
	if err != nil || !bytes.Equal(got, payload) {
		t.Fatalf("recv returned %q, %v; want %q", got, err, payload)
	}
	if err := <-written; err != nil {
		t.Fatal(err)
	}
}

func TestSendReportsLengthWriteError(t *testing.T) {
	client, server := protocolPipe(t)
	if err := server.netConn.Close(); err != nil {
		t.Fatal(err)
	}

	if err := client.send([]byte("request"), 8); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("send error is %v; want io.ErrClosedPipe", err)
	}
}

func TestCancelFramedRequest(t *testing.T) {
	client, server := protocolPipe(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	read := make(chan error, 1)
	go func() {
		_, err := server.recv(8)
		read <- err
		cancel()
	}()

	if err := client.Ping(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Ping error is %v; want context.Canceled", err)
	}
	if err := <-read; err != nil {
		t.Fatal(err)
	}
	if client.IsValid() {
		t.Fatal("a canceled in-flight connection must not be reused")
	}
}

type protocolStep struct {
	match   func(*prism.Request) bool
	replies []*prism.Response
}

func replyTo(match func(*prism.Request) bool, replies ...*prism.Response) protocolStep {
	return protocolStep{match: match, replies: replies}
}

func serveProtocol(t *testing.T, steps ...protocolStep) *PolyphenyConn {
	t.Helper()
	client, server := protocolPipe(t)
	finished := make(chan error, 1)
	go func() {
		defer server.netConn.Close()
		finished <- runProtocolSteps(server, steps)
	}()
	t.Cleanup(func() {
		if err := <-finished; err != nil {
			t.Error(err)
		}
	})
	return client
}

func runProtocolSteps(server *PolyphenyConn, steps []protocolStep) error {
	for index, step := range steps {
		data, err := server.recv(8)
		if err != nil {
			return fmt.Errorf("read Prism request at step %d: %w", index+1, err)
		}
		request := new(prism.Request)
		if err := proto.Unmarshal(data, request); err != nil {
			return fmt.Errorf("decode Prism request at step %d: %w", index+1, err)
		}
		if !step.match(request) {
			return fmt.Errorf("unexpected Prism request at step %d: %v", index+1, request)
		}
		for i, reply := range step.replies {
			if reply.Id == 0 {
				reply.Id = request.GetId()
			}
			reply.Last = i == len(step.replies)-1
			data, err := proto.Marshal(reply)
			if err != nil {
				return fmt.Errorf("encode Prism reply at step %d: %w", index+1, err)
			}
			if err := server.send(data, 8); err != nil {
				return fmt.Errorf("write Prism reply at step %d: %w", index+1, err)
			}
		}
	}
	return nil
}

func serverError(message string, details *prism.ErrorDetails) *prism.Response {
	return &prism.Response{Type: &prism.Response_ErrorResponse{
		ErrorResponse: &prism.ErrorResponse{Message: message, Details: details},
	}}
}

func assertClientError(t *testing.T, err error, message, state string, code int32, hasCode bool) {
	t.Helper()
	var clientErr *ClientError
	if !errors.As(err, &clientErr) {
		t.Fatalf("error = %T %v, want *ClientError", err, err)
	}
	if got := err.Error(); got != message {
		t.Errorf("Error() = %q, want %q", got, message)
	}
	if got := clientErr.SQLState(); got != state {
		t.Errorf("SQLState() = %q, want %q", got, state)
	}
	if got, present := clientErr.ErrorCode(); got != code || present != hasCode {
		t.Errorf("ErrorCode() = (%d, %t), want (%d, %t)", got, present, code, hasCode)
	}
	if errors.Is(err, driver.ErrBadConn) {
		t.Errorf("error %q must not trigger driver.ErrBadConn retry", err)
	}
}

func TestStructuredServerErrorDetails(t *testing.T) {
	const message = "  request failed\nSQLSTATE 23505, code 17  "
	for _, tc := range []struct {
		name    string
		details *prism.ErrorDetails
		state   string
		code    int32
		hasCode bool
	}{
		{name: "legacy"},
		{name: "empty details", details: &prism.ErrorDetails{}},
		{
			name:    "state only",
			details: &prism.ErrorDetails{State: proto.String("XX000"), Message: proto.String("different details message")},
			state:   "XX000",
		},
		{
			name:    "code only",
			details: &prism.ErrorDetails{ErrorCode: proto.Int32(-17)},
			code:    -17,
			hasCode: true,
		},
		{
			name:    "state and explicit zero",
			details: &prism.ErrorDetails{State: proto.String("22012"), ErrorCode: proto.Int32(0)},
			state:   "22012",
			hasCode: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := serveProtocol(t, replyTo(isPing, serverError(message, tc.details)))
			assertClientError(t, client.Ping(context.Background()), message, tc.state, tc.code, tc.hasCode)
			if !client.IsValid() {
				t.Fatal("a complete server error response should preserve the connection")
			}
		})
	}
}

func TestMismatchedResponseIDHasNoServerDetails(t *testing.T) {
	reply := serverError("unrelated error", &prism.ErrorDetails{State: proto.String("22012"), ErrorCode: proto.Int32(17)})
	reply.Id = 2
	client := serveProtocol(t, replyTo(isPing, reply))
	assertClientError(t, client.Ping(context.Background()), "Prism response ID 2 does not match request ID 1", "", 0, false)
	if client.IsValid() {
		t.Fatal("a connection with mismatched response IDs must not be reused")
	}
}

func isPing(request *prism.Request) bool     { return request.GetConnectionCheckRequest() != nil }
func isRollback(request *prism.Request) bool { return request.GetRollbackRequest() != nil }

func closeStatementStep(statementID int32) protocolStep {
	return replyTo(func(r *prism.Request) bool { return r.GetCloseStatementRequest().GetStatementId() == statementID },
		&prism.Response{Type: &prism.Response_CloseStatementResponse{CloseStatementResponse: &prism.CloseStatementResponse{}}})
}

func rollbackStep() protocolStep {
	return replyTo(isRollback, &prism.Response{Type: &prism.Response_RollbackResponse{RollbackResponse: &prism.RollbackResponse{}}})
}

type protocolConnector struct {
	client   *PolyphenyConn
	attempts atomic.Int32
}

func (c *protocolConnector) Connect(context.Context) (driver.Conn, error) {
	if c.attempts.Add(1) != 1 {
		return nil, errors.New("unexpected additional driver connection")
	}
	return c.client, nil
}

func (c *protocolConnector) Driver() driver.Driver { return PolyphenyDriver{} }

func TestSQLStatementErrorAfterAcknowledgment(t *testing.T) {
	const (
		statementID = 42
		query       = "bad SQL"
		message     = "invalid SQL statement"
	)
	client := serveProtocol(t,
		replyTo(func(r *prism.Request) bool {
			q := r.GetExecuteUnparameterizedStatementRequest()
			return q.GetLanguageName() == "sql" && q.GetStatement() == query
		}, &prism.Response{Type: &prism.Response_StatementResponse{StatementResponse: &prism.StatementResponse{StatementId: statementID}}},
			serverError(message, &prism.ErrorDetails{State: proto.String("42601"), ErrorCode: proto.Int32(0)})),
		closeStatementStep(statementID),
		rollbackStep(),
		replyTo(isPing, &prism.Response{Type: &prism.Response_ConnectionCheckResponse{ConnectionCheckResponse: &prism.ConnectionCheckResponse{}}}),
		replyTo(func(r *prism.Request) bool { return r.GetDisconnectRequest() != nil },
			&prism.Response{Type: &prism.Response_DisconnectResponse{DisconnectResponse: &prism.DisconnectResponse{}}}),
	)
	connector := &protocolConnector{client: client}
	db := sql.OpenDB(connector)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	result, err := db.ExecContext(context.Background(), query)
	if result != nil {
		t.Errorf("ExecContext returned %v, want nil result", result)
	}
	assertClientError(t, err, message, "42601", 0, true)
	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("ping after statement cleanup: %v", err)
	}
	if got := connector.attempts.Load(); got != 1 {
		t.Errorf("database/sql made %d connection attempts, want 1", got)
	}
}

func TestFailedImplicitRollbackDiscardsConnection(t *testing.T) {
	const message = "could not start statement"
	client := serveProtocol(t,
		replyTo(func(r *prism.Request) bool { return r.GetExecuteUnparameterizedStatementRequest() != nil },
			serverError(message, &prism.ErrorDetails{State: proto.String("42000"), ErrorCode: proto.Int32(7)})),
		replyTo(isRollback, serverError("could not roll back", &prism.ErrorDetails{State: proto.String("XX000")})),
	)

	_, err := client.ExecContext(context.Background(), "SELECT invalid", nil)
	assertClientError(t, err, message, "42000", 7, true)
	if client.IsValid() {
		t.Fatal("connection with a failed implicit rollback must not be reused")
	}
}
