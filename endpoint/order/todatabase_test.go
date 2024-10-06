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
