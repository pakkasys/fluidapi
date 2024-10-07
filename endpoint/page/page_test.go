package page

import (
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/stretchr/testify/assert"
)

// Validate_ValidLimit tests the Validate function for a valid limit.
func TestValidate_ValidLimit(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  5,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.NoError(t, err, "Expected no error when limit is within maxLimit")
}

// Validate_LimitExceeded tests the Validate function for a limit that exceeds
// the max limit.
func TestValidate_LimitExceeded(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  15,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.Error(t, err, "Expected an error when limit exceeds maxLimit")
	apiErr, ok := err.(*api.Error[MaxPageLimitExceededErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "MAX_PAGE_LIMIT_EXCEEDED", apiErr.ID, "Error ID should")
	assert.Equal(t, maxLimit, apiErr.Data.MaxLimit, "Max limit should match")
}

// Validate_ZeroLimit tests the Validate function for a limit of zero.
func TestValidate_ZeroLimit(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  0,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.NoError(t, err, "Expected no error when limit is zero")
}
