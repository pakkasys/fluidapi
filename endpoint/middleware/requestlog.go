package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/pakkasys/fluidapi/core/api"
)

const RequestLogMiddlewareID = "request_log"

type GetRequestMetadataFunc func(ctx context.Context) *RequestMetadata

type requestLog struct {
	StartTime     time.Time `json:"start_time"`
	RemoteAddress string    `json:"remote_address"`
	Protocol      string    `json:"protocol"`
	HTTPMethod    string    `json:"http_method"`
	URL           string    `json:"url"`
}

func RequestLogMiddlewareWrapper(
	requestLoggerFn func(r *http.Request) func(messages ...any),
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID:         RequestLogMiddlewareID,
		Middleware: RequestLogMiddleware(GetRequestMetadata, requestLoggerFn),
	}
}

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

			if requestLoggerFn != nil {
				requestLoggerFn(r)("Request completed")
			}
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
