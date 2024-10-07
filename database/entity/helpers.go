package entity

import (
	"context"
	"fmt"

	"github.com/pakkasys/fluidapi/database/transaction"
	"github.com/pakkasys/fluidapi/database/util"
)

// UpsertOptions is the options struct used for upserts.
type UpsertOptions struct {
	UpdateProjection []util.Projection
}

// Updates contains a list of updates.
type Updates []Update

// GetByField returns update with the given field.
//
//   - field: the field to search for
func (s Updates) GetByField(field string) *Update {
	for j := range s {
		if s[j].Field == field {
			return &s[j]
		}
	}
	return nil
}

// UpdateHandler is the handler for entity updates.
type UpdateHandler struct {
	// Timestamp field indicating when the entity was last updated
	UpdatedField string
	// Function that returns a new UpdateOptions for the timestamp field
	GetUpdateOptionsFn func() Update
}

// TXHelpers is a convenience struct for working with transactions.
type TXHelpers[T any] struct {
	// GetTxFn is a function that returns a transaction
	GetTxFn func(ctx context.Context) (util.Tx, error)
}

// ExecuteManagedTransaction wraps a transactional function in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - transactionalFunc: The function to execute in a transaction.
func (t *TXHelpers[T]) ExecuteManagedTransaction(
	ctx context.Context,
	transactionalFunc transaction.TransactionalFunc[T],
) (T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		t.GetTxFn,
		transactionalFunc,
	)
}

// Collection of convenience functions for working with entities.
type EntityHelpers[T any] struct {
	// TableName is the name of the table
	TableName string
	// GetTxFn is a function that returns a transaction
	GetTxFn func(ctx context.Context) (util.Tx, error)
	// InserterFn is a function that returns an inserter
	InserterFn Inserter[*T]
	// ScannerFn is a function that returns a row scanner
	ScanRowFn RowScanner[T]
	// ScannerMultipleFn is a function that returns a multiple rows scanner
	ScanRowsFn RowScannerMultiple[T]
	// EntityNotFoundFn is a function that returns an entity not found error
	EntityNotFoundFn func() error
	// UpdateHandler is an optional handler for entity updates
	UpdateHandler *UpdateHandler
	// SQLUtil is a utility for working with SQL
	SQLUtil SQLUtil
}

// CreateEntity is a generic function for creating or upserting an entity.
//
//   - preparer: The preparer used to prepare the query.
//   - object: The entity to create or upsert.
//   - opts: The options struct for upserting the entity.
func (e *EntityHelpers[T]) CreateEntity(
	preparer util.Preparer,
	object *T,
	opts *UpsertOptions,
) (*T, error) {
	if opts != nil {
		_, err := UpsertEntity(
			preparer,
			e.TableName,
			object,
			e.InserterFn,
			opts.UpdateProjection,
			e.SQLUtil,
		)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := CreateEntity(
			object,
			preparer,
			e.TableName,
			e.InserterFn,
			e.SQLUtil,
		)
		if err != nil {
			return nil, err
		}
	}

	return object, nil
}

// CreateEntityWithManagedTransaction wraps entity creation in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - entity: The entity to create or upsert.
//   - opts: The options struct for upserting the entity.
func (e *EntityHelpers[T]) CreateEntityWithManagedTransaction(
	ctx context.Context,
	entity *T,
	opts *UpsertOptions,
) (*T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) (*T, error) {
			return e.CreateEntity(tx, entity, opts)
		},
	)
}

// CreateEntities is a generic function for creating or upserting multiple
// entities.
//
//   - preparer: The preparer used to prepare the query.
//   - entities: The entities to create or upsert.
//   - opts: The options struct for upserting the entities.
func (e *EntityHelpers[T]) CreateEntities(
	preparer util.Preparer,
	entities []*T,
	opts *UpsertOptions,
) ([]*T, error) {
	if opts != nil {
		_, err := UpsertEntities(
			preparer,
			e.TableName,
			entities,
			e.InserterFn,
			opts.UpdateProjection,
			e.SQLUtil,
		)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := CreateEntities(
			entities,
			preparer,
			e.TableName,
			e.InserterFn,
			e.SQLUtil,
		)
		if err != nil {
			return nil, err
		}
	}

	return entities, nil
}

// CreateEntitiesWithManagedTransaction wraps entity creation in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - entities: The entities to create or upsert.
//   - opts: The options struct for upserting the entities.
func (e *EntityHelpers[T]) CreateEntitiesWithManagedTransaction(
	ctx context.Context,
	entities []*T,
	opts *UpsertOptions,
) ([]*T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) ([]*T, error) {
			return e.CreateEntities(tx, entities, opts)
		},
	)
}

