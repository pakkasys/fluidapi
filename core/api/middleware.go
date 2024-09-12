package api

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func ApplyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

type MiddlewareInput struct {
	Input any
}

func NewMiddlewareInput(input any) MiddlewareInput {
	return MiddlewareInput{
		Input: input,
	}
}

type MiddlewareWrapper struct {
	ID               string
	Middleware       Middleware
	MiddlewareInputs []MiddlewareInput
}

type MiddlewareWrapperBuilder struct {
	middlewareWrapper MiddlewareWrapper
}

func NewMiddlewareWrapperBuilder() *MiddlewareWrapperBuilder {
	return &MiddlewareWrapperBuilder{
		middlewareWrapper: MiddlewareWrapper{},
	}
}

func (b *MiddlewareWrapperBuilder) Build() *MiddlewareWrapper {
	return &b.middlewareWrapper
}

func (b *MiddlewareWrapperBuilder) ID(id string) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.ID = id
	return b
}

func (b *MiddlewareWrapperBuilder) Middleware(
	middleware Middleware,
) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.Middleware = middleware
	return b
}

func (b *MiddlewareWrapperBuilder) MiddlewareInputs(
	middlewareInputs []MiddlewareInput,
) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.MiddlewareInputs = middlewareInputs
	return b
}

func (b *MiddlewareWrapperBuilder) AddMiddlewareInput(
	middlewareInput MiddlewareInput,
) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.MiddlewareInputs = append(
		b.middlewareWrapper.MiddlewareInputs,
		middlewareInput,
	)
	return b
}
