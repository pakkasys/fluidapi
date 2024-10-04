package entity

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/pakkasys/fluidapi/endpoint/page"
	endpointutil "github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetByField_Found tests the case where a field is found in the
// UpdateOptionList.
func TestGetByField_Found(t *testing.T) {
	updates := Updates{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
		{Field: "email", Value: "alice@example.com"},
	}

	result := updates.GetByField("age")

	assert.NotNil(t, result)
	assert.Equal(t, "age", result.Field)
	assert.Equal(t, 30, result.Value)
}

// TestGetByField_NotFound tests the case where a field is not found in the
// UpdateOptionList.
func TestGetByField_NotFound(t *testing.T) {
	updates := Updates{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}

	result := updates.GetByField("address")
	assert.Nil(t, result)
}

// TestGetByField_EmptyList tests the case where the UpdateOptionList is empty.
func TestGetByField_EmptyList(t *testing.T) {
	updates := Updates{}
	result := updates.GetByField("name")
	assert.Nil(t, result)
}

// TestGetByField_CaseSensitivity tests the case where field names differ by
// case.
func TestGetByField_CaseSensitivity(t *testing.T) {
	updates := Updates{
		{Field: "Name", Value: "Alice"},
		{Field: "age", Value: 30},
	}

	result := updates.GetByField("name")

	assert.Nil(t, result)

	result = updates.GetByField("Name")

	assert.NotNil(t, result)
	assert.Equal(t, "Name", result.Field)
	assert.Equal(t, "Alice", result.Value)
}

// TestGetByField_MultipleMatches tests the case where multiple fields have the
// same name.
func TestGetByField_MultipleMatches(t *testing.T) {
	updates := Updates{
		{Field: "name", Value: "Alice"},
		{Field: "name", Value: "Bob"},
		{Field: "age", Value: 30},
	}

	result := updates.GetByField("name")

	assert.NotNil(t, result)
	assert.Equal(t, "name", result.Field)
	assert.Equal(t, "Alice", result.Value)
}

// TestExecuteManagedTransaction_SuccessfulTransaction tests the scenario where
// the transactional function is executed successfully.
func TestExecuteManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	txHelpers := TXHelpers[string]{
		GetTxFn: getTxFn,
	}

	transactionalFunc := func(tx util.Tx) (string, error) {
		return "transaction-success", nil
	}

	mockTx.On("Commit").Return(nil)

	result, err := txHelpers.ExecuteManagedTransaction(ctx, transactionalFunc)

	assert.NoError(t, err)
	assert.Equal(t, "transaction-success", result)

	mockTx.AssertExpectations(t)
}

// TestExecuteManagedTransaction_TransactionFunctionError tests the scenario
// where the transactional function returns an error.
func TestExecuteManagedTransaction_TransactionFunctionError(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	txHelpers := TXHelpers[string]{
		GetTxFn: getTxFn,
	}

	transactionalFunc := func(tx util.Tx) (string, error) {
		return "", errors.New("transactional error")
	}

	mockTx.On("Rollback").Return(nil)

	result, err := txHelpers.ExecuteManagedTransaction(ctx, transactionalFunc)

	assert.EqualError(t, err, "transactional error")
	assert.Equal(t, "", result)

	mockTx.AssertExpectations(t)
}

// TestCreateEntity_CreateWithoutUpsertOptions tests creating an entity without
// UpsertOptions.
func TestCreateEntity_CreateWithoutUpsertOptions(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entity := &TestEntity{ID: 1, Name: "Alice"}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	result, err := entityHelpers.CreateEntity(mockPreparer, entity, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity, result)
	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestCreateEntity_UpsertWithOptions tests creating or upserting an entity
// using UpsertOptions.
func TestCreateEntity_UpsertWithOptions(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entity := &TestEntity{ID: 1, Name: "Alice"}

	opts := &UpsertOptions{
		UpdateProjection: []util.Projection{
			{Column: "name", Alias: "test_alias"},
		},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	result, err := entityHelpers.CreateEntity(mockPreparer, entity, opts)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestCreateEntity_UpsertEntityFailure tests the scenario where the UpsertEntity
// function fails.
func TestCreateEntity_UpsertEntityFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entity := &TestEntity{ID: 1, Name: "Alice"}

	opts := &UpsertOptions{
		UpdateProjection: []util.Projection{{Column: "name", Alias: "test_alias"}},
	}

	expectedErr := errors.New("upsert error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)

	result, err := entityHelpers.CreateEntity(mockPreparer, entity, opts)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
}

// TestCreateEntity_EntityCreationFailure tests the scenario where creating the
// entity fails.
func TestCreateEntity_EntityCreationFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entity := &TestEntity{ID: 1, Name: "Alice"}
	expectedErr := errors.New("exec error")

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, expectedErr)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)

	result, err := entityHelpers.CreateEntity(mockPreparer, entity, nil)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestCreateEntityWithManagedTransaction_SuccessfulTransaction tests the
// scenario where the entity creation is executed successfully.
func TestCreateEntityWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entity := &TestEntity{ID: 1, Name: "Alice"}

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	mockTx.On("Commit").Return(nil)

	result, err := entityHelpers.CreateEntityWithManagedTransaction(
		ctx,
		entity,
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity, result)

	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestCreateEntityWithManagedTransaction_CreateError tests the scenario where
// an error occurs during the entity creation.
func TestCreateEntityWithManagedTransaction_CreateError(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entity := &TestEntity{ID: 1, Name: "Alice"}

	expectedErr := errors.New("exec error")
	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, expectedErr)
	mockStmt.On("Close").Return(nil)
	mockTx.On("Rollback").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)

	result, err := entityHelpers.CreateEntityWithManagedTransaction(
		ctx,
		entity,
		nil,
	)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockStmt.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestCreateEntities_CreateSuccess tests the scenario where the entity
