package polypheny

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

func executeCypherIntegration(t *testing.T, ctx context.Context, db *sql.DB, namespace, query string) *prism.StatementResult {
	t.Helper()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var result *prism.StatementResult
	if err := conn.Raw(func(raw any) error {
		var err error
		result, err = raw.(*PolyphenyConn).ExecuteCypherContext(ctx, namespace, query)
		return err
	}); err != nil {
		t.Fatalf("Cypher %q: %v", query, err)
	}
	if result == nil {
		t.Fatal("Cypher did not return a statement result")
	}
	return result
}

func graphFrameIntegration(t *testing.T, result *prism.StatementResult) []*prism.GraphElement {
	t.Helper()
	frame := result.GetFrame()
	if frame == nil || !frame.GetIsLast() || frame.GetGraphFrame() == nil {
		t.Fatalf("expected a collected final graph frame, got %v", result)
	}
	return frame.GetGraphFrame().GetElement()
}

func relationalFrameIntegration(t *testing.T, result *prism.StatementResult, labels ...string) *prism.RelationalFrame {
	t.Helper()
	frame := result.GetFrame()
	if frame == nil || !frame.GetIsLast() || frame.GetRelationalFrame() == nil {
		t.Fatalf("expected a collected final relational frame, got %v", result)
	}
	page := frame.GetRelationalFrame()
	var got []string
	for _, column := range page.GetColumnMeta() {
		got = append(got, column.GetColumnLabel())
	}
	if !reflect.DeepEqual(got, labels) {
		t.Fatalf("column labels = %v, want %v", got, labels)
	}
	for _, row := range page.GetRows() {
		if len(row.GetValues()) != len(labels) {
			t.Fatalf("result row has %d values, want %d", len(row.GetValues()), len(labels))
		}
	}
	return page
}

