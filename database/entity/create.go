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

type Inserter interface {
	GetInserted() (columns []string, values []any)
}

func CreateEntity[T Inserter](
	entity T,
	db util.DB,
	tableName string,
) (int64, error) {
	return checkInsertResult(
		insert(db, entity, tableName),
	)
}

func CreateEntities[T Inserter](
	entities []T,
	db util.DB,
	tableName string,
) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	return checkInsertResult(
		insertMany(db, entities, tableName),
	)
}

func UpsertEntity[T Inserter](
	entity T,
	db util.DB,
	tableName string,
	updateProjections []util.Projection,
) (int64, error) {
	return checkInsertResult(
		upsert(db, entity, tableName, updateProjections),
	)
}

func UpsertEntities[T Inserter](
	entities []T,
	db util.DB,
	tableName string,
	updateProjections []util.Projection,
) (int64, error) {
	return checkInsertResult(
		upsertMany(db, entities, tableName, updateProjections),
	)
}

func checkInsertResult(
	result sql.Result,
	err error,
) (int64, error) {
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if internal.IsDuplicateEntryError(mysqlErr) {
				return 0, errors.DuplicateEntry(mysqlErr)
			} else if internal.IsForeignConstraintError(mysqlErr) {
				return 0, errors.ForeignConstraintError(mysqlErr)
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

func getInsertQueryColumnNames(columns []string) string {
	wrappedColumns := make([]string, len(columns))
	for i, column := range columns {
		wrappedColumns[i] = "`" + column + "`"
	}
	columnNames := strings.Join(wrappedColumns, ", ")

	return columnNames
}

func insertQuery(
	inserter Inserter,
	tableName string,
) (string, []any) {
	columns, values := inserter.GetInserted()
	columnNames := getInsertQueryColumnNames(columns)

	valuePlaceholders := strings.TrimSuffix(
		strings.Repeat("?, ", len(values)),
		", ",
	)

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tableName,
		columnNames,
		valuePlaceholders,
	)

	return query, values
}

func insert(
	db util.DB,
	inserter Inserter,
	tableName string,
) (sql.Result, error) {
	query, values := insertQuery(inserter, tableName)

	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func insertManyQuery[T Inserter](
	entities []T,
	tableName string,
) (string, []any) {
	if len(entities) == 0 {
		return "", nil
	}

	columns, _ := entities[0].GetInserted()
	columnNames := getInsertQueryColumnNames(columns)

	var allValues []any
	valuePlaceholders := make([]string, len(entities))
	for i, entity := range entities {
		_, values := entity.GetInserted()
		placeholders := make([]string, len(values))
		for j := range values {
			placeholders[j] = "?"
		}
		valuePlaceholders[i] = "(" + strings.Join(placeholders, ", ") + ")"
		allValues = append(allValues, values...)
	}

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES %s",
		tableName,
		columnNames,
		strings.Join(valuePlaceholders, ", "),
	)

	return query, allValues
}

func insertMany[T Inserter](
	db util.DB,
	entities []T,
	tableName string,
) (sql.Result, error) {
	query, values := insertManyQuery(entities, tableName)

	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func upsertQuery(
	inserter Inserter,
	tableName string,
	updateProjections []util.Projection,
) (string, []any) {
	query, values := insertQuery(inserter, tableName)

	updateParts := make([]string, len(updateProjections))
	for i, proj := range updateProjections {
		updateParts[i] = fmt.Sprintf(
			"`%s` = %s.`%s`",
			proj.Column,
			proj.Alias,
			proj.Column,
		)
	}

	upsertQuery := fmt.Sprintf(
		"%s AS %s ON DUPLICATE KEY UPDATE %s",
		query,
		updateProjections[0].Alias,
		strings.Join(updateParts, ", "),
	)

	return upsertQuery, values
}

func upsert(
	db util.DB,
	inserter Inserter,
	tableName string,
	updateProjections []util.Projection,
) (sql.Result, error) {
	if len(updateProjections) == 0 {
		return nil, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return nil, fmt.Errorf("must provide update projections alias")
	}

	upsertQuery, values := upsertQuery(inserter, tableName, updateProjections)

	statement, err := db.Prepare(upsertQuery)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func upsertManyQuery[T Inserter](
	entities []T,
	tableName string,
	updateProjections []util.Projection,
) (string, []any) {
	insertQueryPart, allValues := insertManyQuery(entities, tableName)

	updateParts := make([]string, len(updateProjections))
	for i, proj := range updateProjections {
		updateParts[i] = fmt.Sprintf(
			"`%s` = VALUES(`%s`)",
			proj.Column,
			proj.Column,
		)
	}

	upsertQuery := fmt.Sprintf(
		"%s ON DUPLICATE KEY UPDATE %s",
		insertQueryPart,
		strings.Join(updateParts, ", "),
	)

	return upsertQuery, allValues
}

func upsertMany[T Inserter](
	db util.DB,
	entities []T,
	tableName string,
	updateProjections []util.Projection,
) (sql.Result, error) {
	if len(entities) == 0 {
		return nil, fmt.Errorf("must provide entities to upsert")
	}
	if len(updateProjections) == 0 {
		return nil, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return nil, fmt.Errorf("must provide update projections alias")
	}

	query, values := upsertManyQuery(entities, tableName, updateProjections)
	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