// creation is successful.
func TestCreateEntities_CreateSuccess(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	result, err := entityHelpers.CreateEntities(mockPreparer, entities, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entities, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestCreateEntities_UpsertError(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	opts := &UpsertOptions{
		UpdateProjection: []util.Projection{
			{Column: "name", Alias: "test_alias"},
		},
	}

	expectedErr := errors.New("upsert error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)

	result, err := entityHelpers.CreateEntities(mockPreparer, entities, opts)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
}

// TestCreateEntities_PreparerFailure tests the scenario where preparing the
// query fails.
func TestCreateEntities_PreparerFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	expectedErr := errors.New("prepare error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)

	result, err := entityHelpers.CreateEntities(mockPreparer, entities, nil)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
}

// TestCreateEntitiesWithManagedTransaction_SuccessfulTransaction tests the
// scenario where multiple entities are successfully created.
func TestCreateEntitiesWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		InserterFn: func(entity *TestEntity) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		},
		SQLUtil: mockSQLUtil,
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)
	mockTx.On("Commit").Return(nil)

	result, err := entityHelpers.CreateEntitiesWithManagedTransaction(
		ctx,
		entities,
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entities, result)

	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestCreateEntitiesWithManagedTransaction_GetTxError tests the scenario where
// starting the transaction fails.
func TestCreateEntitiesWithManagedTransaction_GetTxError(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return nil, errors.New("failed to start transaction")
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
	}

	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	result, err := entityHelpers.CreateEntitiesWithManagedTransaction(
		ctx,
		entities,
		nil,
	)

	assert.EqualError(t, err, "failed to start transaction")
	assert.Nil(t, result)
}

// TestGetEntity_SuccessfulRetrieval tests the scenario where the entity is
// successfully retrieved.
func TestGetEntity_SuccessfulRetrieval(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		ScanRowFn: func(row util.Row, entity *TestEntity) error {
			entity.ID = 1
			entity.Name = "Alice"
			return nil
		},
	}

	getOpts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{{Field: "id", Predicate: "=", Value: 1}},
		},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)

	result, err := entityHelpers.GetEntity(mockPreparer, getOpts)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Alice", result.Name)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestGetEntity_EntityNotFoundFnDefined tests the scenario where the entity
// not found function is defined.
func TestGetEntity_EntityNotFoundFnDefined(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		ScanRowFn: func(row util.Row, entity *TestEntity) error {
			return nil
		},
		EntityNotFoundFn: func() error {
			return errors.New("custom entity not found")
		},
	}

	getOpts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{{Field: "id", Predicate: "=", Value: 1}},
		},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(sql.ErrNoRows)

	result, err := entityHelpers.GetEntity(mockPreparer, getOpts)

	assert.EqualError(t, err, "custom entity not found")
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestGetEntity_EntityNotFoundFnDefined tests the scenario where the entity
// not found function is not defined.
func TestGetEntity_EntityNotFoundFnNotDefined(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		ScanRowFn: func(row util.Row, entity *TestEntity) error {
			return nil
		},
	}

	getOpts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{{Field: "id", Predicate: "=", Value: 1}},
		},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(sql.ErrNoRows)

	result, err := entityHelpers.GetEntity(mockPreparer, getOpts)

	assert.EqualError(t, err, "entity not found")
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestGetEntity_PreparerFailure tests the scenario where preparing the query
// fails.
func TestGetEntity_PreparerFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		ScanRowFn: func(row util.Row, entity *TestEntity) error {
			entity.ID = 1
			entity.Name = "Alice"
			return nil
		},
	}

	getOpts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{{Field: "id", Predicate: "=", Value: 1}},
		},
	}

	expectedErr := errors.New("prepare error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)

	result, err := entityHelpers.GetEntity(mockPreparer, getOpts)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
}

