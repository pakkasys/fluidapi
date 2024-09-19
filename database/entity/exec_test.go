package entity

import (
	"errors"
	"testing"

	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
)

// TestRowsQuery_NormalOperation tests the normal operation of RowsQuery.
func TestRowsQuery_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	parameters := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", parameters).Return(mockRows, nil)

	rows, stmt, err := RowsQuery(mockDB, query, parameters)

	assert.NoError(t, err)
	assert.NotNil(t, rows)
	assert.NotNil(t, stmt)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

// TestRowsQuery_PrepareError tests the case where there is an error during
// query preparation.
func TestRowsQuery_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	parameters := []any{1}

	// Simulate an error during Prepare
	mockDB.On("Prepare", query).Return(nil, errors.New("prepare error"))

	rows, stmt, err := RowsQuery(mockDB, query, parameters)

	assert.Nil(t, rows)
	assert.Nil(t, stmt)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestRowsQuery_QueryError tests the case where there is an error during query
// execution.
func TestRowsQuery_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	parameters := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	// Simulate an error during Query
	mockStmt.On("Query", parameters).Return(nil, errors.New("query error"))
	mockStmt.On("Close").Return(nil)

	rows, stmt, err := RowsQuery(mockDB, query, parameters)

	assert.Nil(t, rows)
	assert.Nil(t, stmt)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestExecQuery_NormalOperation tests the normal operation of ExecQuery.
func TestExecQuery_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test query and parameters
	query := "UPDATE users SET name = ? WHERE id = ?"
	parameters := []any{"Alice", 1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Exec", parameters).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	// Call the function being tested
	result, err := ExecQuery(mockDB, query, parameters)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestExecQuery_PrepareError tests the case where there is an error during
// query preparation.
func TestExecQuery_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "UPDATE users SET name = ? WHERE id = ?"
	parameters := []any{"Alice", 1}

	// Simulate an error during Prepare
	mockDB.On("Prepare", query).Return(nil, errors.New("prepare error"))

	result, err := ExecQuery(mockDB, query, parameters)

	assert.Nil(t, result)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestExecQuery_ExecError tests the case where there is an error during query
// execution.
func TestExecQuery_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test query and parameters
	query := "UPDATE users SET name = ? WHERE id = ?"
	parameters := []any{"Alice", 1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	// Simulate an error during Exec
	mockStmt.On("Exec", parameters).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	result, err := ExecQuery(mockDB, query, parameters)

	assert.Nil(t, result)
	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
