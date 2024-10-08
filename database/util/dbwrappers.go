package util

import (
	"context"
	"database/sql"
	"time"
)

// NewSQLDB creates a new instance of SQLDB.
//
// Parameters:
//   - driver: The database driver name
//   - dsn: The database connection string
//
// Returns:
//   - A new instance of SQLDB
func NewSQLDB(driver string, dsn string) (DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &SQLDB{DB: db}, nil
}

// SQLDB wraps *sql.DB to implement DBInterface.
type SQLDB struct {
	*sql.DB
}

// Ping sends a ping to the database.
func (db *SQLDB) Ping() error {
	return db.DB.Ping()
}

// SetConnMaxLifetime sets the maximum time a connection may be reused.
func (db *SQLDB) SetConnMaxLifetime(d time.Duration) {
	db.DB.SetConnMaxLifetime(d)
}

// SetConnMaxIdleTime sets the maximum time an idle connection may be reused.
func (db *SQLDB) SetConnMaxIdleTime(d time.Duration) {
	db.DB.SetConnMaxIdleTime(d)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (db *SQLDB) SetMaxOpenConns(n int) {
	db.DB.SetMaxOpenConns(n)
}

// SetMaxIdleConns sets the maximum number of idle connections to the database.
func (db *SQLDB) SetMaxIdleConns(n int) {
	db.DB.SetMaxIdleConns(n)
}

// Prepare creates a prepared statement for later queries or executions.
func (db *SQLDB) Prepare(query string) (Stmt, error) {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &RealStmt{Stmt: stmt}, nil
}

// BeginTx creates a transaction and returns it.
func (db *SQLDB) BeginTx(
	ctx context.Context,
	opts *sql.TxOptions,
) (Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &RealTx{Tx: tx}, nil
}

// Exec executes a query without returning rows.
func (db *SQLDB) Exec(query string, args ...any) (Result, error) {
	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &RealResult{Result: res}, nil
}

// Query executes a query that returns rows.
func (db *SQLDB) Query(query string, args ...any) (Rows, error) {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &RealRows{Rows: rows}, nil
}

// RealStmt wraps *sql.Stmt to implement the Stmt interface.
type RealStmt struct {
	*sql.Stmt
}

// Close closes the statement.
func (s *RealStmt) Close() error {
	return s.Stmt.Close()
}

// QueryRow executes a prepared query statement with the given arguments.
func (s *RealStmt) QueryRow(args ...any) Row {
	return s.Stmt.QueryRow(args...)
}

// Exec executes a prepared statement with the given arguments.
func (s *RealStmt) Exec(args ...any) (Result, error) {
	return s.Stmt.Exec(args...)
}

// Query executes a prepared query statement with the given arguments.
func (s *RealStmt) Query(args ...any) (Rows, error) {
	return s.Stmt.Query(args...)
}

// RealRows wraps *sql.Rows to implement the Rows interface.
type RealRows struct {
	*sql.Rows
}

// Scan scans the rows into dest.
func (r *RealRows) Scan(dest ...any) error {
	return r.Rows.Scan(dest...)
}

// Next advances the rows.
func (r *RealRows) Next() bool {
	return r.Rows.Next()
}

// Close closes the rows.
func (r *RealRows) Close() error {
	return r.Rows.Close()
}

// Err returns the error, if any, that was encountered during iteration.
func (r *RealRows) Err() error {
	return r.Rows.Err()
}

// RealRow wraps *sql.Row to implement the Row interface.
type RealRow struct {
	*sql.Row
}

// Scan scans the row into dest.
func (r *RealRow) Scan(dest ...any) error {
	return r.Row.Scan(dest...)
}

// RealResult wraps sql.Result to implement the Result interface.
type RealResult struct {
	Result sql.Result
}

// LastInsertId returns the last inserted id.
func (r *RealResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

// RowsAffected returns the number of rows affected.
func (r *RealResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

// RealTx wraps *sql.Tx to implement the Tx interface.
type RealTx struct {
	*sql.Tx
}

// Prepare commits the transaction.
func (tx *RealTx) Commit() error {
	return tx.Tx.Commit()
}

// Rollback rollbacks the transaction.
func (tx *RealTx) Rollback() error {
	return tx.Tx.Rollback()
}

// Prepare prepares the statement.
func (tx *RealTx) Prepare(query string) (Stmt, error) {
	stmt, err := tx.Tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &RealStmt{Stmt: stmt}, nil
}
