package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDuplicateEntry tests the DuplicateEntry function.
func TestDuplicateEntry(t *testing.T) {
	// Simulate an error
	simulatedError := errors.New("duplicate entry error")

	apiError := DuplicateEntryError.WithData(simulatedError)

	assert.Equal(t, DuplicateEntryError.ID, apiError.ID)
	assert.Equal(t, simulatedError, apiError.Data)
}

// TestForeignConstraintError tests the ForeignConstraintError function.
func TestForeignConstraintError(t *testing.T) {
	// Simulate an error
	simulatedError := errors.New("foreign constraint error")

	apiError := ForeignConstraintError.WithData(simulatedError)

	assert.Equal(t, ForeignConstraintError.ID, apiError.ID)
	assert.Equal(t, simulatedError, apiError.Data)
}
