package runner

import (
	"fmt"

	"github.com/pakkasys/fluidapi/database/entity"
	databaseutil "github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/selector"
	"github.com/pakkasys/fluidapi/endpoint/update"
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

	inputPage := endpointInput.GetPage()
	if inputPage == nil {
		inputPage = &page.InputPage{
			Offset: 0,
			Limit:  specification.MaxPageCount(),
		}
	}
	if err := inputPage.Validate(specification.MaxPageCount()); err != nil {
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
		Page:              inputPage,
		GetCount:          endpointInput.GetGetCount(),
	}, nil
}

type ParsedUpdateEndpointInput struct {
	DatabaseSelectors []databaseutil.Selector
	DatabaseUpdates   []entity.UpdateOptions
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
		DeleteOpts:        &entity.DeleteOptions{Limit: limit, Orders: orders},
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
