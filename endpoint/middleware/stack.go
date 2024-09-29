package middleware

import "github.com/pakkasys/fluidapi/core/api"

type Stack []api.MiddlewareWrapper

// InsertAfterID inserts a middleware wrapper after the given ID.
//
//   - id: The ID of the middleware to insert after.
//   - wrapper: The middleware wrapper to insert.
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
