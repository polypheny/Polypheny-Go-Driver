package polypheny

import (
	prism "github.com/polypheny/Polypheny-Go-Driver/prism"
)

type PolyphenyTranaction struct {
	conn *PolyphenyConn
}

func (tx *PolyphenyTranaction) Commit() error {
	request := prism.Request{
		Type: &prism.Request_CommitRequest{
			CommitRequest: &prism.CommitRequest{},
		},
	}
	_, err := tx.conn.helperSendAndRecv(&request)
	return err
}

func (tx *PolyphenyTranaction) Rollback() error {
	request := prism.Request{
		Type: &prism.Request_RollbackRequest{
			RollbackRequest: &prism.RollbackRequest{},
		},
	}
	_, err := tx.conn.helperSendAndRecv(&request)
	return err
}
