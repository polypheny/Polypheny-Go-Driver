package polypheny

import (
	"context"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// ExecuteCypherContext runs one Cypher statement in the specified graph
// namespace. An empty namespace uses the default namespace of the connection.
// Writes can return a scalar update count. For queries, the result contains a
// single final frame with all pages in server order. Graph-only projections
// return graph elements. Tabular projections retain column metadata and typed
// Prism values, including nodes, edges, paths and graphs.
//
// Call this method inside (*sql.Conn).Raw. As with other driver.Conn
// operations, do not use the connection concurrently.
func (conn *PolyphenyConn) ExecuteCypherContext(ctx context.Context, namespace, query string) (*prism.StatementResult, error) {
	return conn.ExecuteCypherWithParametersContext(ctx, namespace, query, nil)
}

// ExecuteCypherWithParametersContext runs one Cypher statement with named
// Prism parameters. Map keys are the parameter names without the leading dollar
// sign. A nil map uses unparameterized execution.
//
// Call this method inside (*sql.Conn).Raw, like ExecuteCypherContext.
func (conn *PolyphenyConn) ExecuteCypherWithParametersContext(ctx context.Context, namespace, query string, parameters map[string]*prism.ProtoValue) (result *prism.StatementResult, err error) {
	var namespaceName *string
	if namespace != "" {
		namespaceName = &namespace
	}
	var statementID int32
	if parameters == nil {
		statementID, result, err = conn.execute(ctx, "cypher", query, namespaceName)
	} else {
		statementID, result, err = conn.executeNamedCypher(ctx, namespaceName, query, parameters)
	}
	if statementID != 0 {
		defer conn.finishCypherStatement(statementID, &err)
	}
	if err != nil {
		return nil, err
	}
	frame := result.GetFrame()
	if frame == nil {
		return result, nil
	}
	if err := conn.collectCypherFrames(ctx, statementID, frame); err != nil {
		return nil, err
	}
	frame.IsLast = true
	return result, nil
}

func (conn *PolyphenyConn) executeNamedCypher(ctx context.Context, namespace *string, query string, parameters map[string]*prism.ProtoValue) (int32, *prism.StatementResult, error) {
	for name, parameter := range parameters {
		if name == "" || parameter.GetValue() == nil {
			return 0, nil, &ClientError{message: "each Cypher parameter needs a name and a Prism value"}
		}
	}
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_PrepareNamedStatementRequest{
			PrepareNamedStatementRequest: &prism.PrepareStatementRequest{
				LanguageName:  "cypher",
				Statement:     query,
				NamespaceName: namespace,
			},
		},
	})
	if err != nil {
		return 0, nil, err
	}
	signature := response.GetPreparedStatementSignature()
	if signature == nil || signature.GetStatementId() <= 0 {
		_ = conn.close()
		return 0, nil, &ClientError{message: "Prism server did not return a valid prepared statement signature"}
	}
	id := signature.GetStatementId()
	response, err = conn.call(ctx, &prism.Request{
		Type: &prism.Request_ExecuteNamedStatementRequest{
			ExecuteNamedStatementRequest: &prism.ExecuteNamedStatementRequest{
				StatementId: id,
				Parameters:  &prism.NamedParameters{Parameters: parameters},
			},
		},
	})
	if err != nil {
		return id, nil, err
	}
	result := response.GetStatementResult()
	if result == nil {
		_ = conn.close()
		return id, nil, &ClientError{message: "Prism server did not return a prepared statement result"}
	}
	return id, result, nil
}

func (conn *PolyphenyConn) collectCypherFrames(ctx context.Context, statementID int32, first *prism.Frame) error {
	switch result := first.GetResult().(type) {
	case *prism.Frame_GraphFrame:
		graph := result.GraphFrame
		return conn.readFrames(ctx, statementID, first, func(page *prism.Frame) error {
			next := page.GetGraphFrame()
			if next == nil {
				return &ClientError{message: "expected a Prism graph result frame"}
			}
			if page != first {
				graph.Element = append(graph.Element, next.GetElement()...)
			}
			return nil
		})
	case *prism.Frame_RelationalFrame:
		rows := result.RelationalFrame
		return conn.readFrames(ctx, statementID, first, func(page *prism.Frame) error {
			next := page.GetRelationalFrame()
			if next == nil {
				return &ClientError{message: "expected a Prism relational result frame"}
			}
			if page != first {
				rows.Rows = append(rows.Rows, next.GetRows()...)
			}
			return nil
		})
	default:
		return &ClientError{message: "expected a Prism graph or relational result frame"}
	}
}

func (conn *PolyphenyConn) finishCypherStatement(statementID int32, operationErr *error) {
	if !conn.IsValid() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Closing a partially fetched Cypher statement can commit it. Roll back
	// a failed auto-commit operation before closing the statement.
	if *operationErr != nil && !conn.inTx.Load() {
		if err := conn.endTransaction(ctx, false); err != nil {
			_ = conn.close()
			return
		}
	}
	if err := conn.closeStatement(ctx, statementID); err != nil {
		if *operationErr == nil {
			*operationErr = err
		}
		_ = conn.close()
	}
}
