package polypheny

import (
	"context"
	"fmt"
	"testing"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
	"google.golang.org/protobuf/proto"
)

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
	client, finished := serveCypherPages(t, first, next, false)
	result, err := client.ExecuteCypherContext(context.Background(), "sample", "MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	if err := <-finished; err != nil {
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
	client, finished := serveCypherPages(t, first, next, true)
	result, err := client.ExecuteCypherContext(context.Background(), "sample", "MATCH (n) RETURN n")
	if result != nil || err == nil || err.Error() != "expected a Prism graph result frame" {
		t.Fatalf("unexpected result %v, error %v", result, err)
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
	if !client.IsValid() {
		t.Fatal("expected the connection to remain usable after rollback and close")
	}
}

func serveCypherPages(t *testing.T, first, next *prism.Frame, rollback bool) (*PolyphenyConn, <-chan error) {
	t.Helper()
	client, server := protocolPipe(t)
	finished := make(chan error, 1)
	go func() {
		const statementID = 42
		type step struct {
			match func(*prism.Request) bool
			reply *prism.Response
		}
		steps := []step{
			{
				match: func(r *prism.Request) bool {
					q := r.GetExecuteUnparameterizedStatementRequest()
					return q.GetLanguageName() == "cypher" && q.GetNamespaceName() == "sample" && q.GetStatement() == "MATCH (n) RETURN n"
				},
				reply: &prism.Response{Type: &prism.Response_StatementResponse{
					StatementResponse: &prism.StatementResponse{StatementId: statementID, Result: &prism.StatementResult{Frame: first}},
				}},
			},
			{
				match: func(r *prism.Request) bool { return r.GetFetchRequest().GetStatementId() == statementID },
				reply: &prism.Response{Type: &prism.Response_Frame{Frame: next}},
			},
		}
		if rollback {
			steps = append(steps, step{
				match: func(r *prism.Request) bool { return r.GetRollbackRequest() != nil },
				reply: &prism.Response{Type: &prism.Response_RollbackResponse{RollbackResponse: &prism.RollbackResponse{}}},
			})
		}
		steps = append(steps, step{
			match: func(r *prism.Request) bool { return r.GetCloseStatementRequest().GetStatementId() == statementID },
			reply: &prism.Response{Type: &prism.Response_CloseStatementResponse{CloseStatementResponse: &prism.CloseStatementResponse{}}},
		})
		for i, step := range steps {
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
			if !step.match(request) {
				finished <- fmt.Errorf("unexpected Cypher request at step %d: %v", i, request)
				return
			}
			step.reply.Id, step.reply.Last = request.GetId(), true
			data, err = proto.Marshal(step.reply)
			if err == nil {
				err = server.send(data, 8)
			}
			if err != nil {
				finished <- err
				return
			}
		}
		finished <- nil
	}()
	return client, finished
}
