package entity

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUpsertEntity_NormalOperation tests the normal operation of UpsertEntity.
func TestUpsertEntity_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test entity and projections
	entity := &MockInserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entity.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(0), nil)

	_, err := UpsertEntity(mockDB, "user", entity, projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertEntity_ErrorFromUpsert tests the case where the upsert function
// returns an error.
func TestUpsertEntity_ErrorFromUpsert(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity and projections
	entity := &MockInserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entity.On("GetInserted").Return([]string{"id", "name"}, []any{1, "Alice"})

	// Simulate an error in the upsert call
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := UpsertEntity(mockDB, "user", entity, projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsertEntities_NormalOperation tests the normal operation of
// UpsertEntities.
func TestUpsertEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test entities and projections
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(0), nil)

	_, err := UpsertEntities(mockDB, "user", entities, projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertEntities_EmptyEntities tests the case where the entities list is
// empty.
func TestUpsertEntities_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	_, err := UpsertEntities(
		mockDB,
		"user",
		[]Inserter{},
		[]util.Projection{{Column: "name", Alias: "test"}},
	)

	assert.EqualError(t, err, "must provide entities to upsert")
	mockDB.AssertExpectations(t)
}

// TestUpsertEntities_ErrorFromUpsertMany tests the case where upsertMany
// returns an error.
func TestUpsertEntities_ErrorFromUpsertMany(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities and projections
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	// Simulate an error in the upsertMany call
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := UpsertEntities(mockDB, "user", entities, projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsert_NormalOperation tests the normal operation of upsert.
func TestUpsert_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entity and projections
	entity := &MockInserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
		{Column: "age", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entity.On("GetInserted").
		Return([]string{"id", "name", "age"}, []any{1, "Alice", 30})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := upsert(mockDB, "user", entity, projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsert_MissingProjections tests the case where the update projections
// are missing.
func TestUpsert_MissingProjections(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity
	entity := &MockInserter{}

	_, err := upsert(mockDB, "user", entity, []util.Projection{})

	assert.EqualError(t, err, "must provide update projections")
	mockDB.AssertExpectations(t)
}

// TestUpsert_MissingProjectionAlias tests the case where the update projection
// alias is missing.
func TestUpsert_MissingProjectionAlias(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity and projections with an empty alias
	entity := &MockInserter{}
	projections := []util.Projection{
		{Column: "name", Alias: ""},
	}

	_, err := upsert(mockDB, "user", entity, projections)

	assert.EqualError(t, err, "must provide update projections alias")
	mockDB.AssertExpectations(t)
}

// TestUpsert_PrepareError tests the case where Prepare returns an error.
func TestUpsert_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity and projections
	entity := &MockInserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entity.On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})

	// Simulate an error on Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := upsert(mockDB, "user", entity, projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsert_ExecError tests the case where Exec returns an error.
func TestUpsert_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entity and projections
	entity := &MockInserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entity.On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := upsert(mockDB, "user", entity, projections)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertManyQuery_NormalOperation tests upsertManyQuery with multiple
// entities and projections.
func TestUpsertManyQuery_NormalOperation(t *testing.T) {
	// Test entities and projections
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
		{Column: "age", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name", "age"}, []any{1, "Alice", 30})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name", "age"}, []any{2, "Bob", 25})

	query, values := upsertManyQuery(entities, "user", projections)

	expectedQuery := "INSERT INTO `user` (`id`, `name`, `age`) VALUES (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `age` = VALUES(`age`)"
	expectedValues := []any{1, "Alice", 30, 2, "Bob", 25}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpsertManyQuery_SingleEntity tests upsertManyQuery with a single entity.
func TestUpsertManyQuery_SingleEntity(t *testing.T) {
	// Test entity and projections
	entities := []Inserter{
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})

	query, values := upsertManyQuery(entities, "user", projections)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpsertManyQuery_EmptyEntities tests upsertManyQuery with no entities.
func TestUpsertManyQuery_EmptyEntities(t *testing.T) {
	// Test with no entities
	entities := []Inserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	query, values := upsertManyQuery(entities, "user", projections)

	assert.Equal(t, "", query)
	assert.Equal(t, []any(nil), values)
}

// TestUpsertManyQuery_MissingUpdateProjections tests upsertManyQuery with no
// update projections.
func TestUpsertManyQuery_MissingUpdateProjections(t *testing.T) {
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

	// Call the function with an empty projection list
	query, values := upsertManyQuery(entities, "user", []util.Projection{})

	assert.Equal(
		t,
		"INSERT INTO `user` (`id`, `name`) VALUES (?, ?), (?, ?)",
		query,
	)
	assert.Equal(t, []any{1, "Alice", 2, "Bob"}, values)
}

// TestUpsertMany tests the normal operation of upsertMany.
func TestUpsertMany_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entities and projections
	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
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

	_, err := upsertMany(mockDB, entities, "user", projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertMany_EmptyEntities tests the case where the upserted entity list
// is empty.
func TestUpsertMany_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities
	entities := []Inserter{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	_, err := upsertMany(mockDB, entities, "user", projections)

	assert.EqualError(t, err, "must provide entities to upsert")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_MissingUpdateProjections tests the case where the update
// projections are missing.
func TestUpsertMany_MissingUpdateProjections(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities
	entities := []Inserter{
		&MockInserter{},
	}

	_, err := upsertMany(mockDB, entities, "user", []util.Projection{})

	assert.EqualError(t, err, "must provide update projections")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_MissingAliasInProjections tests the case where the update
// projections alias is missing.
func TestUpsertMany_MissingAliasInProjections(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities
	entities := []Inserter{
		&MockInserter{},
	}
	// Projections with an empty alias
	projections := []util.Projection{
		{Column: "name", Alias: ""},
	}

	_, err := upsertMany(mockDB, entities, "user", projections)

	assert.EqualError(t, err, "must provide update projections alias")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_PrepareError tests the case with the SQL prepare error.
func TestUpsertMany_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	mockDB.On("Prepare", mock.Anything).
		Return(mockStmt, fmt.Errorf("prepare error"))

	_, err := upsertMany(mockDB, entities, "user", projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_ExecError tests the case with the SQL execution error.
func TestUpsertMany_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	entities := []Inserter{
		&MockInserter{},
		&MockInserter{},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Setup the MockInserter to return specific columns and values
	entities[0].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{1, "Alice"})
	entities[1].(*MockInserter).On("GetInserted").
		Return([]string{"id", "name"}, []any{2, "Bob"})

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, fmt.Errorf("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := upsertMany(mockDB, entities, "user", projections)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
}
