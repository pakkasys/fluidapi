// Package runner provides utilities for defining and creating HTTP API endpoints, including middleware configuration, client communication, and data parsing.
package runner

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/core/client"
	"github.com/pakkasys/fluidapi/database/entity"
	databaseutil "github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/definition"
	"github.com/pakkasys/fluidapi/endpoint/middleware"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/selector"
	"github.com/pakkasys/fluidapi/endpoint/update"
)

// InputFactory is a function type used to create new instances of an input
// type.
type InputFactory[T any] func() *T

// InputSpecification defines the URL, HTTP method, and input factory function
// for an endpoint.
type InputSpecification[Input any] struct {
	// URL of the endpoint.
	URL string
	// HTTP method for the endpoint (e.g., GET, POST).
	Method string
	// Factory function to create input objects.
	InputFactory InputFactory[Input]
}

// StackBuilder is an interface for building middleware stacks for endpoints.
type StackBuilder interface {
	// Builds and returns the middleware stack.
	Build() middleware.Stack
	// Adds middleware wrappers to the stack.
	MustAddMiddleware(wrapper ...api.MiddlewareWrapper) StackBuilder
}

// StackBuilderFactory is a function type that creates a new StackBuilder
// instance.
type StackBuilderFactory func() StackBuilder

// SendFunc is a function type used to send requests from a client.
// It returns a response or an error.
type SendFunc[I any, W any] func(
	// The input data to send.
	input *I,
	// The host to send the request to.
	host string,
) (*client.Response[I, W], error)

// Client represents an HTTP client used to communicate with an API endpoint.
type Client[I any, O any, W any] struct {
	// URL of the API endpoint.
	URL string
	// HTTP method to use.
	Method string
	// Function to send the request.
	Send SendFunc[I, W]
}

// Endpoint represents an API endpoint, including its definition and client.
type Endpoint[I any, O any, W any] struct {
	// The definition of the endpoint.
	Definition *definition.EndpointDefinition
	// The client used to interact with the endpoint.
	Client *Client[I, O, W]
}

// EndpointOption is a function type used to configure an Endpoint.
type EndpointOption[I any, O any, W any] func(*Endpoint[I, O, W])

// WithMiddlewareWrapper configures an endpoint to use a specific middleware
// wrapper.
//
// Parameters:
//   - middlewareWrapper: The middleware wrapper to apply to the endpoint.
//
// Returns:
//   - An EndpointOption that applies the specified middleware wrapper.
func WithMiddlewareWrapper[I ValidatedInput, O any, W any](
	middlewareWrapper *api.MiddlewareWrapper,
) EndpointOption[I, O, W] {
	return func(endpoint *Endpoint[I, O, W]) {
		endpoint.Definition.MiddlewareStack = middleware.Stack{
			*middlewareWrapper,
		}
	}
}

