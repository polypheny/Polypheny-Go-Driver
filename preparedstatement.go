package polypheny

import (
	driver "database/sql/driver"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
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
