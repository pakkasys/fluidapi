package definition

import "github.com/pakkasys/fluidapi/endpoint/middleware"

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

// WithMiddlewareWrappers clones an endpoint definition with the provided
func WithMiddlewareWrappers(
	middlewareWrappers middleware.MiddlewareStack,
) Option {
	return func(e *EndpointDefinition) {
		e.MiddlewareStack = middlewareWrappers
	}
}

// WithMiddlewareWrappersFunc clones an endpoint definition with the provided
func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(
		endpointDefinition *EndpointDefinition,
	) middleware.MiddlewareStack,
) Option {
	return func(e *EndpointDefinition) {
		e.MiddlewareStack = middlewareWrappersFunc(e)
	}
}
