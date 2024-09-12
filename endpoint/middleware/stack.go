package middleware

import "github.com/PakkaSys/fluidapi/core/api"

type MiddlewareStack []api.MiddlewareWrapper

func (ms *MiddlewareStack) InsertAfterID(
	id string,
	newMiddleware api.MiddlewareWrapper,
) bool {
	for i, mw := range *ms {
		if mw.ID == id {
			if i == len(*ms)-1 {
				*ms = append(*ms, newMiddleware)
			} else {
				*ms = append(
					(*ms)[:i+1],
					append(
						[]api.MiddlewareWrapper{newMiddleware},
						(*ms)[i+1:]...,
					)...,
				)
			}
			return true
		}
	}
	return false
}
