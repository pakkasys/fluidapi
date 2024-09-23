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

			SetRequestWrapper(r, requestWrapper)
			SetResponseWrapper(r, responseWrapper)

			next.ServeHTTP(responseWrapper, requestWrapper.Request)
		})
	}
}

func SetResponseWrapper(r *http.Request, rw *util.ResponseWrapper) {
	util.SetContextValue(r.Context(), responseDataKey, rw)
}

func GetResponseWrapper(r *http.Request) *util.ResponseWrapper {
	return util.GetContextValue[*util.ResponseWrapper](
		r.Context(),
		responseDataKey,
		nil,
	)
}

func SetRequestWrapper(r *http.Request, rw *util.RequestWrapper) {
	util.SetContextValue(r.Context(), requestDataKey, rw)
}

func GetRequestWrapper(r *http.Request) *util.RequestWrapper {
	return util.GetContextValue[*util.RequestWrapper](
		r.Context(),
		requestDataKey,
		nil,
	)
}
