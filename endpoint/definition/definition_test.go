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

// TestCloneEndpointDefinition tests CloneEndpointDefinition with various
// options
func TestCloneEndpointDefinition(t *testing.T) {
	original := &EndpointDefinition{
		URL:    "/original-url",
		Method: "GET",
		MiddlewareStack: middleware.Stack{
			{ID: "auth"},
			{ID: "logging"},
		},
	}

	cloned := CloneEndpointDefinition(
		original,
		WithURL("/new-url"),
		WithMethod("POST"),
		WithMiddlewareStack(
			middleware.Stack{{ID: "new-middleware"}},
		),
	)

	assert.NotSame(t, original, cloned, "cloned object should be new instance")
	assert.Equal(t, "/new-url", cloned.URL, "URL should be updated")
	assert.Equal(t, "POST", cloned.Method, "Method should be updated")
	assert.Equal(
		t,
		middleware.Stack{{ID: "new-middleware"}},
		cloned.MiddlewareStack,
		"MiddlewareStack should be updated",
	)
	assert.Equal(t, "/original-url", original.URL, "URL should not change")
	assert.Equal(t, "GET", original.Method, "method should not change")
	assert.Equal(
		t, middleware.Stack{{ID: "auth"},
			{ID: "logging"}}, original.MiddlewareStack,
		"original MiddlewareStack should remain unchanged",
	)
}

// TestWithURL tests the WithURL function
func TestWithURL(t *testing.T) {
	original := &EndpointDefinition{URL: "/old-url"}
	option := WithURL("/new-url")
	option(original)

	assert.Equal(t, "/new-url", original.URL, "URL should be /new-url")
}

// TestWithMethod tests the WithMethod function
func TestWithMethod(t *testing.T) {
	original := &EndpointDefinition{Method: "GET"}
	option := WithMethod("POST")
	option(original)

	assert.Equal(t, "POST", original.Method, "Method should be updated to POST")
}

// TestWithMiddlewareWrappers tests the WithMiddlewareWrappers function
func TestWithMiddlewareWrappers(t *testing.T) {
	original := &EndpointDefinition{
		MiddlewareStack: middleware.Stack{{ID: "auth"}},
	}
	newMiddleware := middleware.Stack{{ID: "logging"}}
	option := WithMiddlewareStack(newMiddleware)
	option(original)

	assert.Equal(
		t,
		newMiddleware,
		original.MiddlewareStack,
		"MiddlewareStack should be updated to logging",
	)
}

// TestWithMiddlewareWrappersFunc tests the WithMiddlewareWrappersFunc function
func TestWithMiddlewareWrappersFunc(t *testing.T) {
	original := &EndpointDefinition{
		MiddlewareStack: middleware.Stack{{ID: "auth"}},
	}
	option := WithMiddlewareWrappersFunc(
		func(e *EndpointDefinition) middleware.Stack {
			return middleware.Stack{{ID: "dynamic-middleware"}}
		},
	)
	option(original)

	assert.Equal(
		t,
		middleware.Stack{{ID: "dynamic-middleware"}},
		original.MiddlewareStack,
		"MiddlewareStack should be dynamically updated",
	)
}
