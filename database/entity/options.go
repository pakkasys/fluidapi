package entity

import (
	"github.com/PakkaSys/fluidapi/database/util"
	"github.com/PakkaSys/fluidapi/endpoint/page"
)

type Options struct {
	Selectors   []util.Selector
	Orders      []util.Order
	Page        *page.InputPage
	Joins       []util.Join
	Projections []util.Projection
}

func NewOptions() *Options {
	return &Options{}
}

func (c *Options) WithSelectors(selectors []util.Selector) *Options {
	c.Selectors = selectors
	return c
}

func (c *Options) WithOrders(orders []util.Order) *Options {
	c.Orders = orders
	return c
}

func (c *Options) WithPage(page *page.InputPage) *Options {
	c.Page = page
	return c
}

func (c *Options) WithJoins(joins []util.Join) *Options {
	c.Joins = joins
	return c
}

func (c *Options) WithProjections(projections []util.Projection) *Options {
	c.Projections = projections
	return c
}
