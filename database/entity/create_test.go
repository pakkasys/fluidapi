package entity

import (
	"errors"
	"testing"

	databaseerrors "github.com/pakkasys/fluidapi/database/errors"
	"github.com/pakkasys/fluidapi/database/internal"
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

// MockSQLResult is a mock implementation of the sql.Result interface.
type MockSQLResult struct {
	mock.Mock
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// TestCreateEntity_NormalOperation tests the normal operation of CreateEntity.
func TestCreateEntity_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Mock Inserter for entity
	inserter := &MockInserter{}

	// Setup the Inserter to return columns and values
	inserter.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	// Call CreateEntity
	id, err := CreateEntity(inserter, mockDB, "user")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestCreateEntity_InsertError tests the case where the insert function returns
// an error.
func TestCreateEntity_InsertError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Mock Inserter for entity
	inserter := &MockInserter{}

	// Setup the Inserter to return columns and values
	inserter.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	// Call CreateEntity
	_, err := CreateEntity(inserter, mockDB, "user")

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestCreateEntities_NormalOperation tests the normal operation of
// CreateEntities with multiple entities.
func TestCreateEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Mock Inserter for entities
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}

	// Setup the Inserter for each entity
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	// Call CreateEntities
	id, err := CreateEntities(entities, mockDB, "user")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestCreateEntities_EmptyEntities tests the case where no entities are passed
// to CreateEntities.
func TestCreateEntities_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Call CreateEntities with an empty list
	id, err := CreateEntities([]Inserter{}, mockDB, "user")

	assert.NoError(t, err)
	assert.Equal(t, int64(0), id)
	mockDB.AssertExpectations(t)
}

