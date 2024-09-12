package runner

import (
	"fmt"

	"github.com/PakkaSys/fluidapi/database/entity"
	databaseutil "github.com/PakkaSys/fluidapi/database/util"
	"github.com/PakkaSys/fluidapi/endpoint/dbfield"
	"github.com/PakkaSys/fluidapi/endpoint/order"
	"github.com/PakkaSys/fluidapi/endpoint/page"
	"github.com/PakkaSys/fluidapi/endpoint/selector"
	"github.com/PakkaSys/fluidapi/endpoint/update"
	"github.com/PakkaSys/fluidapi/endpoint/util"
)

type ParsedGetEndpointInput struct {
	Orders            []databaseutil.Order
	DatabaseSelectors []databaseutil.Selector
	Page              *page.InputPage
	GetCount          bool
}

func ParseGetEndpointInput[Input any](
	endpointInput IGetInput,
	specification IGetSpecification[Input],
	apiFields APIFields,
) (*ParsedGetEndpointInput, error) {
	orders, err := order.ValidateAndTranslateToDatabaseOrders(
		endpointInput.GetOrders(),
		specification.AllowedOrderFields(),
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	page, err := page.ValidatePage(
		endpointInput.GetPage(),
		specification.MaxPageCount(),
	)
	if err != nil {
		return nil, err
	}

	databaseSelectors, err := handleDatabaseSelectors(
		endpointInput.GetSelectors(),
		specification.AllowedSelectors(),
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	return &ParsedGetEndpointInput{
		Orders:            orders,
		DatabaseSelectors: databaseSelectors,
		Page:              page,
		GetCount:          endpointInput.GetGetCount(),
	}, nil
}

type ParsedUpdateEndpointInput struct {
	DatabaseSelectors []databaseutil.Selector
	DatabaseUpdates   []entity.Update
	Upsert            bool
}

func ParseUpdateEndpointInput[Input any](
	endpointInput IUpdateInput,
	specification IUpdateSpecification[Input],
	apiFields APIFields,
) (*ParsedUpdateEndpointInput, error) {
	databaseSelectors, err := handleDatabaseSelectors(
		endpointInput.GetSelectors(),
		specification.AllowedSelectors(),
		apiFields,
	)
	if err != nil {
		return nil, err
	}
	if len(databaseSelectors) == 0 {
		return nil, selector.NeedAtLeastOneSelectorError()
	}

	databaseUpdates, err := update.GetDatabaseUpdatesFromUpdates(
		endpointInput.GetUpdates(),
		specification.AllowedUpdates(),
		apiFields,
		util.NewValidation(),
	)
	if err != nil {
		return nil, err
	}
	if len(databaseSelectors) == 0 {
		return nil, update.NeedAtLeastOneUpdateError()
	}

	var upsert bool
	upsertInput, ok := endpointInput.(IUpsertInput)
	if ok {
		upsert = upsertInput.GetUpsert()
	}

	return &ParsedUpdateEndpointInput{
		DatabaseSelectors: databaseSelectors,
		DatabaseUpdates:   databaseUpdates,
		Upsert:            upsert,
	}, nil
}

type ParsedDeleteEndpointInput struct {
	DatabaseSelectors []databaseutil.Selector
	DeleteOpts        *entity.DeleteOptions
}

func ParseDeleteEndpointInput[Input any](
	endpointInput IDeleteInput,
	specification IDeleteSpecification[Input],
	apiFields APIFields,
) (*ParsedDeleteEndpointInput, error) {
	databaseSelectors, err := handleDatabaseSelectors(
		endpointInput.GetSelectors(),
		specification.AllowedSelectors(),
		apiFields,
	)
	if err != nil {
		return nil, err
	}
	if len(databaseSelectors) == 0 {
		return nil, selector.NeedAtLeastOneSelectorError()
	}

	orders := []databaseutil.Order{}
	ordersObj, ok := endpointInput.(order.IOrderable)
	if ok {
		orderableSpec, ok := specification.(IDeleteSpecificationOrderable[Input])
		if !ok {
			return nil, fmt.Errorf("delete specification is not orderable")
		}

		orders, err = order.ValidateAndTranslateToDatabaseOrders(
			ordersObj.GetOrders(),
			orderableSpec.AllowedOrderFields(),
			apiFields,
		)
		if err != nil {
			return nil, err
		}

	}

	limit := 0
	limitObj, ok := endpointInput.(ILimitable)
	if ok {
		limit = limitObj.GetLimit()
	}

	return &ParsedDeleteEndpointInput{
		DatabaseSelectors: databaseSelectors,
		DeleteOpts:        entity.NewDeleteOptions(limit, orders),
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
		util.NewValidation(),
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
