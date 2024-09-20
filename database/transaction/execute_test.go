package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/database/util/mock"
	endpointutil "github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

// TestExecuteTransaction_Success tests the case where a transaction is
// successfully executed.
func TestExecuteTransaction_Success(t *testing.T) {
	mockTx := new(mock.MockTx)

	// Mock the transactional function to return a successful result
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "success", nil
	}

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(nil).Once()

	result, err := ExecuteTransaction(mockTx, transactionalFunc)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	mockTx.AssertExpectations(t)
}

// TestExecuteTransaction_TransactionError tests the case where the
// transactional function returns an error.
func TestExecuteTransaction_TransactionalFnError(t *testing.T) {
	mockTx := new(mock.MockTx)

	// Mock the transactional function to return an error
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "", errors.New("application error")
	}

	// Setup the mock transaction expectations
	mockTx.On("Rollback").Return(nil).Once()

	result, err := ExecuteTransaction(mockTx, transactionalFunc)

	assert.Equal(t, "", result)
	assert.EqualError(t, err, "application error")
	mockTx.AssertExpectations(t)
}

// TestExecuteTransaction_TransactionalFnError tests the case where the
// transactional function returns an error.
func TestExecuteTransaction_FinalizeError(t *testing.T) {
	mockTx := new(mock.MockTx)

	// Mock the transactional function to return an error
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "", nil
	}

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(errors.New("commit error")).Once()

	result, err := ExecuteTransaction(mockTx, transactionalFunc)

	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed to commit transaction: commit error")
	mockTx.AssertExpectations(t)
}

// TestExecuteManagedTransaction_SuccessfulTransaction tests the case where a
// transaction is successfully executed.
func TestExecuteManagedTransaction_SuccessfulTransaction(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return the mock transaction
	getTxFunc := func() (util.Tx, error) {
		return mockTx, nil
	}

	// Mock the transactionalFunc to return a successful result
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "success", nil
	}

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(nil).Once()

	result, err := ExecuteManagedTransaction(ctx, getTxFunc, transactionalFunc)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	mockTx.AssertExpectations(t)
}

// TestExecuteManagedTransaction_GetTxFuncError tests the case where the
// getTxFunc returns an error.
func TestExecuteManagedTransaction_GetTxFuncError(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return an error
	getTxFunc := func() (util.Tx, error) {
		return nil, errors.New("failed to create transaction")
	}

	// Mock the transactionalFunc (should not be called)
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "success", nil
	}

	result, err := ExecuteManagedTransaction(ctx, getTxFunc, transactionalFunc)

	assert.Equal(t, "", result) // Expect empty result due to tx failure
	assert.EqualError(t, err, "failed to create transaction")
}

// TestExecuteManagedTransaction_TransactionalFnError tests the case where an
// error occurs in the transactional function.
func TestExecuteManagedTransaction_TransactionalFnError(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return the mock transaction
	getTxFunc := func() (util.Tx, error) {
		return mockTx, nil
	}

	// Mock the transactionalFunc to return an error
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "", errors.New("application error")
	}

	// Setup the mock transaction expectations
	mockTx.On("Rollback").Return(nil).Once()

	result, err := ExecuteManagedTransaction(ctx, getTxFunc, transactionalFunc)

	assert.Equal(t, "", result)
	assert.EqualError(t, err, "application error")
	mockTx.AssertExpectations(t)
}

// TestExecuteManagedTransaction_CommitError tests the case where the commit
// operation fails.
func TestExecuteManagedTransaction_CommitError(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return the mock transaction
	getTxFunc := func() (util.Tx, error) {
		return mockTx, nil
	}

	// Mock the transactionalFunc to return a successful result
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "success", nil
	}

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(errors.New("commit error")).Once()

	result, err := ExecuteManagedTransaction(ctx, getTxFunc, transactionalFunc)

	assert.Equal(t, "success", result)
	assert.EqualError(t, err, "failed to commit transaction: commit error")
	mockTx.AssertExpectations(t)
}

// TestExecuteManagedTransaction_ExistingTransactionInContext tests the case
// where a transaction already exists in the context.
func TestExecuteManagedTransaction_ExistingTransactionInContext(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Set the transaction in the context
	SetTransactionToContext(ctx, mockTx)

	// Mock the transactionalFunc to return a successful result
	transactionalFunc := func(tx util.Tx) (string, error) {
		return "success", nil
	}

	result, err := ExecuteManagedTransaction(ctx, nil, transactionalFunc)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	mockTx.AssertExpectations(t)
}

