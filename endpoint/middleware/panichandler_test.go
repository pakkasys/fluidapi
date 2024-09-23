package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

// MockReadCloser is a mock implementation of io.ReadCloser that returns an
// error on Read.
type MockReadCloser struct{}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func (m *MockReadCloser) Close() error {
	return nil
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// TestPanicHandlerMiddlewareWrapper tests the PanicHandlerMiddlewareWrapper
// function.
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

// TestPanicHandlerMiddleware tests the PanicHandlerMiddleware function.
func TestPanicHandlerMiddleware(t *testing.T) {
	var loggedMessages []any
	mockLoggerFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
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

// TestPanicHandlerMiddleware_NoPanic tests the PanicHandlerMiddleware function
// when there is no panic.
func TestPanicHandlerMiddleware_NoPanic(t *testing.T) {
	var loggerCalled bool
	mockLoggerFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
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

// TestPanicHandlerMiddleware_NilLogger tests the PanicHandlerMiddleware
// function when the logger is nil.
func TestPanicHandlerMiddleware_NilLogger(t *testing.T) {
	assert.Panics(t, func() {
		PanicHandlerMiddleware(nil)
	}, "Expected panic when passing nil panicHandlerLoggerFn")
}

// TestHandlePanic tests the handlePanic function.
func TestHandlePanic(t *testing.T) {
	var loggedMessages []any
	mockLoggerFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			loggedMessages = append(loggedMessages, messages...)
		}
	}

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	rw := util.NewResponseWrapper(w)

	err := "test panic"
	handlePanic(rw, req, err, mockLoggerFn)

	// Check that the response has a 500 status code
	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
		"Expected 500 status code",
	)

	// Check that the panic was logged
	assert.Len(t, loggedMessages, 2, "Expected Panic and panicData messages")
	assert.Equal(t, "Panic", loggedMessages[0], "Expected 'Panic' message")

	// Validate the logged panic data
	panicDataLogged := loggedMessages[1].(panicData)
	assert.Equal(t, "test panic", panicDataLogged.Err, "Expected panic message")
	assert.Equal(
		t,
		req.URL.String(),
		panicDataLogged.RequestDump.Request.URL,
		"Expected correct URL in request dump",
	)
	assert.NotEmpty(t, panicDataLogged.StackTrace, "Expected a stack trace")
}

// TestHandlePanic_WithNonNilResponseData tests the handlePanic function with
// non nil response data.
func TestHandlePanic_WithNonNilResponseData(t *testing.T) {
	var loggedMessages []any
	mockLoggerFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			loggedMessages = append(loggedMessages, messages...)
		}
	}

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	rw := util.NewResponseWrapper(w)

	req = req.WithContext(util.NewContext(req.Context()))
	SetResponseWrapper(req, rw)
	rw.Body = []byte("test body")

	err := "test panic"
	handlePanic(rw, req, err, mockLoggerFn)

	// Check that the panic was logged and that response data was included
	assert.Len(t, loggedMessages, 2, "Expected Panic and panicData messages")
	assert.Equal(t, "Panic", loggedMessages[0], "Expected 'Panic' message")

	// Validate the logged panic data
	panicDataLogged := loggedMessages[1].(panicData)
	assert.Equal(
		t,
		req.URL.String(),
		panicDataLogged.RequestDump.Request.URL,
		"Expected correct URL in request dump",
	)
	assert.Equal(
		t,
		"test body",
		panicDataLogged.RequestDump.Response.Body,
		"Expected correct response body in request dump",
	)
	assert.NotEmpty(t, panicDataLogged.StackTrace, "Expected a stack trace")
}

// TestStackTraceSlice tests the stackTraceSlice function.
func TestStackTraceSlice(t *testing.T) {
	stackTrace := stackTraceSlice()
	assert.NotEmpty(t, stackTrace, "Expected a non-empty stack trace")

	for _, entry := range stackTrace {
		// Each entry should be in the format: file.go:line functionName
		assert.True(
			t,
			strings.Contains(entry, ":"),
			"Expected file and line number in the stack trace entry",
		)
		assert.True(
			t,
			strings.Contains(entry, " "),
			"Expected function name in the stack trace entry",
		)
	}

	assert.Contains(
		t,
		stackTrace[1],
		"TestStackTraceSlice",
		"Expected TestStackTraceSlice function to appear in the stack trace",
	)
}

