package entity

import (
	"database/sql"

	"github.com/PakkaSys/fluidapi/database/util"
)

// RowsQuery runs the query and returns the rows and statement.
// Caller is responsible for closing the statement and the rows after
// successful execution.
func RowsQuery(
	exec util.Executor,
	query string,
	parameters []any,
) (*sql.Rows, *sql.Stmt, error) {
	statement, err := exec.Prepare(query)
	if err != nil {
		return nil, nil, err
	}

	rows, err := statement.Query(parameters...)
	if err != nil {
		statement.Close()
		return nil, nil, err
	}

	return rows, statement, nil
}

// ExecQuery runs the query and returns the result.
func ExecQuery(
	exec util.Executor,
	query string,
	parameters []any,
) (sql.Result, error) {
	statement, err := exec.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(parameters...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
