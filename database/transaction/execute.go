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
) (result Result, err error) {
	result, err = transactionalFn(tx)
	defer func() {
		err = finalizeTransaction(tx, err)
	}()
	return result, err
}

// ExecuteManagedTransaction executes a TransactionalFunc in a transaction.
// It uses the context to get the transaction from and if not found it creates
// a new one. Is successful, the transaction will be, and rolled back if an
// error occurs.
//
//   - ctx: The context to use when getting and setting the transaction.
//   - newTxFn: The function to create a new transaction with.
//   - transactionalFunc: The function to execute in a transaction.
func ExecuteManagedTransaction[Result any](
	ctx context.Context,
	newTxFn func() (util.Tx, error),
	transactionalFunc TransactionalFunc[Result],
) (result Result, err error) {
	tx, isNewTx, err := handleGetTxFromContext(ctx, newTxFn)
	if err != nil {
		var zero Result
		return zero, err
	}

	// Finalize the transaction once and after the transactional function
	if isNewTx {
		defer func() {
			err = finalizeTransaction(tx, err)
			ClearTransactionFromContext(ctx)
		}()
	}

	return transactionalFunc(tx)
}

func handleGetTxFromContext(
	ctx context.Context,
	newTxFn func() (util.Tx, error),
) (util.Tx, bool, error) {
	tx := GetTransactionFromContext(ctx)
	isNewTx := false

	if tx == nil {
		var err error
		tx, err = newTxFn()
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

		return fmt.Errorf("transaction error: %v", txErr)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
