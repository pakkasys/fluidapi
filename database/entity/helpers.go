package entity

import (
	"context"

	"github.com/PakkaSys/fluidapi/database/transaction"
	"github.com/PakkaSys/fluidapi/database/util"
)

func GetEntitiesWithManagedTransaction[T any](
	ctx context.Context,
	rowScannerMultiple RowScannerMultiple[T],
	tableName string,
	getTxFunc func() (util.Transaction, error),
	dbOptions *GetOptions,
) ([]T, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		getTxFunc,
		func(tx util.Transaction) ([]T, error) {
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
	getTxFunc func() (util.Transaction, error),
	dbOptions *DBOptionsCount,
) (int, error) {
	return transaction.ExecuteManagedTransaction(
		ctx,
		getTxFunc,
		func(tx util.Transaction) (int, error) {
			return CountEntities(tx, tableName, dbOptions)
		},
	)
}
