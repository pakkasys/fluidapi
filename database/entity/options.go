package entity

import (
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/page"
)

// Options is the options struct used for queries.
type Options struct {
	Selectors   []util.Selector
	Orders      []util.Order
	Page        *page.InputPage
	Joins       []util.Join
	Projections []util.Projection
}
