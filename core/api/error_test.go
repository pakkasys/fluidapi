package api

import (
	"testing"
)

// TestError tests the Error method of the Error type.
func TestError(t *testing.T) {
	err := Error{ID: "test_error", Data: nil}

	if err.Error() != "test_error" {
		t.Errorf("expected error string to be test_error, got %s", err.Error())
	}
}
