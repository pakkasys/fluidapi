package order

import (
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/stretchr/testify/assert"
)

// TestValidateAndTranslateToDBOrders tests the
// ValidateAndTranslateToDBOrders function
func TestValidateAndTranslateToDBOrders_ValidInput(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "age", Direction: DIRECTION_DESC},
	}

	allowedFields := []string{"name", "age", "email"}
	fieldTranslations := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.NoError(t, err, "Expected no error for valid orders")
	assert.Len(t, dbOrders, 2, "Expected two translated orders")
	assert.Equal(t, "users", dbOrders[0].Table, "Expected correct table for name order")
	assert.Equal(t, "user_name", dbOrders[0].Field, "Expected correct column for name order")
	assert.Equal(t, util.OrderAsc, dbOrders[0].Direction, "Expected ascending order for name order")
	assert.Equal(t, "users", dbOrders[1].Table, "Expected correct table for age order")
	assert.Equal(t, "user_age", dbOrders[1].Field, "Expected correct column for age order")
	assert.Equal(t, util.OrderDesc, dbOrders[1].Direction, "Expected descending order for age order")
}

// TestValidateAndTranslateToDBOrders tests the scenario where
// an invalid field is passed in
func TestValidateAndTranslateToDBOrders_InvalidField(t *testing.T) {
	orders := []Order{
		{Field: "invalid_field", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}
	fieldTranslations := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Error(t, err, "Expected an error for an invalid field")
	assert.Nil(t, dbOrders, "Expected no database orders for an invalid field")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data.Field, "Error fields should match")
}

// TestValidateAndTranslateToDBOrders tests the scenario where
// an invalid direction is passed in
func TestValidateAndTranslateToDBOrders_InvalidDirection(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: "INVALID_DIRECTION"},
	}

	allowedFields := []string{"name", "age", "email"}
	fieldTranslations := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Error(t, err, "Expected an error for an invalid direction")
	assert.Nil(t, dbOrders, "Expected no database orders for an invalid direction")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data.Field, "Error fields should match")
}

// TestValidateAndTranslateToDBOrders tests the scenario where
// a field not in the translation map is passed in
func TestValidateAndTranslateToDBOrders_FieldNotInTranslationMap(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}
	// Missing "name" field in the translation map to trigger the error
	fieldTranslations := map[string]dbfield.DBField{
		"age": {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Error(t, err, "Expected an error for a field not present in the translation map")
	assert.Nil(t, dbOrders, "Expected no database orders for a field not present in the translation map")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data.Field, "Error fields should match")
}

// TestToDBOrders_ValidOrders tests the scenario where
// valid orders are passed in
func TestToDBOrders_ValidOrders(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "age", Direction: DIRECTION_DESC},
	}

	fieldTranslations := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ToDBOrders(orders, fieldTranslations)

	assert.NoError(t, err, "Expected no error for valid orders")
	assert.Len(t, dbOrders, 2, "Expected two translated orders")
	assert.Equal(t, "users", dbOrders[0].Table, "Expected correct table for name order")
	assert.Equal(t, "user_name", dbOrders[0].Field, "Expected correct column for name order")
	assert.Equal(t, util.OrderAsc, dbOrders[0].Direction, "Expected ascending order for name order")
	assert.Equal(t, "users", dbOrders[1].Table, "Expected correct table for age order")
	assert.Equal(t, "user_age", dbOrders[1].Field, "Expected correct column for age order")
	assert.Equal(t, util.OrderDesc, dbOrders[1].Direction, "Expected descending order for age order")
}

// TestToDBOrders_InvalidField tests the scenario where
// an invalid field is passed in
func TestToDBOrders_InvalidField(t *testing.T) {
	orders := []Order{
		{Field: "invalid_field", Direction: DIRECTION_ASC},
	}

	fieldTranslations := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ToDBOrders(orders, fieldTranslations)

	assert.Error(t, err, "Expected an error for an invalid field")
	assert.Nil(t, dbOrders, "Expected no database orders for an invalid field")
	apiErr, ok := err.(*api.Error[InvalidOrderFieldErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data.Field, "Error fields should match")
}

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
