package transaction

import (
	"context"
	"testing"

	"github.com/pakkasys/fluidapi/database/util/mock"
	endpointutil "github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

// TestSetTransactionToContext tests the SetTransactionToContext function.
func TestSetTransactionToContext(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())
	mockTx := new(mock.MockTx)

	// Set the transaction in the context
	SetTransactionToContext(ctx, mockTx)

	// Retrieve the transaction from the context and verify it matches
	retrievedTx := GetTransactionFromContext(ctx)
	assert.Equal(t, mockTx, retrievedTx)
}

// TestGetTransactionFromContext_NoTransaction tests the
// GetTransactionFromContext function when no transaction is set.
func TestGetTransactionFromContext_NoTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	// Retrieve the transaction from the context when none is set
	retrievedTx := GetTransactionFromContext(ctx)

	// Verify that no transaction is returned
	assert.Nil(t, retrievedTx)
}

// TestClearTransactionFromContext tests the ClearTransactionFromContext
// function.
func TestClearTransactionFromContext(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())
	mockTx := new(mock.MockTx)

	// Set the transaction in the context
	SetTransactionToContext(ctx, mockTx)

	// Clear the transaction from the context
	ClearTransactionFromContext(ctx)

	// Verify that the transaction is cleared
	retrievedTx := GetTransactionFromContext(ctx)
	assert.Nil(t, retrievedTx)
}
