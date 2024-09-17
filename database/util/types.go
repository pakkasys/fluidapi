package util

import "database/sql"

// DB is an interface that wraps the basic methods used from sql.DB.
type DB interface {
	Prepare(query string) (Stmt, error)
}

// Tx is an interface that wraps the basic methods used from sql.Tx.
type Tx interface {
	Prepare(query string) (Stmt, error)
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

// RealStmt wraps *sql.Stmt to implement the Stmt interface.
type RealStmt struct {
	*sql.Stmt
}

func (s *RealStmt) Close() error {
	return s.Stmt.Close()
}

func (s *RealStmt) QueryRow(args ...any) Row {
	return &RealRow{Row: s.Stmt.QueryRow(args...)}
}

func (s *RealStmt) Exec(args ...any) (Result, error) {
	res, err := s.Stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return &RealResult{Result: res}, nil
}

func (s *RealStmt) Query(args ...any) (Rows, error) {
	rows, err := s.Stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return &RealRows{Rows: rows}, nil
}

// RealRows wraps *sql.Rows to implement the Rows interface.
type RealRows struct {
	*sql.Rows
}

func (r *RealRows) Scan(dest ...any) error {
	return r.Rows.Scan(dest...)
}

func (r *RealRows) Next() bool {
	return r.Rows.Next()
}

func (r *RealRows) Close() error {
	return r.Rows.Close()
}

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
