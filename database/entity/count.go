package entity

import (
	"fmt"

	"github.com/PakkaSys/fluidapi/database/util"
)

func NewDBOptionsCount() *DBOptionsCount {
	return &DBOptionsCount{}
}

func (c *DBOptionsCount) WithSelectors(
	selectors []util.Selector,
) *DBOptionsCount {
	c.Selectors = selectors
	return c
}

func (c *DBOptionsCount) WithJoins(joins []util.Join) *DBOptionsCount {
	c.Joins = joins
	return c
}

func CountEntities(
	exec util.Executor,
	tableName string,
	dbOptions *DBOptionsCount,
) (int, error) {
	query, whereValues := buildBaseCountQuery(tableName, dbOptions)

	statement, err := exec.Prepare(query)
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

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM `%s` %s %s",
		tableName,
		joinClause(dbOptions.Joins),
		whereClause,
	)

	return query, whereValues
}
