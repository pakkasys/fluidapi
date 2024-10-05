package inputlogic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockObjectPicker is a mock implementation of the ObjectPicker interface.
type MockObjectPicker[T any] struct {
	mock.Mock
}

func (m *MockObjectPicker[T]) PickObject(
	r *http.Request,
	w http.ResponseWriter,
	obj T,
) (*T, error) {
	args := m.Called(r, w, obj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

// MockOutputHandler is a mock implementation of the OutputHandler interface.
type MockOutputHandler struct {
	mock.Mock
}

func (m *MockOutputHandler) ProcessOutput(
	w http.ResponseWriter,
	r *http.Request,
	out any,
	outError error,
	statusCode int,
) error {
	args := m.Called(w, r, out, outError, statusCode)
	return args.Error(0)
}

// MockLogger is a mock implementation of the ILogger interface.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Trace(messages ...any) {
	m.Called(messages)
}

func (m *MockLogger) Error(messages ...any) {
	m.Called(messages)
}

// MockHelper helps to mock a value receiver object.
type MockHelper struct {
	mock.Mock
}

// MockValidatedInput is a mock implementation of the ValidatedInput interface.
type MockValidatedInput struct {
	helper *MockHelper
}

func (m MockValidatedInput) Validate() []FieldError {
	if m.helper == nil {
		panic("MockHelper is not initialized.")
	}
	args := m.helper.Called()
	return args.Get(0).([]FieldError)
}

// Factory function to create MockValidatedInput with initialized helper
func NewMockValidatedInput(helper *MockHelper) MockValidatedInput {
	return MockValidatedInput{helper: helper}
}

// TestMiddlewareWrapper_Success checks if MiddlewareWrapper correctly wraps
// the callback.
func TestMiddlewareWrapper_Success(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	inputFactory := func() *MockValidatedInput {
		return &input
	}

	opts := Options[MockValidatedInput]{
		ObjectPicker:  mockObjectPicker,
		OutputHandler: mockOutputHandler,
		Logger: func(*http.Request) ILogger {
			return mockLogger
		},
	}

	callback := func(
		w http.ResponseWriter,
		r *http.Request,
		i *MockValidatedInput,
	) (*string, error) {
		result := "Success"
		return &result, nil
	}

	expectedErrors := []ExpectedError{}

	mockHelper.On("Validate").Return([]FieldError{})

	wrapper := MiddlewareWrapper(callback, inputFactory, expectedErrors, opts)

	assert.NotNil(t, wrapper)
	assert.Equal(t, MiddlewareID, wrapper.ID)
	assert.NotNil(t, wrapper.Middleware)
	assert.Equal(t, 1, len(wrapper.Inputs))
	assert.IsType(t, MockValidatedInput{}, wrapper.Inputs[0])
}

// TestMiddleware_Success tests that the middleware handles successful
// callbacks.
func TestMiddleware_Success(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	inputFactory := func() *MockValidatedInput {
		return &input
	}

	opts := Options[MockValidatedInput]{
		Logger: func(*http.Request) ILogger {
			return mockLogger
		},
	}

	expectedOutput := "Success"
	callback := func(
		w http.ResponseWriter,
		r *http.Request,
		i *MockValidatedInput,
	) (*string, error) {
		return &expectedOutput, nil
	}

	mockObjectPicker.On("PickObject", r, w, input).Return(&input, nil)
	mockLogger.On("Trace", mock.Anything)
	mockOutputHandler.
		On("ProcessOutput", w, r, &expectedOutput, nil, http.StatusOK).
		Return(nil)
	mockHelper.On("Validate").Return([]FieldError{})

	middleware := Middleware(
		callback,
		inputFactory,
		nil,
		mockObjectPicker,
		mockOutputHandler,
		opts.Logger,
	)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := middleware(next)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusTeapot, w.Result().StatusCode)

	mockObjectPicker.AssertExpectations(t)
	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestMiddleware_InputValidationError tests that the middleware handles
// input validation errors.
func TestMiddleware_InputValidationError(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	inputFactory := func() *MockValidatedInput {
		return &input
	}

	opts := Options[MockValidatedInput]{
		Logger: func(*http.Request) ILogger {
			return mockLogger
		},
	}

	callback := func(
		w http.ResponseWriter,
		r *http.Request,
		i *MockValidatedInput,
	) (*string, error) {
		return nil, nil
	}

	mockObjectPicker.On("PickObject", r, w, input).Return(&input, nil)
	mockLogger.On("Trace", mock.Anything)
	mockOutputHandler.
		On(
			"ProcessOutput",
			w,
			r,
			mock.Anything,
			mock.Anything,
			http.StatusBadRequest,
		).
		Return(nil)

	input.helper.
		On("Validate").
		Return([]FieldError{{Field: "test", Message: "invalid"}})

	middleware := Middleware(
		callback,
		inputFactory,
		nil,
		mockObjectPicker,
		mockOutputHandler,
		opts.Logger,
	)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(next)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockObjectPicker.AssertExpectations(t)
	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestMiddleware_CallbackError tests that the middleware handles callback
// errors.
func TestMiddleware_CallbackError(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	inputFactory := func() *MockValidatedInput {
		return &input
	}

	opts := Options[MockValidatedInput]{
		Logger: func(*http.Request) ILogger {
			return mockLogger
		},
	}

	expectedError := errors.New("callback failed")
	callback := func(
		w http.ResponseWriter,
		r *http.Request,
		i *MockValidatedInput,
	) (*string, error) {
		return nil, expectedError
	}

	mockObjectPicker.On("PickObject", r, w, input).Return(&input, nil)
	mockLogger.On("Trace", mock.Anything)
	mockOutputHandler.
		On(
			"ProcessOutput",
			w,
			r,
			nil,
			mock.Anything,
			http.StatusInternalServerError,
		).Return(nil)
	mockHelper.On("Validate").Return([]FieldError{})

	middleware := Middleware(
		callback,
		inputFactory,
		nil,
		mockObjectPicker,
		mockOutputHandler,
		opts.Logger,
	)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(next)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockObjectPicker.AssertExpectations(t)
	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestMiddleware_ObjectPickerNil_Panics tests that the middleware panics when
// the objectPicker is nil.
func TestMiddleware_ObjectPickerNil_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Middleware did not panic when objectPicker was nil")
		}
	}()

	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)
	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	inputFactory := func() *MockValidatedInput {
		return &input
	}

	opts := Options[MockValidatedInput]{
		Logger: func(*http.Request) ILogger {
			return mockLogger
		},
	}

	Middleware(
		func(
			w http.ResponseWriter,
			r *http.Request,
			i *MockValidatedInput,
		) (*string, error) {
			return nil, nil
		},
		inputFactory,
		nil,
		nil,
		mockOutputHandler,
		opts.Logger,
	)
}

