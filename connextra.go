package polypheny

import (
	"context"

	"github.com/polypheny/Polypheny-Prism-API/go/v2/prism"
)

// QueryMongoContext executes an MQL query and returns all documents in its
// result. Document keys are strings. Nested documents are maps with string
// keys, and lists are slices of values.
func (conn *PolyphenyConn) QueryMongoContext(ctx context.Context, query string) (rows *[]map[any]any, err error) {
	statementID, result, err := conn.execute(ctx, "mongo", query, nil)
	paginated := result.GetFrame() != nil && !result.GetFrame().GetIsLast()
	if statementID != 0 {
		defer conn.finishStatement(statementID, true, paginated, &err)
	}
	if err != nil {
		return nil, err
	}

	documents := make([]map[any]any, 0)
	err = conn.readFrames(ctx, statementID, result.GetFrame(), func(frame *prism.Frame) error {
		documentFrame := frame.GetDocumentFrame()
		if documentFrame == nil {
			return &ClientError{message: "Prism server did not return a document result"}
		}
		for _, document := range documentFrame.GetDocuments() {
			result, err := convertDocument(document)
			if err != nil {
				return err
			}
			row := make(map[any]any, len(result))
			for key, value := range result {
				row[key] = value
			}
			documents = append(documents, row)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &documents, nil
}
