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
	GetSelectors() []selector.Selector
	GetGetCount() bool
}

type IUpdateInput interface {
	GetSelectors() []selector.Selector
	GetUpdates() []update.Update
}

type IUpsertInput interface {
	GetUpsert() bool
}

type IDeleteInput interface {
	GetSelectors() []selector.Selector
}
