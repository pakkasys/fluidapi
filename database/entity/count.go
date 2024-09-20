package entity

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/util"
)

// DBOptionsCount is the options struct for entity count queries.
type DBOptionsCount struct {
	Selectors []util.Selector
	Joins     []util.Join
}

// CountEntities counts the number of entities in the database.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - dbOptions: The options for the query.
func CountEntities(
	preparer util.Preparer,
	tableName string,
	dbOptions *DBOptionsCount,
) (int, error) {
	query, whereValues := buildBaseCountQuery(tableName, dbOptions)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var count int
	if err := statement.QueryRow(whereValues...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func buildBaseCountQuery(
	tableName string,
	dbOptions *DBOptionsCount,
) (string, []any) {
	whereClause, whereValues := whereClause(dbOptions.Selectors)

	query := strings.Trim(fmt.Sprintf(
		"SELECT COUNT(*) FROM `%s` %s %s",
		tableName,
		joinClause(dbOptions.Joins),
		whereClause,
	), " ")

	return query, whereValues
}
