package transaction

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Rollback() error
	Commit() error
}

type TxKey struct{}

func withTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, TxKey{}, tx)
}

func ExtractTx(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}

type transactionManager struct {
	db *sqlx.DB
}

type TxManager interface {
	BeginTx(ctx context.Context) (context.Context, Transaction, error)
}

func NewTxManager(db *sqlx.DB) TxManager {
	return &transactionManager{db: db}
}

// BeginTx starts a transaction and inject it into the context
func (t *transactionManager) BeginTx(ctx context.Context) (context.Context, Transaction, error) {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return withTx(ctx, tx), tx, nil
}