// TestGetEntityWithManagedTransaction_SuccessfulTransaction tests the
// scenario where the entity is successfully retrieved.
func TestGetEntityWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	// mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		ScanRowFn: func(row util.Row, entity *TestEntity) error {
			entity.ID = 1
			entity.Name = "Alice"
			return nil
		},
	}

	getOpts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{{Field: "id", Predicate: "=", Value: 1}},
		},
	}

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)

	mockTx.On("Commit").Return(nil)

	result, err := entityHelpers.GetEntityWithManagedTransaction(ctx, getOpts)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Alice", result.Name)

	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestGetEntityWithManagedTransaction_GetTxError tests the scenario where
// starting the transaction fails.
func TestGetEntityWithManagedTransaction_GetTxError(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return nil, errors.New("failed to start transaction")
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
	}

	getOpts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{{Field: "id", Predicate: "=", Value: 1}},
		},
	}

	result, err := entityHelpers.GetEntityWithManagedTransaction(ctx, getOpts)

	assert.EqualError(t, err, "failed to start transaction")
	assert.Nil(t, result)
}

// TestGetEntities_SuccessfulRetrieval tests the scenario where multiple
// entities are successfully retrieved.
func TestGetEntities_SuccessfulRetrieval(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		ScanRowsFn: func(rows util.Rows, entity *TestEntity) error {
			entity.ID = 1
			entity.Name = "Alice"
			return nil
		},
	}

	opts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{
				{Field: "active", Predicate: "=", Value: true},
			},
			Orders: []util.Order{
				{Field: "created_at", Direction: util.OrderAsc},
			},
			Page: &page.InputPage{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Err").Return(nil)

	result, err := entityHelpers.GetEntities(mockPreparer, opts)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "Alice", result[0].Name)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

// TestGetEntities_PreparerFailure tests the scenario where preparing the query
// fails.
func TestGetEntities_PreparerFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		SQLUtil:   mockSQLUtil,
	}

	opts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{
				{Field: "active", Predicate: "=", Value: true},
			},
		},
	}

	expectedErr := errors.New("prepare error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)

	result, err := entityHelpers.GetEntities(mockPreparer, opts)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	mockPreparer.AssertExpectations(t)
}

// TestGetEntitiesWithManagedTransaction_SuccessfulTransaction tests the
// scenario where multiple entities are successfully retrieved.
func TestGetEntitiesWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		ScanRowsFn: func(rows util.Rows, entity *TestEntity) error {
			entity.ID = 1
			entity.Name = "Alice"
			return nil
		},
	}

	opts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{
				{Field: "active", Predicate: "=", Value: true},
			},
		},
	}

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	mockTx.On("Commit").Return(nil)

	result, err := entityHelpers.GetEntitiesWithManagedTransaction(ctx, opts)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "Alice", result[0].Name)

	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

// TestGetEntitiesWithManagedTransaction_TransactionFailure tests the scenario
// where the transaction fails during the get process.
func TestGetEntitiesWithManagedTransaction_TransactionFailure(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
	}

	opts := GetOptions{
		Options: Options{
			Selectors: []util.Selector{
				{Field: "active", Predicate: "=", Value: true},
			},
		},
	}

	expectedErr := errors.New("execution error")

	// Mock transaction preparation, execution, and error
	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(nil, expectedErr)
	mockStmt.On("Close").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Run the managed transaction get
	result, err := entityHelpers.GetEntitiesWithManagedTransaction(ctx, opts)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)

	// Assert expectations
	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestGetEntityCount_SuccessfulCountWithJoins tests the scenario where the
// entity count is successfully retrieved with join clauses.
func TestGetEntityCount_SuccessfulCountWithJoins(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}

	joins := []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "related_table",
			OnLeft: util.ColumSelector{
				Table:  "test_table",
				Column: "related_id",
			},
			OnRight: util.ColumSelector{
				Table:  "related_table",
				Column: "id",
			},
		},
	}

	expectedCount := 42
	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		// Get the pointer to the int and set it to the expected value
		arguments := args.Get(0).([]any)
		countPtr, ok := arguments[0].(*int)
		if !ok {
			panic("expected *int argument")
		}
		*countPtr = expectedCount
	}).Return(nil)

	result, err := entityHelpers.GetEntityCount(mockPreparer, selectors, joins)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestGetEntityCount_PreparerFailure tests the scenario where preparing the
