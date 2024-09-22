package definition

import (
	"net/http"
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMiddleware is a test mock middleware.
type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Called()
	})
}

// Test converting endpoint definitions without any middlewares
func TestEndpointDefinitionsToAPIEndpoints_NoMiddlewares(t *testing.T) {
	// Prepare input
	endpointDefinitions := []EndpointDefinition{
		{
			URL:             "/test-url",
			Method:          "GET",
			MiddlewareStack: middleware.Stack{},
		},
	}

	apiEndpoints := EndpointDefinitionsToAPIEndpoints(endpointDefinitions)

	assert.Len(t, apiEndpoints, 1)
	assert.Equal(t, "/test-url", apiEndpoints[0].URL)
	assert.Equal(t, "GET", apiEndpoints[0].Method)
	assert.Empty(t, apiEndpoints[0].Middlewares)
}

// Test converting endpoint definitions with middlewares
func TestEndpointDefinitionsToAPIEndpoints_WithMiddlewares(t *testing.T) {
	mockMiddleware := MockMiddleware{}

	middleware1 := api.MiddlewareWrapper{
		ID:         "auth",
		Middleware: mockMiddleware.Middleware,
	}
	middleware2 := api.MiddlewareWrapper{
		ID:         "logging",
		Middleware: mockMiddleware.Middleware,
	}

	// Prepare input
	endpointDefinitions := []EndpointDefinition{
		{
			URL:    "/test-url",
			Method: "POST",
			MiddlewareStack: middleware.Stack{
				middleware1,
				middleware2,
			},
		},
	}

	apiEndpoints := EndpointDefinitionsToAPIEndpoints(endpointDefinitions)

	assert.Len(t, apiEndpoints, 1)
	assert.Equal(t, "/test-url", apiEndpoints[0].URL)
	assert.Equal(t, "POST", apiEndpoints[0].Method)
	assert.Len(t, apiEndpoints[0].Middlewares, 2)

	// Call middlewares
	mockMiddleware.On("func1").Return(nil)
	apiEndpoints[0].Middlewares[0](nil).ServeHTTP(nil, nil)
	mockMiddleware.AssertCalled(t, "func1")

	mockMiddleware.On("func1").Return(nil)
	apiEndpoints[0].Middlewares[1](nil).ServeHTTP(nil, nil)
	mockMiddleware.AssertCalled(t, "func1")
}

// Test converting multiple endpoint definitions
func TestEndpointDefinitionsToAPIEndpoints_MultipleDefinitions(t *testing.T) {
	mockMiddleware := MockMiddleware{}

	middleware1 := api.MiddlewareWrapper{
		ID:         "auth",
		Middleware: mockMiddleware.Middleware,
	}

	// Prepare input
	endpointDefinitions := []EndpointDefinition{
		{
			URL:    "/test-url-1",
			Method: "GET",
			MiddlewareStack: middleware.Stack{
				middleware1,
			},
		},
		{
			URL:             "/test-url-2",
			Method:          "POST",
			MiddlewareStack: middleware.Stack{},
		},
	}

	apiEndpoints := EndpointDefinitionsToAPIEndpoints(endpointDefinitions)

	assert.Len(t, apiEndpoints, 2)

	// Check endpoints
	assert.Equal(t, "/test-url-1", apiEndpoints[0].URL)
	assert.Equal(t, "GET", apiEndpoints[0].Method)
	assert.Len(t, apiEndpoints[0].Middlewares, 1)

	assert.Equal(t, "/test-url-2", apiEndpoints[1].URL)
	assert.Equal(t, "POST", apiEndpoints[1].Method)
	assert.Empty(t, apiEndpoints[1].Middlewares)
}

// Test converting empty endpoint definitions (no definitions)
func TestEndpointDefinitionsToAPIEndpoints_EmptyDefinitions(t *testing.T) {
	endpointDefinitions := []EndpointDefinition{}
	apiEndpoints := EndpointDefinitionsToAPIEndpoints(endpointDefinitions)
	assert.Len(t, apiEndpoints, 0)
}
