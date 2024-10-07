package runner

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockParseableGetInput is a mock implementation of
// ParseableInput[ParsedGetEndpointInput]
type MockParseableGetInput struct {
	mock.Mock
}

func (m *MockParseableGetInput) Validate() []inputlogic.FieldError {
	return nil
}

func (m *MockParseableGetInput) Parse() (*ParsedGetEndpointInput, error) {
	args := m.Called()
	return args.Get(0).(*ParsedGetEndpointInput), args.Error(1)
}

// TestGetInvoke_ValidInput tests the GetInvoke function with a valid input.
func TestGetInvoke_ValidInput(t *testing.T) {
	mockInput := new(MockParseableGetInput)
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              &page.Page{Offset: 0, Limit: 10},
		GetCount:          false,
	}
	mockInput.On("Parse").Return(parsedInput, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	output, err := GetInvoke(
		rr,
		req,
		mockInput,
		MockGetServiceFunc,
		MockGetCountFunc,
		MockToGetEndpointOutput,
	)

	assert.NoError(t, err, "GetInvoke should not return an error for valid input")
	assert.NotNil(t, output, "GetInvoke should return a valid output")
	assert.Equal(t, "output", *output, "GetInvoke should return the expected output")

	mockInput.AssertExpectations(t)
}

// TestGetInvoke_ParseError tests the GetInvoke function with a parse error.
func TestGetInvoke_ParseError(t *testing.T) {
	mockInput := new(MockParseableGetInput)
	mockInput.On("Parse").
		Return((*ParsedGetEndpointInput)(nil), errors.New("parse error"))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	output, err := GetInvoke(
		rr,
		req,
		mockInput,
		MockGetServiceFunc,
		MockGetCountFunc,
		MockToGetEndpointOutput,
	)

	assert.Error(t, err, "GetInvoke should return an error if parsing fails")
	assert.Nil(t, output, "GetInvoke should return nil output if parsing fails")

	mockInput.AssertExpectations(t)
}

// TestGetInvoke_GetServiceError tests the GetInvoke function with a service
// error.
func TestGetInvoke_GetServiceError(t *testing.T) {
	mockInput := new(MockParseableGetInput)
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              &page.Page{Offset: 0, Limit: 10},
		GetCount:          false,
	}
	mockInput.On("Parse").Return(parsedInput, nil)

	mockGetServiceFn := func(
		ctx context.Context,
		opts entity.GetOptions,
	) ([]string, error) {
		return nil, errors.New("service error")
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	output, err := GetInvoke(
		rr,
		req,
		mockInput,
		mockGetServiceFn,
		MockGetCountFunc,
		MockToGetEndpointOutput,
	)

	assert.Error(t, err, "GetInvoke should return an error if the service function fails")
	assert.Nil(t, output, "GetInvoke should return nil output if the service function fails")

	mockInput.AssertExpectations(t)
}

// TestGetInvoke_GetCount tests the GetInvoke function with GetCount true.
func TestGetInvoke_GetCount(t *testing.T) {
	mockInput := new(MockParseableGetInput)
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              &page.Page{Offset: 0, Limit: 10},
		GetCount:          true,
	}
	mockInput.On("Parse").Return(parsedInput, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	output, err := GetInvoke(
		rr,
		req,
		mockInput,
		MockGetServiceFunc,
		MockGetCountFunc,
		MockToGetEndpointOutput,
	)

	assert.NoError(t, err, "GetInvoke should not return an error for valid input with GetCount true")
	assert.NotNil(t, output, "GetInvoke should return a valid output when GetCount is true")
	assert.Equal(t, "output", *output, "GetInvoke should return the expected output when GetCount is true")

	mockInput.AssertExpectations(t)
}

func (m MockParseableUpdateInputInvoke) Validate() []inputlogic.FieldError {
	if m.helper == nil {
		panic("MockHelper is not initialized.")
	}
	args := m.helper.Called()
	return args.Get(0).([]inputlogic.FieldError)
}

func (m MockParseableUpdateInputInvoke) Parse() (
	*ParsedUpdateEndpointInput,
	error,
) {
	if m.helper == nil {
		panic("MockHelper is not initialized.")
	}
	args := m.helper.Called()
	return args.Get(0).(*ParsedUpdateEndpointInput), args.Error(1)
}

