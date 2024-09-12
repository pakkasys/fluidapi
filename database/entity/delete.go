package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/PakkaSys/fluidapi/database/internal"
	"github.com/PakkaSys/fluidapi/database/util"
)

type DeleteOptions struct {
	Limit  int
	Orders []util.Order
}

func NewDeleteOptions(limit int, orders []util.Order) *DeleteOptions {
	return &DeleteOptions{
		Limit:  limit,
		Orders: orders,
	}
}

type DBOptionsCount struct {
	Selectors []util.Selector
	Joins     []util.Join
}

func DeleteEntities(
	selectors []util.Selector,
	exec util.Executor,
	tableName string,
	opts *DeleteOptions,
) (int64, error) {
	result, err := delete(
		exec,
		tableName,
		selectors,
		opts,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func delete(
	exec util.Executor,
	tableName string,
	selectors []util.Selector,
	opts *DeleteOptions,
) (sql.Result, error) {
	whereColumns, whereValues :=
		internal.ProcessSelectors(selectors)

	whereClause := ""
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}

	query := fmt.Sprintf(
		"DELETE FROM `%s` %s",
		tableName,
		whereClause,
	)

	if opts != nil {
		orderClause := getOrderClauseFromOrders(opts.Orders)
		if orderClause != "" {
			query += " " + orderClause
		}

		limit := opts.Limit
		if limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", limit)
		}
	}

	statement, err := exec.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.Exec(whereValues...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
