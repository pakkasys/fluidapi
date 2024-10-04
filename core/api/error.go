package api

// Error represents a JSON marshalable custom error type with an ID and optional
// data.
type Error[T any] struct {
	ID      string `json:"id"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewError returns a new error with the given ID.
func NewError[T any](id string) *Error[T] {
	return &Error[T]{
		ID: id,
	}
}

// WithData returns a new error with the given data.
func (e *Error[T]) WithData(data T) *Error[T] {
	return &Error[T]{
		ID:   e.ID,
		Data: data,
	}
}

func (e *Error[T]) WithMessage(message string) *Error[T] {
	return &Error[T]{
		ID:      e.ID,
		Data:    e.Data,
		Message: message,
	}
}

// Error returns the error message as a string, which is the ID of the error.
func (e *Error[T]) Error() string {
	if e.Message != "" {
		return e.ID + ": " + e.Message
	}
	return e.ID
}
