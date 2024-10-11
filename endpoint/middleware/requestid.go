package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/util"
)

const RequestIDMiddlewareID = "request_metadata"

var dataKey = util.NewDataKey()

// RequestMetadata represents metadata information associated with an HTTP
// request.
type RequestMetadata struct {
	TimeStart     time.Time // Time when the request started.
	RequestID     string    // Unique identifier for the request.
	RemoteAddress string    // Remote IP address of the request.
	Protocol      string    // Protocol used in the request (e.g., HTTP/1.1).
	HTTPMethod    string    // HTTP method used for the request (e.g., GET).
	URL           string    // URL of the request.
}

// RequestIDMiddlewareWrapper creates a new MiddlewareWrapper for the Request ID
// middleware. This middleware generates a unique ID for each request and stores
// it as metadata in the context.
//
//   - requestIDFn: A function that generates a unique request ID.
func RequestIDMiddlewareWrapper(
	requestIDFn func() string,
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID:         RequestIDMiddlewareID,
		Middleware: RequestIDMiddleware(requestIDFn),
	}
}

// RequestIDMiddleware constructs a middleware function that generates request
// metadata and stores it in the request's context. This metadata can be used
// for logging and tracking purposes.
//
//   - requestIDFn: A function that generates a unique request ID.
func RequestIDMiddleware(requestIDFn func() string) api.Middleware {
	if requestIDFn == nil {
		panic("requestIDFn cannot be nil")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestMetadata := RequestMetadata{
				TimeStart:     time.Now().UTC(),
				RequestID:     requestIDFn(),
				RemoteAddress: util.RequestIPAddress(r),
				Protocol:      r.Proto,
				HTTPMethod:    r.Method,
				URL:           fmt.Sprintf("%s%s", r.Host, r.URL.Path),
			}
			util.SetContextValue(r.Context(), dataKey, &requestMetadata)
			next.ServeHTTP(w, r)
		})
	}
}

// GetRequestMetadata retrieves the RequestMetadata from the given context.
// If no metadata is found, it returns nil.
//
//   - ctx: The context from which to retrieve the metadata.
func GetRequestMetadata(ctx context.Context) *RequestMetadata {
	return util.GetContextValue[*RequestMetadata](ctx, dataKey, nil)
}
