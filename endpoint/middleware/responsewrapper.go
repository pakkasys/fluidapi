package middleware

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/util"
)

const ResponseWrapperMiddlewareID = "response_wrapper"

var (
	responseDataKey = util.NewDataKey()
	requestDataKey  = util.NewDataKey()
)

// ResponseWrapperMiddlewareWrapper creates a new MiddlewareWrapper for the
// Response Wrapper middleware.
func ResponseWrapperMiddlewareWrapper() *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID: ResponseWrapperMiddlewareID,
		Middleware: ResponseWrapperMiddleware(
			util.NewRequestWrapper,
			util.NewResponseWrapper,
		),
	}
}

// ResponseWrapperMiddleware constructs a middleware function that wraps the
// request and response.
//
//   - requestWrapperFn: A function that wraps the HTTP request.
//   - responseWrapperFn: A function that wraps the HTTP response.
func ResponseWrapperMiddleware(
	requestWrapperFn func(*http.Request) (*util.RequestWrapper, error),
	responseWrapperFn func(http.ResponseWriter) *util.ResponseWrapper,
) api.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseWrapper := util.NewResponseWrapper(w)

			requestWrapper, err := requestWrapperFn(r)
			if err != nil {
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}

			setRequestWrapper(r, requestWrapper)
			setResponseWrapper(r, responseWrapper)

			next.ServeHTTP(responseWrapper, requestWrapper.Request)
		})
	}
}

// GetResponseWrapper retrieves the `ResponseWrapper` from the request context.
//
//   - r: The HTTP request from which to retrieve the `ResponseWrapper`.
func GetResponseWrapper(r *http.Request) *util.ResponseWrapper {
	return util.GetContextValue[*util.ResponseWrapper](
		r.Context(),
		responseDataKey,
		nil,
	)
}

// GetRequestWrapper retrieves the `RequestWrapper` from the request context.
//
//   - r: The HTTP request from which to retrieve the `RequestWrapper`.
func GetRequestWrapper(r *http.Request) *util.RequestWrapper {
	return util.GetContextValue[*util.RequestWrapper](
		r.Context(),
		requestDataKey,
		nil,
	)
}

func setResponseWrapper(r *http.Request, rw *util.ResponseWrapper) {
	util.SetContextValue(r.Context(), responseDataKey, rw)
}

func setRequestWrapper(r *http.Request, rw *util.RequestWrapper) {
	util.SetContextValue(r.Context(), requestDataKey, rw)
}
