package definition

import (
	"testing"

	"github.com/pakkasys/fluidapi/endpoint/middleware"
	"github.com/stretchr/testify/assert"
)

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
