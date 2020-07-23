// Package sql provides utilities for working with package database/sql.
package sql

import (
	"context"
	"database/sql"
)

// DB is a minimal interface that is implemented by both *sql.DB and *sql.Tx.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// TxDB extends DB with transactions..
type TxDB interface {
	DB
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Row is a minimal interface around *sql.Row.
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows is a minimal interface around *sql.Rows.
type Rows interface {
	Row
	Next() bool
	Close() error
}
