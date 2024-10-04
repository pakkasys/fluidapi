package runner

import (
	"github.com/pakkasys/fluidapi/database/entity"
	databaseutil "github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/selector"
	"github.com/pakkasys/fluidapi/endpoint/update"
)

type APIFields map[string]dbfield.DBField

type ParsedGetEndpointInput struct {
	Orders            []databaseutil.Order
	DatabaseSelectors databaseutil.Selectors
	Page              *page.InputPage
	GetCount          bool
}

type ParsedUpdateEndpointInput struct {
	DatabaseSelectors databaseutil.Selectors
	DatabaseUpdates   []entity.Update
	Upsert            bool
}

type ParsedDeleteEndpointInput struct {
	DatabaseSelectors databaseutil.Selectors
	DeleteOpts        *entity.DeleteOptions
}

func ParseGetEndpointInput(
	apiFields APIFields,
	selectors []selector.Selector,
	orders []order.Order,
	allowedOrderFields []string,
	inputPage *page.InputPage,
	maxPageCount int,
	getCount bool,
) (*ParsedGetEndpointInput, error) {
	dbOrders, err := order.ValidateAndTranslateToDatabaseOrders(
		orders,
		allowedOrderFields,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	if inputPage == nil {
		inputPage = &page.InputPage{
			Offset: 0,
			Limit:  maxPageCount,
		}
	}
	if err := inputPage.Validate(maxPageCount); err != nil {
		return nil, err
	}

	dbSelectors, err := selector.ToDBSelectors(selectors, apiFields)
	if err != nil {
		return nil, err
	}

	return &ParsedGetEndpointInput{
		Orders:            dbOrders,
		DatabaseSelectors: dbSelectors,
		Page:              inputPage,
		GetCount:          getCount,
	}, nil
}

func ParseUpdateEndpointInput(
	apiFields APIFields,
	selectors []selector.Selector,
	updates []update.Update,
	upsert bool,
) (*ParsedUpdateEndpointInput, error) {
	dbSelectors, err := selector.ToDBSelectors(selectors, apiFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, selector.NeedAtLeastOneSelectorError
	}

	dbUpdates, err := update.ToDBUpdates(updates, apiFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, update.NeedAtLeastOneUpdateError
	}

	return &ParsedUpdateEndpointInput{
		DatabaseSelectors: dbSelectors,
		DatabaseUpdates:   dbUpdates,
		Upsert:            upsert,
	}, nil
}

func ParseDeleteEndpointInput(
	apiFields APIFields,
	selectors []selector.Selector,
	orders []order.Order,
	allowedOrderFields []string,
	limit int,
) (*ParsedDeleteEndpointInput, error) {
	dbSelectors, err := selector.ToDBSelectors(selectors, apiFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, selector.NeedAtLeastOneSelectorError
	}

	dbOrders, err := order.ValidateAndTranslateToDatabaseOrders(
		orders,
		allowedOrderFields,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	return &ParsedDeleteEndpointInput{
		DatabaseSelectors: dbSelectors,
		DeleteOpts: &entity.DeleteOptions{
			Limit:  limit,
			Orders: dbOrders,
		},
	}, nil
}
