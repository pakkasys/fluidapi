package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicHandlerMiddlewareWrapper(t *testing.T) {
	var loggedMessages []any
	mockLoggerFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			loggedMessages = append(loggedMessages, messages...)
		}
	}

	wrapper := PanicHandlerMiddlewareWrapper(mockLoggerFn)
	assert.Equal(t, PanicHandlerMiddlewareID, wrapper.ID)

	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		},
	)

	// Call the middleware
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	handler := wrapper.Middleware(mockHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Len(t, loggedMessages, 2, "Expected Panic and panicData messages")
}

func TestPanicHandlerMiddleware(t *testing.T) {
	var loggedMessages []interface{}
	mockLoggerFn := func(r *http.Request) func(messages ...interface{}) {
		return func(messages ...interface{}) {
			loggedMessages = append(loggedMessages, messages...)
		}
	}

	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		},
	)

	// Call the middleware
	middleware := PanicHandlerMiddleware(mockLoggerFn)
	wrappedHandler := middleware(mockHandler)
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	// Check that the response has a 500 status code
	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
		"Expected 500 status code after panic",
	)

	// Check that the panic was logged
	assert.Len(t, loggedMessages, 2, "Expected Panic and panicData messages")
	assert.Equal(t, "Panic", loggedMessages[0], "Expected panic msg")
	assert.IsType(t, panicData{}, loggedMessages[1], "Expected panicData msg")

	// Check that the panic data includes the correct error and stack trace
	panicDataLogged := loggedMessages[1].(panicData)
	assert.Equal(t, "test panic", panicDataLogged.Err, "Expected panic message")
	assert.NotEmpty(t, panicDataLogged.StackTrace, "Expected a stack trace")
}

func TestPanicHandlerMiddleware_NoPanic(t *testing.T) {
	var loggerCalled bool
	mockLoggerFn := func(r *http.Request) func(messages ...interface{}) {
		return func(messages ...interface{}) {
			loggerCalled = true
		}
	}

	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	// Call the middleware
	middleware := PanicHandlerMiddleware(mockLoggerFn)
	wrappedHandler := middleware(mockHandler)
	req := httptest.NewRequest("GET", "/no-panic", nil)
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 status code")
	assert.False(t, loggerCalled, "Expected logger not to be called")
}
