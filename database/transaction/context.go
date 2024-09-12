package transaction

import (
	"context"

	"github.com/pakkasys/fluidapi/database/util"
	endpointutil "github.com/pakkasys/fluidapi/endpoint/util"
)

var txContextKey = endpointutil.NewDataKey()

func SetTransactionToContext(ctx context.Context, tx util.Transaction) {
	endpointutil.SetContextValue(ctx, txContextKey, tx)
}

func GetTransactionFromContext(ctx context.Context) util.Transaction {
	return endpointutil.GetContextValue[util.Transaction](
		ctx,
		txContextKey,
		nil,
	)
}

func ClearTransactionFromContext(ctx context.Context) {
	endpointutil.SetContextValue(ctx, txContextKey, nil)
}
