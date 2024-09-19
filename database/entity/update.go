package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/errors"
	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"
)

// UpdateOptions is the options struct for entity update queries.
type UpdateOptions struct {
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
	db util.DB,
	tableName string,
	selectors []util.Selector,
	updates []UpdateOptions,
) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	return checkUpdateResult(update(db, tableName, updates, selectors))
}

func checkUpdateResult(result sql.Result, err error) (int64, error) {
	if err != nil {
		if internal.IsMySQLError(
			err,
			internal.MySQLDuplicateEntryErrorCode,
		) {
			return 0, errors.DuplicateEntry(err)
		} else if internal.IsMySQLError(
			err,
			internal.MySQLForeignConstraintErrorCode,
		) {
			return 0, errors.ForeignConstraintError(err)
		}
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func update(
	db util.DB,
	tableName string,
	updates []UpdateOptions,
	selectors []util.Selector,
) (sql.Result, error) {
	query, values := updateQuery(tableName, updates, selectors)

	statement, err := db.Prepare(query)
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
	updates []UpdateOptions,
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

func getSetClause(updates []UpdateOptions) (string, []any) {
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
