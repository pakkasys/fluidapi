package middleware

import (
	"fmt"
	"net/http"

	"github.com/PakkaSys/fluidapi/core/api"
)

type MiddlewareStackBuilder struct {
	middlewareStack MiddlewareStack
}

func NewMiddlewareStackBuilder(
	requestIDFn func() string,
	panicHandlerLoggerFn func(r *http.Request) func(messages ...any),
	requestLoggerFn func(r *http.Request) func(messages ...any),
) *MiddlewareStackBuilder {
	return &MiddlewareStackBuilder{
		[]api.MiddlewareWrapper{
			*ContextMiddlewareWrapper(),
			*ResponseWrapperMiddlewareWrapper(),
			*RequestIDMiddlewareWrapper(requestIDFn),
			*PanicHandlerMiddlewareWrapper(panicHandlerLoggerFn),
			*RequestLogMiddlewareWrapper(requestLoggerFn),
		},
	}
}

func (b *MiddlewareStackBuilder) Build() MiddlewareStack {
	return b.middlewareStack
}

func (b *MiddlewareStackBuilder) MustAddMiddleware(
	middlewareWrapper api.MiddlewareWrapper,
) *MiddlewareStackBuilder {
	success := b.middlewareStack.InsertAfterID(
		RequestLogMiddlewareID,
		middlewareWrapper,
	)
	if !success {
		panic(fmt.Sprintf("Failed to add middleware: %s", middlewareWrapper.ID))
	}

	return b
}
