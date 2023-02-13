// Package dbcontext provides DB transaction support for transactions.
package dbcontext

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/golang-sql/sqlexp"
)

type DB struct {
	db     *sql.DB
	logger logrus.Logger
}

type TransactionFunc func(ctx context.Context, f func(ctx context.Context) error) error

type contextKey int

const (
	txKey contextKey = iota
)

func New(db *sql.DB, logger logrus.Logger) *DB {
	return &DB{db, logger}
}

func (db *DB) DB() *sql.DB {
	return db.db
}

// Return the Querier if transaction object presetn in context, otherwise *sql.DB
func (db *DB) With(ctx context.Context) sqlexp.Querier {
	if anyTx := ctx.Value(txKey); anyTx == nil {
		db.logger.Warn("No *sql.Tx in context.Context. Methods will invoced witout transaction")
		return db.db
	} else {
		return anyTx.(*sql.Tx)
	}
}

// TODO: can add logic by propagation
func (db *DB) Transactional(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context) error) (err error) {
	fail := func(err error) error {
		return fmt.Errorf("Error during transaction: %v", err)
	}

	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return fail(err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = f(context.WithValue(ctx, txKey, tx))
	// and now defer checks this err
	return err
}
