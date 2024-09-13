package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestApplyMiddlewares tests the ApplyMiddlewares function.
func TestApplyMiddlewares(t *testing.T) {
	// Simple middleware that adds a header
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "true")
			next.ServeHTTP(w, r)
		})
	}

	// Simple handler that returns OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Apply middleware
	ApplyMiddlewares(handler, testMiddleware).ServeHTTP(rr, req)

	// Test that middleware worked
	if rr.Header().Get("X-Test") != "true" {
		t.Errorf(
			"expected X-Test header to be true, got %s",
			rr.Header().Get("X-Test"),
		)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code to be 200, got %d", status)
	}
}

// TestMiddlewareWrapperBuilder tests the MiddlewareWrapperBuilder.
func TestMiddlewareWrapperBuilder(t *testing.T) {
	builder := NewMiddlewareWrapperBuilder()
	middlewareWrapper := builder.ID("test_id").Build()

	if middlewareWrapper.ID != "test_id" {
		t.Errorf(
			"expected middleware ID to be test_id, got %s",
			middlewareWrapper.ID,
		)
	}
}

// TestNewMiddlewareInput tests the creation of a new MiddlewareInput.
func TestNewMiddlewareInput(t *testing.T) {
	input := "test_input"
	middlewareInput := NewMiddlewareInput(input)

	if middlewareInput.Input != input {
		t.Errorf(
			"expected middleware input to be %v, got %v",
			input,
			middlewareInput.Input,
		)
	}
}

// TestMiddlewareWrapperBuilderMiddleware tests setting middleware in the
// builder.
func TestMiddlewareWrapperBuilderMiddleware(t *testing.T) {
	builder := NewMiddlewareWrapperBuilder()
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	middlewareWrapper := builder.Middleware(testMiddleware).Build()

	if middlewareWrapper.Middleware == nil {
		t.Error("expected middleware to be set, got nil")
	}
}

// TestMiddlewareWrapperBuilderMiddlewareInputs tests setting middleware inputs
// in the builder.
func TestMiddlewareWrapperBuilderMiddlewareInputs(t *testing.T) {
	builder := NewMiddlewareWrapperBuilder()
	middlewareInputs := []MiddlewareInput{
		NewMiddlewareInput("input1"),
		NewMiddlewareInput("input2")}
	middlewareWrapper := builder.MiddlewareInputs(middlewareInputs).Build()

	if len(middlewareWrapper.MiddlewareInputs) != 2 {
		t.Errorf(
			"expected 2 middleware inputs, got %d",
			len(middlewareWrapper.MiddlewareInputs),
		)
	}
}

// TestMiddlewareWrapperBuilderAddMiddlewareInput tests adding middleware input
// in the builder.
func TestMiddlewareWrapperBuilderAddMiddlewareInput(t *testing.T) {
	builder := NewMiddlewareWrapperBuilder()
	middlewareInput := NewMiddlewareInput("test_input")
	builder.AddMiddlewareInput(middlewareInput)
	middlewareWrapper := builder.Build()

	if len(middlewareWrapper.MiddlewareInputs) != 1 {
		t.Errorf(
			"expected 1 middleware input, got %d",
			len(middlewareWrapper.MiddlewareInputs),
		)
	}

	if middlewareWrapper.MiddlewareInputs[0].Input != "test_input" {
		t.Errorf(
			"expected middleware input to be 'test_input', got %v",
			middlewareWrapper.MiddlewareInputs[0].Input,
		)
	}
}
