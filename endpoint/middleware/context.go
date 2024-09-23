package middleware

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/util"
)

// ContextMiddlewareID is the ID of the context middleware
const ContextMiddlewareID = "context"

// ContextMiddlewareWrapper is the middleware wrapper for the context middleware
func ContextMiddlewareWrapper() *api.MiddlewareWrapper {
	return api.NewMiddlewareWrapperBuilder().
		ID(ContextMiddlewareID).
		Middleware(ContextMiddleware()).
		Build()
}

// ContextMiddleware is the middleware for the context middleware
func ContextMiddleware() api.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(util.NewContext(r.Context())))
		})
	}
}
