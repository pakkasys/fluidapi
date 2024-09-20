package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"
)

// DeleteOptions is the options struct for entity delete queries.
type DeleteOptions struct {
	Limit  int
	Orders []util.Order
}

// DeleteEntities deletes entities from the database.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - selectors: The selectors for the entities to delete.
//   - opts: The options for the query.
func DeleteEntities(
	preparer util.Preparer,
	tableName string,
	selectors []util.Selector,
	opts *DeleteOptions,
) (int64, error) {
	result, err := delete(preparer, tableName, selectors, opts)
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
	preparer util.Preparer,
	tableName string,
	selectors []util.Selector,
	opts *DeleteOptions,
) (sql.Result, error) {
	whereColumns, whereValues := internal.ProcessSelectors(selectors)

	whereClause := ""
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}

	builder := strings.Builder{}
	builder.WriteString(
		fmt.Sprintf("DELETE FROM `%s` %s", tableName, whereClause),
	)

	if opts != nil {
		writeDeleteOptions(&builder, opts)
	}

	statement, err := preparer.Prepare(builder.String())
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

func writeDeleteOptions(
	builder *strings.Builder,
	opts *DeleteOptions,
) {
	orderClause := getOrderClauseFromOrders(opts.Orders)
	if orderClause != "" {
		builder.WriteString(" " + orderClause)
	}

	limit := opts.Limit
	if limit > 0 {
		builder.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}
}
