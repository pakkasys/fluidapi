package inputlogic

import (
	"net/http"
	"testing"

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
	return args.Get(0).(*T), args.Error(1)
}

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

// MockValidatedInput is a mock implementation of the ValidatedInput interface.
type MockValidatedInput struct{}

func (m MockValidatedInput) Validate() []FieldError {
	return []FieldError{}
}

// TestMiddlewareWrapper_Success checks if MiddlewareWrapper correctly wraps
// the callback.
func TestMiddlewareWrapper_Success(t *testing.T) {
	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockOutputHandler := new(MockOutputHandler)
	mockLogger := new(MockLogger)

	inputFactory := func() *MockValidatedInput {
		return &MockValidatedInput{}
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

	wrapper := MiddlewareWrapper(callback, inputFactory, expectedErrors, opts)

	assert.NotNil(t, wrapper)
	assert.Equal(t, MiddlewareID, wrapper.ID)
	assert.NotNil(t, wrapper.Middleware)
	assert.Equal(t, 1, len(wrapper.Inputs))
	assert.IsType(t, MockValidatedInput{}, wrapper.Inputs[0])
}
