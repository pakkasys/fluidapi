package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

// TestResponseWrapperMiddlewareWrapper tests the
// ResponseWrapperMiddlewareWrapper function
func TestResponseWrapperMiddlewareWrapper(t *testing.T) {
	wrapper := ResponseWrapperMiddlewareWrapper()

	expectID := ResponseWrapperMiddlewareID
	assert.NotNil(t, wrapper, "MiddlewareWrapper should not be nil")
	assert.Equal(t, expectID, wrapper.ID, "MiddlewareWrapper IDs should match")
	assert.NotNil(t, wrapper.Middleware, "Middleware func should not be nil")

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))
	w := httptest.NewRecorder()

	// Define a next handler that will be called by the middleware
	nextHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Next handler called"))
			assert.NoError(t, err, "Expected no error when writing response")
		},
	)

	handler := wrapper.Middleware(nextHandler)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status OK")
	assert.Equal(
		t,
		"Next handler called",
		w.Body.String(),
		"Expected 'Next handler called' body",
	)
}

// TestResponseWrapperMiddleware tests the ResponseWrapperMiddleware function
func TestResponseWrapperMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))
	w := httptest.NewRecorder()

	middleware := ResponseWrapperMiddleware(
		util.NewRequestWrapper,
		util.NewResponseWrapper,
	)

	// Define a next handler that will be called by the middleware
	nextHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqW := GetRequestWrapper(r)
			resW := GetResponseWrapper(r)

			assert.NotNil(t, reqW, "Request wrapper should not be nil")
			assert.Equal(t, req, reqW.Request, "Requests should match")

			assert.NotNil(t, resW, "Response wrapper should not be nil")

			// Check that the response writer is an httptest.ResponseRecorder
			recorder, ok := resW.ResponseWriter.(*httptest.ResponseRecorder)
			if ok {
				assert.Equal(
					t,
					resW.ResponseWriter,
					recorder,
					"Writers should match",
				)
			} else {
				t.Errorf(
					"Should be ResponseRecorder, got %T",
					resW.ResponseWriter,
				)
			}

			// Continue the response
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Success"))
			assert.NoError(t, err)
		},
	)

	handler := middleware(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status 200")
	assert.Equal(t, "Success", w.Body.String(), "Expected 'Success' body")
}

// TestResponseWrapperMiddleware_RequestWrapperError tests the scenario where
// util.NewRequestWrapper returns an error.
func TestResponseWrapperMiddleware_RequestWrapperError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	mockCreateRequestWrapper := func(
		r *http.Request,
	) (*util.RequestWrapper, error) {
		return nil, errors.New("mock error")
	}

	middleware := ResponseWrapperMiddleware(
		mockCreateRequestWrapper,
		util.NewResponseWrapper,
	)
	nextHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("Next handler should not be called")
			w.WriteHeader(http.StatusOK)
		},
	)

	handler := middleware(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Result().StatusCode, "Expected status 500",
	)
	assert.Equal(
		t,
		http.StatusText(http.StatusInternalServerError)+"\n", w.Body.String(),
		"Expected error body",
	)
}

// TestGetRequestWrapper verifies that the request wrapper can be correctly
// retrieved from the request context.
func TestGetRequestWrapper(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))

	mockWrapper, err := util.NewRequestWrapper(req)
	assert.NoError(t, err, "Expected no error")

	setRequestWrapper(req, mockWrapper)

	getWrapper := GetRequestWrapper(req)

	assert.NotNil(t, getWrapper, "Wrapper should not be nil")
	assert.Equal(t, mockWrapper, getWrapper, "Wrappers should match")
	assert.Equal(t, req, getWrapper.Request, "Requests should match")
}

// TestGetResponseWrapper verifies that the response wrapper can be correctly
// retrieved from the request context.
func TestGetResponseWrapper(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))

	w := httptest.NewRecorder()

	mockWrapper := util.NewResponseWrapper(w)
	setResponseWrapper(req, mockWrapper)

	getWrapper := GetResponseWrapper(req)

	assert.NotNil(t, getWrapper, "Wrapper should not be nil")
	assert.Equal(t, mockWrapper, getWrapper, "Wrappers should match")
	assert.Equal(t, w, getWrapper.ResponseWriter, "Responses should match")
}

// TestSetResponseWrapper tests the SetResponseWrapper function
func TestSetResponseWrapper(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))

	w := httptest.NewRecorder()

	mockWrapper := util.NewResponseWrapper(w)

	setResponseWrapper(req, mockWrapper)

	getWrapper := GetResponseWrapper(req)

	assert.NotNil(t, getWrapper, "Wrapper should not be nil")
	assert.Equal(t, mockWrapper, getWrapper, "Wrappers should match")
	assert.Equal(t, w, getWrapper.ResponseWriter, "Responses should match")
}

// TestSetRequestWrapper tests the SetRequestWrapper function
func TestSetRequestWrapper(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))
	mockWrapper, err := util.NewRequestWrapper(req)
	assert.NoError(t, err, "Expected no error")

	setRequestWrapper(req, mockWrapper)

	getWrapper := GetRequestWrapper(req)

	assert.NotNil(t, getWrapper, "Request wrapper should not be nil")
	assert.Equal(t, mockWrapper, getWrapper, "Wrappers should match")
	assert.Equal(t, req, getWrapper.Request, "Requests should match")
}