// query fails.
func TestGetEntityCount_PreparerFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}

	joins := []util.Join{
		{
			Type:  util.JoinTypeLeft,
			Table: "related_table",
			OnLeft: util.ColumSelector{
				Table:  "test_table",
				Column: "related_id",
			},
			OnRight: util.ColumSelector{
				Table:  "related_table",
				Column: "id",
			},
		},
	}

	expectedErr := errors.New("prepare error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)

	result, err := entityHelpers.GetEntityCount(mockPreparer, selectors, joins)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, 0, result)

	mockPreparer.AssertExpectations(t)
}

// TestGetEntityCountWithManagedTransaction_SuccessfulTransaction tests the
// scenario where the entity count is success.
func TestGetEntityCountWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}

	joins := []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "related_table",
			OnLeft: util.ColumSelector{
				Table:  "test_table",
				Column: "related_id",
			},
			OnRight: util.ColumSelector{
				Table:  "related_table",
				Column: "id",
			},
		},
	}

	expectedCount := 42

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockStmt.On("Close").Return(nil)

	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		// Get the pointer to the int and set it to the expected value
		arguments := args.Get(0).([]any)
		countPtr, ok := arguments[0].(*int)
		if !ok {
			panic("expected *int argument")
		}
		*countPtr = expectedCount
	}).Return(nil)

	mockTx.On("Commit").Return(nil)

	result, err := entityHelpers.GetEntityCountWithManagedTransaction(
		ctx,
		selectors,
		joins,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)

	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestGetEntityCountWithManagedTransaction_PreparerFailure(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   func(ctx context.Context) (util.Tx, error) { return mockTx, nil },
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}

	joins := []util.Join{
		{
			Type:  util.JoinTypeInner,
			Table: "related_table",
			OnLeft: util.ColumSelector{
				Table:  "test_table",
				Column: "related_id",
			},
			OnRight: util.ColumSelector{
				Table:  "related_table",
				Column: "id",
			},
		},
	}

	expectedErr := errors.New("prepare error")
	mockTx.On("Prepare", mock.Anything).Return(nil, expectedErr)
	mockTx.On("Rollback").Return(nil)

	result, err := entityHelpers.GetEntityCountWithManagedTransaction(
		ctx,
		selectors,
		joins,
	)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, 0, result)

	mockTx.AssertExpectations(t)
}

// TestUpdateEntities_SuccessfulUpdate tests the scenario where entities are
// successfully updated.
func TestUpdateEntities_SuccessfulUpdate(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}
	updates := Updates{
		{Field: "name", Value: "Updated Name"},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("RowsAffected").Return(int64(1), nil)

	result, err :=
		entityHelpers.UpdateEntities(mockPreparer, selectors, updates)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpdateEntities_PreparerFailure tests the scenario where preparing the
// query fails.
func TestUpdateEntities_PreparerFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}
	updates := Updates{
		{Field: "name", Value: "Updated Name"},
	}

	expectedErr := errors.New("prepare error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)

	result, err := entityHelpers.UpdateEntities(mockPreparer, selectors, updates)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, int64(0), result)

	mockPreparer.AssertExpectations(t)
}

// TestUpdateEntitiesWithManagedTransaction_SuccessfulTransaction tests the
// scenario where entities are successfully updated.
func TestUpdateEntitiesWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}
	updates := Updates{
		{Field: "name", Value: "Updated Name"},
	}

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(1), nil)

	mockTx.On("Commit").Return(nil)

	result, err := entityHelpers.UpdateEntitiesWithManagedTransaction(
		ctx,
		selectors,
		updates,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)

	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpdateEntitiesWithManagedTransaction_TransactionFailure tests the
// scenario where the transaction fails during the update process.
func TestUpdateEntitiesWithManagedTransaction_TransactionFailure(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: true},
	}
	updates := Updates{
		{Field: "name", Value: "Updated Name"},
	}

	expectedErr := errors.New("execution error")

	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, expectedErr)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(expectedErr)
	mockTx.On("Rollback").Return(nil)

	result, err := entityHelpers.UpdateEntitiesWithManagedTransaction(
		ctx,
		selectors,
		updates,
	)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, int64(0), result)

	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