// TestExecuteManagedTransaction_ExecuteMultipleTimes tests the case where a
// managed transaction is executed multiple times
func TestExecuteManagedTransaction_ExecuteMultipleTimes(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return the mock transaction
	getTxFunc := func() (util.Tx, error) {
		return mockTx, nil
	}

	// Mock the transactionalFunc to return a successful result
	transactionalFunc := func(tx util.Tx) (string, error) {
		return ExecuteManagedTransaction(
			ctx,
			getTxFunc,
			func(tx util.Tx) (string, error) {
				return "success", nil
			},
		)
	}

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(nil).Once()

	result, err := ExecuteManagedTransaction(ctx, getTxFunc, transactionalFunc)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	mockTx.AssertExpectations(t)
}

// TestHandleGetTxFromContext_WithExistingTx tests the case where a transaction
// already exists in the context.
func TestHandleGetTxFromContext_WithExistingTx(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Set the transaction in the context
	SetTransactionToContext(ctx, mockTx)

	tx, isNewTx, err := handleGetTxFromContext(ctx, nil)

	assert.NoError(t, err)
	assert.False(t, isNewTx)    // The transaction is not new
	assert.Equal(t, mockTx, tx) // Should return the existing transaction
}

// TestHandleGetTxFromContext_WithNewTx tests the case where no transaction
// exists in the context, so a new one is created.
func TestHandleGetTxFromContext_WithNewTx(t *testing.T) {
	mockTx := new(mock.MockTx)
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return the mock transaction
	getTxFunc := func() (util.Tx, error) {
		return mockTx, nil
	}

	tx, isNewTx, err := handleGetTxFromContext(ctx, getTxFunc)

	assert.NoError(t, err)
	assert.True(t, isNewTx)     // The transaction is new
	assert.Equal(t, mockTx, tx) // The new transaction should be returned

	// Transaction should be stored in the context
	assert.Equal(t, mockTx, GetTransactionFromContext(ctx))
}

// TestHandleGetTxFromContext_GetTxError tests the case where an error occurs
// while getting a new
func TestHandleGetTxFromContext_GetTxError(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return an error
	getTxFunc := func() (util.Tx, error) {
		return nil, errors.New("failed to create transaction")
	}

	tx, isNewTx, err := handleGetTxFromContext(ctx, getTxFunc)

	assert.Nil(t, tx)
	assert.False(t, isNewTx) // No new transaction was created
	assert.EqualError(t, err, "failed to create transaction")
}

// TestHandleGetTxFromContext_NilTx tests the case where the getTxFunc returns a
// nil transaction without an error.
func TestHandleGetTxFromContext_NilTx(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	// Mock the getTxFunc to return nil transaction with no error
	getTxFunc := func() (util.Tx, error) {
		return nil, nil
	}

	tx, isNewTx, err := handleGetTxFromContext(ctx, getTxFunc)

	assert.Nil(t, tx)
	assert.True(t, isNewTx) // New transaction was created
	assert.NoError(t, err)  // No error should be returned
}

// TestFinalizeTransaction_SuccessfulCommit tests the successful commit case.
func TestFinalizeTransaction_SuccessfulCommit(t *testing.T) {
	mockTx := new(mock.MockTx)

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(nil).Once()

	err := finalizeTransaction(mockTx, nil)

	assert.NoError(t, err)
	mockTx.AssertExpectations(t)
}

// TestFinalizeTransaction_CommitError tests the case where commit fails.
func TestFinalizeTransaction_CommitError(t *testing.T) {
	mockTx := new(mock.MockTx)

	// Setup the mock transaction expectations
	mockTx.On("Commit").Return(errors.New("commit error")).Once()

	err := finalizeTransaction(mockTx, nil)

	assert.EqualError(t, err, "failed to commit transaction: commit error")
	mockTx.AssertExpectations(t)
}

// TestFinalizeTransaction_RollbackError tests the case where rollback fails.
func TestFinalizeTransaction_RollbackError(t *testing.T) {
	mockTx := new(mock.MockTx)

	// Setup the mock transaction expectations
	mockTx.On("Rollback").Return(errors.New("rollback error")).Once()

	err := finalizeTransaction(mockTx, errors.New("transaction error"))

	assert.EqualError(t, err, "failed to rollback transaction: rollback error")
	mockTx.AssertExpectations(t)
}
