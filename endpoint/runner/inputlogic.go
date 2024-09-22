package runner

import (
	"net/http"

	"github.com/pakkasys/fluidapi/endpoint/definition"
	"github.com/pakkasys/fluidapi/endpoint/middleware"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
)

type Options struct {
	TraceLoggerFn func(r *http.Request) func(messages ...any)
	ErrorLoggerFn func(r *http.Request) func(messages ...any)
}

type StackBuilderFactory func() middleware.StackBuilder

func GenericEndpointDefinition[I any, O any](
	specification IInputSpecification[I],
	callback inputlogic.Callback[I, O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts Options,
) *definition.EndpointDefinition {
	return &definition.EndpointDefinition{
		URL:    specification.URL(),
		Method: specification.HTTPMethod(),
		MiddlewareStack: stackBuilderFactoryFn().
			MustAddMiddleware(
				*inputlogic.MiddlewareWrapper(
					callback,
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsCreate),
					nil,
					nil,
					nil,
					opts.TraceLoggerFn,
					opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

func GetEndpointDefinition[I IGetInput, O any, E any](
	specification IGetSpecification[I],
	apiFields APIFields,
	getEntitiesFn GetServiceFunc[E],
	getCountFn GetCountFunc,
	toOutputFn ToGetEndpointOutput[E, O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts Options,
) *definition.EndpointDefinition {
	return &definition.EndpointDefinition{
		URL:    specification.URL(),
		Method: specification.HTTPMethod(),
		MiddlewareStack: stackBuilderFactoryFn().
			MustAddMiddleware(
				*inputlogic.MiddlewareWrapper(
					func(
						writer http.ResponseWriter,
						request *http.Request,
						input *I,
					) (*O, error) {
						return GetInvoke(
							writer,
							request,
							*input,
							specification,
							apiFields,
							getEntitiesFn,
							getCountFn,
							toOutputFn,
						)
					},
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsGet),
					nil,
					nil,
					nil,
					opts.TraceLoggerFn,
					opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

func UpdateEndpointDefinition[I IUpdateInput, O any](
	specification IUpdateSpecification[I],
	apiFields APIFields,
	updateEntitiesFn UpdateServiceFunc,
	toOutputFn ToUpdateEndpointOutput[O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts Options,
) *definition.EndpointDefinition {
	return &definition.EndpointDefinition{
		URL:    specification.URL(),
		Method: specification.HTTPMethod(),
		MiddlewareStack: stackBuilderFactoryFn().
			MustAddMiddleware(
				*inputlogic.MiddlewareWrapper(
					func(
						writer http.ResponseWriter,
						request *http.Request,
						input *I,
					) (*O, error) {
						return UpdateInvoke(
							writer,
							request,
							*input,
							specification,
							apiFields,
							updateEntitiesFn,
							toOutputFn,
						)
					},
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsUpdate),
					nil,
					nil,
					nil,
					opts.TraceLoggerFn,
					opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

func DeleteEndpointDefinition[I IDeleteInput, O any, E any](
	specification IDeleteSpecification[I],
	apiFields APIFields,
	deleteEntitiesFn DeleteServiceFunc[E],
	toOutputFn ToDeleteEndpointOutput[O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilderFactoryFn StackBuilderFactory,
	opts Options,
) *definition.EndpointDefinition {
	return &definition.EndpointDefinition{
		URL:    specification.URL(),
		Method: specification.HTTPMethod(),
		MiddlewareStack: stackBuilderFactoryFn().
			MustAddMiddleware(
				*inputlogic.MiddlewareWrapper(
					func(
						writer http.ResponseWriter,
						request *http.Request,
						input *I,
					) (*O, error) {
						return DeleteInvoke(
							writer,
							request,
							*input,
							specification,
							apiFields,
							deleteEntitiesFn,
							toOutputFn,
						)
					},
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsDelete),
					nil,
					nil,
					nil,
					opts.TraceLoggerFn,
					opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

func selectExpectedErrors(
	providedErrors []inputlogic.ExpectedError,
	defaultErrors []inputlogic.ExpectedError,
) []inputlogic.ExpectedError {
	if len(providedErrors) > 0 {
		return providedErrors
	}
	return defaultErrors
}