// TestCreateEntities_InsertError tests the case where an error occurs during
// insertion.
func TestCreateEntities_InsertError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Mock Inserter for entities
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}

	// Setup the Inserter for each entity
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	// Call CreateEntities
	_, err := CreateEntities(entities, mockDB, "user")

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestCheckInsertResult_NoError tests the case where there is no error and the
// result returns an ID.
func TestCheckInsertResult_NoError(t *testing.T) {
	mockResult := new(MockSQLResult)
	mockResult.On("LastInsertId").Return(int64(123), nil)

	id, err := checkInsertResult(mockResult, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
	mockResult.AssertExpectations(t)
}

// TestCheckInsertResult_DuplicateEntryError tests the case where there is a
// duplicate entry error.
func TestCheckInsertResult_DuplicateEntryError(t *testing.T) {
	// Create a custom MySQLError simulating a duplicate entry
	mysqlErr := internal.NewMySQLError(
		uint16(internal.MySQLDuplicateEntryErrorCode),
		"duplicate entry",
	)

	_, err := checkInsertResult(nil, mysqlErr)

	assert.Error(t, err)
	assert.EqualError(t, err, databaseerrors.DUPLICATE_ENTRY_ERROR_ID)

}

// TestCheckInsertResult_ForeignKeyConstraintError tests the case where there is
// a foreign key constraint error.
func TestCheckInsertResult_ForeignKeyConstraintError(t *testing.T) {
	// Create a custom MySQLError simulating a foreign key constraint failure
	mysqlErr := internal.NewMySQLError(
		uint16(internal.MySQLForeignConstraintErrorCode),
		"foreign key constraint fails",
	)

	_, err := checkInsertResult(nil, mysqlErr)

	assert.Error(t, err)
	assert.EqualError(t, err, databaseerrors.FOREIGN_CONSTRAINT_ERROR_ID)
}

// TestCheckInsertResult_GeneralError tests the case where a general error is
// passed.
func TestCheckInsertResult_GeneralError(t *testing.T) {
	// Simulate a general error
	generalErr := errors.New("some other error")

	_, err := checkInsertResult(nil, generalErr)

	assert.EqualError(t, err, "some other error")
}

// TestCheckInsertResult_LastInsertIdError tests the case where getting the last
// insert ID returns an error.
func TestCheckInsertResult_LastInsertIdError(t *testing.T) {
	mockResult := new(MockSQLResult)
	mockResult.On("LastInsertId").Return(
		int64(0),
		errors.New("last insert ID error"),
	)

	_, err := checkInsertResult(mockResult, nil)

	assert.EqualError(t, err, "last insert ID error")
	mockResult.AssertExpectations(t)
}

// TestGetInsertQueryColumnNames_MultipleColumns tests getInsertQueryColumnNames
// with multiple columns.
func TestGetInsertQueryColumnNames_MultipleColumns(t *testing.T) {
	// Multiple columns
	columns := []string{"id", "name", "age"}

	result := getInsertQueryColumnNames(columns)

	// Expected result
	expectedResult := "`id`, `name`, `age`"

	assert.Equal(t, expectedResult, result)
}

// TestGetInsertQueryColumnNames_SingleColumn tests getInsertQueryColumnNames
// with a single column.
func TestGetInsertQueryColumnNames_SingleColumn(t *testing.T) {
	// Single column
	columns := []string{"id"}

	result := getInsertQueryColumnNames(columns)

	// Expected result
	expectedResult := "`id`"

	assert.Equal(t, expectedResult, result)
}

// TestGetInsertQueryColumnNames_EmptyColumns tests getInsertQueryColumnNames
// with an empty list of columns.
func TestGetInsertQueryColumnNames_EmptyColumns(t *testing.T) {
	// Empty columns
	columns := []string{}

	result := getInsertQueryColumnNames(columns)

	// Expected result is an empty string
	expectedResult := ""

	assert.Equal(t, expectedResult, result)
}

// TestInsertQuery_NormalOperation tests insertQuery with a standard entity.
func TestInsertQuery_NormalOperation(t *testing.T) {
	// MockInserter that returns two columns and values
	inserter := &MockInserter{}

	// Setup the MockInserter to return specific columns and values
	inserter.On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})

	query, values := insertQuery(inserter, "user")

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertQuery_SingleColumnEntity tests insertQuery with an entity that has
// only one column.
func TestInsertQuery_SingleColumnEntity(t *testing.T) {
	// MockInserter that returns a single column and value
	inserter := &MockInserter{}

	// Setup the MockInserter to return one column and one value
	inserter.On("GetInserted").
		Return([]string{"id"}, []any{1})

	query, values := insertQuery(inserter, "user")

	expectedQuery := "INSERT INTO `user` (`id`) VALUES (?)"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertQuery_NoColumns tests insertQuery with an entity that has no
// columns.
func TestInsertQuery_NoColumns(t *testing.T) {
	// MockInserter that returns no columns and no values
	inserter := &MockInserter{}

	// Setup the MockInserter to return no columns or values
	inserter.On("GetInserted").
		Return([]string{}, []any{})

	query, values := insertQuery(inserter, "user")

	// We expect an empty query here because there are no columns
	expectedQuery := "INSERT INTO `user` () VALUES ()"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsert_NormalOperation tests the normal operation of insert.
func TestInsert_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entity
	inserter := &MockInserter{}

	// Setup the MockInserter to return specific columns and values
	inserter.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := insert(mockDB, inserter, "user")

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsert_PrepareError tests the case where Prepare returns an error.
func TestInsert_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity
	inserter := &MockInserter{}

	// Setup the MockInserter to return specific columns and values
	inserter.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Simulate an error on Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := insert(mockDB, inserter, "user")

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestInsert_ExecError tests the case where Exec returns an error.
func TestInsert_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entity
	inserter := &MockInserter{}

	// Setup the MockInserter to return specific columns and values
	inserter.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := insert(mockDB, inserter, "user")

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsertManyQuery_NormalOperation tests insertManyQuery with multiple
// entities.
func TestInsertManyQuery_NormalOperation(t *testing.T) {
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

	query, values := insertManyQuery(entities, "user")

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?), (?, ?)"
	expectedValues := []any{1, "Alice", 2, "Bob"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertManyQuery_SingleEntity tests insertManyQuery with a single entity.
func TestInsertManyQuery_SingleEntity(t *testing.T) {
	// Test entity
	entities := []Inserter{
		&MockInserter{},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})

	query, values := insertManyQuery(entities, "user")

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertManyQuery_NoEntities tests insertManyQuery with no entities.
func TestInsertManyQuery_NoEntities(t *testing.T) {
	// Test with no entities
	entities := []Inserter{}

	query, values := insertManyQuery(entities, "user")

	// Expected an empty query and nil values since no entities were provided
	assert.Equal(t, "", query)
	assert.Equal(t, []any(nil), values)
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

	_, err := insertMany(mockDB, entities, "user")

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

	_, err := insertMany(mockDB, entities, "user")

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

	_, err := insertMany(mockDB, entities, "user")

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

	_, err := insertMany(mockDB, entities, "user")

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
