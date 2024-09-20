package internal

import (
	"reflect"
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/stretchr/testify/assert"
)

// TestProcessSelectors_NoSelectors tests the case where no selectors are
// provided.
func TestProcessSelectors_NoSelectors(t *testing.T) {
	selectors := []util.Selector{}

	whereColumns, whereValues := ProcessSelectors(selectors)

	// Expect no columns and no values
	assert.Empty(t, whereColumns)
	assert.Empty(t, whereValues)
}

// TestProcessSelectors_SingleSelector tests the case where a single selector is
// provided.
func TestProcessSelectors_SingleSelector(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	whereColumns, whereValues := ProcessSelectors(selectors)

	expectedColumns := []string{"`user`.`id` = ?"}
	expectedValues := []any{1}

	assert.Equal(t, expectedColumns, whereColumns)
	assert.Equal(t, expectedValues, whereValues)
}

// TestProcessSelectors_MultipleSelectors tests the case where multiple
// selectors are provided.
func TestProcessSelectors_MultipleSelectors(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "age", Predicate: ">", Value: 18},
	}

	whereColumns, whereValues := ProcessSelectors(selectors)

	expectedColumns := []string{"`user`.`id` = ?", "`user`.`age` > ?"}
	expectedValues := []any{1, 18}

	assert.Equal(t, expectedColumns, whereColumns)
	assert.Equal(t, expectedValues, whereValues)
}

// TestProcessSelectors_WithInPredicate tests the case where a selector with
// "IN" predicate is provided.
func TestProcessSelectors_WithInPredicate(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "IN", Value: []int{1, 2, 3}},
	}

	whereColumns, whereValues := ProcessSelectors(selectors)

	expectedColumns := []string{"`user`.`id` IN (?,?,?)"}
	expectedValues := []any{1, 2, 3}

	assert.Equal(t, expectedColumns, whereColumns)
	assert.Equal(t, expectedValues, whereValues)
}

// TestProcessSelectors_WithNilValue tests the case where a selector with a nil
// value is provided.
func TestProcessSelectors_WithNilValue(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "deleted_at", Predicate: "=", Value: nil},
	}

	whereColumns, whereValues := ProcessSelectors(selectors)

	expectedColumns := []string{"`user`.`deleted_at` IS NULL"}

	assert.Equal(t, expectedColumns, whereColumns)
	assert.Empty(t, whereValues) // No values since it's a NULL condition
}

// TestProcessSelectors_WithDifferentPredicates tests the case where different
// predicates are provided.
func TestProcessSelectors_WithDifferentPredicates(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "name", Predicate: "LIKE", Value: "%Alice%"},
		{Table: "user", Field: "age", Predicate: "<", Value: 30},
	}

	whereColumns, whereValues := ProcessSelectors(selectors)

	expectedColumns := []string{"`user`.`name` LIKE ?", "`user`.`age` < ?"}
	expectedValues := []any{"%Alice%", 30}

	assert.Equal(t, expectedColumns, whereColumns)
	assert.Equal(t, expectedValues, whereValues)
}

// TestProcessSelectors_EmptyTableField tests the case where a selector with an
// empty table and field is provided.
func TestProcessSelectors_EmptyTableField(t *testing.T) {
	selectors := []util.Selector{
		{Table: "", Field: "", Predicate: "=", Value: 1},
	}

	whereColumns, whereValues := ProcessSelectors(selectors)

	expectedColumns := []string{"`` = ?"}
	expectedValues := []any{1}

	assert.Equal(t, expectedColumns, whereColumns)
	assert.Equal(t, expectedValues, whereValues)
}

// TestProcessSelector_WithInPredicate tests the processSelector function with
// an "IN" predicate.
func TestProcessSelector_WithInPredicate(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "id",
		Predicate: "IN",
		Value:     []int{1, 2, 3},
	}

	column, values := processSelector(selector)

	expectedColumn := "`user`.`id` IN (?,?,?)"
	expectedValues := []any{1, 2, 3}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessSelector_WithDefaultPredicate tests the processSelector function
