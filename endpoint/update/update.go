package update

// InputUpdate represents an update DTO.
type InputUpdate struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

// APIUpdate represents an update object in API.
type APIUpdate struct {
	Validation string
}
