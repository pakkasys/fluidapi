package update

type InputUpdate struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

func NewInputUpdate(field string, value any) *InputUpdate {
	return &InputUpdate{
		Field: field,
		Value: value,
	}
}

type APIUpdate struct {
	Validation string
}
