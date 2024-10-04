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
	err := Error[string]{ID: "test_error", Data: "", Message: ""}

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

// TestEnhancedErrorMethod tests the enhanced Error method with a message.
func TestEnhancedErrorMethod(t *testing.T) {
	err := Error[string]{
		ID:      "test_error",
		Data:    "",
		Message: "Something went wrong",
	}

	assert.Equal(
		t,
		"test_error: Something went wrong",
		err.Error(),
		"Error string should include ID and message",
	)
}
