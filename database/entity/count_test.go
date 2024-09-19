package entity

import (
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCountEntities_NormalOperation tests the CountEntities function.
func TestCountEntities_NormalOperation(t *testing.T) {
	mockExecutor := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	// Setup mock expectations
	mockExecutor.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(nil)

	// Example table name and dbOptions
	tableName := "test_table"
	dbOptions := &DBOptionsCount{}

	// Call the function being tested
	count, err := CountEntities(mockExecutor, tableName, dbOptions)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 0, count) // Adjust as per the test case

	// Verify that all expectations were met
	mockExecutor.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestCountEntities_PrepareError tests the case where an error occurs during
// prepare call.
func TestCountEntities_PrepareError(t *testing.T) {
	mockExecutor := new(utilmock.MockDB)

	// Setup mock expectations for Prepare error
	mockExecutor.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	// Example table name and dbOptions
	tableName := "test_table"
	dbOptions := &DBOptionsCount{}

	// Call the function being tested
	count, err := CountEntities(mockExecutor, tableName, dbOptions)

	// Assertions
	assert.Equal(t, 0, count)
	assert.EqualError(t, err, "prepare error")

	// Verify that all expectations were met
	mockExecutor.AssertExpectations(t)
}

// TestCountEntities_QueryRowError tests the case where an error occurs during
// query row call.
func TestCountEntities_QueryRowError(t *testing.T) {
	mockExecutor := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	// Setup mock expectations
	mockExecutor.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("query row error"))

	// Example table name and dbOptions
	tableName := "test_table"
	dbOptions := &DBOptionsCount{}

	// Call the function being tested
	count, err := CountEntities(mockExecutor, tableName, dbOptions)

	// Assertions
	assert.Equal(t, 0, count)
	assert.EqualError(t, err, "query row error")

	// Verify that all expectations were met
	mockExecutor.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestBuildBaseCountQuery_NoSelectorsNoJoins tests buildBaseCountQuery with no
// selectors or joins.
func TestBuildBaseCountQuery_NoSelectorsNoJoins(t *testing.T) {
	tableName := "test_table"
	dbOptions := &DBOptionsCount{}

	query, whereValues := buildBaseCountQuery(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table`"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}

// TestBuildBaseCountQuery_WithSelectors tests buildBaseCountQuery only
// selectors.
func TestBuildBaseCountQuery_WithSelectors(t *testing.T) {
	tableName := "test_table"
	dbOptions := &DBOptionsCount{
		Selectors: []util.Selector{
			{Table: "test_table", Field: "id", Predicate: "=", Value: 1},
		},
	}

	query, whereValues := buildBaseCountQuery(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table`  WHERE `test_table`.`id` = ?"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}

// TestBuildBaseCountQuery_WithJoins tests buildBaseCountQuery with joins only.
func TestBuildBaseCountQuery_WithJoins(t *testing.T) {
	tableName := "test_table"
	dbOptions := &DBOptionsCount{
		Joins: []util.Join{
			{
				Type:  util.JoinTypeInner,
				Table: "other_table",
				OnLeft: util.ColumSelector{
					Table:  "test_table",
					Column: "id",
				},
				OnRight: util.ColumSelector{
					Table:  "other_table",
					Column: "ref_id",
				},
			},
		},
	}
	query, whereValues := buildBaseCountQuery(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table` INNER JOIN `other_table` ON `test_table`.`id` = `other_table`.`ref_id`"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}

// TestBuildBaseCountQuery_WithSelectorsAndJoins tests buildBaseCountQuery with
// both selectors and joins.
func TestBuildBaseCountQuery_WithSelectorsAndJoins(t *testing.T) {
	tableName := "test_table"
	dbOptions := &DBOptionsCount{
		Selectors: []util.Selector{
			{Table: "test_table", Field: "id", Predicate: "=", Value: 1},
		},
		Joins: []util.Join{
			{
				Type:  util.JoinTypeInner,
				Table: "other_table",
				OnLeft: util.ColumSelector{
					Table:  "test_table",
					Column: "id",
				},
				OnRight: util.ColumSelector{
					Table:  "other_table",
					Column: "ref_id",
				},
			},
		},
	}

	query, whereValues := buildBaseCountQuery(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table` INNER JOIN `other_table` ON `test_table`.`id` = `other_table`.`ref_id` WHERE `test_table`.`id` = ?"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}
