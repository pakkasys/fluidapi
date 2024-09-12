package inputlogic

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core/api"
)

var INTERNAL_SERVER_ERROR_ID = "INTERNAL_SERVER_ERROR"

func InternalServerError() *api.Error {
	return &api.Error{
		ID: INTERNAL_SERVER_ERROR_ID,
	}
}

type ErrorHandler struct{}

type ExpectedError struct {
	ErrorID       string
	MaskedErrorID string
	StatusCode    int
	DataIsVisible bool
}

func NewExpectedError(
	errorID string,
	statusCode int,
) *ExpectedError {
	return &ExpectedError{
		ErrorID:    errorID,
		StatusCode: statusCode,
	}
}

func (e *ExpectedError) WithMaskedErrorID(maskedErrorID string) *ExpectedError {
	e.MaskedErrorID = maskedErrorID
	return e
}

func (e *ExpectedError) WithDataIsVisible(dataIsVisible bool) *ExpectedError {
	e.DataIsVisible = dataIsVisible
	return e
}

func (e *ErrorHandler) Handle(
	handleError error,
	expectedErrors []ExpectedError,
) (int, *api.Error) {
	apiError, ok := handleError.(*api.Error)
	if !ok {
		return http.StatusInternalServerError, InternalServerError()
	} else {
		return e.handleAPIError(apiError, expectedErrors)
	}
}

func (e *ErrorHandler) handleAPIError(
	apiError *api.Error,
	expectedErrors []ExpectedError,
) (int, *api.Error) {
	expectedError := e.getExpectedError(*apiError, expectedErrors)

	if expectedError == nil {
		return http.StatusInternalServerError, nil
	} else {
		var useErrorID string
		if expectedError.MaskedErrorID != "" {
			useErrorID = expectedError.MaskedErrorID
		} else {
			useErrorID = expectedError.ErrorID
		}

		var useData any
		if expectedError.DataIsVisible {
			useData = apiError.Data
		} else {
			useData = nil
		}

		return expectedError.StatusCode, api.NewError(useErrorID, useData)
	}
}

func (e *ErrorHandler) getExpectedError(
	apiError api.Error,
	expectedErrors []ExpectedError,
) *ExpectedError {
	for _, expectedError := range expectedErrors {
		if apiError.ID == expectedError.ErrorID {
			return &expectedError
		}
	}

	return nil
}
