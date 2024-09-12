package api

type Error struct {
	ID   string `json:"id"`
	Data any    `json:"data,omitempty"`
}

func NewError(id string, data any) *Error {
	return &Error{
		ID:   id,
		Data: data,
	}
}

func (e *Error) Error() string {
	return e.ID
}
