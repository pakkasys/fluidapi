package util

import (
	"context"
	"database/sql"
	"time"
)

// Preparer is an interface that allows to prepare a query.
type Preparer interface {
	Prepare(query string) (Stmt, error)
}

// DB is an interface that wraps the basic methods used from sql.DB.
type DB interface {
	Preparer
	Ping() error
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Exec(query string, args ...any) (Result, error)
	Query(query string, args ...any) (Rows, error)
	Close() error
}

// Tx is an interface that wraps the basic methods used from sql.Tx.
type Tx interface {
	Preparer
	Commit() error
	Rollback() error
}

// Stmt is an interface that wraps the basic methods used from sql.Stmt.
type Stmt interface {
	Close() error
	QueryRow(args ...any) Row
	Exec(args ...any) (Result, error)
	Query(args ...any) (Rows, error)
}

// Rows is an interface that wraps the basic methods used from sql.Rows.
type Rows interface {
	Scan(dest ...any) error
	Next() bool
	Close() error
	Err() error
}

// Row is an interface that wraps the basic methods used from sql.Row.
type Row interface {
	Scan(dest ...any) error
	Err() error
}

// Result is an interface that wraps the basic methods used from sql.Result.
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
