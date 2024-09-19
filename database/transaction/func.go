package transaction

import (
	"context"
	"fmt"
	"os"

	"github.com/pakkasys/fluidapi/database/util"
)

// TransactionalFunc is a function that takes a transaction and returns a
// result.
type TransactionalFunc[Result any] func(tx util.Tx) (Result, error)

// ExecuteTransaction executes a TransactionalFunc in a transaction.
//
//   - transactionalFunc: The function to execute in a transaction.
//   - getTxFunc: The function to get a transaction.
func ExecuteTransaction[Result any](
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

// ExecuteManagedTransaction executes a TransactionalFunc in a transaction.
// It uses the context to get the transaction from and if not found it creates
// a new one. Is successful, the transaction will be, and rolled back if an
// error occurs.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - getTxFunc: The function to get a new transaction.
//   - transactionalFunc: The function to execute in a transaction.
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
