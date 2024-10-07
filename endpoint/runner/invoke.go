// Package runner provides utilities for invoking API endpoint functions,
// including the common CRUD operations.
package runner

import (
	"context"
	"fmt"

	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"

	"net/http"
)

// ValidatedInput represents an interface that requires validation of the input
// fields.
type ValidatedInput interface {
	Validate() []inputlogic.FieldError
}

// ParseableInput represents an input that can be parsed to a more specific
// format.
//
// Output: Type of the parsed output.
type ParseableInput[Output any] interface {
	ValidatedInput
	Parse() (*Output, error)
}

// ToGetEndpointOutput represents a function type to convert service output to
// endpoint output.
//
// Parameters:
// - froms: Slice of service output values.
// - count: Pointer to the count of items.
//
// Returns:
// - Converted endpoint output.
type ToGetEndpointOutput[ServiceOutput any, EndpointOutput any] func(
	froms []ServiceOutput,
	count *int,
) *EndpointOutput

// UpdateServiceFunc represents a function type to perform update operations on
// the database.
type UpdateServiceFunc func(
	ctx context.Context,
	databaseSelectors []util.Selector,
	databaseUpdates []entity.Update,
) (int64, error)

// ToUpdateEndpointOutput represents a function type to convert update count to
// endpoint output.
type ToUpdateEndpointOutput[EndpointOutput any] func(
	count int64,
) *EndpointOutput

// DeleteServiceFunc represents a function type to perform delete operations on
// the database.
type DeleteServiceFunc func(
	ctx context.Context,
	databaseSelectors []util.Selector,
	opts *entity.DeleteOptions,
) (int64, error)

// ToDeleteEndpointOutput represents a function type to convert delete count to
// endpoint output.
type ToDeleteEndpointOutput[EndpointOutput any] func(
	count int64,
) *EndpointOutput

// GetServiceFunc represents a function type to retrieve entities from the
// database.
type GetServiceFunc[Output any] func(
	ctx context.Context,
	opts entity.GetOptions,
) ([]Output, error)

// GetCountFunc represents a function type to get the count of entities from the
// database.
type GetCountFunc func(
	ctx context.Context,
	selectors []util.Selector,
	joins []util.Join,
) (int, error)

// GetInvoke handles the invocation of a GET endpoint.
//
// Parameters:
// - writer: The HTTP response writer.
// - request: The HTTP request.
// - input: Input data to the endpoint, implementing ParseableInput.
// - serviceFn: Function to retrieve entities from the database.
// - getCountFn: Function to get the count of entities from the database.
// - toEndpointOutputFn: Function to convert service output to endpoint output.
//
// Returns:
// - Pointer to the output object or an error.
func GetInvoke[I ParseableInput[ParsedGetEndpointInput], O any, E any](
	writer http.ResponseWriter,
	request *http.Request,
	input I,
	serviceFn GetServiceFunc[E],
	getCountFn GetCountFunc,
	toEndpointOutputFn ToGetEndpointOutput[E, O],
) (*O, error) {
	parsedInput, err := input.Parse()
	if err != nil {
		return nil, err
	}

	output, count, err := runGetService(
		request.Context(),
		parsedInput,
		serviceFn,
		getCountFn,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return toEndpointOutputFn(output, &count), nil
}

// UpdateInvoke handles the invocation of an UPDATE endpoint.
//
// Parameters:
//   - writer: The HTTP response writer.
//   - request: The HTTP request.
//   - input: Input data to the endpoint, implementing ParseableInput.
//   - serviceFn: Function to perform update operations in the database.
//   - toEndpointOutputFn: Function to convert the update result to endpoint
//     output.
//
// Returns:
// - Pointer to the output object or an error.
func UpdateInvoke[I ParseableInput[ParsedUpdateEndpointInput], EndpointOutput any](
	writer http.ResponseWriter,
	request *http.Request,
	input ParseableInput[ParsedUpdateEndpointInput],
	serviceFn UpdateServiceFunc,
	toEndpointOutputFn ToUpdateEndpointOutput[EndpointOutput],
) (*EndpointOutput, error) {
	parsedInput, err := input.Parse()
	if err != nil {
		return nil, err
	}

	count, err := serviceFn(
		request.Context(),
		parsedInput.DatabaseSelectors,
		parsedInput.DatabaseUpdates,
	)
	if err != nil {
		return nil, err
	}

	return toEndpointOutputFn(count), nil
}

// DeleteInvoke handles the invocation of a DELETE endpoint.
//
// Parameters:
//   - writer: The HTTP response writer.
//   - request: The HTTP request.
//   - input: Input data to the endpoint, implementing ParseableInput.
//   - serviceFn: Function to perform delete operations in the database.
//   - toEndpointOutputFn: Function to convert the delete result to endpoint
//     output.
//
// Returns:
// - Pointer to the output object or an error.
func DeleteInvoke[EndpointInput ParseableInput[ParsedDeleteEndpointInput], EndpointOutput any](
	writer http.ResponseWriter,
	request *http.Request,
	input ParseableInput[ParsedDeleteEndpointInput],
	serviceFn DeleteServiceFunc,
	toEndpointOutputFn ToDeleteEndpointOutput[EndpointOutput],
) (*EndpointOutput, error) {
	parsedInput, err := input.Parse()
	if err != nil {
		return nil, err
	}

	count, err := serviceFn(
		request.Context(),
		parsedInput.DatabaseSelectors,
		parsedInput.DeleteOpts,
	)
	if err != nil {
		return nil, err
	}

	return toEndpointOutputFn(count), nil
}

// runGetService performs the core logic of a GET request for retrieving
// entities or their count.
//
// Parameters:
// - ctx: The context for the request.
// - parsedEndpoint: Parsed input data for the endpoint.
// - serviceFn: Function to retrieve entities from the database.
// - getCountFn: Function to get the count of entities from the database.
// - joins: Joins to include in the retrieval.
// - projections: Projections to apply when retrieving entities.
//
// Returns:
// - A slice of output entities, the count of entities, or an error.
func runGetService[Output any](
	ctx context.Context,
	parsedEndpoint *ParsedGetEndpointInput,
	serviceFn GetServiceFunc[Output],
	getCountFn GetCountFunc,
	joins []util.Join,
	projections []util.Projection,
) ([]Output, int, error) {
	if parsedEndpoint.GetCount {
		if getCountFn == nil {
			return nil, 0, fmt.Errorf("GetCountFunc is nil")
		}

		count, err := getCountFn(
			ctx,
			parsedEndpoint.DatabaseSelectors,
			nil,
		)
		if err != nil {
			return nil, 0, err
		}

		return nil, count, nil
	} else {
		if serviceFn == nil {
			return nil, 0, fmt.Errorf("GetServiceFunc is nil")
		}

		entities, err := serviceFn(
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
