package runner

import (
	"context"

	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"

	"net/http"
)

type ParseableInput[Output any] interface {
	Parse() (*Output, error)
}
type ToGetEndpointOutput[ServiceOutput any, EndpointOutput any] func(
	froms []ServiceOutput,
	count *int,
) *EndpointOutput

type UpdateServiceFunc func(
	ctx context.Context,
	databaseSelectors []util.Selector,
	databaseUpdates []entity.UpdateOptions,
) (int64, error)

type ToUpdateEndpointOutput[EndpointOutput any] func(
	count int64,
) *EndpointOutput

type DeleteServiceFunc func(
	ctx context.Context,
	databaseSelectors []util.Selector,
	opts *entity.DeleteOptions,
) (int64, error)

type ToDeleteEndpointOutput[EndpointOutput any] func(
	count int64,
) *EndpointOutput

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
