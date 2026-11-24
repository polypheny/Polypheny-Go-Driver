package polypheny

import (
	"context"
	"testing"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
	"google.golang.org/protobuf/proto"
)

func TestCypherNamedParametersAndRelationalFormat(t *testing.T) {
	const (
		id    = 17
		query = "RETURN $1 AS value"
	)
	parameter := &prism.ProtoValue{Value: &prism.ProtoValue_Integer{Integer: &prism.ProtoInteger{Integer: 42}}}
	frame := &prism.Frame{IsLast: true, Result: &prism.Frame_RelationalFrame{RelationalFrame: &prism.RelationalFrame{
		ColumnMeta: []*prism.ColumnMeta{{ColumnLabel: "value"}},
		Rows:       []*prism.Row{{Values: []*prism.ProtoValue{parameter}}},
	}}}
	client := serveProtocol(t,
		replyTo(func(r *prism.Request) bool {
			return r.GetSetClientInfoPropertiesRequest().GetProperties()["CypherResultFormat"] == "relational"
		}, &prism.Response{Type: &prism.Response_SetClientInfoPropertiesResponse{SetClientInfoPropertiesResponse: &prism.ClientInfoPropertiesResponse{}}}),
		replyTo(func(r *prism.Request) bool {
			prepare := r.GetPrepareNamedStatementRequest()
			return prepare.GetLanguageName() == "cypher" && prepare.GetNamespaceName() == "sample" && prepare.GetStatement() == query
		}, &prism.Response{Type: &prism.Response_PreparedStatementSignature{PreparedStatementSignature: &prism.PreparedStatementSignature{StatementId: id}}}),
		replyTo(func(r *prism.Request) bool {
			execute := r.GetExecuteNamedStatementRequest()
			return execute.GetStatementId() == id && proto.Equal(execute.GetParameters().GetParameters()["1"], parameter)
		}, &prism.Response{Type: &prism.Response_StatementResult{StatementResult: &prism.StatementResult{Frame: frame}}}),
		closeStatementStep(id),
	)
	ctx := context.Background()
	if err := client.SetClientInfoPropertiesContext(ctx, map[string]string{"CypherResultFormat": "relational"}); err != nil {
		t.Fatal(err)
	}
	result, err := client.ExecuteCypherWithParametersContext(ctx, "sample", query, map[string]*prism.ProtoValue{"1": parameter})
	if err != nil || !proto.Equal(result.GetFrame(), frame) {
		t.Fatalf("ExecuteCypherWithParametersContext returned %v, %v", result, err)
	}
}

func TestCypherPreservesTypedGraphCellsAcrossPages(t *testing.T) {
	node := &prism.ProtoNode{Id: "node", Labels: []string{"Person"}}
	edge := &prism.ProtoEdge{Id: "edge", Left: "node", Right: "other", Direction: prism.ProtoEdge_LEFT_TO_RIGHT}
	other := &prism.ProtoNode{Id: "other"}
	path := &prism.ProtoPath{Nodes: []*prism.ProtoNode{node, other}, Edges: []*prism.ProtoEdge{edge}}
	whole := &prism.ProtoGraph{Nodes: []*prism.ProtoNode{node, other, {Id: "isolated"}}, Edges: []*prism.ProtoEdge{edge}}
	graphCells := []*prism.ProtoValue{
		{Value: &prism.ProtoValue_Node{Node: node}},
		{Value: &prism.ProtoValue_Edge{Edge: edge}},
		{Value: &prism.ProtoValue_Path{Path: path}},
		{Value: &prism.ProtoValue_Graph{Graph: whole}},
	}
	item := &prism.ProtoValue{Value: &prism.ProtoValue_List{List: &prism.ProtoList{Values: graphCells}}}
	document := &prism.ProtoDocument{Entries: map[string]*prism.ProtoValue{"elements": item}}
	firstRow := &prism.Row{Values: []*prism.ProtoValue{{Value: &prism.ProtoValue_Document{Document: document}}}}
	secondRow := &prism.Row{Values: []*prism.ProtoValue{{Value: &prism.ProtoValue_Null{Null: &prism.ProtoNull{}}}}}
	first := &prism.Frame{Result: &prism.Frame_RelationalFrame{RelationalFrame: &prism.RelationalFrame{
		ColumnMeta: []*prism.ColumnMeta{{ColumnLabel: "graph_values"}},
		Rows:       []*prism.Row{firstRow},
	}}}
	next := &prism.Frame{IsLast: true, Result: &prism.Frame_RelationalFrame{RelationalFrame: &prism.RelationalFrame{
		Rows: []*prism.Row{secondRow},
	}}}
	client := serveCypherFetch(t, first, &prism.Response{Type: &prism.Response_Frame{Frame: next}}, false)
	result, err := client.ExecuteCypherContext(context.Background(), "sample", "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	first.IsLast = true
	first.GetRelationalFrame().Rows = append(first.GetRelationalFrame().Rows, secondRow)
	if !proto.Equal(result.GetFrame(), first) {
		t.Fatalf("got %v, want %v", result.GetFrame(), first)
	}
}

func TestCypherRollsBackBeforeClosingUnexpectedPage(t *testing.T) {
	first := &prism.Frame{Result: &prism.Frame_GraphFrame{GraphFrame: &prism.GraphFrame{}}}
	next := &prism.Frame{Result: &prism.Frame_RelationalFrame{RelationalFrame: &prism.RelationalFrame{}}}
	client := serveCypherFetch(t, first, &prism.Response{Type: &prism.Response_Frame{Frame: next}}, true)
	result, err := client.ExecuteCypherContext(context.Background(), "sample", "MATCH (n) RETURN n")
	if result != nil || err == nil || err.Error() != "expected a Prism graph result frame" {
		t.Fatalf("unexpected result %v, error %v", result, err)
	}
	if !client.IsValid() {
		t.Fatal("expected the connection to remain usable after rollback and close")
	}
}

func TestCypherFetchPreservesServerError(t *testing.T) {
	const message = "cannot fetch the next Cypher result page"
	first := &prism.Frame{Result: &prism.Frame_GraphFrame{GraphFrame: &prism.GraphFrame{}}}
	next := serverError(message, &prism.ErrorDetails{State: proto.String("XX000"), ErrorCode: proto.Int32(17)})
	client := serveCypherFetch(t, first, next, true)
	result, err := client.ExecuteCypherContext(context.Background(), "sample", "MATCH (n) RETURN n")
	if result != nil {
		t.Errorf("ExecuteCypherContext returned %v, want nil result", result)
	}
	assertClientError(t, err, message, "XX000", 17, true)
	if !client.IsValid() {
		t.Fatal("the connection should remain usable after rollback and close")
	}
}

func serveCypherFetch(t *testing.T, first *prism.Frame, next *prism.Response, rollback bool) *PolyphenyConn {
	t.Helper()
	const statementID = 42
	steps := []protocolStep{
		replyTo(func(r *prism.Request) bool {
			q := r.GetExecuteUnparameterizedStatementRequest()
			return q.GetLanguageName() == "cypher" && q.GetNamespaceName() == "sample" && q.GetStatement() == "MATCH (n) RETURN n"
		}, &prism.Response{Type: &prism.Response_StatementResponse{
			StatementResponse: &prism.StatementResponse{StatementId: statementID, Result: &prism.StatementResult{Frame: first}},
		}}),
		replyTo(func(r *prism.Request) bool { return r.GetFetchRequest().GetStatementId() == statementID }, next),
	}
	if rollback {
		steps = append(steps, rollbackStep())
	}
	return serveProtocol(t, append(steps, closeStatementStep(statementID))...)
}
