package middleware

import "github.com/pakkasys/fluidapi/core/api"

type Stack []api.MiddlewareWrapper

// InsertAfterID inserts a middleware wrapper after the given ID.
//
//   - id: The ID of the middleware to insert after.
//   - wrapper: The middleware wrapper to insert.
func (ms *Stack) InsertAfterID(
	id string,
	wrapper api.MiddlewareWrapper,
) bool {
	for i, mw := range *ms {
		if mw.ID == id {
			if i == len(*ms)-1 {
				*ms = append(*ms, wrapper)
			} else {
				*ms = append(
					(*ms)[:i+1],
					append(
						[]api.MiddlewareWrapper{wrapper},
						(*ms)[i+1:]...,
					)...,
				)
			}
			return true
		}
	}
	return false
}
