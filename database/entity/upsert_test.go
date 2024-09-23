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
	entity := &TestEntity{ID: 1, Name: "Alice"}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(0), nil)

	_, err := UpsertEntity(mockDB, "user", entity, inserter, projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertEntities_NormalOperation tests the normal operation of
// UpsertEntities.
func TestUpsertEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for multiple entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(0), nil)

	_, err := UpsertEntities(mockDB, "user", entities, inserter, projections)

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
		[]*TestEntity{}, // Updated to use []*TestEntity
		func(e *TestEntity) ([]string, []any) {
			return []string{}, []any{}
		},
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
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for multiple entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Simulate an error in the upsertMany call
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := UpsertEntities(mockDB, "user", entities, inserter, projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsert_NormalOperation tests the normal operation of upsert.
func TestUpsert_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entity and projections
	entity := &TestEntity{ID: 1, Name: "Alice", Age: 30}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
		{Column: "age", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name", "age"}, []any{e.ID, e.Name, e.Age}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := upsert(mockDB, "user", entity, inserter, projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsert_MissingProjections tests the case where the update projections are
// missing.
func TestUpsert_MissingProjections(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity
	entity := &TestEntity{ID: 1, Name: "Alice"}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	_, err := upsert(mockDB, "user", entity, inserter, []util.Projection{})

	assert.EqualError(t, err, "must provide update projections")
	mockDB.AssertExpectations(t)
}

// TestUpsert_MissingProjectionAlias tests the case where the update projection
// alias is missing.
func TestUpsert_MissingProjectionAlias(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity and projections with an empty alias
	entity := &TestEntity{ID: 1, Name: "Alice"}
	projections := []util.Projection{
		{Column: "name", Alias: ""},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	_, err := upsert(mockDB, "user", entity, inserter, projections)

	assert.EqualError(t, err, "must provide update projections alias")
	mockDB.AssertExpectations(t)
}

// TestUpsert_PrepareError tests the case where Prepare returns an error.
func TestUpsert_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entity and projections
	entity := &TestEntity{ID: 1, Name: "Alice"}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Simulate an error on Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := upsert(mockDB, "user", entity, inserter, projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsert_ExecError tests the case where Exec returns an error.
func TestUpsert_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entity and projections
	entity := &TestEntity{ID: 1, Name: "Alice"}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := upsert(mockDB, "user", entity, inserter, projections)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertManyQuery_NormalOperation tests upsertManyQuery with multiple
// entities and projections.
func TestUpsertManyQuery_NormalOperation(t *testing.T) {
	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
		{Column: "age", Alias: "test"},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name", "age"}, []any{e.ID, e.Name, e.Age}
	}

	query, values := upsertManyQuery(entities, "user", inserter, projections)

	expectedQuery := "INSERT INTO `user` (`id`, `name`, `age`) VALUES (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `age` = VALUES(`age`)"
	expectedValues := []any{1, "Alice", 30, 2, "Bob", 25}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpsertManyQuery_SingleEntity tests upsertManyQuery with a single entity.
func TestUpsertManyQuery_SingleEntity(t *testing.T) {
	// Test entity and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	query, values := upsertManyQuery(entities, "user", inserter, projections)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpsertManyQuery_EmptyEntities tests upsertManyQuery with no entities.
func TestUpsertManyQuery_EmptyEntities(t *testing.T) {
	// Test with no entities
	entities := []*TestEntity{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function (not used here since entities is empty)
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{}, []any{}
	}

	query, values := upsertManyQuery(entities, "user", inserter, projections)

	assert.Equal(t, "", query)
	assert.Equal(t, []any(nil), values)
}

// TestUpsertManyQuery_MissingUpdateProjections tests upsertManyQuery with no
// update projections.
func TestUpsertManyQuery_MissingUpdateProjections(t *testing.T) {
	// Test entities
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Call the function with an empty projection list
	query, values := upsertManyQuery(
		entities,
		"user",
		inserter,
		[]util.Projection{},
	)

	assert.Equal(
		t,
		"INSERT INTO `user` (`id`, `name`) VALUES (?, ?), (?, ?)",
		query,
	)
	assert.Equal(t, []any{1, "Alice", 2, "Bob"}, values)
}

// TestUpsertMany_NormalOperation tests the normal operation of upsertMany.
func TestUpsertMany_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := upsertMany(mockDB, entities, "user", inserter, projections)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpsertMany_EmptyEntities tests the case where the upserted entity list is
// empty.
func TestUpsertMany_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test with no entities
	entities := []*TestEntity{}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function (not used since entities is empty)
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{}, []any{}
	}

	_, err := upsertMany(mockDB, entities, "user", inserter, projections)

	assert.EqualError(t, err, "must provide entities to upsert")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_MissingUpdateProjections tests the case where the update
// projections are missing.
func TestUpsertMany_MissingUpdateProjections(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	_, err := upsertMany(
		mockDB,
		entities,
		"user",
		inserter,
		[]util.Projection{},
	)

	assert.EqualError(t, err, "must provide update projections")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_MissingAliasInProjections tests the case where the update
// projections alias is missing.
func TestUpsertMany_MissingAliasInProjections(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test entities
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
	}
	// Projections with an empty alias
	projections := []util.Projection{
		{Column: "name", Alias: ""},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	_, err := upsertMany(mockDB, entities, "user", inserter, projections)

	assert.EqualError(t, err, "must provide update projections alias")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_PrepareError tests the case with the SQL prepare error.
func TestUpsertMany_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	mockDB.On("Prepare", mock.Anything).
		Return(mockStmt, fmt.Errorf("prepare error"))

	_, err := upsertMany(mockDB, entities, "user", inserter, projections)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpsertMany_ExecError tests the case with the SQL execution error.
func TestUpsertMany_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, fmt.Errorf("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := upsertMany(mockDB, entities, "user", inserter, projections)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
}
