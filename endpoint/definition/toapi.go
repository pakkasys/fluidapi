package definition

import (
	"github.com/pakkasys/fluidapi/core/api"
)

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
