package selector

import (
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/predicate"
	"github.com/stretchr/testify/assert"
)

// TestSelectorString tests the String() method of the Selector struct.
func TestSelectorString(t *testing.T) {
	sel := Selector{
		Field:     "name",
		Predicate: predicate.EQUAL,
		Value:     "Alice",
	}

	expected := "name = Alice"
	assert.Equal(t, expected, sel.String(), "String() method should return the correct string representation")
}

// TestSelectors_GetByFields_SingleField tests the GetByFields method of
// Selectors when searching for a single field.
func TestSelectors_GetByFields_SingleField(t *testing.T) {
	selectors := Selectors{
		{Field: "name", Predicate: predicate.EQUAL, Value: "Alice"},
		{Field: "age", Predicate: predicate.GREATER, Value: 25},
		{Field: "email", Predicate: predicate.EQUAL, Value: "alic@example.com"},
	}

	result := selectors.GetByFields("name")

	assert.Len(t, result, 1, "Expected 1 selector for field 'name'")
	assert.Equal(t, "name", result[0].Field, "Field should match")
	assert.Equal(t, predicate.EQUAL, result[0].Predicate, "Predicate should match")
	assert.Equal(t, "Alice", result[0].Value, "Value should match")
}

// TestSelectors_GetByFields_MultipleFields tests the GetByFields method of
// Selectors when searching for multiple fields.
func TestSelectors_GetByFields_MultipleFields(t *testing.T) {
	selectors := Selectors{
		{Field: "name", Predicate: predicate.EQUAL, Value: "Alice"},
		{Field: "age", Predicate: predicate.GREATER, Value: 25},
		{Field: "email", Predicate: predicate.EQUAL, Value: "alic@example.com"},
	}

	result := selectors.GetByFields("name", "email")

	assert.Len(t, result, 2, "Expected 2 selectors for fields 'name' and 'email'")
	assert.Equal(t, "name", result[0].Field, "First result field should be 'name'")
	assert.Equal(t, "email", result[1].Field, "Second result field should be 'email'")
}

// TestSelectors_GetByFields_NoMatch tests the GetByFields method of Selectors
// when there are no matching fields.
func TestSelectors_GetByFields_NoMatch(t *testing.T) {
	selectors := Selectors{
		{Field: "name", Predicate: predicate.EQUAL, Value: "Alice"},
		{Field: "age", Predicate: predicate.GREATER, Value: 25},
	}

	result := selectors.GetByFields("email")

	assert.Len(t, result, 0, "Expected no selectors for field 'email'")
}

// TestToDBSelectors_ValidSelectors tests the successful translation of API
// selectors to DB selectors.
func TestToDBSelectors_ValidSelectors(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "name",
			Predicate:         predicate.EQUAL,
			Value:             "Alice",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL, predicate.NOT_EQUAL},
		},
		{
			Field:             "age",
			Predicate:         predicate.GREATER,
			Value:             25,
			AllowedPredicates: []predicate.Predicate{predicate.GREATER, predicate.LESS},
		},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.NoError(t, err, "Expected no error for valid selectors")
	assert.Len(t, dbSelectors, 2, "Expected two translated DB selectors")

	assert.Equal(t, "users", dbSelectors[0].Table, "Expected correct table for 'name' selector")
	assert.Equal(t, "user_name", dbSelectors[0].Field, "Expected correct column for 'name' selector")
	assert.Equal(t, util.Predicate("="), dbSelectors[0].Predicate, "Expected correct predicate for 'name' selector")
	assert.Equal(t, "Alice", dbSelectors[0].Value, "Expected correct value for 'name' selector")

	assert.Equal(t, "users", dbSelectors[1].Table, "Expected correct table for 'age' selector")
	assert.Equal(t, "user_age", dbSelectors[1].Field, "Expected correct column for 'age' selector")
	assert.Equal(t, util.Predicate(">"), dbSelectors[1].Predicate, "Expected correct predicate for 'age' selector")
	assert.Equal(t, 25, dbSelectors[1].Value, "Expected correct value for 'age' selector")
}

// TestToDBSelectors_InvalidPredicate tests the case when a selector has an
// invalid predicate.
func TestToDBSelectors_InvalidPredicate(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "name",
			Predicate:         predicate.EQUAL,
			Value:             "Alice",
			AllowedPredicates: []predicate.Predicate{predicate.NOT_EQUAL},
		},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Error(t, err, "Expected error for invalid predicate")
	assert.Nil(t, dbSelectors, "Expected no database selectors when predicate is not allowed")
	predicateErr, ok := err.(*api.Error[PredicateNotAllowedErrorData])
	assert.True(t, ok, "Expected error to be PREDICATE_NOT_ALLOWED")
	assert.Equal(t, predicate.EQUAL, predicateErr.Data.Predicate, "Expected error predicate to match the disallowed predicate")
}

// TestToDBSelectors_InvalidField tests the case when a selector field cannot be
// translated.
func TestToDBSelectors_InvalidField(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "unknown_field",
			Predicate:         predicate.EQUAL,
			Value:             "Alice",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for an unknown field")
	assert.Nil(t, dbSelectors, "Expected no database selectors for an unknown field")
	fieldErr, ok := err.(*api.Error[InvalidDatabaseSelectorTranslationErrorData])
	assert.True(t, ok, "Expected error to be INVALID_DATABASE_SELECTOR_TRANSLATION")
	assert.Equal(t, "unknown_field", fieldErr.Data.Field, "Expected error field to match the unknown field")
}

// TestToDBSelectors_InvalidDBPredicate tests the case when a predicate cannot
// be translated to a DB predicate.
func TestToDBSelectors_InvalidDBPredicate(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "name",
			Predicate:         "NONE",
			Value:             "Alice",
			AllowedPredicates: []predicate.Predicate{"NONE"},
		},
	}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for an invalid DB predicate")
	assert.Nil(t, dbSelectors, "Expected no database selectors for an invalid DB predicate")
	dbPredicateErr, ok := err.(*api.Error[InvalidPredicateErrorData])
	assert.True(t, ok, "Expected error to be INVALID_PREDICATE")
	assert.Equal(t, predicate.Predicate("NONE"), dbPredicateErr.Data.Predicate, "Expected error predicate to match the invalid predicate")
}

// TestToDBSelectors_EmptySelectors tests the case when the input selectors list
// is empty.
func TestToDBSelectors_EmptySelectors(t *testing.T) {
	apiSelectors := []Selector{}

	apiToDBFieldMap := map[string]dbfield.DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.NoError(t, err, "Expected no error for empty selectors")
	assert.Empty(t, dbSelectors, "Expected no database selectors for empty input")
}
