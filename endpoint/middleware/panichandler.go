package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/util"
)

const (
	PanicHandlerMiddlewareID = "panic_handler"

	maxDumpPartSize = 1024 * 1024
)

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

func PanicHandlerMiddlewareWrapper(
	loggerFn func(r *http.Request) func(messages ...any),
) *api.MiddlewareWrapper {
	return api.NewMiddlewareWrapperBuilder().
		ID(PanicHandlerMiddlewareID).
		Middleware(PanicHandlerMiddleware(loggerFn)).
		Build()
}

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
	panicHandlerLoggerFn(r)(
		"Panic",
		panicData{
			Err: err,
			RequestDump: *createRequestDumpData(
				GetResponseWrapper(r),
				r,
			),
			StackTrace: stackTraceSlice(),
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
	rw *util.ResponseWrapper,
	r *http.Request,
) *requestDumpData {
	requestBody, err := readBodyWithLimit(r.Body, maxDumpPartSize)
	if err != nil {
		requestBody = "Error reading request body"
	}

	responseData := getResponseData(rw)

	return &requestDumpData{
		StatusCode: responseData.StatusCode,
		Request: struct {
			URL     string
			Params  string
			Headers map[string][]string
			Body    string
		}{
			URL:     r.URL.String(),
			Params:  limitQueryParameters(r.URL.RawQuery, maxDumpPartSize),
			Headers: limitHeaders(r.Header, maxDumpPartSize),
			Body:    requestBody,
		},
		Response: struct {
			Headers map[string][]string
			Body    string
		}{
			Headers: limitHeaders(responseData.Headers, maxDumpPartSize),
			Body:    responseData.Body,
		},
	}
}

func getResponseData(rw *util.ResponseWrapper) responseData {
	if rw == nil {
		return responseData{
			StatusCode: 0,
			Headers:    nil,
			Body:       "",
		}
	}

	return responseData{
		StatusCode: rw.StatusCode(),
		Headers:    limitHeaders(rw.Header(), maxDumpPartSize),
		Body:       string(rw.Body()),
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

func limitHeaders(headers map[string][]string, maxHeaderSize int) map[string][]string {
	limitedHeaders := make(map[string][]string)
	for key, values := range headers {
		var limitedValues []string
		for _, value := range values {
			if len(value) > maxHeaderSize {
				limitedValues = append(
					limitedValues,
					value[:maxHeaderSize]+"... (truncated)",
				)
			} else {
				limitedValues = append(limitedValues, value)
			}
		}
		limitedHeaders[key] = limitedValues
	}
	return limitedHeaders
}

func limitQueryParameters(params string, maxParamSize int) string {
	if len(params) > maxParamSize {
		return params[:maxParamSize] + "... (truncated)"
	}
	return params
}
