package definition

import (
	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/middleware"
)

type EndpointDefinition struct {
	URL             string
	Method          string
	MiddlewareStack middleware.Stack
}

// EndpointDefinitionsToAPIEndpoints converts a list of endpoint definitions to
// a list of API endpoints
//
//   - endpointDefinitions: A list of endpoint definitions to convert
func EndpointDefinitionsToAPIEndpoints(
	endpointDefinitions []EndpointDefinition,
) []api.Endpoint {
	endpoints := []api.Endpoint{}

	for _, endpointDefinition := range endpointDefinitions {
		middlewares := []api.Middleware{}
		for _, mw := range endpointDefinition.MiddlewareStack {
			middlewares = append(middlewares, mw.Middleware)
		}

		endpoints = append(
			endpoints,
			api.Endpoint{
				URL:         endpointDefinition.URL,
				Method:      endpointDefinition.Method,
				Middlewares: middlewares,
			},
		)
	}

	return endpoints
}

// Option is a function that modifies an endpoint definition when it is cloned
type Option func(*EndpointDefinition)

// CloneEndpointDefinition clones an endpoint definition with options
func CloneEndpointDefinition(
	original *EndpointDefinition,
	options ...Option,
) *EndpointDefinition {
	cloned := *original
	for _, option := range options {
		option(&cloned)
	}
	return &cloned
}

// WithURL clones an endpoint definition with the provided URL
func WithURL(url string) Option {
	return func(e *EndpointDefinition) {
		e.URL = url
	}
}

// WithMethod clones an endpoint definition with the provided HTTP method
func WithMethod(method string) Option {
	return func(e *EndpointDefinition) {
		e.Method = method
	}
}

// WithMiddlewareStack clones an endpoint definition with the provided
// middleware stack.
func WithMiddlewareStack(
	stack middleware.Stack,
) Option {
	return func(e *EndpointDefinition) {
		e.MiddlewareStack = stack
	}
}

// WithMiddlewareWrappersFunc clones an endpoint definition with the provided
// middleware wrappers
func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(
		endpointDefinition *EndpointDefinition,
	) middleware.Stack,
) Option {
	return func(e *EndpointDefinition) {
		e.MiddlewareStack = middlewareWrappersFunc(e)
	}
}
