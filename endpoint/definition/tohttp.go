package definition

import (
	"github.com/pakkasys/fluidapi/core/api"
)

func EndpointDefinitionsToHTTPEndpoints(
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
			*api.NewEndpoint(
				endpointDefinition.URL,
				endpointDefinition.Method,
				middlewares,
			),
		)
	}

	return endpoints
}