// GetEntity is a generic function for getting an entity.
//
//   - preparer: The preparer used to prepare the query.
//   - opts: The options struct for getting the entity.
func (e *EntityHelpers[T]) GetEntity(
	preparer util.Preparer,
	opts GetOptions,
) (*T, error) {
	entity, err := GetEntity(e.TableName, e.ScanRowFn, preparer, &opts)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		if e.EntityNotFoundFn != nil {
			return nil, e.EntityNotFoundFn()
		} else {
			return nil, fmt.Errorf("entity not found")
		}
	}

	return entity, nil
}

// GetEntityWithManagedTransaction wraps entity get in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - opts: The options struct for getting the entity.
func (e *EntityHelpers[T]) GetEntityWithManagedTransaction(
	ctx context.Context,
	opts GetOptions,
) (*T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) (*T, error) {
			return e.GetEntity(tx, opts)
		},
	)
}

// GetEntities is a generic function for getting multiple entities.
//
//   - preparer: The preparer used to prepare the query.
//   - opts: The options struct for getting the entities.
func (e *EntityHelpers[T]) GetEntities(
	preparer util.Preparer,
	opts GetOptions,
) ([]T, error) {
	return GetEntities(e.TableName, e.ScanRowsFn, preparer, &opts)
}

// GetEntitiesWithManagedTransaction wraps entity get in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - opts: The options struct for getting the entities.
func (e *EntityHelpers[T]) GetEntitiesWithManagedTransaction(
	ctx context.Context,
	opts GetOptions,
) ([]T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) ([]T, error) {
			return e.GetEntities(tx, opts)
		},
	)
}

// GetEntityCount is a generic function for getting the count of entities.
//
//   - preparer: The preparer used to prepare the query.
//   - selectors: The selectors for the query.
//   - joins: The joins for the query.
func (e *EntityHelpers[T]) GetEntityCount(
	preparer util.Preparer,
	selectors []util.Selector,
	joins []util.Join,
) (int, error) {
	return CountEntities(
		preparer,
		e.TableName,
		&DBOptionsCount{
			Selectors: selectors,
			Joins:     joins,
		},
	)
}

// GetEntityCountWithManagedTransaction wraps entity count in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - selectors: The selectors for the query.
//   - joins: The joins for the query.
func (e *EntityHelpers[T]) GetEntityCountWithManagedTransaction(
	ctx context.Context,
	selectors []util.Selector,
	joins []util.Join,
) (int, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) (int, error) {
			return e.GetEntityCount(tx, selectors, joins)
		},
	)
}

// UpdateEntities updates entities and returns the number of updated
// rows. If update options are set they will be used to update the "updated"
// timestamp field only if that field update options is not explicitly set.
//
//   - preparer: The preparer used to prepare the query.
//   - selectors: The selectors for the query.
//   - updates: The update options for the query.
func (e *EntityHelpers[T]) UpdateEntities(
	preparer util.Preparer,
	selectors []util.Selector,
	updates Updates,
) (int64, error) {
	if e.UpdateHandler != nil {
		update := updates.GetByField(e.UpdateHandler.UpdatedField)
		// Add update options if not explicitly set.
		if update == nil {
			updates = append(updates, e.UpdateHandler.GetUpdateOptionsFn())
		}
	}

	return UpdateEntities(
		preparer,
		e.TableName,
		selectors,
		updates,
		e.SQLUtil,
	)
}

// UpdateEntitiesWithManagedTransaction wraps entity update in a transaction.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - selectors: The selectors for the query.
//   - updates: The update options for the query.
func (e *EntityHelpers[T]) UpdateEntitiesWithManagedTransaction(
	ctx context.Context,
	selectors []util.Selector,
	updates []Update,
) (int64, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) (int64, error) {
			return e.UpdateEntities(tx, selectors, updates)
		},
	)
}

// DeleteEntities is a generic function for deleting multiple entities.
//
//   - preparer: The preparer used to prepare the query.
//   - opts: The options struct for getting the entity.
func (e *EntityHelpers[T]) DeleteEntities(
	preparer util.Preparer,
	selectors []util.Selector,
	opts *DeleteOptions,
) (int64, error) {
	count, err := DeleteEntities(preparer, e.TableName, selectors, opts)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// DeleteEntitiesWithManagedTransaction deletes entities and returns the number
// of deleted rows.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - selectors: The selectors for the query.
//   - opts: The options struct for deleting the entities.
func (e *EntityHelpers[T]) DeleteEntitiesWithManagedTransaction(
	ctx context.Context,
	selectors []util.Selector,
	opts *DeleteOptions,
) (int64, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		e.GetTxFn,
		func(tx util.Tx) (int64, error) {
			return e.DeleteEntities(tx, selectors, opts)
		},
	)
}
