package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProjectionString_NoTableNoAlias tests the case where the Projection has
// no table and no alias.
func TestProjectionString_NoTableNoAlias(t *testing.T) {
	projection := Projection{
		Table:  "",
		Column: "column_name",
		Alias:  "",
	}

	result := projection.String()

	expected := "`column_name`"
	assert.Equal(t, expected, result)
}

// TestProjectionString_WithTableNoAlias tests the case where the Projection has
// a table but no alias.
func TestProjectionString_WithTableNoAlias(t *testing.T) {
	projection := Projection{
		Table:  "table_name",
		Column: "column_name",
		Alias:  "",
	}

	result := projection.String()

	expected := "`table_name`.`column_name`"
	assert.Equal(t, expected, result)
}

// TestProjectionString_WithTableAndAlias tests the case where the Projection
// has both a table and an alias.
func TestProjectionString_WithTableAndAlias(t *testing.T) {
	projection := Projection{
		Table:  "table_name",
		Column: "column_name",
		Alias:  "alias_name",
	}

	result := projection.String()

	expected := "`table_name`.`column_name` AS `alias_name`"
	assert.Equal(t, expected, result)
}

// TestProjectionString_NoTableWithAlias tests the case where the Projection has
// no table but has an alias.
func TestProjectionString_NoTableWithAlias(t *testing.T) {
	projection := Projection{
		Table:  "",
		Column: "column_name",
		Alias:  "alias_name",
	}

	result := projection.String()

	expected := "`column_name` AS `alias_name`"
	assert.Equal(t, expected, result)
}

// TestProjectionString_EmptyFields tests the case where all fields are empty.
func TestProjectionString_EmptyFields(t *testing.T) {
	projection := Projection{
		Table:  "",
		Column: "",
		Alias:  "",
	}

	result := projection.String()

	expected := "``"
	assert.Equal(t, expected, result)
}
