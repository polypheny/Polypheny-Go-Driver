package polypheny

import (
	"context"
	"database/sql/driver"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// PreparedStatement implements driver.Stmt and its context-aware extensions.
type PreparedStatement struct {
	conn   *PolyphenyConn
	id     int32
	args   []ArgSpec
	closed atomic.Bool
}

type ArgSpec struct {
	precision     int32
	scale         int32
	typeName      string
	parameterName string
	name          string
}

// Close releases the prepared statement on the server.
func (stmt *PreparedStatement) Close() error {
	if !stmt.closed.CompareAndSwap(false, true) || !stmt.conn.IsValid() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := stmt.conn.closeStatement(ctx, stmt.id); err != nil {
		_ = stmt.conn.close()
		return err
	}
	return nil
}

func (stmt *PreparedStatement) recordCloseError(operationErr *error) {
	if err := stmt.Close(); *operationErr == nil {
		*operationErr = err
	}
}

// NumInput returns the number of positional parameters.
func (stmt *PreparedStatement) NumInput() int {
	return len(stmt.args)
}

// Exec executes a prepared statement.
func (stmt *PreparedStatement) Exec(args []driver.Value) (driver.Result, error) {
	return stmt.ExecContext(context.Background(), namedValues(args))
}

// ExecContext executes a prepared statement with the supplied parameters.
func (stmt *PreparedStatement) ExecContext(ctx context.Context, args []driver.NamedValue) (result driver.Result, err error) {
	values, err := stmt.parameters(args)
	if err != nil {
		return nil, err
	}
	response, err := stmt.execute(ctx, values)
	if err != nil {
		defer stmt.conn.finishStatement(stmt.id, false, false, &err)
		return nil, err
	}
	if response.GetFrame() != nil {
		defer stmt.conn.finishStatement(stmt.id, false, !response.GetFrame().GetIsLast(), &err)
	}
	return helperExtractResultFromStatementResult(response)
}

// Query executes a prepared query.
func (stmt *PreparedStatement) Query(args []driver.Value) (driver.Rows, error) {
	return stmt.QueryContext(context.Background(), namedValues(args))
}

// QueryContext executes a prepared query with the supplied parameters.
func (stmt *PreparedStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	values, err := stmt.parameters(args)
	if err != nil {
		return nil, err
	}
	response, err := stmt.execute(ctx, values)
	paginated := response.GetFrame() != nil && !response.GetFrame().GetIsLast()
	defer stmt.conn.finishStatement(stmt.id, false, paginated, &err)
	if err != nil {
		return nil, err
	}
	return stmt.conn.readRows(ctx, stmt.id, response.GetFrame())
}

func (stmt *PreparedStatement) parameters(args []driver.NamedValue) ([]*prism.ProtoValue, error) {
	if stmt.closed.Load() {
		return nil, &ClientError{message: "the prepared statement is closed"}
	}
	if len(args) != stmt.NumInput() {
		return nil, &ClientError{message: fmt.Sprintf("expected %d SQL arguments, got %d", stmt.NumInput(), len(args))}
	}
	return helperConvertNamedvalueToProto(args)
}

func (stmt *PreparedStatement) execute(ctx context.Context, values []*prism.ProtoValue) (*prism.StatementResult, error) {
	response, err := stmt.conn.call(ctx, &prism.Request{
		Type: &prism.Request_ExecuteIndexedStatementRequest{
			ExecuteIndexedStatementRequest: &prism.ExecuteIndexedStatementRequest{
				StatementId: stmt.id,
				Parameters:  &prism.IndexedParameters{Parameters: values},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	result := response.GetStatementResult()
	if result == nil {
		return nil, &ClientError{message: "Prism server did not return a prepared statement result"}
	}
	return result, nil
}
