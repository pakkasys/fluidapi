package api

// Error represents a JSON marshalable custom error type with an ID and optional
// data.
type Error struct {
	ID   string `json:"id"`
	Data any    `json:"data,omitempty"`
}

// NewError creates a new Error instance.
// - id: A unique identifier for the error.
// - data: Additional data related to the error (optional).
func NewError(id string, data any) *Error {
	return &Error{
		ID:   id,
		Data: data,
	}
}

// Error returns the error message as a string, which is the ID of the error.
func (e *Error) Error() string {
	return e.ID
}
