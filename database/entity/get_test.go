package entity

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRowScanner is a mock implementation of RowScanner.
type MockRowScanner[T any] struct {
	mock.Mock
}

func (m *MockRowScanner[T]) Scan(row util.Row, entity *T) error {
	args := m.Called(row, entity)
	return args.Error(0)
}

// MockRowScannerMultiple is a mock implementation of RowScannerMultiple.
type MockRowScannerMultiple[T any] struct {
	mock.Mock
}

func (m *MockRowScannerMultiple[T]) Scan(rows util.Rows, entity *T) error {
	args := m.Called(rows, entity)
	return args.Error(0)
}

type TestEntity struct {
	ID   int
	Name string
}

// TestGetEntity_NormalOperation tests the normal operation where a single entity is successfully retrieved.
func TestGetEntity_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test table name and GetOptions
	tableName := "users"
	dbOptions := &GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)

	// Setup row scanning
	mockScanner.On("Scan", mockRow, &TestEntity{}).Return(nil).Once()

	// Call the function being tested
	entity, err := GetEntity(tableName, mockScanner.Scan, mockDB, dbOptions)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, entity)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntity_QueryError tests the case where an error occurs during query
// execution.
func TestGetEntity_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test table name and GetOptions
	tableName := "users"
	dbOptions := &GetOptions{}

	// Simulate an error during query execution
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("query error"))

	// Call the function being tested
	entity, err := GetEntity[TestEntity](tableName, nil, mockDB, dbOptions)

	// Assertions
	assert.Nil(t, entity)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestGetEntity_NoRows tests the case where sql.ErrNoRows is returned.
func TestGetEntity_NoRows(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test table name and GetOptions
	tableName := "users"
	dbOptions := &GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(sql.ErrNoRows).Once()
	mockScanner.On("Scan", mockRow, &TestEntity{}).Return(nil).Once()

	// Call the function being tested
	entity, err := GetEntity(tableName, mockScanner.Scan, mockDB, dbOptions)

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, entity) // No entity should be returned
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntity_RowScannerError tests the case where an error occurs during row
// scanning.
func TestGetEntity_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test table name and GetOptions
	tableName := "users"
	dbOptions := &GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)

	// Simulate an error during row scanning
	mockScanner.On("Scan", mockRow, &TestEntity{}).
		Return(errors.New("row scanner error")).Once()

	// Call the function being tested
	entity, err := GetEntity(tableName, mockScanner.Scan, mockDB, dbOptions)

	// Assertions
	assert.Nil(t, entity)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntityWithQuery_NormalOperation tests the normal operation where a
// single entity is successfully retrieved.
func TestGetEntityWithQuery_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)
	mockScanner.On("Scan", mockRow, &TestEntity{}).Return(nil).Once()

	// Call the function being tested
	entity, err := GetEntityWithQuery(
		"users",
		mockScanner.Scan,
		mockDB,
		query,
		params,
	)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, entity)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntityWithQuery_QueryError tests the case where an error occurs during
// query execution.
func TestGetEntityWithQuery_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	params := []any{1}

	// Simulate an error during query execution
	mockDB.On("Prepare", query).Return(nil, errors.New("query error"))

	// Call the function being tested
	entity, err := GetEntityWithQuery[TestEntity](
		"users",
		nil,
		mockDB,
		query,
		params,
	)

	// Assertions
	assert.Nil(t, entity)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestGetEntityWithQuery_NoRows tests the case where sql.ErrNoRows is returned.
