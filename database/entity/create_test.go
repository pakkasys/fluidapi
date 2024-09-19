package entity

import (
	"errors"
	"testing"

	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInserter is a mock implementation of the Inserter interface.
type MockInserter struct {
	mock.Mock
}

func (m *MockInserter) GetInserted() (columns []string, values []any) {
	args := m.Called()
	return args.Get(0).([]string), args.Get(1).([]any)
}

// TestInsertMany_NormalOperation tests the normal operation of insertMany.
func TestInsertMany_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entities
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := insertMany(mockDB, entities, "users")

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsertMany_EmptyEntities tests the case where the inserted entity list
// is empty.
func TestInsertMany_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test with no entities
	entities := []Inserter{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)

	_, err := insertMany(mockDB, entities, "users")

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

// TestInsertMany_PrepareError tests the case where Prepare returns an error.
func TestInsertMany_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Simulate an error on Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := insertMany(mockDB, entities, "users")

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestInsertMany_ExecError tests the case where Exec returns an error.
func TestInsertMany_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entities
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := insertMany(mockDB, entities, "users")

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
