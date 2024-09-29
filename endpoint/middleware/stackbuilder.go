package middleware

import (
	"github.com/pakkasys/fluidapi/core/api"
)

type StackBuilder interface {
	Build() Stack
	MustAddMiddleware(wrapper ...api.MiddlewareWrapper) StackBuilder
}
