package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"
)

// Update is the options struct for entity update queries.
type Update struct {
	Field string
	Value any
}

// UpdateEntities updates entities in the database.
//
//   - db: The database connection to use.
//   - tableName: The name of the database table.
//   - selectors: The selectors of the entities to update.
//   - updates: The updates to apply to the entities.
func UpdateEntities(
	preparer util.Preparer,
	tableName string,
	selectors []util.Selector,
	updates []Update,
	sqlUtil SQLUtil,
) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}
	res, err := update(preparer, tableName, updates, selectors)
	return checkUpdateResult(res, err, sqlUtil)
}

func checkUpdateResult(
	result sql.Result,
	err error,
	sqlUtil SQLUtil,
) (int64, error) {
	if err != nil {
		return 0, sqlUtil.CheckDBError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func update(
	preparer util.Preparer,
	tableName string,
	updates []Update,
	selectors []util.Selector,
) (sql.Result, error) {
	query, values := updateQuery(tableName, updates, selectors)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func updateQuery(
	tableName string,
	updates []Update,
	selectors []util.Selector,
) (string, []any) {
	whereColumns, whereValues := internal.ProcessSelectors(selectors)

	setClause, values := getSetClause(updates)
	values = append(values, whereValues...)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(
		"UPDATE `%s` SET %s",
		tableName,
		setClause,
	))
	if len(whereColumns) != 0 {
		builder.WriteString(" " + getWhereClause(whereColumns))
	}

	return builder.String(), values
}

func getWhereClause(whereColumns []string) string {
	whereClause := ""
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}
	return whereClause
}

func getSetClause(updates []Update) (string, []any) {
	setClauseParts := make([]string, len(updates))
	values := make([]any, len(updates))

	for i, update := range updates {
		setClauseParts[i] = fmt.Sprintf(
			"%s = ?",
			update.Field,
		)
		values[i] = update.Value
	}

	return strings.Join(setClauseParts, ", "), values
}