func TestUpdateEntities_WithUpdateHandler_ExistingUpdateOption(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	updateHandler := &UpdateHandler{
		UpdatedField: "last_updated",
		GetUpdateOptionsFn: func() Update {
			return Update{Field: "last_updated", Value: "2024-01-01"}
		},
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName:     "test_table",
		UpdateHandler: updateHandler,
		SQLUtil:       mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}

	updates := Updates{
		{Field: "last_updated", Value: "2024-01-01"}, // Explicitly set field
		{Field: "name", Value: "Alice"},
	}

	expectedRowsAffected := int64(1)
	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("RowsAffected").Return(expectedRowsAffected, nil)

	result, err := entityHelpers.UpdateEntities(
		mockPreparer,
		selectors,
		updates,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedRowsAffected, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestUpdateEntities_WithUpdateHandler_NoExistingUpdateOption(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	updateHandler := &UpdateHandler{
		UpdatedField: "last_updated",
		GetUpdateOptionsFn: func() Update {
			return Update{Field: "last_updated", Value: "2024-01-01"}
		},
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName:     "test_table",
		UpdateHandler: updateHandler,
		SQLUtil:       mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}

	updates := Updates{
		{Field: "name", Value: "Alice"},
	}

	expectedRowsAffected := int64(1)
	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockSQLUtil.On("CheckDBError", mock.Anything).Return(nil)
	mockResult.On("RowsAffected").Return(expectedRowsAffected, nil)

	result, err := entityHelpers.UpdateEntities(
		mockPreparer,
		selectors,
		updates,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedRowsAffected, result)

	mockPreparer.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDeleteEntities_SuccessfulDeletion tests the scenario where entities are
// successfully deleted.
func TestDeleteEntities_SuccessfulDeletion(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: false},
	}
	opts := &DeleteOptions{
		Limit:  10,
		Orders: []util.Order{{Field: "created_at", Direction: util.OrderAsc}},
	}

	mockPreparer.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(5), nil)

	result, err := entityHelpers.DeleteEntities(mockPreparer, selectors, opts)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), result)

	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDeleteEntities_PreparerFailure tests the scenario where preparing the
// query fails.
func TestDeleteEntities_PreparerFailure(t *testing.T) {
	mockPreparer := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: false},
	}
	opts := &DeleteOptions{
		Limit:  10,
		Orders: []util.Order{{Field: "created_at", Direction: util.OrderAsc}},
	}

	expectedErr := errors.New("prepare error")
	mockPreparer.On("Prepare", mock.Anything).Return(nil, expectedErr)

	// Run the delete function
	result, err := entityHelpers.DeleteEntities(mockPreparer, selectors, opts)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, int64(0), result)

	// Assert expectations
	mockPreparer.AssertExpectations(t)
}

// TestDeleteEntitiesWithManagedTransaction_SuccessfulTransaction tests the
// where entities are successfully deleted.
func TestDeleteEntitiesWithManagedTransaction_SuccessfulTransaction(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
		SQLUtil:   mockSQLUtil,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: false},
	}
	opts := &DeleteOptions{
		Limit:  10,
		Orders: []util.Order{{Field: "created_at", Direction: util.OrderDesc}},
	}

	// Mock the transaction preparation, execution, and the result
	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(5), nil)

	// Mock the commit of the transaction
	mockTx.On("Commit").Return(nil)

	// Run the managed transaction delete
	result, err := entityHelpers.DeleteEntitiesWithManagedTransaction(
		ctx,
		selectors,
		opts,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), result)

	// Assert expectations
	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDeleteEntitiesWithManagedTransaction_TransactionFailure tests the
// scenario where the transaction fails.
func TestDeleteEntitiesWithManagedTransaction_TransactionFailure(t *testing.T) {
	ctx := endpointutil.NewContext(context.Background())

	mockTx := new(utilmock.MockTx)
	mockStmt := new(utilmock.MockStmt)

	getTxFn := func(ctx context.Context) (util.Tx, error) {
		return mockTx, nil
	}

	entityHelpers := EntityHelpers[TestEntity]{
		TableName: "test_table",
		GetTxFn:   getTxFn,
	}

	selectors := []util.Selector{
		{Field: "active", Predicate: "=", Value: false},
	}
	opts := &DeleteOptions{
		Limit:  10,
		Orders: []util.Order{{Field: "created_at", Direction: util.OrderDesc}},
	}

	expectedErr := errors.New("execution error")

	// Mock the transaction preparation, execution, and the error
	mockTx.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, expectedErr)
	mockStmt.On("Close").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Run the managed transaction delete
	result, err := entityHelpers.DeleteEntitiesWithManagedTransaction(
		ctx,
		selectors,
		opts,
	)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, int64(0), result)

	// Assert expectations
	mockTx.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