func TestGetEntityWithQuery_NoRows(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(sql.ErrNoRows).Once()
	mockScanner.On("Scan", mockRow, &TestEntity{}).Return(nil).Once()

	// Call the function being tested
	entity, err := GetEntityWithQuery(
		"users",
		mockScanner.Scan,
		mockDB,
		query,
		params,
	)

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, entity) // No entity should be returned
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntityWithQuery_RowScannerError tests the case where an error occurs
// during row scanning.
func TestGetEntityWithQuery_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM users WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockStmt.On("Close").Return(nil)

	// Simulate an error during row scanning
	mockScanner.On("Scan", mockRow, &TestEntity{}).
		Return(errors.New("row scanner error")).Once()

	// Call the function being tested
	entity, err := GetEntityWithQuery(
		"users",
		mockScanner.Scan,
		mockDB,
		query,
		params,
	)

	// Assertions
	assert.Nil(t, entity)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntities_NormalOperation tests normal operation where multiple
// entities are successfully retrieved.
func TestGetEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and options
	tableName := "users"
	dbOptions := &GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()  // Simulate a row read
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)

	// Setup row scanning
	mockScanner.On("Scan", mockRows, &TestEntity{}).Return(nil).Once()

	entities, err := GetEntities(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.NoError(t, err)
	assert.NotNil(t, entities)
	assert.Len(t, entities, 1)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntities_QueryError tests the case where an error occurs during query
// execution.
func TestGetEntities_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and options
	tableName := "users"
	dbOptions := &GetOptions{}

	// Simulate an error during RowsQuery
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("query error"))

	entities, err := GetEntities[TestEntity](tableName, nil, mockDB, dbOptions)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestGetEntities_NoRows tests the case where sql.ErrNoRows is returned.