// TestCreateRequestDumpData_ValidRequest tests the case where the request
// body and headers are valid.
func TestCreateRequestDumpData_ValidRequest(t *testing.T) {
	bodyContent := "test body content"
	req := httptest.NewRequest(
		"POST",
		"/test",
		io.NopCloser(strings.NewReader(bodyContent)),
	)

	rd := responseData{
		Body: "response body content",
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		StatusCode: http.StatusOK,
	}

	dumpData := createRequestDumpData(rd, req)

	// Assertions for request part
	h := dumpData.Request.Headers
	assert.Equal(t, "/test", dumpData.Request.URL, "Expected request URL /test")
	assert.Equal(t, bodyContent, dumpData.Request.Body, "Expected body match")
	assert.ElementsMatch(t, req.Header, h, "Expected headers match")
	assert.Empty(t, dumpData.Request.Params, "Expected empty query params")

	// Assertions for response part
	h = dumpData.Response.Headers
	assert.Equal(t, rd.StatusCode, dumpData.StatusCode, "Expected status match")
	assert.Equal(t, rd.Body, dumpData.Response.Body, "Expected body match")
	assert.Equal(t, rd.Headers, h, "Expected headers match")
}

// TestCreateRequestDumpData_ErrorReadingBody tests the case where there is an
// error while reading the request body.
func TestCreateRequestDumpData_ErrorReadingBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", errReader{})

	rd := responseData{
		Body: "response body content",
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		StatusCode: http.StatusOK,
	}

	dumpData := createRequestDumpData(rd, req)

	// Assertions for request part
	assert.Equal(
		t,
		"Error reading request body",
		dumpData.Request.Body,
		"Expected request body to show read error",
	)
}

// TestReadBodyWithLimit_ValidBody tests the case where the body is successfully
// read within the size limit.
func TestReadBodyWithLimit_ValidBody(t *testing.T) {
	bodyContent := "test body content"
	body := io.NopCloser(strings.NewReader(bodyContent))

	result, err := readBodyWithLimit(body, int64(len(bodyContent)+10))

	assert.NoError(t, err, "Expected no error when reading within limit")
	assert.Equal(t, bodyContent, result, "Expected body content to match")
}

// TestReadBodyWithLimit_ExceedsLimit tests the case where the body exceeds the
// size limit and gets truncated.
func TestReadBodyWithLimit_ExceedsLimit(t *testing.T) {
	bodyContent := "this content is too long"
	body := io.NopCloser(strings.NewReader(bodyContent))

	result, err := readBodyWithLimit(body, 10)

	expectedResult := "this conte... (truncated)"
	assert.NoError(t, err, "Expected no error when reading body")
	assert.Equal(t, expectedResult, result, "Expected body to be truncated")
}

// TestReadBodyWithLimit_NilBody tests the case where the body is nil.
func TestReadBodyWithLimit_NilBody(t *testing.T) {
	result, err := readBodyWithLimit(nil, 10)

	assert.NoError(t, err, "Expected no error when reading nil body")
	assert.Equal(t, "", result, "Expected empty string for nil body")
}

// TestReadBodyWithLimit_ZeroSizeLimit tests the case where the size limit is
// zero.
func TestReadBodyWithLimit_ZeroSizeLimit(t *testing.T) {
	bodyContent := "this content should not be read"
	body := io.NopCloser(strings.NewReader(bodyContent))

	result, err := readBodyWithLimit(body, 0)

	expectedResult := "... (truncated)"
	assert.NoError(t, err, "Expected no error when reading body")
	assert.Equal(t, expectedResult, result, "Expected body to be truncated")
}

// TestReadBodyWithLimit_ReadError tests the case where an error occurs while
// reading the body.
func TestReadBodyWithLimit_ReadError(t *testing.T) {
	mockBody := &MockReadCloser{}

	result, err := readBodyWithLimit(mockBody, 10)

	assert.EqualError(t, err, "read error", "Expected a read error")
	assert.Equal(t, "", result, "Expected empty result on read error")
}

// TestLimitHeaders_NoTruncation tests that headers that do not exceed the
// maxHeaderSize are not truncated.
func TestLimitHeaders_NoTruncation(t *testing.T) {
	headers := map[string][]string{
		"X-Test-Header":  {"TestValue"},
		"X-Other-Header": {"ShortValue"},
	}

	limitedHeaders := limitHeaders(headers, 100)

	assert.Equal(
		t,
		headers,
		limitedHeaders,
		"Expected headers to remain unchanged",
	)
}

