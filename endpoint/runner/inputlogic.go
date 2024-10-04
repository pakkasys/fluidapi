package runner

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/core/client"
	"github.com/pakkasys/fluidapi/endpoint/definition"
	"github.com/pakkasys/fluidapi/endpoint/middleware"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
)

type StackBuilder interface {
	Build() middleware.Stack
	MustAddMiddleware(wrapper ...api.MiddlewareWrapper) StackBuilder
}

type StackBuilderFactory func() StackBuilder

type SendFunc[I any, W any] func(
	input *I,
	host string,
) (*client.Response[I, W], error)

type Client[I any, O any, W any] struct {
	URL    string
	Method string
	Send   SendFunc[I, W]
}

type Endpoint[I any, O any, W any] struct {
	Definition *definition.EndpointDefinition
	Client     *Client[I, O, W]
}

func GenericEndpointDefinition[I ValidatedInput, O any, W any](
	specification InputSpecification[I],
	callback inputlogic.Callback[I, O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
) *Endpoint[I, O, W] {
	definition := &definition.EndpointDefinition{
		URL:    specification.URL,
		Method: specification.Method,
		MiddlewareStack: stackBuilderFactoryFn().
			MustAddMiddleware(
				*inputlogic.MiddlewareWrapper(
					callback,
					specification.InputFactory,
					expectedErrors,
					opts,
				),
			).
			Build(),
	}

	client := &Client[I, O, W]{
		URL:    specification.URL,
		Method: specification.Method,
		Send:   sendFn,
	}

	return &Endpoint[I, O, W]{
		Definition: definition,
		Client:     client,
	}
}

func GetEndpointDefinition[I ParseableInput[ParsedGetEndpointInput], O any, E any, W any](
	specification InputSpecification[I],
	getEntitiesFn GetServiceFunc[E],
	getCountFn GetCountFunc,
	toOutputFn ToGetEndpointOutput[E, O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
) *Endpoint[I, O, W] {
	callback := func(
		writer http.ResponseWriter,
		request *http.Request,
		input *I,
	) (*O, error) {
		return GetInvoke(
			writer,
			request,
			*input,
			getEntitiesFn,
			getCountFn,
			toOutputFn,
		)
	}

	return GenericEndpointDefinition(
		specification,
		callback,
		expectedErrors,
		stackBuilderFactoryFn,
		opts,
		sendFn,
	)
}

func UpdateEndpointDefinition[I ParseableInput[ParsedUpdateEndpointInput], O any, W any](
	specification InputSpecification[I],
	updateEntitiesFn UpdateServiceFunc,
	toOutputFn ToUpdateEndpointOutput[O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
) *Endpoint[I, O, W] {
	callback := func(
		writer http.ResponseWriter,
		request *http.Request,
		input *I,
	) (*O, error) {
		return UpdateInvoke[I](
			writer,
			request,
			*input,
			updateEntitiesFn,
			toOutputFn,
		)
	}

	return GenericEndpointDefinition(
		specification,
		callback,
		expectedErrors,
		stackBuilderFactoryFn,
		opts,
		sendFn,
	)
}

func DeleteEndpointDefinition[I ParseableInput[ParsedDeleteEndpointInput], O any, W any](
	specification InputSpecification[I],
	deleteEntitiesFn DeleteServiceFunc,
	toOutputFn ToDeleteEndpointOutput[O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
) *Endpoint[I, O, W] {
	callback := func(
		writer http.ResponseWriter,
		request *http.Request,
		input *I,
	) (*O, error) {
		return DeleteInvoke[I](
			writer,
			request,
			*input,
			deleteEntitiesFn,
			toOutputFn,
		)
	}

	return GenericEndpointDefinition(
		specification,
		callback,
		expectedErrors,
		stackBuilderFactoryFn,
		opts,
		sendFn,
	)
}
