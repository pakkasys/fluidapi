package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/errors"
	"github.com/pakkasys/fluidapi/database/internal"
	"github.com/pakkasys/fluidapi/database/util"
)

// Inserter is a function used to insert an entity into the database.
type Inserter[T any] func(entity T) (columns []string, values []any)

// CreateEntity creates an entity in the database.
//
//   - entity: The entity to insert.
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - inserter: The function used to get the columns and values to insert.
func CreateEntity[T any](
	entity *T,
	preparer util.Preparer,
	tableName string,
	inserter Inserter[*T],
) (int64, error) {
	return checkInsertResult(insert(preparer, entity, tableName, inserter))
}

// CreateEntities creates entities in the database.
//
//   - entities: The entities to insert.
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - inserter: The function used to get the columns and values to insert.
func CreateEntities[T any](
	entities []*T,
	preparer util.Preparer,
	tableName string,
	inserter Inserter[*T],
) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	return checkInsertResult(insertMany(
		preparer,
		entities,
		tableName,
		inserter,
	))
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

func insertQuery[T any](
	entity *T,
	tableName string,
	inserter Inserter[*T],
) (string, []any) {
	columns, values := inserter(entity)
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

func insert[T any](
	preparer util.Preparer,
	entity *T,
	tableName string,
	inserter Inserter[*T],
) (sql.Result, error) {
	query, values := insertQuery(entity, tableName, inserter)

	statement, err := preparer.Prepare(query)
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

func insertManyQuery[T any](
	entities []*T,
	tableName string,
	inserter Inserter[*T],
) (string, []any) {
	if len(entities) == 0 {
		return "", nil
	}

	columns, _ := inserter(entities[0])
	columnNames := getInsertQueryColumnNames(columns)

	var allValues []any
	valuePlaceholders := make([]string, len(entities))
	for i, entity := range entities {
		_, values := inserter(entity)
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

func insertMany[T any](
	preparer util.Preparer,
	entities []*T,
	tableName string,
	inserter Inserter[*T],
) (sql.Result, error) {
	query, values := insertManyQuery(entities, tableName, inserter)

	statement, err := preparer.Prepare(query)
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
