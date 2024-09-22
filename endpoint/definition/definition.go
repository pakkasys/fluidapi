package definition

import "github.com/pakkasys/fluidapi/endpoint/middleware"

type EndpointDefinition struct {
	URL             string
	Method          string
	MiddlewareStack middleware.Stack
}
