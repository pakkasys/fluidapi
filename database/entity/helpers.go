package entity

import (
	"context"

	"github.com/pakkasys/fluidapi/database/transaction"
	"github.com/pakkasys/fluidapi/database/util"
)

func GetEntitiesWithManagedTransaction[T any](
	ctx context.Context,
	rowScannerMultiple RowScannerMultiple[T],
	tableName string,
	getTxFunc func() (util.Tx, error),
	dbOptions *GetOptions,
) ([]T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		getTxFunc,
		func(tx util.Tx) ([]T, error) {
			return GetEntities(
				tableName,
				rowScannerMultiple,
				tx,
				dbOptions,
			)
		},
	)
}

func GetCountWithManagedTransaction(
	ctx context.Context,
	tableName string,
	getTxFunc func() (util.Tx, error),
	dbOptions *DBOptionsCount,
) (int, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		getTxFunc,
		func(tx util.Tx) (int, error) {
			return CountEntities(tx, tableName, dbOptions)
		},
	)
}