// TestUpdateInvoke_Success tests UpdateInvoke with valid input and successful
// update.
func TestUpdateInvoke_Success(t *testing.T) {
	mockHelper := new(MockHelper)
	mockInput := MockParseableUpdateInputInvoke{helper: mockHelper}

	parsedInput := &ParsedUpdateEndpointInput{
		DatabaseSelectors: util.Selectors{
			{
				Table:     "table1",
				Field:     "column1",
				Predicate: util.EQUAL,
				Value:     "value1",
			},
		},
		DatabaseUpdates: []entity.Update{
			{Field: "column1", Value: "newValue"},
		},
	}
	mockInput.helper.On("Parse").Return(parsedInput, nil)

	req := httptest.NewRequest(http.MethodPut, "/update", nil)
	rr := httptest.NewRecorder()

	mockUpdateServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		databaseUpdates []entity.Update,
	) (int64, error) {
		return 1, nil
	}

	output, err := UpdateInvoke[MockParseableUpdateInput](
		rr,
		req,
		mockInput,
		mockUpdateServiceFunc,
		MockToUpdateEndpointOutput,
	)

	assert.NoError(t, err, "UpdateInvoke should not return an error for valid input")
	assert.NotNil(t, output, "UpdateInvoke should return a valid output")
	assert.Equal(t, "update successful", *output, "UpdateInvoke should return the correct output")

	mockInput.helper.AssertExpectations(t)
}

// TestUpdateInvoke_ParseError tests UpdateInvoke when parsing fails.
func TestUpdateInvoke_ParseError(t *testing.T) {
	mockHelper := new(MockHelper)
	mockInput := MockParseableUpdateInputInvoke{helper: mockHelper}

	mockInput.helper.On("Parse").
		Return((*ParsedUpdateEndpointInput)(nil), errors.New("parse error"))

	req := httptest.NewRequest(http.MethodPut, "/update", nil)
	rr := httptest.NewRecorder()

	mockUpdateServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		databaseUpdates []entity.Update,
	) (int64, error) {
		return 1, nil
	}

	output, err := UpdateInvoke[MockParseableUpdateInput](
		rr,
		req,
		mockInput,
		mockUpdateServiceFunc,
		MockToUpdateEndpointOutput,
	)

	assert.Error(t, err, "UpdateInvoke should return an error if parsing fails")
	assert.Nil(t, output, "UpdateInvoke should return nil output if parsing fails")

	mockInput.helper.AssertExpectations(t)
}

// TestUpdateInvoke_ServiceError tests UpdateInvoke when the update service
// function returns an error.
func TestUpdateInvoke_ServiceError(t *testing.T) {
	mockHelper := new(MockHelper)
	mockInput := MockParseableUpdateInputInvoke{helper: mockHelper}

	parsedInput := &ParsedUpdateEndpointInput{
		DatabaseSelectors: util.Selectors{},
		DatabaseUpdates: []entity.Update{
			{Field: "column1", Value: "newValue"},
		},
	}
	mockInput.helper.On("Parse").Return(parsedInput, nil)

	req := httptest.NewRequest(http.MethodPut, "/update", nil)
	rr := httptest.NewRecorder()

	mockUpdateServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		databaseUpdates []entity.Update,
	) (int64, error) {
		return 0, errors.New("service error")
	}

	output, err := UpdateInvoke[MockParseableUpdateInput](
		rr,
		req,
		mockInput,
		mockUpdateServiceFunc,
		MockToUpdateEndpointOutput,
	)

	assert.Error(t, err, "UpdateInvoke should return an error if the service function fails")
	assert.Nil(t, output, "UpdateInvoke should return nil output if the service function fails")

	mockInput.helper.AssertExpectations(t)
}

// TODO: Get working without helper
// MockHelper helps to mock a value receiver object.
type MockHelper struct {
	mock.Mock
}

// MockParseableUpdateInputInvoke is a mock implementation of the
// ParseableInput[ParsedUpdateEndpointInput] interface.
// interface.
type MockParseableUpdateInputInvoke struct {
	helper *MockHelper
}

// MockParseableDeleteInputInvoke is a mock implementation of
// ParseableInput[ParsedDeleteEndpointInput].
type MockParseableDeleteInputInvoke struct {
	helper *MockHelper
}