// TestLimitHeaders_WithTruncation tests that headers exceeding the
// maxHeaderSize are truncated correctly.
func TestLimitHeaders_WithTruncation(t *testing.T) {
	headers := map[string][]string{
		"X-Test-Header":  {"Too long value"},
		"X-Other-Header": {"ShortValue"},
	}

	limitedHeaders := limitHeaders(headers, 10)

	expectedHeaders := map[string][]string{
		"X-Test-Header":  {"Too long v... (truncated)"},
		"X-Other-Header": {"ShortValue"},
	}

	assert.Equal(
		t,
		expectedHeaders,
		limitedHeaders,
		"Expected headers to be truncated",
	)
}

// TestLimitHeaders_MultipleValues tests that multiple values are truncated
// correctly.
func TestLimitHeaders_MultipleValues(t *testing.T) {
	headers := map[string][]string{
		"X-Test-Header": {"Value1", "Value2", "Too long value"},
	}

	limitedHeaders := limitHeaders(headers, 10)

	expectedHeaders := map[string][]string{
		"X-Test-Header": {"Value1", "Value2", "Too long v... (truncated)"},
	}

	assert.Equal(
		t,
		expectedHeaders,
		limitedHeaders,
		"Expected one value to be truncated",
	)
}

// TestLimitHeaders_EmptyHeaders tests that empty headers are not truncated.
func TestLimitHeaders_EmptyHeaders(t *testing.T) {
	headers := map[string][]string{}

	limitedHeaders := limitHeaders(headers, 10)

	assert.Equal(
		t,
		headers,
		limitedHeaders,
		"Expected headers to be unchanged for empty input",
	)
}

// TestLimitHeaders_NoValues tests that keys without values are not truncated.
func TestLimitHeaders_NoValues(t *testing.T) {
	headers := map[string][]string{
		"X-Test-Header": {},
	}

	limitedHeaders := limitHeaders(headers, 10)

	assert.Equal(
		t,
		headers,
		limitedHeaders,
		"Expected headers to remain unchanged for keys without values",
	)
}

// TestLimitQueryParameters_NoTruncation tests that query parameters that
// do not exceed the maxHeaderSize are not truncated.
func TestLimitQueryParameters_NoTruncation(t *testing.T) {
	queryParams := "name=alice&age=30"

	limitedParams := limitQueryParameters(queryParams, 100)

	assert.Equal(
		t,
		queryParams,
		limitedParams,
		"Expected query parameters to remain unchanged",
	)
}

// TestLimitQueryParameters_WithTruncation tests that query parameters
// exceeding the maxHeaderSize are truncated correctly.
func TestLimitQueryParameters_WithTruncation(t *testing.T) {
	queryParams := "name=alice&age=30&location=someverylonglocationstring"

	limitedParams := limitQueryParameters(queryParams, 20)

	expectedParams := "name=alice&age=30&lo... (truncated)"
	assert.Equal(
		t,
		expectedParams,
		limitedParams,
		"Expected query parameters to be truncated",
	)
}

// TestLimitQueryParameters_SmallSizeLimit tests that query parameters
// exactly the maxHeaderSize are truncated correctly.
func TestLimitQueryParameters_ExactSize(t *testing.T) {
	queryParams := "name=alice&age=30"

	limitedParams := limitQueryParameters(queryParams, len(queryParams))

	assert.Equal(
		t,
		queryParams,
		limitedParams,
		"Expected query parameters to remain unchanged for exact size",
	)
}

// TestLimitQueryParameters_EmptyQuery tests that empty query parameters
// are not truncated.
func TestLimitQueryParameters_EmptyQuery(t *testing.T) {
	queryParams := ""

	limitedParams := limitQueryParameters(queryParams, 10)

	assert.Equal(
		t,
		queryParams,
		limitedParams,
		"Expected query parameters to remain unchanged for empty input",
	)
}

// TestLimitQueryParameters_ZeroSizeLimit tests that query parameters
// with zero size limit are truncated correctly.
func TestLimitQueryParameters_ZeroSizeLimit(t *testing.T) {
	queryParams := "name=alice"

	limitedParams := limitQueryParameters(queryParams, 0)

	expectedParams := "... (truncated)"
	assert.Equal(
		t,
		expectedParams,
		limitedParams,
		"Expected query parameters to be truncated even for small size limit",
	)
}