func graphIntegerIntegration(t *testing.T, value *prism.ProtoValue) int64 {
	t.Helper()
	switch v := value.GetValue().(type) {
	case *prism.ProtoValue_Integer:
		return int64(v.Integer.GetInteger())
	case *prism.ProtoValue_Long:
		return v.Long.GetLong()
	case *prism.ProtoValue_String_:
		n, err := strconv.ParseInt(v.String_.GetString_(), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		return n
	default:
		t.Fatalf("expected integer graph property, got %v", value)
		return 0
	}
}

func assertGraphNodeIntegration(t *testing.T, node *prism.ProtoNode, name string, id int64) {
	t.Helper()
	if node == nil || node.GetId() == "" {
		t.Fatalf("missing node identity: %v", node)
	}
	if !reflect.DeepEqual(node.GetLabels(), []string{"Person"}) {
		t.Fatalf("node labels = %v", node.GetLabels())
	}
	if node.GetProperties()["name"].GetString_().GetString_() != name {
		t.Fatalf("node name property = %v, want %q", node.GetProperties()["name"], name)
	}
	if got := graphIntegerIntegration(t, node.GetProperties()["id"]); got != id {
		t.Fatalf("node id property = %d, want %d", got, id)
	}
}

func assertGraphEdgeIntegration(t *testing.T, edge *prism.ProtoEdge, from, to *prism.ProtoNode, since int64) {
	t.Helper()
	if edge == nil || edge.GetId() == "" {
		t.Fatalf("missing edge identity: %v", edge)
	}
	if !reflect.DeepEqual(edge.GetLabels(), []string{"KNOWS"}) {
		t.Fatalf("edge labels = %v", edge.GetLabels())
	}
	if edge.GetLeft() != from.GetId() || edge.GetRight() != to.GetId() || edge.GetDirection() != prism.ProtoEdge_LEFT_TO_RIGHT {
		t.Fatalf("unexpected edge endpoints or direction: %v", edge)
	}
	if got := graphIntegerIntegration(t, edge.GetProperties()["since"]); got != since {
		t.Fatalf("edge since property = %d, want %d", got, since)
	}
}

func assertCypherWriteIntegration(t *testing.T, result *prism.StatementResult) {
	t.Helper()
	if result.GetFrame() == nil {
		if result.GetScalar() < 0 {
			t.Fatalf("negative graph update count: %v", result)
		}
		return
	}
	// Some Cypher plans expose their update count as a tabular ROWCOUNT.
	page := relationalFrameIntegration(t, result, "ROWCOUNT")
	if len(page.GetRows()) != 1 || graphIntegerIntegration(t, page.GetRows()[0].GetValues()[0]) < 0 {
		t.Fatalf("unexpected graph update result: %v", result)
	}
}

func TestIntegrationCypher(t *testing.T) {
	db := integrationDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	namespace := fmt.Sprintf("go_graph_%x", time.Now().UnixNano())
	if _, err := db.ExecContext(ctx, "CREATE GRAPH NAMESPACE "+namespace); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if _, err := db.ExecContext(cleanupCtx, "DROP NAMESPACE IF EXISTS "+namespace); err != nil {
			t.Errorf("drop graph fixture: %v", err)
		}
	})

	created := executeCypherIntegration(t, ctx, db, namespace,
		"CREATE (:Person {id:1,name:'Alice'})-[:KNOWS {since:2020}]->"+
			"(:Person {id:2,name:'Bob'})-[:KNOWS {since:2021}]->"+
			"(:Person {id:3,name:'Charlie'}), (:Isolated {name:'alone'})")
	if created.GetFrame() != nil {
		t.Fatalf("Cypher write returned a frame: %v", created)
	}
	if created.GetScalar() < 0 {
		t.Fatalf("negative update count: %v", created)
	}

	t.Run("flat nodes edges and paths", func(t *testing.T) {
		result := executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (a:Person {name:'Alice'})-[r:KNOWS]->(b:Person) RETURN a,r,b")
		elements := graphFrameIntegration(t, result)
		if len(elements) != 3 {
			t.Fatalf("graph element count = %d, want 3", len(elements))
		}
		alice, bob := elements[0].GetNode(), elements[2].GetNode()
		assertGraphNodeIntegration(t, alice, "Alice", 1)
		assertGraphNodeIntegration(t, bob, "Bob", 2)
		assertGraphEdgeIntegration(t, elements[1].GetEdge(), alice, bob, 2020)

		path := graphFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH p=(:Person {name:'Alice'})-[:KNOWS*1..3]->(:Person {name:'Charlie'}) RETURN p"))
		if len(path) != 5 {
			t.Fatalf("path element count = %d, want 5", len(path))
		}
		assertGraphNodeIntegration(t, path[0].GetNode(), "Alice", 1)
		assertGraphNodeIntegration(t, path[2].GetNode(), "Bob", 2)
		assertGraphNodeIntegration(t, path[4].GetNode(), "Charlie", 3)
		assertGraphEdgeIntegration(t, path[1].GetEdge(), path[0].GetNode(), path[2].GetNode(), 2020)
		assertGraphEdgeIntegration(t, path[3].GetEdge(), path[2].GetNode(), path[4].GetNode(), 2021)
	})

	t.Run("whole graph preserves isolated node", func(t *testing.T) {
		elements := graphFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace, "MATCH (*) RETURN *"))
		nodes, edges := map[string]*prism.ProtoNode{}, []*prism.ProtoEdge{}
		for _, element := range elements {
			if node := element.GetNode(); node != nil {
				nodes[node.GetProperties()["name"].GetString_().GetString_()] = node
			} else if edge := element.GetEdge(); edge != nil {
				edges = append(edges, edge)
			} else {
				t.Fatalf("unrecognized graph element: %v", element)
			}
		}
		if len(nodes) != 4 || len(edges) != 2 || nodes["alone"] == nil {
			t.Fatalf("whole graph returned %d nodes and %d edges: %v", len(nodes), len(edges), elements)
		}
	})

	t.Run("mixed nodes and scalars in both orders", func(t *testing.T) {
		names := []string{"Alice", "Bob", "Charlie"}
		for _, nodeFirst := range []bool{true, false} {
			projection, labels, nodeColumn := "n,n.name AS name", []string{"n", "name"}, 0
			if !nodeFirst {
				projection, labels, nodeColumn = "n.name AS name,n", []string{"name", "n"}, 1
			}
			result := executeCypherIntegration(t, ctx, db, namespace,
				"MATCH (n:Person) RETURN "+projection+" ORDER BY name")
			page := relationalFrameIntegration(t, result, labels...)
			if len(page.GetRows()) != len(names) {
				t.Fatalf("mixed row count = %d, want 3", len(page.GetRows()))
			}
			if page.GetColumnMeta()[nodeColumn].GetTypeMeta().GetProtoValueType() != prism.ProtoPolyType_NODE {
				t.Fatalf("missing NODE column metadata: %v", page.GetColumnMeta())
			}
			for i, row := range page.GetRows() {
				values := row.GetValues()
				assertGraphNodeIntegration(t, values[nodeColumn].GetNode(), names[i], int64(i+1))
				if got := values[1-nodeColumn].GetString_().GetString_(); got != names[i] {
					t.Fatalf("mixed scalar = %q, want %q", got, names[i])
				}
			}
		}
	})

	t.Run("scalar projections retain rows and aliases", func(t *testing.T) {
		page := relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Person) RETURN n.name AS name ORDER BY name"), "name")
		var names []string
		for _, row := range page.GetRows() {
			names = append(names, row.GetValues()[0].GetString_().GetString_())
		}
		if !reflect.DeepEqual(names, []string{"Alice", "Bob", "Charlie"}) {
			t.Fatalf("scalar names = %v", names)
		}
	})

	t.Run("mixed edge and path cells", func(t *testing.T) {
		page := relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (a:Person {name:'Alice'})-[r:KNOWS]->(b:Person) RETURN a,r,b,b.name AS name"),
			"a", "r", "b", "name")
		if len(page.GetRows()) != 1 {
			t.Fatalf("edge row count = %d, want 1", len(page.GetRows()))
		}
		values := page.GetRows()[0].GetValues()
		assertGraphNodeIntegration(t, values[0].GetNode(), "Alice", 1)
		assertGraphNodeIntegration(t, values[2].GetNode(), "Bob", 2)
		assertGraphEdgeIntegration(t, values[1].GetEdge(), values[0].GetNode(), values[2].GetNode(), 2020)
		if page.GetColumnMeta()[1].GetTypeMeta().GetProtoValueType() != prism.ProtoPolyType_EDGE {
			t.Fatalf("missing EDGE metadata: %v", page.GetColumnMeta())
		}

		page = relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH p=(:Person {name:'Alice'})-[:KNOWS*1..3]->(:Person {name:'Charlie'}) RETURN p,'route' AS title"),
			"p", "title")
		if len(page.GetRows()) != 1 {
			t.Fatalf("path row count = %d, want 1", len(page.GetRows()))
		}
		path := page.GetRows()[0].GetValues()[0].GetPath()
		if path == nil || len(path.GetNodes()) != 3 || len(path.GetEdges()) != 2 {
			t.Fatalf("path cell lost structure: %v", path)
		}
		assertGraphNodeIntegration(t, path.GetNodes()[0], "Alice", 1)
		assertGraphNodeIntegration(t, path.GetNodes()[2], "Charlie", 3)
		assertGraphEdgeIntegration(t, path.GetEdges()[0], path.GetNodes()[0], path.GetNodes()[1], 2020)
	})

	t.Run("reverse path preserves graph edge orientation", func(t *testing.T) {
		page := relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH p=(:Person {name:'Charlie'})<-[:KNOWS*1..3]-(:Person {name:'Alice'}) RETURN p,'reverse' AS title"),
			"p", "title")
		if len(page.GetRows()) != 1 {
			t.Fatalf("reverse path rows = %d, want 1", len(page.GetRows()))
		}
		path := page.GetRows()[0].GetValues()[0].GetPath()
		if path == nil || len(path.GetNodes()) != 3 || len(path.GetEdges()) != 2 {
			t.Fatalf("reverse path lost structure: %v", path)
		}
		assertGraphNodeIntegration(t, path.GetNodes()[0], "Charlie", 3)
		assertGraphNodeIntegration(t, path.GetNodes()[1], "Bob", 2)
		assertGraphNodeIntegration(t, path.GetNodes()[2], "Alice", 1)
		assertGraphEdgeIntegration(t, path.GetEdges()[0], path.GetNodes()[1], path.GetNodes()[0], 2021)
		assertGraphEdgeIntegration(t, path.GetEdges()[1], path.GetNodes()[2], path.GetNodes()[1], 2020)
	})

	t.Run("empty graph and mixed frames", func(t *testing.T) {
		if got := graphFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Person {name:'missing'}) RETURN n")); len(got) != 0 {
			t.Fatalf("empty graph returned %d elements", len(got))
		}
		page := relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Person {name:'missing'}) RETURN n,n.name AS name"), "n", "name")
		if len(page.GetRows()) != 0 || page.GetColumnMeta()[0].GetTypeMeta().GetProtoValueType() != prism.ProtoPolyType_NODE {
			t.Fatalf("empty mixed frame lost rows or schema: %v", page)
		}
	})

	t.Run("namespace and graph metadata", func(t *testing.T) {
		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if err := conn.Raw(func(raw any) error {
			driverConn := raw.(*PolyphenyConn)
			kind := "GRAPH"
			namespaces, err := driverConn.SearchNamespaces(ctx, &namespace, &kind)
			if err != nil {
				return err
			}
			if len(namespaces.GetNamespaces()) != 1 || namespaces.GetNamespaces()[0].GetNamespaceName() != namespace {
				return fmt.Errorf("unexpected graph namespaces: %v", namespaces)
			}
			for _, filter := range []*string{nil, &namespace} {
				entities, err := driverConn.SearchEntities(ctx, namespace, filter)
				if err != nil {
					return err
				}
				if len(entities.GetEntities()) != 1 || entities.GetEntities()[0].GetGraph() == nil {
					return fmt.Errorf("unexpected graph entities: %v", entities)
				}
			}
			missing := "missing%"
			entities, err := driverConn.SearchEntities(ctx, namespace, &missing)
			if err != nil {
				return err
			}
			if len(entities.GetEntities()) != 0 {
				return fmt.Errorf("nonmatching graph filter returned %v", entities)
			}
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("create set and delete", func(t *testing.T) {
		for _, query := range []string{
			"CREATE (:Mutation {name:'before'})",
			"MATCH (n:Mutation {name:'before'}) SET n.name='after'",
		} {
			assertCypherWriteIntegration(t, executeCypherIntegration(t, ctx, db, namespace, query))
		}
		page := relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Mutation) RETURN n.name AS name"), "name")
		if len(page.GetRows()) != 1 || page.GetRows()[0].GetValues()[0].GetString_().GetString_() != "after" {
			t.Fatalf("updated graph property was not persisted: %v", page)
		}
		assertCypherWriteIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Mutation) DELETE n"))
		if elements := graphFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Mutation) RETURN n")); len(elements) != 0 {
			t.Fatalf("deleted graph nodes are still visible: %v", elements)
		}
	})

	t.Run("pagination preserves node and scalar order", func(t *testing.T) {
		const count = 110
		nodes := make([]string, count)
		for i := range nodes {
			nodes[i] = fmt.Sprintf("(:Page {sequence:'%03d'})", i)
		}
		executeCypherIntegration(t, ctx, db, namespace, "CREATE "+strings.Join(nodes, ","))
		elements := graphFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Page) RETURN n ORDER BY n.sequence"))
		if len(elements) != count {
			t.Fatalf("paged graph contains %d elements, want %d", len(elements), count)
		}
		for i, element := range elements {
			if got := graphIntegerIntegration(t, element.GetNode().GetProperties()["sequence"]); got != int64(i) {
				t.Fatalf("graph element %d has sequence %d", i, got)
			}
		}

		page := relationalFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			"MATCH (n:Page) RETURN n,n.sequence AS sequence ORDER BY sequence"), "n", "sequence")
		if len(page.GetRows()) != count {
			t.Fatalf("paged mixed result contains %d rows, want %d", len(page.GetRows()), count)
		}
		for i, row := range page.GetRows() {
			values := row.GetValues()
			if got := graphIntegerIntegration(t, values[0].GetNode().GetProperties()["sequence"]); got != int64(i) {
				t.Fatalf("mixed node %d has sequence %d", i, got)
			}
			if got := graphIntegerIntegration(t, values[1]); got != int64(i) {
				t.Fatalf("mixed scalar %d has sequence %d", i, got)
			}
		}
	})

	t.Run("rollback and commit over a shared raw connection", func(t *testing.T) {
		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		for _, commit := range []bool{false, true} {
			name := "rolled_back"
			if commit {
				name = "committed"
			}
			tx, err := conn.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()
			err = conn.Raw(func(raw any) error {
				result, err := raw.(*PolyphenyConn).ExecuteCypherContext(ctx, namespace,
					fmt.Sprintf("CREATE (:Transaction {name:'%s'})", name))
				if err != nil {
					return err
				}
				if result.GetFrame() != nil {
					return fmt.Errorf("write returned a frame")
				}
				return nil
			})
			if err != nil {
				_ = tx.Rollback()
				t.Fatal(err)
			}
			if commit {
				err = tx.Commit()
			} else {
				err = tx.Rollback()
			}
			if err != nil {
				t.Fatal(err)
			}
		}
	})
	for _, entry := range []struct {
		name string
		want int
	}{{"rolled_back", 0}, {"committed", 1}} {
		elements := graphFrameIntegration(t, executeCypherIntegration(t, ctx, db, namespace,
			fmt.Sprintf("MATCH (n:Transaction {name:'%s'}) RETURN n", entry.name)))
		if len(elements) != entry.want {
			t.Errorf("transaction %s produced %d nodes, want %d", entry.name, len(elements), entry.want)
		}
	}

	t.Run("server errors and SQL reuse", func(t *testing.T) {
		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if err := conn.Raw(func(raw any) error {
			driverConn := raw.(*PolyphenyConn)
			canceled, cancel := context.WithCancel(ctx)
			cancel()
			if _, err := driverConn.ExecuteCypherContext(canceled, namespace, "MATCH (n) RETURN n"); !errors.Is(err, context.Canceled) {
				return fmt.Errorf("pre-canceled query returned %v", err)
			}
			if !driverConn.IsValid() {
				return fmt.Errorf("pre-canceled query invalidated a usable connection")
			}
			if _, err := driverConn.ExecuteCypherContext(ctx, namespace, "MATCH ("); err == nil {
				return fmt.Errorf("malformed Cypher must return a server error")
			}
			if _, err := driverConn.ExecuteCypherContext(ctx, namespace, "MATCH (n) RETURN missing_variable"); err == nil {
				return fmt.Errorf("invalid Cypher must return a server error")
			}
			if _, err := driverConn.ExecuteCypherContext(ctx, namespace, "MATCH (n:Person) RETURN n; MATCH (n:Person) RETURN n"); err == nil {
				return fmt.Errorf("multiple Cypher statements must be rejected")
			}
			if _, err := driverConn.ExecuteCypherContext(ctx, namespace+"_missing", "MATCH (n) RETURN n"); err == nil {
				return fmt.Errorf("missing namespace must return a server error")
			}
			return driverConn.Ping(ctx)
		}); err != nil {
			t.Fatal(err)
		}
		var one int
		if err := conn.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil || one != 1 {
			t.Fatalf("SQL reuse returned %d, %v", one, err)
		}
	})
}
