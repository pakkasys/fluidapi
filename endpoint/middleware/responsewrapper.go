package middleware

import (
	"net/http"

	"github.com/PakkaSys/fluidapi/core/api"
	"github.com/PakkaSys/fluidapi/endpoint/util"
)

const ResponseWrapperMiddlewareID = "response_wrapper"

var (
	responseDataKey = util.NewDataKey()
	requestDataKey  = util.NewDataKey()
)

func ResponseWrapperMiddlewareWrapper() *api.MiddlewareWrapper {
	return api.NewMiddlewareWrapperBuilder().
		ID(ResponseWrapperMiddlewareID).
		Middleware(ResponseWrapperMiddleware()).
		Build()
}

func ResponseWrapperMiddleware() api.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseWrapper := util.NewResponseWrapper(w)

			requestWrapper, err := util.NewRequestWrapper(r)
			if err != nil {
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}

			util.SetContextValue(
				r.Context(),
				responseDataKey,
				responseWrapper,
			)
			util.SetContextValue(
				r.Context(),
				requestDataKey,
				requestWrapper,
			)

			next.ServeHTTP(responseWrapper, requestWrapper.Request)
		})
	}
}

func GetResponseWrapper(r *http.Request) *util.ResponseWrapper {
	return util.GetContextValue[*util.ResponseWrapper](
		r.Context(),
		responseDataKey,
		nil,
	)
}

func GetRequestWrapper(r *http.Request) *util.RequestWrapper {
	return util.GetContextValue[*util.RequestWrapper](
		r.Context(),
		requestDataKey,
		nil,
	)
}