func (m MockParseableDeleteInputInvoke) Validate() []inputlogic.FieldError {
	if m.helper == nil {
		panic("MockHelper is not initialized.")
	}
	args := m.helper.Called()
	return args.Get(0).([]inputlogic.FieldError)
}

func (m MockParseableDeleteInputInvoke) Parse() (
	*ParsedDeleteEndpointInput,
	error,
) {
	if m.helper == nil {
		panic("MockHelper is not initialized.")
	}
	args := m.helper.Called()
	return args.Get(0).(*ParsedDeleteEndpointInput), args.Error(1)
}

// TestDeleteInvoke_Success tests DeleteInvoke with valid input and successful
// delete operation.
func TestDeleteInvoke_Success(t *testing.T) {
	mockHelper := new(MockHelper)
	mockInput := MockParseableDeleteInputInvoke{helper: mockHelper}

	parsedInput := &ParsedDeleteEndpointInput{
		DatabaseSelectors: util.Selectors{
			{
				Table:     "table1",
				Field:     "column1",
				Predicate: util.EQUAL,
				Value:     "value1",
			},
		},
		DeleteOpts: &entity.DeleteOptions{
			Limit: 1,
		},
	}
	mockInput.helper.On("Parse").Return(parsedInput, nil)

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	rr := httptest.NewRecorder()

	mockDeleteServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		opts *entity.DeleteOptions,
	) (int64, error) {
		return 1, nil
	}

	output, err := DeleteInvoke[MockParseableDeleteInputInvoke, string](
		rr,
		req,
		mockInput,
		mockDeleteServiceFunc,
		MockToDeleteEndpointOutput,
	)

	assert.NoError(t, err, "DeleteInvoke should not return an error for valid input")
	assert.NotNil(t, output, "DeleteInvoke should return a valid output")
	assert.Equal(t, "mock delete output", *output, "DeleteInvoke should return the correct output")

	mockInput.helper.AssertExpectations(t)
}

// TestDeleteInvoke_ParseError tests DeleteInvoke when parsing fails.
func TestDeleteInvoke_ParseError(t *testing.T) {
	mockHelper := new(MockHelper)
	mockInput := MockParseableDeleteInputInvoke{helper: mockHelper}

	mockInput.helper.On("Parse").
		Return((*ParsedDeleteEndpointInput)(nil), errors.New("parse error"))

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	rr := httptest.NewRecorder()

	mockDeleteServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		opts *entity.DeleteOptions,
	) (int64, error) {
		return 1, nil
	}

	output, err := DeleteInvoke[MockParseableDeleteInputInvoke, string](
		rr,
		req,
		mockInput,
		mockDeleteServiceFunc,
		MockToDeleteEndpointOutput,
	)

	assert.Error(t, err, "DeleteInvoke should return an error if parsing fails")
	assert.Nil(t, output, "DeleteInvoke should return nil output if parsing fails")

	mockInput.helper.AssertExpectations(t)
}

// TestDeleteInvoke_ServiceError tests DeleteInvoke when the delete service
// function returns an error.
func TestDeleteInvoke_ServiceError(t *testing.T) {
	mockHelper := new(MockHelper)
	mockInput := MockParseableDeleteInputInvoke{helper: mockHelper}

	parsedInput := &ParsedDeleteEndpointInput{
		DatabaseSelectors: util.Selectors{
			{
				Table:     "table1",
				Field:     "column1",
				Predicate: util.EQUAL,
				Value:     "value1",
			},
		},
		DeleteOpts: &entity.DeleteOptions{
			Limit: 1,
		},
	}
	mockInput.helper.On("Parse").Return(parsedInput, nil)

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	rr := httptest.NewRecorder()

	mockDeleteServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		opts *entity.DeleteOptions,
	) (int64, error) {
		return 0, errors.New("service error")
	}

	output, err := DeleteInvoke[MockParseableDeleteInputInvoke, string](
		rr,
		req,
		mockInput,
		mockDeleteServiceFunc,
		MockToDeleteEndpointOutput,
	)

	assert.Error(t, err, "DeleteInvoke should return an error if the service function fails")
	assert.Nil(t, output, "DeleteInvoke should return nil output if the service function fails")

	mockInput.helper.AssertExpectations(t)
}

