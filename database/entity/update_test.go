package entity

import (
	"errors"
	"testing"

	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUpdateEntities_NormalOperation tests the normal operation where updates
// are successfully applied.
func TestUpdateEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test table name, updates, and selectors
	tableName := "user"
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(1), nil)

	rowsAffected, err :=
		UpdateEntities(mockDB, tableName, selectors, updates, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpdateEntities_NoUpdates tests the case where no updates are provided.
func TestUpdateEntities_NoUpdates(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test table name and selectors
	tableName := "user"
	updates := []Update{}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	rowsAffected, err :=
		UpdateEntities(mockDB, tableName, selectors, updates, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), rowsAffected)
	mockDB.AssertExpectations(t)
}

// TestUpdateEntities_Error tests the case where an error occurs during the
// update process.
func TestUpdateEntities_Error(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test table name, updates, and selectors
	tableName := "user"
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))
	mockSQLUtil.On("CheckDBError", mock.Anything).
		Return(errors.New("prepare error"))

	rowsAffected, err :=
		UpdateEntities(mockDB, tableName, selectors, updates, mockSQLUtil)

	assert.Equal(t, int64(0), rowsAffected)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
	mockSQLUtil.AssertExpectations(t)
}

// TestCheckUpdateResult_NormalOperation tests the normal operation where rows
// are affected.
func TestCheckUpdateResult_NormalOperation(t *testing.T) {
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Setup mock expectations
	mockResult.On("RowsAffected").Return(int64(1), nil)

	rowsAffected, err := checkUpdateResult(mockResult, nil, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
	mockResult.AssertExpectations(t)
}

// TestCheckUpdateResult_OtherError tests the case where a non-MySQL error
// occurs.
func TestCheckUpdateResult_OtherError(t *testing.T) {
	mockSQLUtil := new(entitymock.MockSQLUtil)

	otherErr := errors.New("some other error")
	mockSQLUtil.On("CheckDBError", otherErr).Return(otherErr)

	rowsAffected, err := checkUpdateResult(nil, otherErr, mockSQLUtil)

	assert.Equal(t, int64(0), rowsAffected)
	assert.EqualError(t, err, "some other error")
	mockSQLUtil.AssertExpectations(t)
}

// TestCheckUpdateResult_RowsAffectedError tests the case where an error occurs
// when retrieving rows affected.
func TestCheckUpdateResult_RowsAffectedError(t *testing.T) {
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Simulate an error when retrieving RowsAffected
	mockResult.On("RowsAffected").
		Return(int64(0), errors.New("rows affected error"))

	rowsAffected, err := checkUpdateResult(mockResult, nil, mockSQLUtil)

	assert.Equal(t, int64(0), rowsAffected)
	assert.EqualError(t, err, "rows affected error")
	mockResult.AssertExpectations(t)
}

// TestUpdate_NormalOperation tests the normal operation of the update function.
func TestUpdate_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test table name, updates, and selectors
	tableName := "user"
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := update(mockDB, tableName, updates, selectors)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpdate_PrepareError tests the case where an error occurs during the
// preparation of the statement.
func TestUpdate_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test table name, updates, and selectors
	tableName := "user"
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	result, err := update(mockDB, tableName, updates, selectors)

	assert.Nil(t, result)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpdate_ExecError tests the case where an error occurs during the
// execution of the statement.
func TestUpdate_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test table name, updates, and selectors
	tableName := "user"
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	result, err := update(mockDB, tableName, updates, selectors)

	assert.Nil(t, result)
	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpdate_EmptyUpdates tests the case where no updates are provided.
func TestUpdate_EmptyUpdates(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test table name and selectors
	tableName := "user"
	updates := []Update{}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := update(mockDB, tableName, updates, selectors)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpdateQuery_SingleUpdate tests the case where a single update is
// provided.
func TestUpdateQuery_SingleUpdate(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, values := updateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET name = ? WHERE `user`.`id` = ?"
	expectedValues := []any{"Alice", 1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_MultipleUpdates tests the case where multiple updates are
// provided.
func TestUpdateQuery_MultipleUpdates(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, values := updateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET name = ?, age = ? WHERE `user`.`id` = ?"
	expectedValues := []any{"Alice", 30, 1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_NoUpdates tests the case where no updates are provided.
func TestUpdateQuery_NoUpdates(t *testing.T) {
	updates := []Update{}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, values := updateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET  WHERE `user`.`id` = ?"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_NoSelectors tests the case where no selectors are provided.
func TestUpdateQuery_NoSelectors(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}

	selectors := []util.Selector{} // No selectors

	query, values := updateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET name = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_EmptyFields tests the case where updates and selectors have
// empty fields.
func TestUpdateQuery_EmptyFields(t *testing.T) {
	updates := []Update{
		{Field: "", Value: "Unknown"},
	}
	selectors := []util.Selector{
		{Table: "", Field: "", Predicate: "=", Value: nil},
	}

	query, values := updateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET  = ? WHERE `` IS NULL"
	expectedValues := []any{"Unknown"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestGetWhereClause_NoConditions tests the case where no conditions are
// provided.
func TestGetWhereClause_NoConditions(t *testing.T) {
	whereColumns := []string{}

	whereClause := getWhereClause(whereColumns)

	expectedWhereClause := ""
	assert.Equal(t, expectedWhereClause, whereClause)
}

// TestGetWhereClause_SingleCondition tests the case where a single condition is
// provided.
func TestGetWhereClause_SingleCondition(t *testing.T) {
	whereColumns := []string{"`user`.`id` = ?"}

	whereClause := getWhereClause(whereColumns)

	expectedWhereClause := "WHERE `user`.`id` = ?"
	assert.Equal(t, expectedWhereClause, whereClause)
}

// TestGetWhereClause_MultipleConditions tests the case where multiple
// conditions are provided.
func TestGetWhereClause_MultipleConditions(t *testing.T) {
	whereColumns := []string{"`user`.`id` = ?", "`user`.`age` > 18"}

	whereClause := getWhereClause(whereColumns)

	expectedWhereClause := "WHERE `user`.`id` = ? AND `user`.`age` > 18"
	assert.Equal(t, expectedWhereClause, whereClause)
}

// TestGetSetClause_SingleUpdate tests the case where a single update is
// provided.
func TestGetSetClause_SingleUpdate(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
	}

	setClause, values := getSetClause(updates)

	expectedSetClause := "name = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}

// TestGetSetClause_MultipleUpdates tests the case where multiple updates are
// provided.
func TestGetSetClause_MultipleUpdates(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}

	setClause, values := getSetClause(updates)

	expectedSetClause := "name = ?, age = ?"
	expectedValues := []any{"Alice", 30}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}

// TestGetSetClause_NoUpdates tests the case where no updates are provided.
func TestGetSetClause_NoUpdates(t *testing.T) {
	updates := []Update{}

	setClause, values := getSetClause(updates)

	expectedSetClause := ""
	expectedValues := []any{}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}

// TestGetSetClause_EmptyField tests the case where an update has an empty
// field.
func TestGetSetClause_EmptyField(t *testing.T) {
	updates := []Update{
		{Field: "", Value: "Unknown"},
	}

	setClause, values := getSetClause(updates)

	expectedSetClause := " = ?"
	expectedValues := []any{"Unknown"}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}
