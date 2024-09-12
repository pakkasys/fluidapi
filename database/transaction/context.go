package transaction

import (
	"context"

	"github.com/PakkaSys/fluidapi/database/util"
	endpointutil "github.com/PakkaSys/fluidapi/endpoint/util"
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