// TestRunGetService_GetCount_Success tests runGetService when GetCount is true
// and the count function succeeds.
func TestRunGetService_GetCount_Success(t *testing.T) {
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              nil,
		GetCount:          true,
	}

	mockGetCountFunc := func(
		ctx context.Context,
		selectors []util.Selector,
		joins []util.Join,
	) (int, error) {
		return 42, nil
	}

	output, count, err := runGetService[any](
		context.Background(),
		parsedInput,
		nil,
		mockGetCountFunc,
		nil,
		nil,
	)

	assert.NoError(t, err, "runGetService should not return an error for valid GetCount operation")
	assert.Nil(t, output, "Output should be nil when counting entities")
	assert.Equal(t, 42, count, "Count should match the expected value")
}

// TestRunGetService_GetCount_Error tests runGetService when GetCount is true
// but the count function returns an error.
func TestRunGetService_GetCount_Error(t *testing.T) {
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              nil,
		GetCount:          true,
	}

	mockGetCountFunc := func(ctx context.Context, selectors []util.Selector, joins []util.Join) (int, error) {
		return 0, errors.New("count function error")
	}

	output, count, err := runGetService[any](
		context.Background(),
		parsedInput,
		nil,
		mockGetCountFunc,
		nil,
		nil,
	)

	assert.Error(t, err, "runGetService should return an error if the count function fails")
	assert.Nil(t, output, "Output should be nil if counting fails")
	assert.Equal(t, 0, count, "Count should be zero if counting fails")
}

// TestRunGetService_FetchEntities_Success tests runGetService when GetCount is
// false and the service function succeeds.
func TestRunGetService_FetchEntities_Success(t *testing.T) {
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              &page.Page{Offset: 0, Limit: 10},
		GetCount:          false,
	}

	mockGetServiceFunc := func(ctx context.Context, opts entity.GetOptions) ([]string, error) {
		return []string{"entity1", "entity2"}, nil
	}

	output, count, err := runGetService(context.Background(), parsedInput, mockGetServiceFunc, nil, nil, nil)

	assert.NoError(t, err, "runGetService should not return an error for valid service function operation")
	assert.NotNil(t, output, "Output should not be nil when fetching entities")
	assert.Equal(t, []string{"entity1", "entity2"}, output, "Output should match the expected entities")
	assert.Equal(t, 2, count, "Count should match the number of entities fetched")
}

// TestRunGetService_FetchEntities_Error tests runGetService when GetCount is
// false but the service function returns an error.
func TestRunGetService_FetchEntities_Error(t *testing.T) {
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              &page.Page{Offset: 0, Limit: 10},
		GetCount:          false,
	}

	mockGetServiceFunc := func(ctx context.Context, opts entity.GetOptions) ([]string, error) {
		return nil, errors.New("service function error")
	}

	output, count, err := runGetService(
		context.Background(),
		parsedInput,
		mockGetServiceFunc,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err, "runGetService should return an error if the service function fails")
	assert.Nil(t, output, "Output should be nil if the service function fails")
	assert.Equal(t, 0, count, "Count should be zero if fetching entities fails")
}

// TestRunGetService_GetCountFuncNil tests runGetService when GetCount is true
// but GetCountFunc is nil.
func TestRunGetService_GetCountFuncNil(t *testing.T) {
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              nil,
		GetCount:          true,
	}

	output, count, err := runGetService[any](
		context.Background(),
		parsedInput,
		nil,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err, "runGetService should return an error if GetCountFunc is nil")
	assert.Nil(t, output, "Output should be nil if GetCountFunc is nil")
	assert.Equal(t, 0, count, "Count should be zero if GetCountFunc is nil")
}

// TestRunGetService_GetServiceFuncNil tests runGetService when GetCount is
// false but GetServiceFunc is nil.
func TestRunGetService_GetServiceFuncNil(t *testing.T) {
	parsedInput := &ParsedGetEndpointInput{
		Orders:            []util.Order{},
		DatabaseSelectors: util.Selectors{},
		Page:              &page.Page{Offset: 0, Limit: 10},
		GetCount:          false,
	}

	output, count, err := runGetService[any](
		context.Background(),
		parsedInput,
		nil,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err, "runGetService should return an error if GetServiceFunc is nil")
	assert.Nil(t, output, "Output should be nil if GetServiceFunc is nil")
	assert.Equal(t, 0, count, "Count should be zero if GetServiceFunc is nil")
}
