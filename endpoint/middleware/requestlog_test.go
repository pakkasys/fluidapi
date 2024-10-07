package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the logger function.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Log(messages ...any) {
	m.Called(messages)
}

// TestRequestLogMiddlewareWrapper tests the RequestLogMiddlewareWrapper
// function
func TestRequestLogMiddlewareWrapper(t *testing.T) {
	mockLogger := new(MockLogger)

	requestLoggerFn := func(r *http.Request) func(messages ...any) {
		return mockLogger.Log
	}

	middlewareWrapper := RequestLogMiddlewareWrapper(requestLoggerFn)

	assert.NotNil(t, middlewareWrapper, "Wrapper should not be nil")
	assert.Equal(
		t,
		RequestLogMiddlewareID,
		middlewareWrapper.ID,
		"Middleware ID should match",
	)
	assert.NotNil(
		t,
		middlewareWrapper.Middleware,
		"Middleware func should not be nil",
	)
}

// TestRequestLogMiddleware tests the RequestLogMiddleware function.
func TestRequestLogMiddleware(t *testing.T) {
	mockLogger := new(MockLogger)

	requestLoggerFn := func(r *http.Request) func(messages ...any) {
		return mockLogger.Log
	}

	middleware := RequestLogMiddleware(GetRequestMetadata, requestLoggerFn)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Request passed through"))
			assert.NoError(t, err)
		},
	)

	mockLogger.On("Log", mock.AnythingOfType("[]interface {}")).Return()

	handler := middleware(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status 200")
	assert.Equal(
		t,
		"Request passed through",
		w.Body.String(),
		"Expected 'Request passed through' body",
	)

	mockLogger.AssertNumberOfCalls(t, "Log", 2)
	mockLogger.AssertExpectations(t)
}

// TestRequestLogMiddleware_NilLogger tests the scenario where the
// requestLoggerFn is nil.
func TestRequestLogMiddleware_NilLogger(t *testing.T) {
	assert.Panics(t, func() {
		_ = RequestLogMiddleware(GetRequestMetadata, nil)
	}, "Expected panic with nil requestLoggerFn")
}

// TestLogRequest_MetadataNotFound tests the scenario where the request metadata
// is not found in the context.
func TestLogRequest_MetadataNotFound(t *testing.T) {
	mockLogger := new(MockLogger)

	mockLogger.On("Log", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(context.Background())

	getMetadataFunc := func(ctx context.Context) *RequestMetadata {
		return nil
	}

	logRequest(
		req,
		getMetadataFunc,
		func(r *http.Request) func(messages ...any) {
			return mockLogger.Log
		},
	)

	mockLogger.AssertExpectations(t)
}

// TestLogRequest_MetadataFound tests the scenario where the request metadata
// is found in the context.
func TestLogRequest_MetadataFound(t *testing.T) {
	mockLogger := new(MockLogger)

	expectedMetadata := &RequestMetadata{
		TimeStart:     time.Now().UTC(),
		RemoteAddress: "127.0.0.1",
		Protocol:      "HTTP/1.1",
		HTTPMethod:    "GET",
		URL:           "/test",
	}

	mockLogger.On(
		"Log",
		mock.MatchedBy(func(args []any) bool {
			if len(args) != 2 {
				return false
			}
			if args[0] != "Request started" {
				return false
			}
			_, ok := args[1].(requestLog)
			return ok
		}),
	).Return()

	getMetadataFunc := func(ctx context.Context) *RequestMetadata {
		return expectedMetadata
	}

	req := httptest.NewRequest("GET", "/test", nil)

	logRequest(
		req,
		getMetadataFunc, func(r *http.Request) func(messages ...any) {
			return mockLogger.Log
		},
	)

	mockLogger.AssertExpectations(t)
}
