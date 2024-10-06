package order

import (
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/stretchr/testify/assert"
)

// TestValidateAndDeduplicateOrders tests the ValidateAndDeduplicateOrders
// function with valid input.
func TestValidateAndDeduplicateOrders_ValidInput(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "age", Direction: DIRECTION_DESC},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, orders, result, "Orders should match the input")
}

// TestValidateAndDeduplicateOrders_DuplicateFields tests the case where
// duplicate fields are provided.
func TestValidateAndDeduplicateOrders_DuplicateFields(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "name", Direction: DIRECTION_DESC},
		{Field: "age", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.NoError(t, err)
	assert.Len(t, result, 2, "Duplicate fields should be removed")
	assert.Equal(t, result[0], orders[0], "First order should be preserved")
	assert.Equal(t, result[1], orders[2], "Order should be preserved")
}

// TestValidateAndDeduplicateOrders_InvalidDirection tests the case where an
// invalid direction is provided.
func TestValidateAndDeduplicateOrders_InvalidDirection(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: "INVALID_DIRECTION"},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Error(t, err)
	assert.Nil(t, result, "Result should be nil")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data.Field, "Error fields should match")
}

// TestValidateAndDeduplicateOrders_InvalidField tests the case where an
// invalid field is provided.
func TestValidateAndDeduplicateOrders_InvalidField(t *testing.T) {
	orders := []Order{
		{Field: "invalid_field", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Error(t, err)
	assert.Nil(t, result, "Result should be nil")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data.Field, "Error fields should match")
}

// TestValidateAndDeduplicateOrders_EmptyOrders tests the case where an
// empty list of orders is provided.
func TestValidateAndDeduplicateOrders_EmptyOrders(t *testing.T) {
	orders := []Order{}
	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.NoError(t, err)
	assert.Empty(t, result, "Result should be empty")
}

// TestValidate_ValidOrder tests the case where a valid order is provided.
func TestValidate_ValidOrder(t *testing.T) {
	order := Order{
		Field:     "name",
		Direction: DIRECTION_ASC,
	}

	allowedFields := []string{"name", "age", "email"}

	err := validate(order, allowedFields)

	assert.NoError(t, err, "Expected no error for a valid order")
}

// TestValidate_InvalidDirection tests the case where an invalid direction
// is provided.
func TestValidate_InvalidDirection(t *testing.T) {
	order := Order{
		Field:     "name",
		Direction: "INVALID_DIRECTION",
	}

	allowedFields := []string{"name", "age", "email"}

	err := validate(order, allowedFields)

	assert.Error(t, err, "Expected an error for an invalid direction")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data.Field, "Error fields should match")
}

// TestValidate_InvalidField tests the case where an invalid field is provided.
func TestValidate_InvalidField(t *testing.T) {
	order := Order{
		Field:     "invalid_field",
		Direction: DIRECTION_ASC,
	}

	allowedFields := []string{"name", "age", "email"}

	err := validate(order, allowedFields)

	assert.Error(t, err, "Expected an error for an invalid field")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data.Field, "Errror fields should match")
}
