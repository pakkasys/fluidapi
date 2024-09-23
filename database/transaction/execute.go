package transaction

import (
	"context"
	"fmt"

	"github.com/pakkasys/fluidapi/database/util"
)

// TransactionalFunc is a function that takes a transaction and returns a
// result.
type TransactionalFunc[Result any] func(tx util.Tx) (Result, error)

// ExecuteTransaction executes a TransactionalFunc in a transaction.
//
//   - tx: The transaction to use.
//   - transactionalFn: The function to execute in a transaction.
func ExecuteTransaction[Result any](
	tx util.Tx,
	transactionalFn TransactionalFunc[Result],
) (result Result, txErr error) {
	result, txErr = transactionalFn(tx)
	defer func() {
		if err := finalizeTransaction(tx, txErr); err != nil {
			txErr = err
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
//   - getTxFn: The function to get a new transaction.
//   - transactionalFn: The function to execute in a transaction.
func ExecuteManagedTransaction[Result any](
	ctx context.Context,
	getTxFn func(ctx context.Context) (util.Tx, error),
	transactionalFn TransactionalFunc[Result],
) (result Result, txErr error) {
	tx, isNewTx, txErr := handleGetTxFromContext(ctx, getTxFn)
	if txErr != nil {
		var zero Result
		return zero, txErr
	}

	if isNewTx {
		defer func() {
			if err := finalizeTransaction(tx, txErr); err != nil {
				txErr = err
			}
			ClearTransactionFromContext(ctx)
		}()
	}

	return transactionalFn(tx)
}

func handleGetTxFromContext(
	ctx context.Context,
	getTxFunc func(ctx context.Context) (util.Tx, error),
) (util.Tx, bool, error) {
	tx := GetTransactionFromContext(ctx)
	isNewTx := false

	if tx == nil {
		var err error
		tx, err = getTxFunc(ctx)
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
		return nil
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
