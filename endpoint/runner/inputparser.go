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
	DatabaseSelectors []databaseutil.Selector
	Page              *page.InputPage
	GetCount          bool
}

type ParsedUpdateEndpointInput struct {
	DatabaseSelectors []databaseutil.Selector
	DatabaseUpdates   []entity.UpdateOptions
	Upsert            bool
}

type ParsedDeleteEndpointInput struct {
	DatabaseSelectors []databaseutil.Selector
	DeleteOpts        *entity.DeleteOptions
}

func ParseGetEndpointInput(
	apiFields APIFields,
	selectors []selector.InputSelector,
	allowedSelectors map[string]selector.APISelector,
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

	databaseSelectors, err := handleDatabaseSelectors(
		selectors,
		allowedSelectors,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	return &ParsedGetEndpointInput{
		Orders:            dbOrders,
		DatabaseSelectors: databaseSelectors,
		Page:              inputPage,
		GetCount:          getCount,
	}, nil
}

func ParseUpdateEndpointInput(
	apiFields APIFields,
	selectors []selector.InputSelector,
	allowedSelectors map[string]selector.APISelector,
	updates []update.InputUpdate,
	allowedUpdates map[string]update.APIUpdate,
	upsert bool,
) (*ParsedUpdateEndpointInput, error) {
	databaseSelectors, err := handleDatabaseSelectors(
		selectors,
		allowedSelectors,
		apiFields,
	)
	if err != nil {
		return nil, err
	}
	if len(databaseSelectors) == 0 {
		return nil, selector.NeedAtLeastOneSelectorError()
	}

	databaseUpdates, err := update.GetDatabaseUpdatesFromUpdates(
		updates,
		allowedUpdates,
		apiFields,
	)
	if err != nil {
		return nil, err
	}
	if len(databaseSelectors) == 0 {
		return nil, update.NeedAtLeastOneUpdateError()
	}

	return &ParsedUpdateEndpointInput{
		DatabaseSelectors: databaseSelectors,
		DatabaseUpdates:   databaseUpdates,
		Upsert:            upsert,
	}, nil
}

func ParseDeleteEndpointInput(
	apiFields APIFields,
	selectors []selector.InputSelector,
	allowedSelectors map[string]selector.APISelector,
	orders []order.Order,
	allowedOrderFields []string,
	limit int,
) (*ParsedDeleteEndpointInput, error) {
	databaseSelectors, err := handleDatabaseSelectors(
		selectors,
		allowedSelectors,
		apiFields,
	)
	if err != nil {
		return nil, err
	}
	if len(databaseSelectors) == 0 {
		return nil, selector.NeedAtLeastOneSelectorError()
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
		DatabaseSelectors: databaseSelectors,
		DeleteOpts: &entity.DeleteOptions{
			Limit:  limit,
			Orders: dbOrders,
		},
	}, nil
}

func handleDatabaseSelectors(
	inputSelectors []selector.InputSelector,
	allowedSelectors map[string]selector.APISelector,
	apiFields map[string]dbfield.DBField,
) ([]databaseutil.Selector, error) {
	matchedSelectors, err := selector.MatchAndValidateInputSelectors(
		inputSelectors,
		allowedSelectors,
	)
	if err != nil {
		return nil, err
	}

	databaseSelectors, err := selector.ToDatabaseSelectors(
		apiFields,
		matchedSelectors,
	)
	if err != nil {
		return nil, err
	}

	return databaseSelectors, nil
}
