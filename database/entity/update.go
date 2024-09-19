package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/errors"
	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"

	"github.com/go-sql-driver/mysql"
)

type Update struct {
	Field string
	Value any
}

func NewUpdate(field string, value any) *Update {
	return &Update{
		Field: field,
		Value: value,
	}
}

func UpdateEntities(
	selectors []util.Selector,
	updates []Update,
	db util.DB,
	tableName string,
) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	return checkUpdateResult(
		update(
			db,
			tableName,
			updates,
			selectors,
		),
	)
}

func checkUpdateResult(
	result sql.Result,
	err error,
) (int64, error) {
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if internal.IsMySQLError(
				mysqlErr,
				internal.MySQLDuplicateEntryErrorCode,
			) {
				return 0, errors.DuplicateEntry(mysqlErr)
			} else if internal.IsMySQLError(
				mysqlErr,
				internal.MySQLForeignConstraintErrorCode,
			) {
				return 0, errors.ForeignConstraintError(mysqlErr)
			}
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
	updates []Update,
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
	updates []Update,
	selectors []util.Selector,
) (string, []any) {
	whereColumns, whereValues := internal.ProcessSelectors(selectors)

	setClause, values := getSetClause(updates)
	values = append(values, whereValues...)

	return fmt.Sprintf(
		"UPDATE `%s` SET %s %s",
		tableName,
		setClause,
		getWhereClause(whereColumns),
	), values
}

func getWhereClause(whereColumns []string) string {
	whereClause := ""
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}
	return whereClause
}

func getSetClause(
	updates []Update,
) (string, []any) {
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
