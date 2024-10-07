package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestColumSelectorString_NormalCase tests the normal operation of the String
// method.
func TestColumSelectorString_NormalCase(t *testing.T) {
	selector := ColumSelector{
		Table:  "users",
		Column: "id",
	}

	result := selector.String()
	expected := "`users`.`id`"

	assert.Equal(t, expected, result)
}

// TestColumSelectorString_EmptyTable tests the case where the Table is empty.
func TestColumSelectorString_EmptyTable(t *testing.T) {
	selector := ColumSelector{
		Table:  "",
		Column: "id",
	}

	result := selector.String()
	expected := "``.`id`"

	assert.Equal(t, expected, result)
}

// TestColumSelectorString_EmptyColumn tests the case where the Column is empty.
func TestColumSelectorString_EmptyColumn(t *testing.T) {
	selector := ColumSelector{
		Table:  "users",
		Column: "",
	}

	result := selector.String()
	expected := "`users`.``"

	assert.Equal(t, expected, result)
}

// TestColumSelectorString_EmptyTableAndColumn tests the case where both the
// Table and Column are empty.
func TestColumSelectorString_EmptyTableAndColumn(t *testing.T) {
	selector := ColumSelector{
		Table:  "",
		Column: "",
	}

	result := selector.String()
	expected := "``.``"

	assert.Equal(t, expected, result)
}
