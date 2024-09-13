package api

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// ApplyMiddlewares applies a chain of middlewares to an http.Handler.
// - h: The http.Handler to wrap with middlewares.
// - middlewares: A variadic parameter of Middleware functions to apply.
func ApplyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// MiddlewareInput represents the input used by a middleware.
type MiddlewareInput struct {
	Input any
}

// NewMiddlewareInput creates a new MiddlewareInput with the provided input.
// - input: The input to be used by a middleware.
func NewMiddlewareInput(input any) MiddlewareInput {
	return MiddlewareInput{
		Input: input,
	}
}

// MiddlewareWrapper wraps a middleware function with additional metadata.
type MiddlewareWrapper struct {
	ID               string
	Middleware       Middleware
	MiddlewareInputs []MiddlewareInput
}

// MiddlewareWrapper wraps a middleware function with additional metadata.
type MiddlewareWrapperBuilder struct {
	middlewareWrapper MiddlewareWrapper
}

// MiddlewareWrapperBuilder provides a builder pattern for constructing
// MiddlewareWrapper instances.
func NewMiddlewareWrapperBuilder() *MiddlewareWrapperBuilder {
	return &MiddlewareWrapperBuilder{
		middlewareWrapper: MiddlewareWrapper{},
	}
}

// Build finalizes the builder and returns the constructed MiddlewareWrapper.
func (b *MiddlewareWrapperBuilder) Build() *MiddlewareWrapper {
	return &b.middlewareWrapper
}

// ID sets the ID for the MiddlewareWrapper being built.
// - id: The identifier for the middleware.
func (b *MiddlewareWrapperBuilder) ID(id string) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.ID = id
	return b
}

// Middleware sets the Middleware function for the MiddlewareWrapper being
// built.
// - middleware: The Middleware function to be set.
func (b *MiddlewareWrapperBuilder) Middleware(
	middleware Middleware,
) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.Middleware = middleware
	return b
}

// MiddlewareInputs sets the slice of MiddlewareInput for the MiddlewareWrapper
// being built.
// - middlewareInputs: A slice of MiddlewareInput to be set.
func (b *MiddlewareWrapperBuilder) MiddlewareInputs(
	middlewareInputs []MiddlewareInput,
) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.MiddlewareInputs = middlewareInputs
	return b
}

// AddMiddlewareInput appends a MiddlewareInput to the MiddlewareWrapper being
// built.
// - middlewareInput: The MiddlewareInput to add.
func (b *MiddlewareWrapperBuilder) AddMiddlewareInput(
	middlewareInput MiddlewareInput,
) *MiddlewareWrapperBuilder {
	b.middlewareWrapper.MiddlewareInputs = append(
		b.middlewareWrapper.MiddlewareInputs,
		middlewareInput,
	)
	return b
}
