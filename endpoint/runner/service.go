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

func RunGetService[Output any](
	ctx context.Context,
	parsedEndpoint *ParsedGetEndpointInput,
	serviceFunc GetServiceFunc[Output],
	getCountFunc GetCountFunc,
	joins []util.Join,
	projections []util.Projection,
) ([]Output, int, error) {
	if parsedEndpoint.GetCount {
		if getCountFunc == nil {
			return nil, 0, fmt.Errorf("GetCountFunc is nil")
		}

		count, err := getCountFunc(
			ctx,
			parsedEndpoint.DatabaseSelectors,
			nil,
		)
		if err != nil {
			return nil, 0, err
		}

		return nil, count, nil
	} else {
		if serviceFunc == nil {
			return nil, 0, fmt.Errorf("GetServiceFunc is nil")
		}

		entities, err := serviceFunc(
			ctx,
			entity.GetOptions{
				Options: entity.Options{
					Selectors:   parsedEndpoint.DatabaseSelectors,
					Orders:      parsedEndpoint.Orders,
					Page:        parsedEndpoint.Page,
					Joins:       joins,
					Projections: projections,
				},
			},
		)
		if err != nil {
			return nil, 0, err
		}

		return entities, len(entities), nil
	}
}
