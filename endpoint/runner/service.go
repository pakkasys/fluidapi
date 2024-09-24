package runner

import (
	"context"
	"fmt"

	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"
)

type GetServiceFunc[Output any] func(
	ctx context.Context,
	opts entity.GetOptions,
) ([]Output, error)

type GetCountFunc func(
	ctx context.Context,
	selectors []util.Selector,
	joins []util.Join,
) (int, error)

type GetServiceOutput[Output any] struct {
	Entities []Output
	Count    int
}

func RunGetService[Output any](
	ctx context.Context,
	parsedEndpoint *ParsedGetEndpointInput,
	serviceFunc GetServiceFunc[Output],
	getCountFunc GetCountFunc,
) (*GetServiceOutput[Output], error) {
	if parsedEndpoint.GetCount {
		if getCountFunc == nil {
			return nil, fmt.Errorf("GetCountFunc is nil")
		}

		count, err := getCountFunc(
			ctx,
			parsedEndpoint.DatabaseSelectors,
			nil,
		)
		if err != nil {
			return nil, err
		}

		return &GetServiceOutput[Output]{
			Count: count,
		}, nil
	} else {
		if serviceFunc == nil {
			return nil, fmt.Errorf("GetServiceFunc is nil")
		}

		entities, err := serviceFunc(
			ctx,
			entity.GetOptions{
				Options: entity.Options{
					Selectors:   parsedEndpoint.DatabaseSelectors,
					Orders:      parsedEndpoint.Orders,
					Page:        parsedEndpoint.Page,
					Joins:       parsedEndpoint.Joins,
					Projections: parsedEndpoint.Projections,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return &GetServiceOutput[Output]{
			Entities: entities,
			Count:    len(entities),
		}, nil
	}
}
