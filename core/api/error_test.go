package api

import (
	"testing"
)

// TestNewError tests the creation of a new Error instance.
func TestNewError(t *testing.T) {
	err := NewError("test_error", nil)

	if err.ID != "test_error" {
		t.Errorf("expected error ID to be test_error, got %s", err.ID)
	}

	if err.Data != nil {
		t.Errorf("expected error data to be nil, got %v", err.Data)
	}
}

// TestError tests the Error method of the Error type.
func TestError(t *testing.T) {
	err := NewError("test_error", nil)

	if err.Error() != "test_error" {
		t.Errorf("expected error string to be test_error, got %s", err.Error())
	}
}
