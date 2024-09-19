package transaction

import (
	"context"

	"github.com/pakkasys/fluidapi/database/util"
	endpointutil "github.com/pakkasys/fluidapi/endpoint/util"
)

var txContextKey = endpointutil.NewDataKey()

// SetTransactionToContext sets the transaction to the context.
//
//   - ctx: The context to set the transaction to.
//   - tx: The transaction to set.
func SetTransactionToContext(ctx context.Context, tx util.Tx) {
	endpointutil.SetContextValue(ctx, txContextKey, tx)
}

// GetTransactionFromContext returns the transaction from the context.
//
//   - ctx: The context to get the transaction from.
func GetTransactionFromContext(ctx context.Context) util.Tx {
	return endpointutil.GetContextValue[util.Tx](
		ctx,
		txContextKey,
		nil,
	)
}

// ClearTransactionFromContext clears the transaction from the context.
//
//   - ctx: The context to clear the transaction from.
func ClearTransactionFromContext(ctx context.Context) {
	endpointutil.SetContextValue(ctx, txContextKey, nil)
}
