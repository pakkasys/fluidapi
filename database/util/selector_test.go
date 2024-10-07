package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetByField_Found tests the case where the selector with the given field
// is found.
func TestGetByField_Found(t *testing.T) {
	selectors := Selectors{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "name", Predicate: "=", Value: "Alice"},
	}

	selector := selectors.GetByField("name")

	assert.NotNil(t, selector)
	assert.Equal(t, "user", selector.Table)
	assert.Equal(t, "name", selector.Field)
	assert.Equal(t, "=", string(selector.Predicate))
	assert.Equal(t, "Alice", selector.Value)
}

// TestGetByField_NotFound tests the case where the selector with the given
// field is not found.
func TestGetByField_NotFound(t *testing.T) {
	selectors := Selectors{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "name", Predicate: "=", Value: "Alice"},
	}

	selector := selectors.GetByField("age")

	assert.Nil(t, selector)
}

// TestGetByFields_Found tests the case where the selectors with the given
// fields are found.
func TestGetByFields_Found(t *testing.T) {
	selectors := Selectors{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "name", Predicate: "=", Value: "Alice"},
		{Table: "user", Field: "age", Predicate: ">", Value: 25},
	}

	resultSelectors := selectors.GetByFields("name", "age")

	assert.Len(t, resultSelectors, 2)
	assert.Equal(t, "name", resultSelectors[0].Field)
	assert.Equal(t, "age", resultSelectors[1].Field)
}

// TestGetByFields_NotFound tests the case where none of the selectors with the
// given fields are found.
func TestGetByFields_NotFound(t *testing.T) {
	selectors := Selectors{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "name", Predicate: "=", Value: "Alice"},
	}

	resultSelectors := selectors.GetByFields("age", "address")

	assert.Len(t, resultSelectors, 0)
}

// TestGetByFields_PartialFound tests the case where some selectors with the
// given fields are found.
func TestGetByFields_PartialFound(t *testing.T) {
	selectors := Selectors{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "name", Predicate: "=", Value: "Alice"},
		{Table: "user", Field: "age", Predicate: ">", Value: 25},
	}

	resultSelectors := selectors.GetByFields("name", "address")

	assert.Len(t, resultSelectors, 1)
	assert.Equal(t, "name", resultSelectors[0].Field)
}
