package order

import "github.com/pakkasys/fluidapi/core/api"

type InvalidOrderFieldErrorData struct {
	Field string `json:"field"`
}

var INVALID_ORDER_FIELD_ERROR_ID = "INVALID_ORDER_FIELD"

func InvalidOrderFieldError(field string) *api.Error {
	return &api.Error{
		ID: INVALID_ORDER_FIELD_ERROR_ID,
		Data: InvalidOrderFieldErrorData{
			Field: field,
		},
	}
}
