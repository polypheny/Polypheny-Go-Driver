package polypheny

import (
	context "context"
	driver "database/sql/driver"

	prism "github.com/polypheny/Polypheny-Go-Driver/prism"
)

// PreparedStatement implments Stmt, StmtExecContext and StmtQueryContext
// Deprecated functions, such as Exec and Query, are also implemented for compatibility, but they are deprecated.
type PreparedStatement struct {
	conn *PolyphenyConn
	id   int32
	args []ArgSpec
}

type ArgSpec struct {
	precision     int32
	scale         int32
	typeName      string
	parameterName string
	name          string
}

// Close closes the statement
func (stmt *PreparedStatement) Close() error {
	request := prism.Request{
		Type: &prism.Request_CloseStatementRequest{
			CloseStatementRequest: &prism.CloseStatementRequest{
				StatementId: stmt.id,
			},
		},
	}
	_, err := stmt.conn.helperSendAndRecv(&request)
	return err
}

// NumInput returns the number of parameters
func (stmt *PreparedStatement) NumInput() int {
	return len(stmt.args)
}

// Exec executes a query that doesn't return data
// TODO: add support to FetchSize
// Deprecated
func (stmt *PreparedStatement) Exec(args []driver.Value) (driver.Result, error) {
	pvs, err := helperConvertValueToProto(args)
	if err != nil {
		return nil, err
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteIndexedStatementRequest{
			ExecuteIndexedStatementRequest: &prism.ExecuteIndexedStatementRequest{
				StatementId: stmt.id,
				Parameters: &prism.IndexedParameters{
					Parameters: pvs,
				},
				FetchSize: nil,
			},
		},
	}
	response, err := stmt.conn.helperSendAndRecv(&request)
	if err != nil {
		return nil, err
	}
	return helperExtractResultFromStatementResult(response.GetStatementResult())
}

func (stmt *PreparedStatement) execContextInternal(args []driver.NamedValue, resultChan chan *PolyphenyResult, errChan chan error) {
	pvs, err := helperConvertNamedvalueToProto(args)
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteIndexedStatementRequest{
			ExecuteIndexedStatementRequest: &prism.ExecuteIndexedStatementRequest{
				StatementId: stmt.id,
				Parameters: &prism.IndexedParameters{
					Parameters: pvs,
				},
				FetchSize: nil,
			},
		},
	}
	response, err := stmt.conn.helperSendAndRecv(&request)
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	result, err := helperExtractResultFromStatementResult(response.GetStatementResult())
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
func (stmt *PreparedStatement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	errChan := make(chan error)
	resultChan := make(chan *PolyphenyResult)
	var err error
	var result *PolyphenyResult
	go stmt.execContextInternal(args, resultChan, errChan)
	select {
	case <-ctx.Done():
		// context timeout or cancelled
		return nil, ctx.Err()
	case result = <-resultChan:
		err = <-errChan
		return result, err
	}
}

// Query executes a query that returns data
// TODO: add support to FetchSize
// Deprecated
func (stmt *PreparedStatement) Query(args []driver.Value) (driver.Rows, error) {
	pvs, err := helperConvertValueToProto(args)
	if err != nil {
		return nil, err
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteIndexedStatementRequest{
			ExecuteIndexedStatementRequest: &prism.ExecuteIndexedStatementRequest{
				StatementId: stmt.id,
				Parameters: &prism.IndexedParameters{
					Parameters: pvs,
				},
				FetchSize: nil,
			},
		},
	}
	response, err := stmt.conn.helperSendAndRecv(&request)
	if err != nil {
		return nil, err
	}
	return helperExtractRowsFromStatementResult(response.GetStatementResult())
}

func (stmt *PreparedStatement) queryContextInternal(args []driver.NamedValue, rowsChan chan *PolyphenyRows, errChan chan error) {
	pvs, err := helperConvertNamedvalueToProto(args)
	if err != nil {
		rowsChan <- nil
		errChan <- err
		return
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteIndexedStatementRequest{
			ExecuteIndexedStatementRequest: &prism.ExecuteIndexedStatementRequest{
				StatementId: stmt.id,
				Parameters: &prism.IndexedParameters{
					Parameters: pvs,
				},
				FetchSize: nil,
			},
		},
	}
	response, err := stmt.conn.helperSendAndRecv(&request)
	if err != nil {
		rowsChan <- nil
		errChan <- err
		return
	}
	result, err := helperExtractRowsFromStatementResult(response.GetStatementResult())
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
func (stmt *PreparedStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	errChan := make(chan error)
	rowsChan := make(chan *PolyphenyRows)
	var err error
	var result *PolyphenyRows
	go stmt.queryContextInternal(args, rowsChan, errChan)
	select {
	case <-ctx.Done():
		// context timeout or cancelled
		return nil, ctx.Err()
	case result = <-rowsChan:
		err = <-errChan
		return result, err
	}
}
