package runner

import (
	"context"

	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"

	"net/http"
)

type ToGetEndpointOutput[ServiceOutput any, EndpointOutput any] func(
	froms []ServiceOutput,
	count *int,
) *EndpointOutput

type APIFields map[string]dbfield.DBField

func GetInvoke[EndpointInput IGetInput, EndpointOutput any, ServiceOutput any](
	writer http.ResponseWriter,
	request *http.Request,
	input EndpointInput,
	specification IGetSpecification[EndpointInput],
	apiFields APIFields,
	serviceFunc GetServiceFunc[ServiceOutput],
	getCountFunc GetCountFunc,
	toEndpointOutputFunc ToGetEndpointOutput[ServiceOutput, EndpointOutput],
) (*EndpointOutput, error) {
	parsedInput, err := ParseGetEndpointInput(input, specification, apiFields)
	if err != nil {
		return nil, err
	}

	output, count, err := RunGetService(
		request.Context(),
		parsedInput,
		serviceFunc,
		getCountFunc,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return toEndpointOutputFunc(output, &count), nil
}

type UpdateServiceFunc func(
	ctx context.Context,
	databaseSelectors []util.Selector,
	databaseUpdates []entity.Update,
) (int64, error)

type ToUpdateEndpointOutput[EndpointOutput any] func(
	count int64,
) *EndpointOutput

func UpdateInvoke[EndpointInput IUpdateInput, EndpointOutput any](
	writer http.ResponseWriter,
	request *http.Request,
	input IUpdateInput,
	specification IUpdateSpecification[EndpointInput],
	apiFields APIFields,
	serviceFunc UpdateServiceFunc,
	toEndpointOutputFunc ToUpdateEndpointOutput[EndpointOutput],
) (*EndpointOutput, error) {
	parsedInput, err := ParseUpdateEndpointInput(
		input,
		specification,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	count, err := serviceFunc(
		request.Context(),
		parsedInput.DatabaseSelectors,
		parsedInput.DatabaseUpdates,
	)
	if err != nil {
		return nil, err
	}

	return toEndpointOutputFunc(count), nil
}

type DeleteServiceFunc[ServiceOutput any] func(
	ctx context.Context,
	databaseSelectors []util.Selector,
	opts *entity.DeleteOptions,
) (int64, error)

type ToDeleteEndpointOutput[EndpointOutput any] func(
	count int64,
) *EndpointOutput

func DeleteInvoke[EndpointInput IDeleteInput, EndpointOutput any, ServiceOutput any](
	writer http.ResponseWriter,
	request *http.Request,
	input IDeleteInput,
	specification IDeleteSpecification[EndpointInput],
	apiFields APIFields,
	serviceFunc DeleteServiceFunc[ServiceOutput],
	toEndpointOutputFunc ToDeleteEndpointOutput[EndpointOutput],
) (*EndpointOutput, error) {
	parsedInput, err := ParseDeleteEndpointInput(
		input,
		specification,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	count, err := serviceFunc(
		request.Context(),
		parsedInput.DatabaseSelectors,
		parsedInput.DeleteOpts,
	)
	if err != nil {
		return nil, err
	}

	return toEndpointOutputFunc(count), nil
}
