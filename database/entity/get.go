package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/page"
)

const lockSQL = "FOR UPDATE"

type RowScanner[T any] func(row util.Row, entity *T) error
type RowScannerMultiple[T any] func(rows util.Rows, entity *T) error

type GetOptions struct {
	Options
	lock bool
}

func NewGetOptions() *GetOptions {
	return &GetOptions{
		Options: *NewOptions(),
	}
}

func GetOptionsFromDBOptions(options Options) *GetOptions {
	return &GetOptions{
		Options: options,
	}
}

func (c *GetOptions) Lock() bool {
	return c.lock
}

func (c *GetOptions) WithLock(lock bool) *GetOptions {
	c.lock = lock
	return c
}

func GetEntity[T any](
	tableName string,
	rowScanner RowScanner[T],
	db util.DB,
	dbOptions *GetOptions,
) (*T, error) {
	query, whereValues := buildBaseGetQuery(
		tableName,
		GetOptionsFromDBOptions(
			*NewGetOptions().
				WithSelectors(dbOptions.Selectors).
				WithOrders(dbOptions.Orders).
				WithPage(page.NewInputPage(0, 1)).
				WithJoins(dbOptions.Joins).
				WithProjections(dbOptions.Projections),
		).WithLock(dbOptions.Lock()),
	)

	return GetEntityWithQuery(
		tableName,
		rowScanner,
		db,
		query,
		whereValues,
	)
}

func GetEntityWithQuery[T any](
	tableName string,
	rowScanner RowScanner[T],
	db util.DB,
	query string,
	parameters []any,
) (*T, error) {
	entity, err := querySingle(
		db,
		query,
		parameters,
		rowScanner,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return entity, nil
}

func GetEntities[T any](
	tableName string,
	rowScannerMultiple RowScannerMultiple[T],
	db util.DB,
	dbOptions *GetOptions,
) ([]T, error) {
	query, whereValues := buildBaseGetQuery(
		tableName,
		dbOptions,
	)

	return GetEntitiesWithQuery(
		tableName,
		rowScannerMultiple,
		db,
		query,
		whereValues,
	)
}

func GetEntitiesWithQuery[T any](
	tableName string,
	rowScannerMultiple RowScannerMultiple[T],
	db util.DB,
	query string,
	parameters []any,
) ([]T, error) {
	entities, err := queryMultiple(db, query, parameters, rowScannerMultiple)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}, nil
		}
		return nil, err
	}

	return entities, nil
}

func queryMultiple[T any](
	db util.DB,
	query string,
	parameters []any,
	rowScannerMultiple RowScannerMultiple[T],
) ([]T, error) {
	rows, statement, err := RowsQuery(db, query, parameters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer statement.Close()

	entities, err := rowsToEntities(rows, rowScannerMultiple)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func querySingle[T any](
	db util.DB,
	query string,
	params []any,
	rowScanner RowScanner[T],
) (*T, error) {
	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var entity T
	err = rowToEntity(statement.QueryRow(params...), &entity, rowScanner)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func projectionsToStrings(projections []util.Projection) []string {
	if len(projections) == 0 {
		return []string{"*"}
	}

	projectionStrings := make([]string, len(projections))
	for i, projection := range projections {
		projectionStrings[i] = projection.String()
	}
	return projectionStrings
}

func joinClause(joins []util.Join) string {
	var joinClause string
	for _, join := range joins {
		joinClause += fmt.Sprintf(
			"%s JOIN `%s` ON %s = %s",
			join.Type,
			join.Table,
			join.OnLeft.String(),
			join.OnRight.String(),
		)
	}
	return joinClause
}

func whereClause(selectors []util.Selector) (string, []any) {
	whereColumns, whereValues := internal.ProcessSelectors(selectors)

	var whereClause string
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}

	return whereClause, whereValues
}

func buildBaseGetQuery(
	tableName string,
	dbOptions *GetOptions,
) (string, []any) {
	whereClause, whereValues := whereClause(dbOptions.Selectors)

	projectionClause := fmt.Sprintf(
		"SELECT %s",
		strings.Join(projectionsToStrings(dbOptions.Projections), ","),
	)

	query := fmt.Sprintf(
		"%s FROM `%s` %s %s %s %s",
		projectionClause,
		tableName,
		joinClause(dbOptions.Joins),
		whereClause,
		getOrderClauseFromOrders(dbOptions.Orders),
		getLimitOffsetClauseFromPage(dbOptions.Page),
	)

	if dbOptions.Lock() {
		query += " " + lockSQL
	}

	return query, whereValues
}

func getLimitOffsetClauseFromPage(page *page.InputPage) string {
	if page == nil {
		return ""
	}

	return fmt.Sprintf(
		"LIMIT %d OFFSET %d",
		page.Limit,
		page.Offset,
	)
}

func getOrderClauseFromOrders(
	orders []util.Order,
) string {
	orderClause := ""

	if len(orders) != 0 {
		orderClause = "ORDER BY"

		for _, readOrder := range orders {
			orderClause += fmt.Sprintf(
				" `%s`.`%s` %s,",
				readOrder.Table,
				readOrder.Field,
				readOrder.Direction,
			)
		}

		orderClause = strings.TrimSuffix(
			orderClause,
			",",
		)
	}

	return orderClause
}

func rowsToEntities[T any](
	rows util.Rows,
	rowScannerMultiple RowScannerMultiple[T],
) ([]T, error) {
	if rowScannerMultiple == nil {
		return nil, fmt.Errorf("must provide rowScannerMultiple")
	}

	var entities []T

	for rows.Next() {
		var entity T
		err := rowScannerMultiple(rows, &entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func rowToEntity[T any](
	row util.Row,
	entity *T,
	rowScanner RowScanner[T],
) error {
	err := rowScanner(row, entity)
	if err != nil {
		return err
	}
	return row.Err()
}
