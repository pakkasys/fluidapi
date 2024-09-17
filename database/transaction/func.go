package transaction

import (
	"context"
	"fmt"
	"os"

	"github.com/pakkasys/fluidapi/database/util"
)

type TransactionalFunc[Result any] func(tx util.Tx) (Result, error)

func ExecuteTransaction[Result any](
	ctx context.Context,
	tx util.Tx,
	transactionalFunc TransactionalFunc[Result],
) (Result, error) {
	result, txErr := transactionalFunc(tx)
	defer func() {
		if err := finalizeTransaction(tx, txErr); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to finalize transaction: %v\n", err)
		}
	}()

	return result, txErr
}

func ExecuteManagedTransaction[Result any](
	ctx context.Context,
	getTxFunc func() (util.Tx, error),
	transactionalFunc TransactionalFunc[Result],
) (result Result, txErr error) {
	tx, isNewTx, txErr := handleGetTransactionFromContext(ctx, getTxFunc)
	if txErr != nil {
		var zero Result
		return zero, txErr
	}

	if isNewTx {
		defer func() {
			if err := finalizeTransaction(tx, txErr); err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Failed to finalize transaction: %v",
					err,
				)
			}
			ClearTransactionFromContext(ctx)
		}()
	}

	return transactionalFunc(tx)
}

func handleGetTransactionFromContext(
	ctx context.Context,
	getTxFunc func() (util.Tx, error),
) (util.Tx, bool, error) {
	tx := GetTransactionFromContext(ctx)
	isNewTx := false

	if tx == nil {
		var err error
		tx, err = getTxFunc()
		if err != nil {
			return nil, false, err
		}
		SetTransactionToContext(ctx, tx)
		isNewTx = true
	}

	return tx, isNewTx, nil
}

func finalizeTransaction(tx util.Tx, txErr error) error {
	if txErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v", rollbackErr)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
