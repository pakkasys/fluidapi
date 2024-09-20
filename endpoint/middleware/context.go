package middleware

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/util"
)

const ContextMiddlewareID = "context"

func ContextMiddlewareWrapper() *api.MiddlewareWrapper {
	return api.NewMiddlewareWrapperBuilder().
		ID(ContextMiddlewareID).
		Middleware(ContextMiddleware()).
		Build()
}

func ContextMiddleware() api.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(util.NewContext(r.Context())))
		})
	}
}