func TestGetEntities_NoRows(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and options
	tableName := "users"
	dbOptions := &GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(false).Once() // No rows found
	mockRows.On("Err").Return(sql.ErrNoRows).Once()

	entities, err := GetEntities(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.NoError(t, err)
	assert.Len(t, entities, 0) // No entities should be returned
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntities_RowScannerError tests the case where an error occurs during
// row scanning.
func TestGetEntities_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and options
	tableName := "users"
	dbOptions := &GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read

	// Simulate an error during row scanning
	mockScanner.On("Scan", mockRows, &TestEntity{}).
		Return(errors.New("row scanner error")).Once()

	entities, err := GetEntities(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntitiesWithQuery_NormalOperation tests the normal operation where
// multiple entities are successfully retrieved.
func TestGetEntitiesWithQuery_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()  // Simulate a row read
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)

	// Setup row scanning
	mockScanner.On("Scan", mockRows, &TestEntity{}).Return(nil).Once()

	entities, err := GetEntitiesWithQuery(
		"user",
		mockScanner.Scan,
		mockDB,
		query,
		params,
	)

	assert.NoError(t, err)
	assert.NotNil(t, entities)
	assert.Len(t, entities, 1)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntitiesWithQuery_QueryError tests the case where an error occurs
// during the query execution.
func TestGetEntitiesWithQuery_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Simulate an error during RowsQuery
	mockDB.On("Prepare", query).Return(nil, errors.New("query error"))

	entities, err := GetEntitiesWithQuery[TestEntity](
		"user",
		nil,
		mockDB,
		query,
		params,
	)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestGetEntitiesWithQuery_NoRows tests the case where no rows are found.
func TestGetEntitiesWithQuery_NoRows(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM users WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)

	// Simulate the `sql.ErrNoRows` case
	mockRows.On("Next").Return(false).Once() // No rows found
	mockRows.On("Err").Return(sql.ErrNoRows).Once()

	entities, err := GetEntitiesWithQuery(
		"users",
		mockScanner.Scan,
		mockDB,
		query,
		params,
	)

	assert.NoError(t, err)
	assert.Len(t, entities, 0) // No entities should be returned
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQueryMultiple_NormalOperation tests normal operation where multiple
// entities are successfully retrieved.
func TestQueryMultiple_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()  // Simulate a row read
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)
	mockScanner.On("Scan", mockRows, &TestEntity{}).Return(nil).Once()

	entities, err := queryMultiple(mockDB, query, params, mockScanner.Scan)

	assert.NoError(t, err)
	assert.NotNil(t, entities)
	assert.Len(t, entities, 1)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQueryMultiple_RowsQueryError tests the case where an error occurs during
// the query execution.
func TestQueryMultiple_RowsQueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Simulate an error during RowsQuery
	mockDB.On("Prepare", query).Return(nil, errors.New("query error"))

	entities, err := queryMultiple[TestEntity](mockDB, query, params, nil)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestQueryMultiple_RowScannerError tests the case where an error occurs during
// row scanning.
func TestQueryMultiple_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read

	// Simulate an error in the row scanner
	mockScanner.On("Scan", mockRows, &TestEntity{}).
		Return(errors.New("row scanner error"))

	entities, err := queryMultiple(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQueryMultiple_RowsErr tests the case where rows.Err() returns an error.
func TestQueryMultiple_RowsErr(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()  // Simulate a row read
	mockRows.On("Next").Return(false).Once() // No more rows
	mockScanner.On("Scan", mockRows, &TestEntity{}).Return(nil).Once()

	// Simulate an error in rows.Err()
	mockRows.On("Err").Return(errors.New("rows error"))

	entities, err := queryMultiple(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "rows error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQuerySingle_NormalOperation tests normal operation where a single entity
// is successfully retrieved.
func TestQuerySingle_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)
	mockScanner.On("Scan", mockRow, &TestEntity{}).Return(nil)

	entity, err := querySingle(mockDB, query, params, mockScanner.Scan)

	assert.NoError(t, err)
	assert.NotNil(t, entity)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQuerySingle_PrepareError tests the case where an error occurs during
// query preparation.
func TestQuerySingle_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Simulate an error during Prepare
	mockDB.On("Prepare", query).Return(nil, errors.New("prepare error"))

	entity, err := querySingle[TestEntity](mockDB, query, params, nil)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestQuerySingle_RowScannerError tests the case where the row scanner returns
// an error.
func TestQuerySingle_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockStmt.On("Close").Return(nil)

	// Simulate an error in the row scanner
	mockScanner.On("Scan", mockRow, &TestEntity{}).
		Return(errors.New("row scanner error"))

	entity, err := querySingle(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQuerySingle_RowErr tests the case where the row.Err() method returns an
// error.
func TestQuerySingle_RowErr(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockStmt.On("Close").Return(nil)

	// Setup the row scanner to succeed
	mockScanner.On("Scan", mockRow, &TestEntity{}).Return(nil)

	// Simulate an error returned by row.Err()
	mockRow.On("Err").Return(errors.New("row error"))

	entity, err := querySingle(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "row error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestProjectionsToStrings_NoProjections tests the case where no projections
// are provided.
func TestProjectionsToStrings_NoProjections(t *testing.T) {
	projections := []util.Projection{}
	projectionStrings := projectionsToStrings(projections)
	assert.Equal(t, []string{"*"}, projectionStrings)
}

// TestProjectionsToStrings_SingleProjection tests the case where a single
// projection is provided.
func TestProjectionsToStrings_SingleProjection(t *testing.T) {
	projections := []util.Projection{
		{Table: "user", Column: "name"},
	}

	projectionStrings := projectionsToStrings(projections)

	expected := []string{"`user`.`name`"}
	assert.Equal(t, expected, projectionStrings)
}

// TestProjectionsToStrings_MultipleProjections tests the case where multiple
// projections are provided.
func TestProjectionsToStrings_MultipleProjections(t *testing.T) {
	projections := []util.Projection{
		{Table: "user", Column: "name"},
		{Table: "orders", Column: "order_id"},
	}

	projectionStrings := projectionsToStrings(projections)

	expected := []string{"`user`.`name`", "`orders`.`order_id`"}
	assert.Equal(t, expected, projectionStrings)
}

// TestProjectionsToStrings_EmptyFields tests the case where a projection has
// empty fields.
func TestProjectionsToStrings_EmptyFields(t *testing.T) {
	projections := []util.Projection{
		{Table: "", Column: ""},
	}

	projectionStrings := projectionsToStrings(projections)

	expected := []string{"``"}
	assert.Equal(t, expected, projectionStrings)
}

// TestJoinClause_NoJoins tests the case where no joins are provided.
func TestJoinClause_NoJoins(t *testing.T) {
	joins := []util.Join{}
	joinClause := joinClause(joins)
	assert.Equal(t, "", joinClause)
}

// TestJoinClause_SingleJoin tests the case where a single join is provided.
func TestJoinClause_SingleJoin(t *testing.T) {
	joins := []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "orders",
			OnLeft: util.ColumSelector{
				Table:  "user",
				Column: "id",
			},
			OnRight: util.ColumSelector{
				Table:  "orders",
				Column: "user_id",
			},
		},
	}

	joinClause := joinClause(joins)

	expected := "INNER JOIN `orders` ON `user`.`id` = `orders`.`user_id`"
	assert.Equal(t, expected, joinClause)
}

// TestJoinClause_MultipleJoins tests the case where multiple joins are
// provided.
func TestJoinClause_MultipleJoins(t *testing.T) {
	joins := []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "order",
			OnLeft: util.ColumSelector{
				Table:  "user",
				Column: "id",
			},
			OnRight: util.ColumSelector{
				Table:  "order",
				Column: "user_id",
			},
		},
		{
			Type:  util.JoinTypeLeft,
			Table: "payments",
			OnLeft: util.ColumSelector{
				Table:  "user",
				Column: "id",
			},
			OnRight: util.ColumSelector{
				Table:  "payments",
				Column: "user_id",
			},
		},
	}

	joinClause := joinClause(joins)

	// Expect multiple JOIN clauses
	expected := "INNER JOIN `order` ON `user`.`id` = `order`.`user_id` LEFT JOIN `payments` ON `user`.`id` = `payments`.`user_id`"
	assert.Equal(t, expected, joinClause)
}

// TestJoinClause_EmptyFields tests the case where a join has empty fields.
func TestJoinClause_EmptyFields(t *testing.T) {
	joins := []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "",
			OnLeft: util.ColumSelector{
				Table:  "",
				Column: "",
			},
			OnRight: util.ColumSelector{
				Table:  "",
				Column: "",
			},
		},
	}

	joinClause := joinClause(joins)

	// Expect a malformed JOIN clause with empty fields
	expected := "INNER JOIN `` ON ``.`` = ``.``"
	assert.Equal(t, expected, joinClause)
}