// GenericEndpointDefinition creates a generic endpoint definition with
// specified options.
//
// Parameters:
//   - specification: The input specification containing URL, method, and input
//     factory.
//   - callback: The callback function to process the endpoint request.
//   - expectedErrors: A list of expected errors to handle.
//   - stackBuilder: A factory function to create a middleware stack
//     builder.
//   - opts: Options for configuring the input logic.
//   - sendFn: Function to send the request.
//   - options: Additional endpoint options to configure.
//
// Returns:
//   - An Endpoint instance.
func GenericEndpointDefinition[I ValidatedInput, O any, W any](
	specification InputSpecification[I],
	callback inputlogic.Callback[I, O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilder StackBuilder,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
	options ...EndpointOption[I, O, W],
) *Endpoint[I, O, W] {
	middlewareWrapper := inputlogic.MiddlewareWrapper(
		callback,
		specification.InputFactory,
		expectedErrors,
		opts,
	)

	definition := &definition.EndpointDefinition{
		URL:    specification.URL,
		Method: specification.Method,
		MiddlewareStack: stackBuilder.
			MustAddMiddleware(*middlewareWrapper).Build(),
	}

	client := &Client[I, O, W]{
		URL:    specification.URL,
		Method: specification.Method,
		Send:   sendFn,
	}

	endpoint := &Endpoint[I, O, W]{
		Definition: definition,
		Client:     client,
	}

	for _, opt := range options {
		opt(endpoint)
	}

	return endpoint
}

// GetEndpointDefinition creates an endpoint definition for a GET request.
//
// Parameters:
//   - specification: The input specification for the GET request.
//   - getEntitiesFn: Function to get entities from the database.
//   - getCountFn: Function to get the count of entities.
//   - toOutputFn: Function to convert entities to output format.
//   - expectedErrors: A list of expected errors to handle.
//   - stackBuilder: A factory function to create a middleware stack
//     builder.
//   - opts: Options for configuring the input logic.
//   - sendFn: Function to send the request.
//   - options: Additional endpoint options to configure.
//
// Returns:
//   - An Endpoint instance.
func GetEndpointDefinition[I ParseableInput[ParsedGetEndpointInput], O any, E any, W any](
	specification InputSpecification[I],
	getEntitiesFn GetServiceFunc[E],
	getCountFn GetCountFunc,
	toOutputFn ToGetEndpointOutput[E, O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilder StackBuilder,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
	options ...EndpointOption[I, O, W],
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
		stackBuilder,
		opts,
		sendFn,
		options...,
	)
}

// UpdateEndpointDefinition creates an endpoint definition for an UPDATE request.
//
// Parameters:
//   - specification: The input specification for the UPDATE request.
//   - updateEntitiesFn: Function to update entities in the database.
//   - toOutputFn: Function to convert the update result to output format.
//   - expectedErrors: A list of expected errors to handle.
//   - stackBuilder: A factory function to create a middleware stack
//     builder.
//   - opts: Options for configuring the input logic.
//   - sendFn: Function to send the request.
//   - options: Additional endpoint options to configure.
//
// Returns:
//   - An Endpoint instance.
func UpdateEndpointDefinition[I ParseableInput[ParsedUpdateEndpointInput], O any, W any](
	specification InputSpecification[I],
	updateEntitiesFn UpdateServiceFunc,
	toOutputFn ToUpdateEndpointOutput[O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilder StackBuilder,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
	options ...EndpointOption[I, O, W],
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
		stackBuilder,
		opts,
		sendFn,
		options...,
	)
}

// DeleteEndpointDefinition creates an endpoint definition for a DELETE request.
//
// Parameters:
//   - specification: The input specification for the DELETE request.
//   - deleteEntitiesFn: Function that deletes entities from the database.
//   - toOutputFn: Function that converts the deletion result to the desired
//     output format.
//   - expectedErrors: A list of errors that are expected and handled during
//     request processing.
//   - stackBuilder: A factory function used to create a middleware
//     stack builder for the endpoint.
//   - opts: Options to configure the input logic, such as the object picker and
//     output handler.
//   - sendFn: Function used to send the request, typically wrapped to include
//     additional headers or settings.
//   - options: Additional options for configuring the endpoint, such as
//     middleware wrappers.
//
// Returns:
//   - An Endpoint instance configured for DELETE operations.
func DeleteEndpointDefinition[I ParseableInput[ParsedDeleteEndpointInput], O any, W any](
	specification InputSpecification[I],
	deleteEntitiesFn DeleteServiceFunc,
	toOutputFn ToDeleteEndpointOutput[O],
	expectedErrors []inputlogic.ExpectedError,
	stackBuilder StackBuilder,
	opts inputlogic.Options[I],
	sendFn SendFunc[I, W],
	options ...EndpointOption[I, O, W],
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
		stackBuilder,
		opts,
		sendFn,
		options...,
	)
}

var NeedAtLeastOneUpdateError = api.NewError[any]("NEED_AT_LEAST_ONE_UPDATE")
var NeedAtLeastOneSelectorError = api.NewError[any]("NEED_AT_LEAST_ONE_SELECTOR")

type APIFields map[string]dbfield.DBField

type ParsedGetEndpointInput struct {
	Orders            []databaseutil.Order
	DatabaseSelectors databaseutil.Selectors
	Page              *page.Page
	GetCount          bool
}

type ParsedUpdateEndpointInput struct {
	DatabaseSelectors databaseutil.Selectors
	DatabaseUpdates   []entity.Update
	Upsert            bool
}

type ParsedDeleteEndpointInput struct {
	DatabaseSelectors databaseutil.Selectors
	DeleteOpts        *entity.DeleteOptions
}

