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
