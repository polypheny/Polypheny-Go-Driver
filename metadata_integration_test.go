package polypheny

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestIntegrationDocumentMetadata(t *testing.T) {
	db := integrationDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	namespace := fmt.Sprintf("go_metadata_%d", time.Now().UnixNano())
	if _, err := db.ExecContext(ctx, "CREATE DOCUMENT NAMESPACE "+namespace); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := db.ExecContext(cleanupCtx, "DROP NAMESPACE "+namespace); err != nil {
			t.Errorf("drop test namespace: %v", err)
		}
	}()

	connection, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()

	err = connection.Raw(func(raw any) error {
		conn := raw.(*PolyphenyConn)
		for _, query := range []string{
			`db.createCollection("symptoms")`,
			`db.createCollection("empty")`,
			`db.symptoms.insertOne({"patient_id":"p1","complications":{"renal":true}})`,
			`db.symptoms.insertOne({"patient_id":"p2"})`,
		} {
			statementID, _, err := conn.execute(ctx, "mongo", query, &namespace)
			if err != nil {
				return err
			}
			if err := conn.closeStatement(ctx, statementID); err != nil {
				return err
			}
		}

		pattern := "symptoms"
		response, err := conn.SearchEntities(ctx, namespace, &pattern)
		if err != nil {
			return err
		}
		entities := response.GetEntities()
		if len(entities) != 1 {
			return fmt.Errorf("expected one matching collection: %v", response)
		}
		document := entities[0].GetDocument()
		if document == nil || document.GetNamespaceName() != namespace || document.GetCollectionName() != pattern {
			return fmt.Errorf("unexpected document collection metadata: %v", document)
		}
		if document.GetSampleDocumentJson() == "" {
			return fmt.Errorf("sample is missing for %s", pattern)
		}
		for _, field := range document.GetFields() {
			if field.GetFieldPath() == "complications.renal" && field.GetTypeName() == "bool" && field.GetIsOptional() {
				return nil
			}
		}
		return fmt.Errorf("expected optional nested field complications.renal: %v", document.GetFields())
	})
	if err != nil {
		t.Fatal(err)
	}
}
