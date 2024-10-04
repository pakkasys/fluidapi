package inputlogic

import (
	"github.com/pakkasys/fluidapi/core/api"
)

const VALIDATION_ERROR_ID = "VALIDATION_ERROR"

// FieldError represents a field-level validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorData contains a list of field-level validation errors
type ValidationErrorData struct {
	Errors []FieldError `json:"errors"`
}

// ValidationError function that accepts a list of field-level validation errors
func ValidationError(fieldErrors []FieldError) *api.Error {
	return &api.Error{
		ID: VALIDATION_ERROR_ID,
		Data: ValidationErrorData{
			Errors: fieldErrors,
		},
	}
}
