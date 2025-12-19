package repository

import (
	"context"
	"database/sql"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionManagerImpl struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &TransactionManagerImpl{
		db: db,
	}
}

type txKey struct{}

// GetTx retrieves sql.Tx from context
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

const TRANSACTION_KEY = "TRANSACTION_KEY"

func (u *TransactionManagerImpl) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var fnErr error
	defer func() {
		if fnErr != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	fnErr = fn(context.WithValue(ctx, txKey{}, tx))

	return fnErr
}

type DbInterface interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
