package runner

import (
	"context"

	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/definition"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"

	"net/http"
)

// ServiceHandler abstracts away the specific service function logic for each endpoint
type ServiceHandler[ParsedInput any, ServiceOutput any] func(ctx context.Context, parsedInput *ParsedInput) (*ServiceOutput, error)

// EndpointBuilder builds the common logic for all endpoints
type EndpointBuilder[ParsedInput any, ServiceOutput any, EndpointOutput any] struct {
	specification   IInputSpecification[ParsedInput]
	stackBuilderFn  StackBuilderFactory
	opts            Options
	serviceHandler  ServiceHandler[ParsedInput, ServiceOutput]
	outputConverter func(serviceOutput *ServiceOutput) *EndpointOutput
	expectedErrors  []inputlogic.ExpectedError
}

// NewEndpointBuilder creates a new builder for an endpoint
func NewEndpointBuilder[ParsedInput any, ServiceOutput any, EndpointOutput any](
	spec IInputSpecification[ParsedInput],
	stackBuilderFn StackBuilderFactory,
	opts Options,
	serviceHandler ServiceHandler[ParsedInput, ServiceOutput],
	outputConverter func(serviceOutput *ServiceOutput) *EndpointOutput,
	expectedErrors []inputlogic.ExpectedError,
) *EndpointBuilder[ParsedInput, ServiceOutput, EndpointOutput] {
	return &EndpointBuilder[ParsedInput, ServiceOutput, EndpointOutput]{
		specification:   spec,
		stackBuilderFn:  stackBuilderFn,
		opts:            opts,
		serviceHandler:  serviceHandler,
		outputConverter: outputConverter,
		expectedErrors:  expectedErrors,
	}
}

// Build builds the final endpoint definition
func (b *EndpointBuilder[ParsedInput, ServiceOutput, EndpointOutput]) Build() *definition.EndpointDefinition {
	return &definition.EndpointDefinition{
		URL:    b.specification.URL(),
		Method: b.specification.HTTPMethod(),
		MiddlewareStack: b.stackBuilderFn().
			MustAddMiddleware(
				*inputlogic.MiddlewareWrapper(
					func(writer http.ResponseWriter, request *http.Request, input *ParsedInput) (*EndpointOutput, error) {
						serviceOutput, err := UnifiedInvoke(
							writer,
							request,
							input,
							b.serviceHandler,
						)
						if err != nil {
							return nil, err
						}
						return b.outputConverter(serviceOutput), nil
					},
					b.specification.InputFactory(),
					selectExpectedErrors(b.expectedErrors, []inputlogic.ExpectedError{}),
					nil, nil, nil, // No custom errorHandler, objectPicker, validator
					b.opts.TraceLoggerFn,
					b.opts.ErrorLoggerFn,
				),
			).
			Build(),
	}
}

type ParseableInput[Output any] interface {
	Parse() (*Output, error)
}

type UnifiedServiceFunc[Input any, Output any] func(
	ctx context.Context,
	input Input,
) (Output, error)

type UnifiedToOutputFunc[Input any, Output any] func(from Input) *Output

type ToGetEndpointOutput[ServiceOutput any, EndpointOutput any] func(
	froms []ServiceOutput,
	count *int,
) *EndpointOutput

type APIFields map[string]dbfield.DBField

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

func UnifiedInvoke[ServiceOutput any, ParsedInput any](
	writer http.ResponseWriter,
	request *http.Request,
	input ParseableInput[ParsedInput],
	serviceFunc UnifiedServiceFunc[*ParsedInput, *ServiceOutput],
) (*ServiceOutput, error) {
	parsedInput, err := input.Parse()
	if err != nil {
		return nil, err
	}

	serviceOutput, err := serviceFunc(request.Context(), parsedInput)
	if err != nil {
		return nil, err
	}

	return serviceOutput, nil
}