// ParseGetEndpointInput parses input for a GET endpoint, translating
// API-specific fields into database selectors and orders.
//
// Parameters:
//   - apiFields: A mapping of API fields to corresponding database fields.
//   - selectors: The list of selectors provided by the client, to filter results.
//   - orders: The list of order specifications, specifying how results should
//     be sorted.
//   - allowedOrderFields: The fields that are allowed to be used for ordering.
//   - inputPage: The pagination information, specifying offset and limit for
//     results.
//   - maxPageCount: The maximum number of results that can be retrieved per
//     page.
//   - getCount: Boolean flag indicating whether the count of results should be
//     retrieved.
//
// Returns:
//   - A pointer to a ParsedGetEndpointInput containing the translated
//     selectors, orders, and pagination information.
//   - An error if parsing fails or if the input does not meet requirements.
func ParseGetEndpointInput(
	apiFields APIFields,
	selectors []selector.Selector,
	orders []order.Order,
	allowedOrderFields []string,
	inputPage *page.Page,
	maxPageCount int,
	getCount bool,
) (*ParsedGetEndpointInput, error) {
	dbOrders, err := order.ValidateAndTranslateToDBOrders(
		orders,
		allowedOrderFields,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	if inputPage == nil {
		inputPage = &page.Page{
			Offset: 0,
			Limit:  maxPageCount,
		}
	}
	if err := inputPage.Validate(maxPageCount); err != nil {
		return nil, err
	}

	dbSelectors, err := selector.ToDBSelectors(selectors, apiFields)
	if err != nil {
		return nil, err
	}

	return &ParsedGetEndpointInput{
		Orders:            dbOrders,
		DatabaseSelectors: dbSelectors,
		Page:              inputPage,
		GetCount:          getCount,
	}, nil
}

// ParseUpdateEndpointInput parses input for an UPDATE endpoint, translating
// API-specific fields into database selectors and updates.
//
// Parameters:
//   - apiFields: A mapping of API fields to corresponding database fields.
//   - selectors: The list of selectors provided by the client, used to filter
//     the entities to be updated.
//   - updates: The list of update operations to be applied to the selected
//     entities.
//   - upsert: Boolean flag indicating if the operation should be an "upsert"
//     (insert if not existing).
//
// Returns:
//   - A pointer to a ParsedUpdateEndpointInput containing the translated
//     selectors and updates.
//   - An error if parsing fails, such as if no valid selectors or updates are
//     provided.
func ParseUpdateEndpointInput(
	apiFields APIFields,
	selectors []selector.Selector,
	updates []update.Update,
	upsert bool,
) (*ParsedUpdateEndpointInput, error) {
	dbSelectors, err := selector.ToDBSelectors(selectors, apiFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, NeedAtLeastOneSelectorError
	}

	dbUpdates, err := update.ToDBUpdates(updates, apiFields)
	if err != nil {
		return nil, err
	}
	if len(dbUpdates) == 0 {
		return nil, NeedAtLeastOneUpdateError
	}

	return &ParsedUpdateEndpointInput{
		DatabaseSelectors: dbSelectors,
		DatabaseUpdates:   dbUpdates,
		Upsert:            upsert,
	}, nil
}

// ParseDeleteEndpointInput parses input for a DELETE endpoint, translating
// API-specific fields into database selectors and orders.
//
// Parameters:
//   - apiFields: A mapping of API fields to corresponding database fields.
//   - selectors: The list of selectors provided by the client, used to filter
//     the entities to be deleted.
//   - orders: The list of order specifications, specifying the order in which
//     entities should be deleted.
//   - allowedOrderFields: The fields that are allowed to be used for ordering
//     the deletions.
//   - limit: The maximum number of entities to delete.
//
// Returns:
//   - A pointer to a ParsedDeleteEndpointInput containing the translated
//     selectors and deletion options.
//   - An error if parsing fails or if the input does not meet requirements.

func ParseDeleteEndpointInput(
	apiFields APIFields,
	selectors []selector.Selector,
	orders []order.Order,
	allowedOrderFields []string,
	limit int,
) (*ParsedDeleteEndpointInput, error) {
	dbSelectors, err := selector.ToDBSelectors(selectors, apiFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, NeedAtLeastOneSelectorError
	}

	dbOrders, err := order.ValidateAndTranslateToDBOrders(
		orders,
		allowedOrderFields,
		apiFields,
	)
	if err != nil {
		return nil, err
	}

	return &ParsedDeleteEndpointInput{
		DatabaseSelectors: dbSelectors,
		DeleteOpts: &entity.DeleteOptions{
			Limit:  limit,
			Orders: dbOrders,
		},
	}, nil
}
