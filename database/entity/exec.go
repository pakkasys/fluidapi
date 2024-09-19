package entity

import (
	"database/sql"

	"github.com/pakkasys/fluidapi/database/util"
)

// RowsQuery runs a row returning query. The returned rows and stmt objects must
// be closed by the caller.
//
//   - db: The database connection.
//   - query: The query string.
//   - parameters: The parameters for the query.
func RowsQuery(
	db util.DB,
	query string,
	parameters []any,
) (util.Rows, util.Stmt, error) {
	stmt, err := db.Prepare(query)
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
func ExecQuery(db util.DB, query string, parameters []any) (sql.Result, error) {
	stmt, err := db.Prepare(query)
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