// TestWhereClause_NoSelectors tests the case where no selectors are provided.
func TestWhereClause_NoSelectors(t *testing.T) {
	selectors := []util.Selector{}

	whereClause, whereValues := whereClause(selectors)

	// Expect an empty string and no values since there are no selectors
	assert.Equal(t, "", whereClause)
	assert.Empty(t, whereValues)
}

// TestWhereClause_SingleSelector tests the case where a single selector is
// provided.
func TestWhereClause_SingleSelector(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	whereClause, whereValues := whereClause(selectors)

	expectedClause := "WHERE `user`.`id` = ?"
	assert.Equal(t, expectedClause, whereClause)
	assert.Equal(t, []any{1}, whereValues)
}

// TestWhereClause_MultipleSelectors tests the case where multiple selectors are
// provided.
func TestWhereClause_MultipleSelectors(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "age", Predicate: ">", Value: 18},
	}

	whereClause, whereValues := whereClause(selectors)

	expectedClause := "WHERE `user`.`id` = ? AND `user`.`age` > ?"
	assert.Equal(t, expectedClause, whereClause)
	assert.Equal(t, []any{1, 18}, whereValues)
}

// TestWhereClause_DifferentPredicates tests the case where different predicates
// are provided.
func TestWhereClause_DifferentPredicates(t *testing.T) {
	selectors := []util.Selector{
		{Table: "user", Field: "name", Predicate: "LIKE", Value: "%Alice%"},
		{Table: "user", Field: "age", Predicate: "<", Value: 30},
	}

	whereClause, whereValues := whereClause(selectors)

	// Expect a WHERE clause with different predicates
	expectedClause := "WHERE `user`.`name` LIKE ? AND `user`.`age` < ?"
	assert.Equal(t, expectedClause, whereClause)
	assert.Equal(t, []any{"%Alice%", 30}, whereValues)
}

