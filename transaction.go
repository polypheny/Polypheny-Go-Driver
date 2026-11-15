package polypheny

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// PolyphenyTransaction implements driver.Tx.
type PolyphenyTransaction struct {
	conn *PolyphenyConn
	done atomic.Bool
}

// Commit commits the transaction and restores auto-commit.
func (tx *PolyphenyTransaction) Commit() error {
	return tx.finish(true)
}

// Rollback rolls back the transaction and restores auto-commit.
func (tx *PolyphenyTransaction) Rollback() error {
	return tx.finish(false)
}

func (tx *PolyphenyTransaction) finish(commit bool) error {
	if !tx.done.CompareAndSwap(false, true) {
		return &ClientError{message: "the transaction is already closed"}
	}
	defer tx.conn.inTx.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tx.conn.endTransaction(ctx, commit); err != nil {
		_ = tx.conn.close()
		return err
	}
	if err := tx.conn.setAutoCommit(ctx, true); err != nil {
		_ = tx.conn.close()
		return err
	}
	return nil
}

func (conn *PolyphenyConn) endTransaction(ctx context.Context, commit bool) error {
	request := &prism.Request{
		Type: &prism.Request_RollbackRequest{RollbackRequest: &prism.RollbackRequest{}},
	}
	if commit {
		request.Type = &prism.Request_CommitRequest{CommitRequest: &prism.CommitRequest{}}
	}
	_, err := conn.call(ctx, request)
	return err
}

func (conn *PolyphenyConn) setAutoCommit(ctx context.Context, autoCommit bool) error {
	_, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_ConnectionPropertiesUpdateRequest{
			ConnectionPropertiesUpdateRequest: &prism.ConnectionPropertiesUpdateRequest{
				ConnectionProperties: &prism.ConnectionProperties{IsAutoCommit: &autoCommit},
			},
		},
	})
	return err
}
