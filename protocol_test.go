package polypheny

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strings"
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

func TestStatementErrorAfterAcknowledgment(t *testing.T) {
	client, server := protocolPipe(t)
	finished := make(chan error, 1)
	go func() {
		data, err := server.recv(8)
		if err != nil {
			finished <- err
			return
		}
		request := new(prism.Request)
		if err := proto.Unmarshal(data, request); err != nil {
			finished <- err
			return
		}

		responses := []*prism.Response{
			{
				Id: request.GetId(),
				Type: &prism.Response_StatementResponse{
					StatementResponse: &prism.StatementResponse{StatementId: 42},
				},
			},
			{
				Id:   request.GetId(),
				Last: true,
				Type: &prism.Response_ErrorResponse{
					ErrorResponse: &prism.ErrorResponse{Message: "invalid SQL statement"},
				},
			},
		}
		for _, response := range responses {
			data, err := proto.Marshal(response)
			if err != nil {
				finished <- err
				return
			}
			if err := server.send(data, 8); err != nil {
				finished <- err
				return
			}
		}
		finished <- nil
	}()

	id, result, err := client.execute(context.Background(), "sql", "bad SQL", nil)
	if id != 42 || result != nil || err == nil || !strings.Contains(err.Error(), "invalid SQL statement") {
		t.Fatalf("execute returned statement %d, result %v, error %v", id, result, err)
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
	if !client.IsValid() {
		t.Fatal("a complete server error response should preserve the connection")
	}
}

func TestFailedImplicitRollbackDiscardsConnection(t *testing.T) {
	client, server := protocolPipe(t)
	finished := make(chan error, 1)
	go func() {
		for _, reply := range []struct {
			check   func(*prism.Request) bool
			message string
		}{
			{func(r *prism.Request) bool { return r.GetExecuteUnparameterizedStatementRequest() != nil }, "could not start statement"},
			{func(r *prism.Request) bool { return r.GetRollbackRequest() != nil }, "could not roll back"},
		} {
			data, err := server.recv(8)
			if err != nil {
				finished <- err
				return
			}
			request := new(prism.Request)
			if err := proto.Unmarshal(data, request); err != nil {
				finished <- err
				return
			}
			if !reply.check(request) {
				finished <- errors.New("unexpected Prism request during rollback")
				return
			}
			data, err = proto.Marshal(&prism.Response{
				Id:   request.GetId(),
				Last: true,
				Type: &prism.Response_ErrorResponse{ErrorResponse: &prism.ErrorResponse{Message: reply.message}},
			})
			if err != nil {
				finished <- err
				return
			}
			if err := server.send(data, 8); err != nil {
				finished <- err
				return
			}
		}
		finished <- nil
	}()

	_, err := client.ExecContext(context.Background(), "SELECT invalid", nil)
	if err == nil || err.Error() != "could not start statement" {
		t.Fatalf("ExecContext error is %v; want original statement error", err)
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
	if client.IsValid() {
		t.Fatal("connection with a failed implicit rollback must not be reused")
	}
}
