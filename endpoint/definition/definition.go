package definition

import "github.com/pakkasys/fluidapi/endpoint/middleware"

type EndpointDefinition struct {
	URL             string
	Method          string
	MiddlewareStack middleware.MiddlewareStack
}

func NewEndpointDefinition(
	url string,
	method string,
	middlewareStack middleware.MiddlewareStack,
) *EndpointDefinition {
	return &EndpointDefinition{
		URL:             url,
		Method:          method,
		MiddlewareStack: middlewareStack,
	}
}

func (e *EndpointDefinition) WithURL(url string) *EndpointDefinition {
	e.URL = url
	return e
}

func (e *EndpointDefinition) WithMethod(method string) *EndpointDefinition {
	e.Method = method
	return e
}

func (e *EndpointDefinition) WithMiddlewareWrappers(
	middlewareWrappers middleware.MiddlewareStack,
) *EndpointDefinition {
	e.MiddlewareStack = middlewareWrappers
	return e
}
