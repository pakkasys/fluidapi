package runner

import (
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/selector"
	"github.com/pakkasys/fluidapi/endpoint/update"
)

type IGetInput interface {
	GetOrders() []order.Order
	GetPage() *page.InputPage
	GetSelectors() []selector.InputSelector
	GetGetCount() bool
}

type IUpdateInput interface {
	GetSelectors() []selector.InputSelector
	GetUpdates() []update.InputUpdate
}

type IUpsertInput interface {
	GetUpsert() bool
}

type IExtraInput[T any] interface {
	GetExtra() T
}

type IDeleteInput interface {
	GetSelectors() []selector.InputSelector
}

type IAPIUpdateInput interface {
	ToAPIUpdates() []update.APIUpdate
}

type IAPISelectorInput interface {
	ToAPISelectors() []selector.InputSelector
}
