package entity

import (
	"database/sql"

	"github.com/pakkasys/fluidapi/database/util"
)

// RowsQuery runs a row returning query. The returned rows and stmt objects must
// be closed by the caller after successful query execution.
//
//   - db: The database connection.
//   - query: The query string.
//   - parameters: The parameters for the query.
func RowsQuery(
	preparer util.Preparer,
	query string,
	parameters []any,
) (util.Rows, util.Stmt, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, nil, err
	}

	rows, err := stmt.Query(parameters...)
	if err != nil {
		stmt.Close()
		return nil, nil, err
	}

	return rows, stmt, nil
}

// ExecQuery runs a query and returns the result.
//
//   - db: The database connection.
//   - query: The query string.
//   - parameters: The parameters for the query.
func ExecQuery(
	preparer util.Preparer,
	query string,
	parameters []any,
) (sql.Result, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(parameters...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
