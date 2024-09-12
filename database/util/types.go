package util

import (
	"database/sql"
)

type Executor interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Transaction interface {
	Executor
	Commit() error
	Rollback() error
}
