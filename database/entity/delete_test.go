package entity

import (
	"errors"
	"strings"
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestDeleteEntities_NormalOperation tests the normal operation of
// DeleteEntities.
func TestDeleteEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test selectors and delete options
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}
	opts := DeleteOptions{
		Limit:  5,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(2), nil)

	rowsAffected, err := DeleteEntities(mockDB, "user", selectors, &opts)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), rowsAffected)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDeleteEntities_DeleteError tests the case where an error occurs during
// the delete operation.
func TestDeleteEntities_DeleteError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test selectors and delete options
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}
	opts := DeleteOptions{
		Limit:  5,
		Orders: nil,
	}

	// Simulate an error during the delete operation
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("delete error"))

	_, err := DeleteEntities(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "delete error")
	mockDB.AssertExpectations(t)
}

// TestDeleteEntities_RowsAffectedError tests the case where an error occurs
// when getting the number of rows affected.
func TestDeleteEntities_RowsAffectedError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test selectors and delete options
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}
	opts := DeleteOptions{
		Limit:  5,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	// Simulate an error when calling RowsAffected
	mockResult.On("RowsAffected").Return(int64(0), errors.New("rows affected error"))

	_, err := DeleteEntities(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "rows affected error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDelete_NormalOperation tests the normal operation of the delete function.
func TestDelete_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test selectors and options
	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}
	opts := DeleteOptions{
		Limit: 10,
		Orders: []util.Order{
			{Table: "user", Field: "name", Direction: "ASC"},
		},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := delete(mockDB, "user", selectors, &opts)

	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDelete_NoSelectors tests the case where no selectors are provided.
func TestDelete_NoSelectors(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Empty selectors and options
	selectors := []util.Selector{}
	opts := DeleteOptions{
		Limit:  0,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := delete(mockDB, "user", selectors, &opts)

	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDelete_PrepareError tests the case where an error occurs during SQL
// preparation.
func TestDelete_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test selectors and options
	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}
	opts := DeleteOptions{
		Limit:  0,
		Orders: nil,
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := delete(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestDelete_ExecError tests the case where an error occurs during SQL
// execution.
func TestDelete_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test selectors and options
	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}
	opts := DeleteOptions{
		Limit:  0,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error during Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := delete(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestWriteDeleteOptions_WithLimitAndOrders tests writeDeleteOptions with both
// limit and orders.
func TestWriteDeleteOptions_WithLimitAndOrders(t *testing.T) {
	// Create a DeleteOptions with a limit and orders
	orders := []util.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
		{Table: "user", Field: "age", Direction: "DESC"},
	}
	opts := DeleteOptions{Limit: 10, Orders: orders}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1 ORDER BY `user`.`name` ASC, `user`.`age` DESC LIMIT 10"

	assert.Equal(t, expectedSQL, builder.String())
}

// TestWriteDeleteOptions_WithOnlyOrders tests writeDeleteOptions with only
// orders and no limit.
func TestWriteDeleteOptions_WithOnlyOrders(t *testing.T) {
	// Create a DeleteOptions with only orders
	orders := []util.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
	}
	opts := DeleteOptions{Limit: 0, Orders: orders}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1 ORDER BY `user`.`name` ASC"

	assert.Equal(t, expectedSQL, builder.String())
}

// TestWriteDeleteOptions_WithOnlyLimit tests writeDeleteOptions with only a
// limit and no orders.
func TestWriteDeleteOptions_WithOnlyLimit(t *testing.T) {
	// Create a DeleteOptions with only a limit
	opts := DeleteOptions{Limit: 5, Orders: nil}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1 LIMIT 5"

	assert.Equal(t, expectedSQL, builder.String())
}

// TestWriteDeleteOptions_WithNoOptions tests writeDeleteOptions with no limit
// and no orders.
func TestWriteDeleteOptions_WithNoOptions(t *testing.T) {
	// Create an empty DeleteOptions with no limit and no orders
	opts := DeleteOptions{Limit: 0, Orders: nil}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1"

	assert.Equal(t, expectedSQL, builder.String())
}
