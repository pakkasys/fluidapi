package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewError tests the NewError function.
func TestNewError(t *testing.T) {
	err := NewError[string]("test_error")

	assert.Equal(t, "test_error", err.ID, "Error ID should be 'test_error'")
	assert.Empty(t, err.Message, "Error Message should be empty")
	assert.Equal(t, "", err.Data, "Error Data should be empty string")
}

// TestError tests the Error method of the generic Error type.
func TestError(t *testing.T) {
	data := ""
	message := ""
	err := Error[string]{ID: "test_error", Data: &data, Message: &message}

	assert.Equal(t, "test_error", err.Error(), "Expected 'test_error'")
}

// TestWithData tests the WithData method of the generic Error type.
func TestWithData(t *testing.T) {
	originalErr := NewError[string]("test_error")
	data := "additional_info"
	newErr := originalErr.WithData(data)

	// Check that the new error has the same ID
	assert.Equal(t, "test_error", newErr.ID, "Expected 'test_error'")

	// Check that the new error has the correct Data
	assert.Equal(t, data, newErr.Data, "New Data should match provided data")

	// Ensure the original error's Data is still empty
	assert.Equal(t, "", originalErr.Data, "Original Data should be empty")
}

// TestWithMessage tests the WithMessage method of the generic Error type.
func TestWithMessage(t *testing.T) {
	originalErr := NewError[string]("test_error")
	msg := "Something went wrong"
	newErr := originalErr.WithMessage(msg)

	// Check that the new error has the same ID
	assert.Equal(t, "test_error", newErr.ID, "New ID should be 'test_error'")

	// Check that the new error has the correct Message
	assert.Equal(t, msg, newErr.Message, "New Message should match provided")

	// Ensure the original error's Message is still empty
	assert.Empty(t, originalErr.Message, "Original Message should be empty")
}

// TestErrorFunc tests the Error method.
func TestErrorFunc(t *testing.T) {
	data := ""
	message := "Something went wrong"
	err := Error[string]{
		ID:      "test_error",
		Data:    &data,
		Message: &message,
	}

	assert.Equal(
		t,
		"test_error: Something went wrong",
		err.Error(),
		"Error string should include ID and message",
	)
}

// TestGetID tests the GetID method of the generic Error type.
func TestGetID(t *testing.T) {
	err := NewError[string]("test_error")

	assert.Equal(t, "test_error", err.GetID(), "Should return 'test_error'")
}

// TestGetData tests the GetData method of the generic Error type.
func TestGetData(t *testing.T) {
	err := NewError[string]("test_error")
	data := "some_data"
	errWithData := err.WithData(data)

	assert.Equal(t, data, errWithData.GetData(), "Should return correct data")
}

// TestGetMessage tests the GetMessage method of the generic Error type.
func TestGetMessage(t *testing.T) {
	err := NewError[string]("test_error")
	msg := "Detailed message"
	errWithMessage := err.WithMessage(msg)

	assert.Equal(t, msg, errWithMessage.GetMessage(), "No correct message")
}

// TestAPIErrorInterface tests if Error satisfies the APIError interface.
func TestAPIErrorInterface(t *testing.T) {
	var apiErr APIError = NewError[string]("test_error")
	assert.Equal(t, "test_error", apiErr.GetID(), "Should return 'test_error'")
}
