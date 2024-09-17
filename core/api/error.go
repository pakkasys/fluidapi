package api

// Error represents a JSON marshalable custom error type with an ID and optional
// data.
type Error struct {
	ID   string `json:"id"`
	Data any    `json:"data,omitempty"`
}

// Error returns the error message as a string, which is the ID of the error.
func (e *Error) Error() string {
	return e.ID
}
