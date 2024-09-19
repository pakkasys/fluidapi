package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/errors"
	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"
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

func checkInsertResult(result sql.Result, err error) (int64, error) {
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

func insertQuery(inserter Inserter, tableName string) (string, []any) {
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