// with a default predicate.
func TestProcessSelector_WithDefaultPredicate(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "name",
		Predicate: "=",
		Value:     "Alice",
	}

	column, values := processSelector(selector)

	expectedColumn := "`user`.`name` = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessInSelector_WithSliceValue tests the processInSelector function
// with a slice value.
func TestProcessInSelector_WithSliceValue(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "id",
		Predicate: "IN",
		Value:     []int{1, 2, 3},
	}

	column, values := processInSelector(selector)

	expectedColumn := "`user`.`id` IN (?,?,?)"
	expectedValues := []any{1, 2, 3}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessInSelector_WithNonSliceValue tests the processInSelector function
// with a non-slice value.
func TestProcessInSelector_WithNonSliceValue(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "id",
		Predicate: "IN",
		Value:     1,
	}

	column, values := processInSelector(selector)

	expectedColumn := "`user`.`id` IN (?)"
	expectedValues := []any{1}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessInSelector_EmptySlice tests the processInSelector function with an
// empty slice value.
func TestProcessInSelector_EmptySlice(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "id",
		Predicate: "IN",
		Value:     []int{},
	}

	column, values := processInSelector(selector)

	expectedColumn := "`user`.`id` IN ()"
	expectedValues := []any{}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessInSelector_WithStringSliceValue tests the processInSelector
// function with a slice of strings.
func TestProcessInSelector_WithStringSliceValue(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "name",
		Predicate: "IN",
		Value:     []string{"Alice", "Bob", "Charlie"},
	}

	column, values := processInSelector(selector)

	expectedColumn := "`user`.`name` IN (?,?,?)"
	expectedValues := []any{"Alice", "Bob", "Charlie"}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessInSelector_NilValue tests the processInSelector function with a
// nil value.
func TestProcessInSelector_NilValue(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "id",
		Predicate: "IN",
		Value:     nil,
	}

	column, values := processInSelector(selector)

	expectedColumn := "`user`.`id` IN (?)"
	expectedValues := []any{nil}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessDefaultSelector_WithValue tests the processDefaultSelector
// function with a non-nil value.
func TestProcessDefaultSelector_WithValue(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "name",
		Predicate: "=",
		Value:     "Alice",
	}

	column, values := processDefaultSelector(selector)

	expectedColumn := "`user`.`name` = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessDefaultSelector_WithoutTable tests the processDefaultSelector
// function without an empty field value.
func TestProcessDefaultSelector_WithoutTable(t *testing.T) {
	selector := util.Selector{
		Table:     "",
		Field:     "name",
		Predicate: "=",
		Value:     "Alice",
	}

	column, values := processDefaultSelector(selector)

	expectedColumn := "`name` = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedColumn, column)
	assert.Equal(t, expectedValues, values)
}

// TestProcessDefaultSelector_NilValue tests the processDefaultSelector function
// with a nil value and "=" predicate.
func TestProcessDefaultSelector_NilValue(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "name",
		Predicate: "=",
		Value:     nil,
	}

	column, values := processDefaultSelector(selector)

	expectedColumn := "`user`.`name` IS NULL"
	assert.Equal(t, expectedColumn, column)
	assert.Nil(t, values)
}

// TestProcessNullSelector_NotEquals tests the processNullSelector function with
// a "!=" predicate.
func TestProcessNullSelector_NotEquals(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "name",
		Predicate: "!=",
		Value:     nil,
	}

	column, values := processNullSelector(selector)

	expectedColumn := "`user`.`name` IS NOT NULL"
	assert.Equal(t, expectedColumn, column)
	assert.Nil(t, values)
}

// TestProcessNullSelector_InvalidPredicate tests the processNullSelector
// function with an unsupported predicate.
func TestProcessNullSelector_InvalidPredicate(t *testing.T) {
	selector := util.Selector{
		Table:     "user",
		Field:     "name",
		Predicate: ">",
		Value:     nil,
	}

	column, values := processNullSelector(selector)

	// Expect an empty column and nil values because of non-supported predicate
	assert.Equal(t, "", column)
	assert.Nil(t, values)
}

