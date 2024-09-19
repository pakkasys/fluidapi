package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/util"
)

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
	if len(entities) == 0 {
		return "", nil
	}

	updateParts := make([]string, len(updateProjections))
	for i, proj := range updateProjections {
		updateParts[i] = fmt.Sprintf(
			"`%s` = VALUES(`%s`)",
			proj.Column,
			proj.Column,
		)
	}

	insertQueryPart, allValues := insertManyQuery(entities, tableName)

	builder := strings.Builder{}
	builder.WriteString(insertQueryPart)
	if len(updateParts) != 0 {
		builder.WriteString(" ON DUPLICATE KEY UPDATE ")
		builder.WriteString(strings.Join(updateParts, ", "))
	}
	upsertQuery := builder.String()

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
