package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/pakkasys/fluidapi/core/api"
)

const (
	PanicHandlerMiddlewareID = "panic_handler"

	maxDumpSize = 1024 * 1024
)

// ResponseWrapper is an interface that wraps an http.ResponseWriter.
type ResponseWrapper interface {
	StatusCode() int
	Body() []byte
	Header() http.Header
}

type requestDumpData struct {
	StatusCode int
	Request    struct {
		URL     string
		Params  string
		Headers map[string][]string
		Body    string
	}
	Response struct {
		Headers map[string][]string
		Body    string
	}
}

type responseData struct {
	StatusCode int
	Headers    map[string][]string
	Body       string
}

type panicData struct {
	Err         any             `json:"err"`
	RequestDump requestDumpData `json:"request_dump"`
	StackTrace  []string        `json:"stack_trace"`
}

// PanicHandlerMiddlewareWrapper creates a new MiddlewareWrapper for
// the Panic Handler middleware. This middleware catches and logs any panics
// during the request lifecycle.
//
//   - loggerFn: A function that logs panic information for the request.
func PanicHandlerMiddlewareWrapper(
	loggerFn func(r *http.Request) func(messages ...any),
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID:         PanicHandlerMiddlewareID,
		Middleware: PanicHandlerMiddleware(loggerFn),
	}
}

// PanicHandlerMiddleware constructs a middleware that captures and logs any
// panic events during request handling. It uses the provided panic handler
// logger function to log the details.
//
//   - panicHandlerLoggerFn: A function that logs messages in the event of a panic.
func PanicHandlerMiddleware(
	panicHandlerLoggerFn func(r *http.Request) func(messages ...any),
) api.Middleware {
	if panicHandlerLoggerFn == nil {
		panic("panicHandlerLoggerFn cannot be nil")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					handlePanic(w, r, err, panicHandlerLoggerFn)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func handlePanic(
	w http.ResponseWriter,
	r *http.Request,
	err any,
	panicHandlerLoggerFn func(r *http.Request) func(messages ...any),
) {
	var rd responseData
	rw := GetResponseWrapper(r)
	if rw != nil {
		rd = responseData{
			StatusCode: rw.StatusCode,
			Headers:    limitHeaders(rw.Header(), maxDumpSize),
			Body:       string(rw.Body),
		}
	}

	panicHandlerLoggerFn(r)(
		"Panic",
		panicData{
			Err:         err,
			RequestDump: *createRequestDumpData(rd, r),
			StackTrace:  stackTraceSlice(),
		},
	)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func stackTraceSlice() []string {
	var stackTrace []string
	var skip int

	for {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		// Get the function name and format entry.
		fn := runtime.FuncForPC(pc)
		entry := fmt.Sprintf("%s:%d %s", file, line, fn.Name())
		stackTrace = append(stackTrace, entry)

		skip++
	}

	return stackTrace
}

func createRequestDumpData(
	rd responseData,
	r *http.Request,
) *requestDumpData {
	requestBody, err := readBodyWithLimit(r.Body, maxDumpSize)
	if err != nil {
		requestBody = "Error reading request body"
	}

	return &requestDumpData{
		StatusCode: rd.StatusCode,
		Request: struct {
			URL     string
			Params  string
			Headers map[string][]string
			Body    string
		}{
			URL:     r.URL.String(),
			Params:  limitQueryParameters(r.URL.RawQuery, maxDumpSize),
			Headers: limitHeaders(r.Header, maxDumpSize),
			Body:    requestBody,
		},
		Response: struct {
			Headers map[string][]string
			Body    string
		}{
			Headers: limitHeaders(rd.Headers, maxDumpSize),
			Body:    rd.Body,
		},
	}
}

func readBodyWithLimit(body io.ReadCloser, maxSize int64) (string, error) {
	if body == nil {
		return "", nil
	}
	defer body.Close()

	// Limit the reader to the max size
	limitedReader := io.LimitReader(body, maxSize)

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(limitedReader)
	if err != nil {
		return "", err
	}

	// Check if the body was truncated
	if buf.Len() == int(maxSize) {
		return buf.String() + "... (truncated)", nil
	}

	return buf.String(), nil
}

func limitHeaders(
	headers map[string][]string,
	maxSize int,
) map[string][]string {
	limitedHeaders := make(map[string][]string)
	for key, values := range headers {
		var limitedValues []string
		if len(values) == 0 {
			limitedHeaders[key] = values
			continue
		}
		for _, value := range values {
			if len(value) > maxSize {
				limitedValues = append(
					limitedValues,
					value[:maxSize]+"... (truncated)",
				)
			} else {
				limitedValues = append(limitedValues, value)
			}
		}
		limitedHeaders[key] = limitedValues
	}
	return limitedHeaders
}

func limitQueryParameters(params string, maxSize int) string {
	if len(params) > maxSize {
		return params[:maxSize] + "... (truncated)"
	}
	return params
}