// TestProcessNullSelector_WithoutTable tests the processNullSelector function
// without a table.
func TestProcessNullSelector_WithoutTable(t *testing.T) {
	selector := util.Selector{
		Table:     "",
		Field:     "name",
		Predicate: "=",
		Value:     nil,
	}

	column, values := processNullSelector(selector)

	expectedColumn := "`name` IS NULL"
	assert.Equal(t, expectedColumn, column)
	assert.Nil(t, values)
}

// TestBuildNullClause_WithTable tests the buildNullClause function when a table
// is provided.
func TestBuildNullClause_WithTable(t *testing.T) {
	// Test case where a table is provided
	selector := util.Selector{
		Table: "user",
		Field: "name",
	}

	clause := buildNullClause(selector, "IS")
	expectedClause := "`user`.`name` IS NULL"

	assert.Equal(t, expectedClause, clause)
}

// TestBuildNullClause_WithoutTable tests the buildNullClause function when no
// table is provided.
func TestBuildNullClause_WithoutTable(t *testing.T) {
	// Test case where no table is provided
	selector := util.Selector{
		Table: "",
		Field: "name",
	}

	clause := buildNullClause(selector, "IS NOT")
	expectedClause := "`name` IS NOT NULL"

	assert.Equal(t, expectedClause, clause)
}

// TestBuildNullClause_EmptyField tests the buildNullClause function with an
// empty field.
func TestBuildNullClause_EmptyField(t *testing.T) {
	// Test case with an empty field
	selector := util.Selector{
		Table: "user",
		Field: "",
	}

	clause := buildNullClause(selector, "IS")
	expectedClause := "`user`.`` IS NULL"

	assert.Equal(t, expectedClause, clause)
}

// TestCreatePlaceholdersAndValues tests the createPlaceholdersAndValues
// function.
func TestCreatePlaceholdersAndValues(t *testing.T) {
	// Test case 1: Slice with multiple values
	values := []int{1, 2, 3}
	value := reflect.ValueOf(values)

	placeholders, actualValues := createPlaceholdersAndValues(value)
	expectedPlaceholders := "?,?,?"
	expectedValues := []any{1, 2, 3}

	assert.Equal(t, expectedPlaceholders, placeholders)
	assert.Equal(t, expectedValues, actualValues)

	// Test case 2: Empty slice
	emptyValues := []int{}
	value = reflect.ValueOf(emptyValues)

	placeholders, actualValues = createPlaceholdersAndValues(value)
	expectedPlaceholders = ""
	expectedValues = []any{}

	assert.Equal(t, expectedPlaceholders, placeholders)
	assert.Equal(t, expectedValues, actualValues)

	// Test case 3: Slice with a single value
	singleValue := []string{"test"}
	value = reflect.ValueOf(singleValue)

	placeholders, actualValues = createPlaceholdersAndValues(value)
	expectedPlaceholders = "?"
	expectedValues = []any{"test"}

	assert.Equal(t, expectedPlaceholders, placeholders)
	assert.Equal(t, expectedValues, actualValues)
}

// TestCreatePlaceholders tests the createPlaceholders function.
func TestCreatePlaceholders(t *testing.T) {
	// Test case 1: Creating placeholders for 3 values
	placeholders := createPlaceholders(3)
	expectedPlaceholders := "?,?,?"

	assert.Equal(t, expectedPlaceholders, placeholders)

	// Test case 2: Creating placeholders for 1 value
	placeholders = createPlaceholders(1)
	expectedPlaceholders = "?"

	assert.Equal(t, expectedPlaceholders, placeholders)

	// Test case 3: Creating placeholders for 0 values
	placeholders = createPlaceholders(0)
	expectedPlaceholders = ""

	assert.Equal(t, expectedPlaceholders, placeholders)
}
