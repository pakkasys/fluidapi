package definition

import "github.com/PakkaSys/fluidapi/endpoint/middleware"

type Option func(*EndpointDefinition)

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

func WithURL(url string) Option {
	return func(e *EndpointDefinition) {
		e.URL = url
	}
}

func WithMethod(method string) Option {
	return func(e *EndpointDefinition) {
		e.Method = method
	}
}

func WithMiddlewareWrappers(
	middlewareWrappers middleware.MiddlewareStack,
) Option {
	return func(e *EndpointDefinition) {
		e.MiddlewareStack = middlewareWrappers
	}
}

func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(
		endpointDefinition *EndpointDefinition,
	) middleware.MiddlewareStack,
) Option {
	return func(e *EndpointDefinition) {
		e.MiddlewareStack = middlewareWrappersFunc(e)
	}
}