// TestBuildBaseGetQuery_NoOptions tests the case where no options are provided.
func TestBuildBaseGetQuery_NoOptions(t *testing.T) {
	getOptions := GetOptions{}

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT * FROM `user`"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithSelectors tests the case where selectors are
// provided.
func TestBuildBaseGetQuery_WithSelectors(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Selectors = []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` WHERE `user`.`id` = ?"
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, []any{1}, whereValues)
}

// TestBuildBaseGetQuery_WithOrders tests the case where orders are provided.
func TestBuildBaseGetQuery_WithOrders(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Orders = []util.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
	}

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` ORDER BY `user`.`name` ASC"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithProjections tests the case where projections are
// provided.
func TestBuildBaseGetQuery_WithProjections(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Projections = []util.Projection{
		{Table: "user", Column: "name", Alias: "user_name"},
	}

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT `user`.`name` AS `user_name` FROM `user`"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithJoins tests the case where joins are provided.
func TestBuildBaseGetQuery_WithJoins(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Joins = []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "order",
			OnLeft: util.ColumSelector{
				Table:  "user",
				Column: "id",
			},
			OnRight: util.ColumSelector{
				Table:  "order",
				Column: "user_id",
			},
		},
	}

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` INNER JOIN `order` ON `user`.`id` = `order`.`user_id`"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithLock tests the case where the lock option is set.
func TestBuildBaseGetQuery_WithLock(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Lock = true

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` FOR UPDATE"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithPage tests the case where pagination is provided.
func TestBuildBaseGetQuery_WithPage(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Page = &page.InputPage{Offset: 10, Limit: 20}

	query, whereValues := buildBaseGetQuery("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` LIMIT 20 OFFSET 10"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestGetLimitOffsetClauseFromPage_NoPage tests the case where no page is
// provided.
func TestGetLimitOffsetClauseFromPage_NoPage(t *testing.T) {
	var p *page.InputPage = nil
	limitOffsetClause := getLimitOffsetClauseFromPage(p)
	assert.Equal(t, "", limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_WithPage tests the case where a page with
// limit and offset is provided.
func TestGetLimitOffsetClauseFromPage_WithPage(t *testing.T) {
	p := &page.InputPage{Limit: 10, Offset: 20}

	limitOffsetClause := getLimitOffsetClauseFromPage(p)

	expected := "LIMIT 10 OFFSET 20"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_ZeroLimit tests the case where limit is 0.
func TestGetLimitOffsetClauseFromPage_ZeroLimit(t *testing.T) {
	p := &page.InputPage{Limit: 0, Offset: 20}

	limitOffsetClause := getLimitOffsetClauseFromPage(p)

	expected := "LIMIT 0 OFFSET 20"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_ZeroOffset tests the case where offset is 0.
func TestGetLimitOffsetClauseFromPage_ZeroOffset(t *testing.T) {
	p := &page.InputPage{Limit: 10, Offset: 0}

	limitOffsetClause := getLimitOffsetClauseFromPage(p)

	expected := "LIMIT 10 OFFSET 0"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_ZeroLimitAndOffset tests the case where both
// limit and offset are 0.
func TestGetLimitOffsetClauseFromPage_ZeroLimitAndOffset(t *testing.T) {
	p := &page.InputPage{Limit: 0, Offset: 0}

	limitOffsetClause := getLimitOffsetClauseFromPage(p)

	expected := "LIMIT 0 OFFSET 0"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetOrderClauseFromOrders_NoOrders tests the case where no orders are
// provided.
func TestGetOrderClauseFromOrders_NoOrders(t *testing.T) {
	orders := []util.Order{}
	orderClause := getOrderClauseFromOrders(orders)
	assert.Equal(t, "", orderClause)
}

// TestGetOrderClauseFromOrders_WithoutTable tests the case where there is no
// table in the order.
func TestGetOrderClauseFromOrders_WithoutTable(t *testing.T) {
	orders := []util.Order{
		{Field: "name", Direction: "ASC"},
	}

	orderClause := getOrderClauseFromOrders(orders)

	expected := "ORDER BY `name` ASC"
	assert.Equal(t, expected, orderClause)
}

// TestGetOrderClauseFromOrders_SingleOrder tests the case where a single order
// is provided.
func TestGetOrderClauseFromOrders_SingleOrder(t *testing.T) {
	orders := []util.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
	}

	orderClause := getOrderClauseFromOrders(orders)

	expected := "ORDER BY `user`.`name` ASC"
	assert.Equal(t, expected, orderClause)
}

// TestGetOrderClauseFromOrders_MultipleOrders tests the case where multiple
// orders are provided.
func TestGetOrderClauseFromOrders_MultipleOrders(t *testing.T) {
	orders := []util.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
		{Table: "user", Field: "age", Direction: "DESC"},
	}

	orderClause := getOrderClauseFromOrders(orders)

	expected := "ORDER BY `user`.`name` ASC, `user`.`age` DESC"
	assert.Equal(t, expected, orderClause)
}

// TestGetOrderClauseFromOrders_EmptyFields tests the case where orders have
// empty fields.
func TestGetOrderClauseFromOrders_EmptyFields(t *testing.T) {
	orders := []util.Order{
		{Table: "", Field: "", Direction: "ASC"},
	}

	orderClause := getOrderClauseFromOrders(orders)

	// Expect an ORDER BY clause with empty table and field
	expected := "ORDER BY `` ASC"
	assert.Equal(t, expected, orderClause)
}

// TestRowsToEntities_NoRowScannerMultiple tests the case where there is no
// RowScannerMultiple provided.
func TestRowsToEntities_NoRowScannerMultiple(t *testing.T) {
	entities, err := rowsToEntities[any](nil, nil)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "must provide rowScannerMultiple")
}

// TestRowsToEntities_NormalOperation tests normal operation with multiple rows.
func TestRowsToEntities_NormalOperation(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(true).Once()  // Simulate the first row read
	mockRows.On("Next").Return(true).Once()  // Simulate the second row read
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)           // No error in rows
	mockRowScanner.On("Scan", mockRows, &TestEntity{}).Return(nil).Twice()

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.NoError(t, err)
	assert.Len(t, entities, 2)
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

// TestRowsToEntities_NoRows tests the case where there are no rows.
func TestRowsToEntities_NoRows(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(false).Once() // No rows
	mockRows.On("Err").Return(nil)           // No error in rows

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.NoError(t, err)
	assert.Len(t, entities, 0)
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

// TestRowsToEntities_RowScannerError tests the case where the row scanner
// returns an error.
func TestRowsToEntities_RowScannerError(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(true).Once() // Simulate the first row read

	// Simulate an error in the row scanner
	mockRowScanner.On("Scan", mockRows, &TestEntity{}).
		Return(errors.New("row scanner error")).Once()

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "row scanner error")
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

// TestRowsToEntities_RowsError tests the case where rows return an error.
func TestRowsToEntities_RowsError(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(true).Once()  // Simulate the first row read
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(errors.New("rows error")).Once()
	mockRowScanner.On("Scan", mockRows, &TestEntity{}).Return(nil).Once()

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "rows error")
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}
