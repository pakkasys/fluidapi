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
	ErrorID       string  // The ID of the expected error.
	MaskedErrorID *string // An optional ID to mask the original error ID in the response.
	StatusCode    int     // The HTTP status code to return for this error.
	DataIsVisible bool    // Whether to include the error data in the response.
}

// Handle processes an error and returns the corresponding HTTP status code and
// API error. It checks if the error is an *api.Error[any] and handles it
// accordingly.
func (e ErrorHandler) Handle(
	handleError error,
	expectedErrors []ExpectedError,
) (int, *api.Error[any]) {
	apiError, ok := handleError.(*api.Error[any])
	if !ok {
		return http.StatusInternalServerError, InternalServerError
	}
	return e.handleAPIError(apiError, expectedErrors)
}

func (e *ErrorHandler) handleAPIError(
	apiError *api.Error[any],
	expectedErrors []ExpectedError,
) (int, *api.Error[any]) {
	expectedError := e.getExpectedError(apiError, expectedErrors)

	if expectedError == nil {
		return http.StatusInternalServerError, InternalServerError
	}
	return expectedError.maskAPIError(apiError)
}

func (e *ErrorHandler) getExpectedError(
	apiError *api.Error[any],
	expectedErrors []ExpectedError,
) *ExpectedError {
	for i := range expectedErrors {
		if apiError.ID == expectedErrors[i].ErrorID {
			return &expectedErrors[i]
		}
	}
	return nil
}

func (expectedError *ExpectedError) maskAPIError(
	apiError *api.Error[any],
) (int, *api.Error[any]) {
	var useErrorID string
	if expectedError.MaskedErrorID != nil {
		useErrorID = *expectedError.MaskedErrorID
	} else {
		useErrorID = expectedError.ErrorID
	}

	var useData any
	if expectedError.DataIsVisible {
		useData = apiError.Data
	} else {
		useData = nil
	}

	return expectedError.StatusCode, &api.Error[any]{ID: useErrorID, Data: useData}
}
