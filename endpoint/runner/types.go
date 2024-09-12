package runner

import (
	"github.com/PakkaSys/fluidapi/endpoint/selector"
	"github.com/PakkaSys/fluidapi/endpoint/update"
)

type ISpecification interface {
	URL() string
	HTTPMethod() string
}

type InputFactory[T any] func() *T

type IInputSpecification[Input any] interface {
	ISpecification
	InputFactory() InputFactory[Input]
}

type IGetSpecification[Input any] interface {
	IInputSpecification[Input]
	MaxPageCount() int
	AllowedSelectors() map[string]selector.APISelector
	AllowedOrderFields() []string
}

type IUpdateSpecification[Input any] interface {
	IInputSpecification[Input]
	AllowedSelectors() map[string]selector.APISelector
	AllowedUpdates() map[string]update.APIUpdate
}

type IDeleteSpecification[Input any] interface {
	IInputSpecification[Input]
	AllowedSelectors() map[string]selector.APISelector
}

type IDeleteSpecificationOrderable[Input any] interface {
	IInputSpecification[Input]
	AllowedOrderFields() []string
}

type ILimitable interface {
	GetLimit() int
}

type InputSpecification[Input any] struct {
	url          string
	httpMethod   string
	inputFactory InputFactory[Input]
}

func NewInputSpecification[Input any](
	url string,
	httpMethod string,
	inputFactory InputFactory[Input],
) IInputSpecification[Input] {
	return &InputSpecification[Input]{
		url:          url,
		httpMethod:   httpMethod,
		inputFactory: inputFactory,
	}
}

func (s *InputSpecification[Input]) URL() string {
	return s.url
}

func (s *InputSpecification[Input]) HTTPMethod() string {
	return s.httpMethod
}

func (s *InputSpecification[Input]) InputFactory() InputFactory[Input] {
	return s.inputFactory
}

type GetSpecification[Input any] struct {
	IInputSpecification[Input]
	maxPageCount       int
	allowedSelectors   map[string]selector.APISelector
	allowedOrderFields []string
}

func NewGetSpecification[Input any](
	url string,
	httpMethod string,
	inputFactory InputFactory[Input],
	maxPageCount int,
	allowedSelectors map[string]selector.APISelector,
	allowedOrderFields []string,
) IGetSpecification[Input] {
	return &GetSpecification[Input]{
		IInputSpecification: NewInputSpecification(
			url,
			httpMethod,
			inputFactory,
		),
		maxPageCount:       maxPageCount,
		allowedSelectors:   allowedSelectors,
		allowedOrderFields: allowedOrderFields,
	}
}

func (s *GetSpecification[Input]) MaxPageCount() int {
	return s.maxPageCount
}

func (s *GetSpecification[Input]) AllowedSelectors() map[string]selector.APISelector {
	return s.allowedSelectors
}

func (s *GetSpecification[Input]) AllowedOrderFields() []string {
	return s.allowedOrderFields
}

type UpdateSpecification[Input any] struct {
	IInputSpecification[Input]
	allowedSelectors map[string]selector.APISelector
	allowedUpdates   map[string]update.APIUpdate
}

func NewUpdateSpecification[Input any](
	url string,
	httpMethod string,
	inputFactory InputFactory[Input],
	allowedSelectors map[string]selector.APISelector,
	allowedUpdates map[string]update.APIUpdate,
) IUpdateSpecification[Input] {
	return &UpdateSpecification[Input]{
		IInputSpecification: NewInputSpecification(
			url,
			httpMethod,
			inputFactory,
		),
		allowedSelectors: allowedSelectors,
		allowedUpdates:   allowedUpdates,
	}
}

func (s *UpdateSpecification[Input]) AllowedSelectors() map[string]selector.APISelector {
	return s.allowedSelectors
}

func (s *UpdateSpecification[Input]) AllowedUpdates() map[string]update.APIUpdate {
	return s.allowedUpdates
}

type DeleteSpecification[Input any] struct {
	IInputSpecification[Input]
	allowedSelectors map[string]selector.APISelector
}

func NewDeleteSpecification[Input any](
	url string,
	httpMethod string,
	inputFactory InputFactory[Input],
	allowedSelectors map[string]selector.APISelector,
) IDeleteSpecification[Input] {
	return &DeleteSpecification[Input]{
		IInputSpecification: NewInputSpecification(
			url,
			httpMethod,
			inputFactory,
		),
		allowedSelectors: allowedSelectors,
	}
}

func (s *DeleteSpecification[Input]) AllowedSelectors() map[string]selector.APISelector {
	return s.allowedSelectors
}
