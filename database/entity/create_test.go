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

	// Create an Inserter function
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	// Call CreateEntity
	entity := &TestEntity{ID: 1, Name: "Alice"}
	id, err := CreateEntity(entity, mockDB, "user", inserter)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestCreateEntity_InsertError tests the case where the insert function returns an error.
func TestCreateEntity_InsertError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Create an Inserter function
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	// Call CreateEntity
	entity := &TestEntity{ID: 1, Name: "Alice"}
	_, err := CreateEntity(entity, mockDB, "user", inserter)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestCreateEntities_NormalOperation tests the normal operation of CreateEntities with multiple entities.
func TestCreateEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Create an Inserter function
	inserter := func(entity *TestEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	// Call CreateEntities
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	id, err := CreateEntities(entities, mockDB, "user", inserter)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestCreateEntities_EmptyEntities tests the case where no entities are passed to CreateEntities.
func TestCreateEntities_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Call CreateEntities with an empty list
	id, err := CreateEntities([]*TestEntity{}, mockDB, "user", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), id)
	mockDB.AssertExpectations(t)
}

// TestCreateEntities_InsertError tests the case where an error occurs during insertion.
func TestCreateEntities_InsertError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Create an Inserter function
	inserter := func(entity *TestEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	// Call CreateEntities
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	_, err := CreateEntities(entities, mockDB, "user", inserter)

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
	// Inserter function that returns two columns and values
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	query, values := insertQuery(
		&TestEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertQuery_SingleColumnEntity tests insertQuery with an entity that has
// only one column.
func TestInsertQuery_SingleColumnEntity(t *testing.T) {
	// Inserter function that returns a single column and value
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id"}, []any{1}
	}

	query, values := insertQuery(&TestEntity{ID: 1}, "user", inserter)

	expectedQuery := "INSERT INTO `user` (`id`) VALUES (?)"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertQuery_NoColumns tests insertQuery with an entity that has no columns.
func TestInsertQuery_NoColumns(t *testing.T) {
	// Inserter function that returns no columns or values
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{}, []any{}
	}

	query, values := insertQuery(&TestEntity{}, "user", inserter)

	expectedQuery := "INSERT INTO `user` () VALUES ()"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsert_NormalOperation tests the normal operation of insert.
func TestInsert_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for the entity
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := insert(
		mockDB,
		&TestEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsert_PrepareError tests the case where Prepare returns an error.
func TestInsert_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Inserter function for the entity
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Simulate an error on Prepare
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))

	_, err := insert(
		mockDB,
		&TestEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestInsert_ExecError tests the case where Exec returns an error.
func TestInsert_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for the entity
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := insert(
		mockDB,
		&TestEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsertManyQuery_NormalOperation tests insertManyQuery with multiple
// entities.
func TestInsertManyQuery_NormalOperation(t *testing.T) {
	// Inserter function for multiple entities
	inserter := func(entity *TestEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	query, values := insertManyQuery(entities, "user", inserter)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?), (?, ?)"
	expectedValues := []any{1, "Alice", 2, "Bob"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertManyQuery_SingleEntity tests insertManyQuery with a single entity.
func TestInsertManyQuery_SingleEntity(t *testing.T) {
	// Inserter function for a single entity
	inserter := func(entity *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
	}

	query, values := insertManyQuery(entities, "user", inserter)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertManyQuery_NoEntities tests insertManyQuery with no entities.
func TestInsertManyQuery_NoEntities(t *testing.T) {
	entities := []*TestEntity{}

	query, values := insertManyQuery(entities, "user", nil)

	assert.Equal(t, "", query)
	assert.Equal(t, []any(nil), values)
}

// TestInsertMany_NormalOperation tests the normal operation of insertMany.
func TestInsertMany_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for multiple entities
	inserter := func(entity *TestEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := insertMany(mockDB, entities, "user", inserter)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsertMany_EmptyEntities tests the case where the inserted entity list is
// empty.
func TestInsertMany_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	entities := []*TestEntity{}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)

	_, err := insertMany(mockDB, entities, "user", nil)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

// TestInsertMany_PrepareError tests the case where Prepare returns an error.
func TestInsertMany_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Inserter function for multiple entities
	inserter := func(entity *TestEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := insertMany(mockDB, entities, "user", inserter)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestInsertMany_ExecError tests the case where Exec returns an error.
func TestInsertMany_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for multiple entities
	inserter := func(entity *TestEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := insertMany(mockDB, entities, "user", inserter)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
