package runner

type InputFactory[T any] func() *T

type InputSpecification[Input any] struct {
	URL          string
	Method       string
	InputFactory InputFactory[Input]
}
