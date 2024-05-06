// connextra defines extra functions which are not part of the sql/driver interface
package polypheny

import (
	context "context"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
	proto "google.golang.org/protobuf/proto"
)

func (conn *PolyphenyConn) queryMongoContextInternal(query string, resultChan chan *[]map[any]any, errChan chan error) {
	request := prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName:  query,
				Statement:     "mongo",
				FetchSize:     nil,
				NamespaceName: nil,
			},
		},
	}
	response, err := conn.helperSendAndRecv(&request)
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	requestID := response.GetStatementResponse().GetStatementId()
	buf, err := conn.recv(8) // this is the query result
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	err = proto.Unmarshal(buf, response)
	if err != nil {
		resultChan <- nil
		errChan <- err
		return
	}
	// is this an error?
	if requestID != response.GetStatementResponse().GetStatementId() {
		resultChan <- nil
		errChan <- nil
		return
	}
	if response.GetStatementResponse().GetResult().GetFrame() == nil || response.GetStatementResponse().GetResult().GetFrame().GetDocumentFrame() == nil {
		resultChan <- nil
		errChan <- &ClientError{
			message: "Query should return document data, however the result is empty",
		}
	}
	documentData := response.GetStatementResponse().GetResult().GetFrame().GetDocumentFrame().GetDocuments()
	result := make([]map[any]any, len(documentData))
	var currentRow map[any]any
	for i, entries := range documentData {
		currentRow = make(map[any]any)
		for _, v := range entries.GetEntries() {
			key, err := convertProtoValue(v.GetKey())
			if err != nil {
				resultChan <- nil
				errChan <- err
				return
			}
			value, err := convertProtoValue(v.GetValue())
			if err != nil {
				resultChan <- nil
				errChan <- err
				return
			}
			currentRow[key] = value
		}
		result[i] = currentRow
	}
	resultChan <- &result
	errChan <- nil
}

func (conn *PolyphenyConn) QueryMongoContext(ctx context.Context, query string) (*[]map[any]any, error) {
	errChan := make(chan error)
	rowsChan := make(chan *[]map[any]any)
	var err error
	var result *[]map[any]any
	go conn.queryMongoContextInternal(query, rowsChan, errChan)
	select {
	case <-ctx.Done():
		// context timeout or cancelled
		return nil, ctx.Err()
	case result = <-rowsChan:
		err = <-errChan
		return result, err
	}
}
