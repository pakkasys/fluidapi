package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

// TestRequestIDMiddlewareWrapper tests the RequestIDMiddlewareWrapper function
func TestRequestIDMiddlewareWrapper(t *testing.T) {
	requestIDFn := func() string {
		return "test-request-id"
	}

	wrapper := RequestIDMiddlewareWrapper(requestIDFn)
	assert.NotNil(t, wrapper, "MiddlewareWrapper should not be nil")
	assert.Equal(t, RequestIDMiddlewareID, wrapper.ID, "Middleware ID should match")
	assert.NotNil(t, wrapper.Middleware, "Middleware should not be nil")
}

// TestRequestIDMiddleware_Success tests the RequestIDMiddleware function
func TestRequestIDMiddleware_Success(t *testing.T) {
	requestIDFn := func() string {
		return "test-request-id"
	}

	middleware := RequestIDMiddleware(requestIDFn)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	req = req.WithContext(util.NewContext(req.Context()))
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadata := GetRequestMetadata(r.Context())

		assert.NotNil(t, metadata, "Request metadata should not be nil")
		assert.Equal(t, "test-request-id", metadata.RequestID, "Request ID should match")
		assert.Equal(t, "127.0.0.1", metadata.RemoteAddress, "Remote address should match")
		assert.Equal(t, "HTTP/1.1", metadata.Protocol, "Protocol should match")
		assert.Equal(t, "GET", metadata.HTTPMethod, "HTTP method should match")
		assert.Equal(t, "example.com/test", metadata.URL, "URL should match")
	})

	handler := middleware(nextHandler)
	req.Host = "example.com"
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status 200")
}

// TestRequestIDMiddleware_NilRequestIDFn tests that RequestIDMiddleware panics
// when the requestIDFn is nil.
func TestRequestIDMiddleware_NilRequestIDFn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when requestIDFn is nil")
		}
	}()

	_ = RequestIDMiddleware(nil)
}

// TestGetRequestMetadata_MetadataExists tests the scenario where metadata is
// found in the context.
func TestGetRequestMetadata_MetadataExists(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := util.NewContext(req.Context())
	req = req.WithContext(ctx)

	expectedMetadata := &RequestMetadata{
		TimeStart:     time.Now().UTC(),
		RequestID:     "test-request-id",
		RemoteAddress: "127.0.0.1",
		Protocol:      "HTTP/1.1",
		HTTPMethod:    "GET",
		URL:           "example.com/test",
	}

	util.SetContextValue(req.Context(), dataKey, expectedMetadata)

	metadata := GetRequestMetadata(ctx)

	assert.NotNil(t, metadata, "Metadata should not be nil")
	assert.Equal(t, expectedMetadata, metadata, "Metadata should match")
}

// TestGetRequestMetadata_NoMetadata tests the scenario where no metadata exists
// in the context.
func TestGetRequestMetadata_NoMetadata(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))

	metadata := GetRequestMetadata(req.Context())

	assert.Nil(t, metadata, "Metadata should be nil")
}
