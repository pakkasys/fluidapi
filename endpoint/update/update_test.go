package update

import (
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/stretchr/testify/assert"
)

// TestToDBUpdates_ValidInput tests the successful translation of updates to
// database updates.
func TestToDBUpdates_ValidInput(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.NoError(t, err, "Expected no error for valid updates")
	assert.Len(t, dbUpdates, 2, "Expected two translated updates")
	assert.Equal(t, "user_name", dbUpdates[0].Field, "Expected correct column for name update")
	assert.Equal(t, "Alice", dbUpdates[0].Value, "Expected correct value for name update")
	assert.Equal(t, "user_age", dbUpdates[1].Field, "Expected correct column for age update")
	assert.Equal(t, 30, dbUpdates[1].Value, "Expected correct value for age update")
}

// TestToDBUpdates_InvalidField tests the case when an update field cannot be
// translated.
func TestToDBUpdates_InvalidField(t *testing.T) {
	updates := []Update{
		{Field: "unknown_field", Value: "value"},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for an unknown field")
	assert.Nil(t, dbUpdates, "Expected no database updates for an unknown field")
	updateErr, ok := err.(*api.Error[InvalidDatabaseUpdateTranslationErrorData])
	assert.True(t, ok, "Expected error to be INVALID_DATABASE_UPDATE_TRANSLATION")
	assert.Equal(t, "unknown_field", updateErr.Data.Field, "Expected error field to match the unknown field")
}

// TestToDBUpdates_EmptyUpdates tests the case when the input updates list is
// empty.
func TestToDBUpdates_EmptyUpdates(t *testing.T) {
	updates := []Update{}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.NoError(t, err, "Expected no error for empty updates")
	assert.Empty(t, dbUpdates, "Expected no database updates for empty input")
}

// TestToDBUpdates_MultipleInvalidFields tests the case when multiple update
// fields cannot be translated.
func TestToDBUpdates_MultipleInvalidFields(t *testing.T) {
	updates := []Update{
		{Field: "unknown_field_1", Value: "value1"},
		{Field: "unknown_field_2", Value: "value2"},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for unknown fields")
	assert.Nil(t, dbUpdates, "Expected no database updates for unknown fields")
	updateErr, ok := err.(*api.Error[InvalidDatabaseUpdateTranslationErrorData])
	assert.True(t, ok, "Expected error to be of type InvalidDatabaseUpdateTranslationError")
	assert.Contains(t, []string{"unknown_field_1", "unknown_field_2"}, updateErr.Data.Field, "Expected error field to match one of the unknown fields")
}
