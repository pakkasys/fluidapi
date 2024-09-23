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
//   - inserter: The function used to get the columns and values to insert.
//   - updateProjections: The projections of the entity to update.
func UpsertEntity[T any](
	preparer util.Preparer,
	tableName string,
	entity *T,
	inserter Inserter[*T],
	updateProjections []util.Projection,
) (int64, error) {
	return checkInsertResult(
		upsert(preparer, tableName, entity, inserter, updateProjections),
	)
}

// UpsertEntities upserts a multiple entities.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - entities: The entities to upsert.
//   - inserter: The function used to get the columns and values to insert.
//   - updateProjections: The projections of the entities to update.
func UpsertEntities[T any](
	preparer util.Preparer,
	tableName string,
	entities []*T,
	inserter Inserter[*T],
	updateProjections []util.Projection,
) (int64, error) {
	return checkInsertResult(
		upsertMany(preparer, entities, tableName, inserter, updateProjections),
	)
}

func upsertQuery[T any](
	entity *T,
	tableName string,
	inserter Inserter[*T],
	updateProjections []util.Projection,
) (string, []any) {
	query, values := insertQuery(entity, tableName, inserter)

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

func upsert[T any](
	preparer util.Preparer,
	tableName string,
	entity *T,
	inserter Inserter[*T],
	updateProjections []util.Projection,
) (sql.Result, error) {
	if len(updateProjections) == 0 {
		return nil, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return nil, fmt.Errorf("must provide update projections alias")
	}

	upsertQuery, values := upsertQuery(
		entity,
		tableName,
		inserter,
		updateProjections,
	)

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

func upsertManyQuery[T any](
	entities []*T,
	tableName string,
	inserter Inserter[*T],
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

	insertQueryPart, allValues := insertManyQuery(entities, tableName, inserter)

	builder := strings.Builder{}
	builder.WriteString(insertQueryPart)
	if len(updateParts) != 0 {
		builder.WriteString(" ON DUPLICATE KEY UPDATE ")
		builder.WriteString(strings.Join(updateParts, ", "))
	}
	upsertQuery := builder.String()

	return upsertQuery, allValues
}

func upsertMany[T any](
	preparer util.Preparer,
	entities []*T,
	tableName string,
	inserter Inserter[*T],
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

	query, values := upsertManyQuery(
		entities,
		tableName,
		inserter,
		updateProjections,
	)
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
