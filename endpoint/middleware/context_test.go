package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

func TestContextMiddlewareWrapper(t *testing.T) {
	// Call the function
	wrapper := ContextMiddlewareWrapper()

	// Check that the wrapper has the correct ID
	assert.Equal(
		t,
		ContextMiddlewareID,
		wrapper.ID,
		"Expected MiddlewareWrapper to have the correct ID",
	)

	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		},
	)

	// Call the middleware
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler := wrapper.Middleware(mockHandler)
	handler.ServeHTTP(w, req)

	// Check that the middleware chain executed correctly
	assert.Equal(
		t,
		http.StatusTeapot,
		w.Code,
		"Expected the middleware to return specified status",
	)
}

func TestContextMiddleware(t *testing.T) {
	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			newCtx := util.IsContextSet(r.Context())
			assert.NotNil(t, newCtx, "Expected the context to be set")
			w.WriteHeader(http.StatusOK)
		},
	)

	// Call the middleware
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler := ContextMiddleware()(mockHandler)

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status OK")
}
