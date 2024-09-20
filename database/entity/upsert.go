package entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/util"
)

// UpsertEntity upserts an entity.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - entity: The entity to upsert.
//   - updateProjections: The projections of the entity to update.
func UpsertEntity[T Inserter](
	preparer util.Preparer,
	tableName string,
	entity T,
	updateProjections []util.Projection,
) (int64, error) {
	return checkInsertResult(
		upsert(preparer, tableName, entity, updateProjections),
	)
}

// UpsertEntities upserts a multiple entities.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - entities: The entities to upsert.
func UpsertEntities[T Inserter](
	preparer util.Preparer,
	tableName string,
	entities []T,
	updateProjections []util.Projection,
) (int64, error) {
	return checkInsertResult(
		upsertMany(preparer, entities, tableName, updateProjections),
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
	preparer util.Preparer,
	tableName string,
	inserter Inserter,
	updateProjections []util.Projection,
) (sql.Result, error) {
	if len(updateProjections) == 0 {
		return nil, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return nil, fmt.Errorf("must provide update projections alias")
	}

	upsertQuery, values := upsertQuery(inserter, tableName, updateProjections)

	statement, err := preparer.Prepare(upsertQuery)
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
	preparer util.Preparer,
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
