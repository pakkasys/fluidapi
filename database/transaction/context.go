package transaction

import (
	"context"

	"github.com/pakkasys/fluidapi/database/util"
	endpointutil "github.com/pakkasys/fluidapi/endpoint/util"
)

var txContextKey = endpointutil.NewDataKey()

func SetTransactionToContext(ctx context.Context, tx util.Tx) {
	endpointutil.SetContextValue(ctx, txContextKey, tx)
}

func GetTransactionFromContext(ctx context.Context) util.Tx {
	return endpointutil.GetContextValue[util.Tx](
		ctx,
		txContextKey,
		nil,
	)
}

func ClearTransactionFromContext(ctx context.Context) {
	endpointutil.SetContextValue(ctx, txContextKey, nil)
}
