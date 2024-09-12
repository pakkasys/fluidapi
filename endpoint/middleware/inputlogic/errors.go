package inputlogic

import "github.com/pakkasys/fluidapi/core/api"

type ErrorData struct {
	Errors []string `json:"errors"`
}

var VALIDATION_ERROR_ID = "VALIDATION_ERROR_ID"

func ValidationError(errorMessages []string) *api.Error {
	return &api.Error{
		ID: VALIDATION_ERROR_ID,
		Data: ErrorData{
			Errors: errorMessages,
		},
	}
}
