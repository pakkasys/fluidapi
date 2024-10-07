package api

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// ApplyMiddlewares applies a chain of middlewares to an http.Handler.
//   - h: The http.Handler to wrap with middlewares.
//   - middlewares: A variadic parameter of Middleware functions to apply.
func ApplyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// MiddlewareWrapper wraps a middleware function with additional metadata.
type MiddlewareWrapper struct {
	ID         string
	Middleware Middleware
	Inputs     []any
}
