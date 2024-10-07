package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/stretchr/testify/assert"
)

// MockMiddleware is a simple middleware for testing.
func MockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Middleware", "Mocked")
		next.ServeHTTP(w, r)
	})
}

// TestMiddlewares_EmptyStack tests the Middlewares function when the stack is
// empty.
func TestMiddlewares_EmptyStack(t *testing.T) {
	mwStack := Stack{}

	middlewares := mwStack.Middlewares()

	assert.Empty(t, middlewares, "Expected middlewares to be empty for an empty stack")
}

// TestMiddlewares_StackWithMiddlewares tests the Middlewares function when the
// stack has middlewares.
func TestMiddlewares_StackWithMiddlewares(t *testing.T) {
	// Create some mock middleware wrappers.
	mw1 := api.MiddlewareWrapper{
		ID:         "auth",
		Middleware: MockMiddleware,
	}
	mw2 := api.MiddlewareWrapper{
		ID:         "logging",
		Middleware: MockMiddleware,
	}

	// Create a middleware stack with these wrappers.
	mwStack := Stack{mw1, mw2}

	middlewares := mwStack.Middlewares()

	assert.Equal(t, 2, len(middlewares), "Middleware stack should have 2 middlewares")
	assert.NotNil(t, middlewares[0], "First middleware should not be nil")
	assert.NotNil(t, middlewares[1], "Second middleware should not be nil")
}

// TestMiddlewares_Order tests the order of the middlewares returned by
// Middlewares function.
func TestMiddlewares_Order(t *testing.T) {
	callOrder := []string{}

	// Middleware 1: Add "first" to a shared slice
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "first")
			next.ServeHTTP(w, r)
		})
	}

	// Middleware 2: Add "second" to the shared slice
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "second")
			next.ServeHTTP(w, r)
		})
	}

	// Create the middleware stack
	mwStack := Stack{
		api.MiddlewareWrapper{ID: "first", Middleware: mw1},
		api.MiddlewareWrapper{ID: "second", Middleware: mw2},
	}

	// Get the middleware functions from the stack
	middlewares := mwStack.Middlewares()

	// Define a final handler
	finalHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {},
	)

	// Apply middlewares to the final handler
	wrappedHandler := api.ApplyMiddlewares(finalHandler, middlewares...)

	// Create a slice to track the middleware execution order
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req = req.WithContext(context.Background())

	// Perform the request
	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	// Verify the expected order of middleware execution
	assert.Equal(t, []string{"first", "second"}, callOrder, "Middlewares should be executed in the correct order")
}

// TestInsertAfterID_Success tests the InsertAfterID function when the
// middleware is inserted successfully.
func TestInsertAfterID_Success(t *testing.T) {
	mw1 := api.MiddlewareWrapper{ID: "auth"}
	mw2 := api.MiddlewareWrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	newMiddleware := api.MiddlewareWrapper{ID: "metrics"}

	inserted := mwStack.InsertAfterID("auth", newMiddleware)

	assert.True(t, inserted, "Middleware should be inserted")
	assert.Equal(t, 3, len(mwStack), "Middleware stack should have 3 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "metrics", mwStack[1].ID, "Middleware not in 2nd position")
	assert.Equal(t, "logging", mwStack[2].ID)
}

// TestInsertAfterID_AppendToEnd tests the InsertAfterID function when the
// middleware is appended to the end.
func TestInsertAfterID_AppendToEnd(t *testing.T) {
	mw1 := api.MiddlewareWrapper{ID: "auth"}
	mw2 := api.MiddlewareWrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	newMiddleware := api.MiddlewareWrapper{ID: "metrics"}

	inserted := mwStack.InsertAfterID("logging", newMiddleware)

	assert.True(t, inserted, "Middleware should be inserted")
	assert.Equal(t, 3, len(mwStack), "Middleware stack should have 3 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "logging", mwStack[1].ID)
	assert.Equal(t, "metrics", mwStack[2].ID, "New middleware not in the end")
}

// TestInsertAfterID_IDNotFound tests the InsertAfterID function when the
// middleware ID is not found.
func TestInsertAfterID_IDNotFound(t *testing.T) {
	mw1 := api.MiddlewareWrapper{ID: "auth"}
	mw2 := api.MiddlewareWrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	newMiddleware := api.MiddlewareWrapper{ID: "metrics"}

	inserted := mwStack.InsertAfterID("non-existent-id", newMiddleware)

	assert.False(t, inserted, "Middleware should not be inserted")
	assert.Equal(t, 2, len(mwStack), "Middleware stack should have 2 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "logging", mwStack[1].ID)
}
