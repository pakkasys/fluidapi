package middleware

import "github.com/pakkasys/fluidapi/core/api"

type Stack []api.MiddlewareWrapper

// Middlewares returns the middlewares in the stack.
//
// Parameters:
//   - s: The middleware stack.
//
// Returns:
//   - The middlewares in the stack.
func (s Stack) Middlewares() []api.Middleware {
	middlewares := []api.Middleware{}
	for _, mw := range s {
		middlewares = append(middlewares, mw.Middleware)
	}
	return middlewares
}

// InsertAfterID inserts a middleware wrapper after the given ID.
//
// Parameters:
//   - id: The ID of the middleware to insert after.
//   - wrapper: The middleware wrapper to insert.
//
// Returns:
//   - True if the middleware was inserted, false otherwise.
func (s *Stack) InsertAfterID(id string, wrapper api.MiddlewareWrapper) bool {
	for i, mw := range *s {
		if mw.ID == id {
			if i == len(*s)-1 {
				*s = append(*s, wrapper)
			} else {
				*s = append(
					(*s)[:i+1],
					append(
						[]api.MiddlewareWrapper{wrapper},
						(*s)[i+1:]...,
					)...,
				)
			}
			return true
		}
	}
	return false
}
