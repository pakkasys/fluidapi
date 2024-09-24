package runner

import (
	"context"
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

func GetEndpointDefinition[I ParseableInput[ParsedGetEndpointInput], O any, E any](
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
						serviceOutput, err := UnifiedInvoke(
							writer,
							request,
							*input,
							func(
								ctx context.Context,
								parsedInput *ParsedGetEndpointInput,
							) (*GetServiceOutput[E], error) {
								return RunGetService(
									ctx,
									parsedInput,
									getEntitiesFn,
									getCountFn,
								)
							},
						)
						if err != nil {
							return nil, err
						}

						return toOutputFn(
							serviceOutput.Entities,
							&serviceOutput.Count,
						), nil
					},
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsGet),
					nil, // errorHandler
					nil, // objectPicker
					nil, // validator
					opts.TraceLoggerFn,
					opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

func UpdateEndpointDefinition[I ParseableInput[ParsedUpdateEndpointInput], O any](
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
						// Use UnifiedInvoke for Update operation
						serviceOutput, err := UnifiedInvoke(
							writer,
							request,
							*input,
							func(ctx context.Context, parsedInput *ParsedUpdateEndpointInput) (*int64, error) {
								// Call update service function
								updatedCount, err := updateEntitiesFn(
									ctx,
									parsedInput.DatabaseSelectors,
									parsedInput.DatabaseUpdates,
								)
								return &updatedCount, err
							},
						)
						if err != nil {
							return nil, err
						}

						// Convert serviceOutput to endpoint output using toOutputFn
						return toOutputFn(*serviceOutput), nil
					},
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsUpdate),
					nil, // errorHandler
					nil, // objectPicker
					nil, // validator
					opts.TraceLoggerFn,
					opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

func DeleteEndpointDefinition[I ParseableInput[ParsedDeleteEndpointInput], O any](
	specification IDeleteSpecification[I],
	apiFields APIFields,
	deleteEntitiesFn DeleteServiceFunc,
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
						// Use UnifiedInvoke for Delete operation
						serviceOutput, err := UnifiedInvoke(
							writer,
							request,
							*input,
							func(ctx context.Context, parsedInput *ParsedDeleteEndpointInput) (*int64, error) {
								// Call delete service function
								deletedCount, err := deleteEntitiesFn(
									ctx,
									parsedInput.DatabaseSelectors,
									parsedInput.DeleteOpts,
								)
								return &deletedCount, err
							},
						)
						if err != nil {
							return nil, err
						}

						// Convert serviceOutput to endpoint output using toOutputFn
						return toOutputFn(*serviceOutput), nil
					},
					specification.InputFactory(),
					selectExpectedErrors(expectedErrors, ExpectedErrorsDelete),
					nil, // errorHandler
					nil, // objectPicker
					nil, // validator
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