// TestMiddleware_OutputHandlerNil_Panics tests that the middleware panics when
// the outputHandler is nil.
func TestMiddleware_OutputHandlerNil_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Middleware did not panic when outputHandler was nil")
		}
	}()

	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockLogger := new(MockLogger)
	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	inputFactory := func() *MockValidatedInput {
		return &input
	}

	opts := Options[MockValidatedInput]{
		Logger: func(*http.Request) ILogger {
			return mockLogger
		},
	}

	Middleware(
		func(
			w http.ResponseWriter,
			r *http.Request,
			i *MockValidatedInput,
		) (*string, error) {
			return nil, nil
		},
		inputFactory,
		nil,
		mockObjectPicker,
		nil,
		opts.Logger,
	)
}

// TestHandleError tests the handleError function by providing an expected
// error to the function.
func TestHandleError_ExpectedError(t *testing.T) {
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	expectedError := &api.Error[any]{ID: "EXPECTED_ERROR"}
	expectedStatusCode := http.StatusBadRequest

	expectedErrors := []ExpectedError{
		{
			ID:         "EXPECTED_ERROR",
			Status:     expectedStatusCode,
			PublicData: true,
		},
	}

	mockLogger.On("Trace", mock.Anything)
	mockOutputHandler.
		On("ProcessOutput", w, r, nil, expectedError, expectedStatusCode).
		Return(nil)

	handleError(
		w,
		r,
		expectedError,
		mockOutputHandler,
		expectedErrors,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestHandleError tests the handleError function by providing an unexpected
// error to the function.
func TestHandleError_UnexpectedError(t *testing.T) {
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	unexpectedError := errors.New("unexpected error")
	expectedStatusCode := http.StatusInternalServerError
	expectedApiError := InternalServerError

	mockLogger.On("Trace", mock.Anything)
	mockOutputHandler.
		On("ProcessOutput", w, r, nil, expectedApiError, expectedStatusCode).
		Return(nil)

	handleError(
		w,
		r,
		unexpectedError,
		mockOutputHandler,
		nil,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestHandleInput_Success tests that handleInput successfully picks and
// validates an input.
func TestHandleInput_Success(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Set up mock expectations
	mockObjectPicker.On("PickObject", r, w, input).Return(&input, nil)
	mockLogger.On("Trace", mock.Anything)
	mockHelper.On("Validate").Return([]FieldError{})

	// Call handleInput function
	returnedInput, err := handleInput(
		w,
		r,
		input,
		mockObjectPicker,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, &input, returnedInput)

	mockObjectPicker.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockHelper.AssertExpectations(t)
}

// TestHandleInput_ObjectPickerFailure tests that handleInput handles object
// picker failure.
func TestHandleInput_ObjectPickerFailure(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Set up mock expectations
	expectedError := errors.New("failed to pick object")
	mockObjectPicker.On("PickObject", r, w, input).Return(nil, expectedError)

	// Call handleInput function
	returnedInput, err := handleInput(
		w,
		r,
		input,
		mockObjectPicker,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	// Assertions
	assert.Nil(t, returnedInput)
	assert.EqualError(t, err, "failed to pick object")

	mockObjectPicker.AssertExpectations(t)
}

// TestHandleInput_ValidationError tests that handleInput handles input
// validation errors.
func TestHandleInput_ValidationError(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockLogger := new(MockLogger)

	mockHelper := new(MockHelper)
	input := NewMockValidatedInput(mockHelper)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mockObjectPicker.On("PickObject", r, w, input).Return(&input, nil)
	mockLogger.On("Trace", mock.Anything)
	validationErrors := []FieldError{
		{Field: "testField", Message: "invalid value"},
	}
	mockHelper.On("Validate").Return(validationErrors)

	returnedInput, err := handleInput(
		w,
		r,
		input,
		mockObjectPicker,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	assert.Nil(t, returnedInput)
	assert.IsType(t, ValidationError, err)

	mockObjectPicker.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockHelper.AssertExpectations(t)
}

// TestHandleOutput_Success tests that handleOutput successfully handles output.
func TestHandleOutput_Success(t *testing.T) {
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	output := "ExpectedOutput"
	statusCode := http.StatusOK

	mockOutputHandler.On("ProcessOutput", w, r, output, nil, statusCode).
		Return(nil)

	handleOutput(
		w,
		r,
		output,
		nil,
		statusCode,
		mockOutputHandler,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertNotCalled(t, "Error")
}

// TestHandleOutput_Failure tests that handleOutput handles output processing
// errors.
func TestHandleOutput_Failure(t *testing.T) {
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	output := "ExpectedOutput"
	statusCode := http.StatusOK
	expectedError := errors.New("output processing failed")

	mockOutputHandler.On("ProcessOutput", w, r, output, nil, statusCode).
		Return(expectedError)
	mockLogger.On("Error", mock.Anything)

	handleOutput(
		w,
		r,
		output,
		nil,
		statusCode,
		mockOutputHandler,
		func(*http.Request) ILogger {
			return mockLogger
		},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	mockOutputHandler.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
