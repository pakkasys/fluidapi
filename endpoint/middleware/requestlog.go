package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/pakkasys/fluidapi/core/api"
)

const RequestLogMiddlewareID = "request_log"

// GetRequestMetadataFunc is a function type used for retrieving RequestMetadata
// from a context.
type GetRequestMetadataFunc func(ctx context.Context) *RequestMetadata

type requestLog struct {
	StartTime     time.Time `json:"start_time"`     // Start time of the request.
	RemoteAddress string    `json:"remote_address"` // Remote IP address of the client making the request.
	Protocol      string    `json:"protocol"`       // Protocol used in the request (e.g., HTTP/1.1).
	HTTPMethod    string    `json:"http_method"`    // HTTP method used for the request.
	URL           string    `json:"url"`            // Full URL of the request.
}

// RequestLogMiddlewareWrapper creates a new MiddlewareWrapper for the
// Request Log middleware. This middleware logs request information.
//
//   - requestLoggerFn: A function that logs messages for the request.
func RequestLogMiddlewareWrapper(
	requestLoggerFn func(r *http.Request) func(messages ...any),
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID:         RequestLogMiddlewareID,
		Middleware: RequestLogMiddleware(GetRequestMetadata, requestLoggerFn),
	}
}

// RequestLogMiddleware constructs a middleware that logs information about
// incoming requests. It uses the given request metadata retrieval function and
// the request logger function.
//
//   - getMetadataFn: A function that gets metadata from the request context.
//   - requestLoggerFn: A function that logs messages for the request.
func RequestLogMiddleware(
	getMetadataFn GetRequestMetadataFunc,
	requestLoggerFn func(r *http.Request) func(messages ...any),
) api.Middleware {
	if requestLoggerFn == nil {
		panic("requestLoggerFn cannot be nil")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logRequest(r, getMetadataFn, requestLoggerFn)
			next.ServeHTTP(w, r)
			requestLoggerFn(r)("Request completed")
		})
	}
}

func logRequest(
	r *http.Request,
	getMetadataFn GetRequestMetadataFunc,
	requestLoggerFn func(r *http.Request) func(messages ...any),
) {
	requestMetadata := getMetadataFn(r.Context())
	if requestMetadata == nil {
		requestLoggerFn(r)("Request started", "Request metadata not found")
	} else {
		requestLoggerFn(r)(
			"Request started",
			requestLog{
				StartTime:     time.Now().UTC(),
				RemoteAddress: requestMetadata.RemoteAddress,
				Protocol:      requestMetadata.Protocol,
				HTTPMethod:    requestMetadata.HTTPMethod,
				URL:           requestMetadata.URL,
			},
		)
	}
}
