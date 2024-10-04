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

type RequestMetadata struct {
	TimeStart     time.Time
	RequestID     string
	RemoteAddress string
	Protocol      string
	HTTPMethod    string
	URL           string
}

func RequestIDMiddlewareWrapper(
	requestIDFn func() string,
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID:         RequestIDMiddlewareID,
		Middleware: RequestIDMiddleware(requestIDFn),
	}
}

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
				URL:           fmt.Sprintf("%s%s", r.Host, r.URL),
			}
			util.SetContextValue(
				r.Context(),
				dataKey,
				&requestMetadata,
			)
			next.ServeHTTP(w, r)
		})
	}
}

func GetRequestMetadata(ctx context.Context) *RequestMetadata {
	return util.GetContextValue[*RequestMetadata](
		ctx,
		dataKey,
		nil,
	)
}
