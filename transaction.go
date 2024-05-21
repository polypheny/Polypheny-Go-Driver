package polypheny

import (
	prism "github.com/polypheny/Polypheny-Go-Driver/prism"
)

type PolyphenyTransaction struct {
	conn *PolyphenyConn
}

func (tx *PolyphenyTransaction) Commit() error {
	request := prism.Request{
		Type: &prism.Request_CommitRequest{
			CommitRequest: &prism.CommitRequest{},
		},
	}
	_, err := tx.conn.helperSendAndRecv(&request)
	return err
}

func (tx *PolyphenyTransaction) Rollback() error {
	request := prism.Request{
		Type: &prism.Request_RollbackRequest{
			RollbackRequest: &prism.RollbackRequest{},
		},
	}
	_, err := tx.conn.helperSendAndRecv(&request)
	return err
}
