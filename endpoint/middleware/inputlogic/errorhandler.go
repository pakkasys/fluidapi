package inputlogic

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
)

var InternalServerError = api.NewError[any]("INTERNAL_SERVER_ERROR")

// ErrorHandler handles errors and maps them to appropriate HTTP responses.
type ErrorHandler struct{}

// ExpectedError represents an expected error configuration.
// It defines how to handle specific errors that are anticipated.
type ExpectedError struct {
	ID         string  // The ID of the expected error.
	MaskedID   *string // An optional ID to mask the original error ID in the response.
	Status     int     // The HTTP status code to return for this error.
	PublicData bool    // Whether to include the error data in the response.
}

// Handle processes an error and returns the corresponding HTTP status code and
// API error. It checks if the error is an *api.Error[any] and handles it
// accordingly.
func (e ErrorHandler) Handle(
	handleError error,
	expectedErrors []ExpectedError,
) (int, *api.Error[any]) {
	apiError, ok := handleError.(api.APIError)
	if !ok {
		return http.StatusInternalServerError, InternalServerError
	}
	return e.handleAPIError(apiError, expectedErrors)
}

func (e *ErrorHandler) handleAPIError(
	apiError api.APIError,
	expectedErrors []ExpectedError,
) (int, *api.Error[any]) {
	expectedError := e.getExpectedError(apiError, expectedErrors)
	if expectedError == nil {
		return http.StatusInternalServerError, InternalServerError
	}
	return expectedError.maskAPIError(apiError)
}

func (e *ErrorHandler) getExpectedError(
	apiError api.APIError,
	expectedErrors []ExpectedError,
) *ExpectedError {
	for i := range expectedErrors {
		if apiError.GetID() == expectedErrors[i].ID {
			return &expectedErrors[i]
		}
	}
	return nil
}

func (expectedError *ExpectedError) maskAPIError(
	apiError api.APIError,
) (int, *api.Error[any]) {
	var useErrorID string
	if expectedError.MaskedID != nil {
		useErrorID = *expectedError.MaskedID
	} else {
		useErrorID = expectedError.ID
	}

	var useData any
	if expectedError.PublicData {
		useData = apiError.GetData()
	} else {
		useData = nil
	}

	return expectedError.Status, &api.Error[any]{ID: useErrorID, Data: &useData}
}
