package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/core/client"
	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
	"github.com/pakkasys/fluidapi/endpoint/definition"
	"github.com/pakkasys/fluidapi/endpoint/middleware"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/predicate"
	"github.com/pakkasys/fluidapi/endpoint/selector"
	"github.com/pakkasys/fluidapi/endpoint/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStackBuilder is a mock implementation of the StackBuilder interface.
type MockStackBuilder struct {
	mock.Mock
	middlewares []api.MiddlewareWrapper
}

// MustAddMiddleware adds the provided middleware wrappers to the stack.
func (m *MockStackBuilder) MustAddMiddleware(
	wrappers ...api.MiddlewareWrapper,
) StackBuilder {
	m.Called(wrappers)
	m.middlewares = append(m.middlewares, wrappers...)
	return m
}

func (m *MockStackBuilder) Build() middleware.Stack {
	m.Called()
	return middleware.Stack(m.middlewares)
}

// MockValidatedInput is a mock implementation of the ValidatedInput interface.
type MockValidatedInput struct{}

func (m MockValidatedInput) Validate() []inputlogic.FieldError {
	return nil
}

// MockParseableInput is a mock implementation of the ParseableInput interface.
type MockParseableInput struct{}

func (m MockParseableInput) Validate() []inputlogic.FieldError {
	return nil
}

func (m MockParseableInput) Parse() (*ParsedGetEndpointInput, error) {
	return &ParsedGetEndpointInput{}, nil
}

// MockGetServiceFunc represents a mock implementation of GetServiceFunc.
func MockGetServiceFunc(
	ctx context.Context,
	opts entity.GetOptions,
) ([]string, error) {
	return []string{"entity1", "entity2"}, nil
}

// MockGetCountFunc represents a mock implementation of GetCountFunc.
func MockGetCountFunc(
	ctx context.Context,
	selectors []util.Selector,
	joins []util.Join,
) (int, error) {
	return 2, nil
}

// MockToGetEndpointOutput is a mock implementation of ToGetEndpointOutput.
func MockToGetEndpointOutput(froms []string, count *int) *string {
	output := "output"
	return &output
}

// MockObjectPicker is a mock implementation of the IObjectPicker interface.
type MockObjectPicker[Input any] struct {
	mock.Mock
}

func (m *MockObjectPicker[Input]) PickObject(
	r *http.Request,
	w http.ResponseWriter,
	obj Input,
) (*Input, error) {
	args := m.Called(r, w, obj)
	return args.Get(0).(*Input), args.Error(1)
}

// MockOutputHandler is a mock implementation of the IOutputHandler interface.
type MockOutputHandler struct {
	mock.Mock
}

func (m *MockOutputHandler) ProcessOutput(
	w http.ResponseWriter,
	r *http.Request,
	out any,
	outError error,
	statusCode int,
) error {
	args := m.Called(w, r, out, outError, statusCode)
	return args.Error(0)
}

// MockToUpdateEndpointOutput is a mock implementation of ToUpdateEndpointOutput
func MockToUpdateEndpointOutput(count int64) *string {
	output := "update successful"
	return &output
}

// MockParseableUpdateInput is a mock implementation of the ParseableInput
// interface.
type MockParseableUpdateInput struct{}

func (m MockParseableUpdateInput) Validate() []inputlogic.FieldError {
	return nil
}

func (m MockParseableUpdateInput) Parse() (*ParsedUpdateEndpointInput, error) {
	return &ParsedUpdateEndpointInput{}, nil
}

// MockParseableDeleteInput is a mock implementation of the ParseableInput
// interface
type MockParseableDeleteInput struct{}

func (m MockParseableDeleteInput) Validate() []inputlogic.FieldError {
	return nil
}

func (m MockParseableDeleteInput) Parse() (*ParsedDeleteEndpointInput, error) {
	return &ParsedDeleteEndpointInput{}, nil
}

// MockDeleteServiceFunc represents a mock implementation of DeleteServiceFunc
func MockDeleteServiceFunc(
	ctx context.Context,
	databaseSelectors []util.Selector,
	opts *entity.DeleteOptions,
) (int64, error) {
	return 1, nil
}

// MockToDeleteEndpointOutput is a mock implementation of ToDeleteEndpointOutput
func MockToDeleteEndpointOutput(count int64) *string {
	output := "mock delete output"
	return &output
}

// MockMiddlewareWrapper creates a mock middleware wrapper.
func MockMiddlewareWrapper[I ValidatedInput, O any](
	callback inputlogic.Callback[I, O],
	inputFactory func() *I,
	expectedErrors []inputlogic.ExpectedError,
	opts inputlogic.Options[I],
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID: inputlogic.MiddlewareID,
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := callback(w, r, inputFactory())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("Error"))
					return
				}
				_, _ = w.Write([]byte("Success"))
			})
		},
		Inputs: []any{*inputFactory()},
	}
}

// TestWithMiddlewareWrapper tests the WithMiddlewareWrapper function.
func TestWithMiddlewareWrapper(t *testing.T) {
	mockMiddlewareWrapper := &api.MiddlewareWrapper{
		ID: "mockMiddleware",
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("Mock Middleware Executed"))
					next.ServeHTTP(w, r)
				},
			)
		},
	}

	endpoint := &Endpoint[MockValidatedInput, string, any]{
		Definition: &definition.EndpointDefinition{
			URL:             "/test-endpoint",
			Method:          "GET",
			MiddlewareStack: middleware.Stack{},
		},
		Client: nil,
	}

	middlewareOption := WithMiddlewareWrapper[MockValidatedInput, string, any](
		mockMiddlewareWrapper,
	)
	middlewareOption(endpoint)

	assert.Equal(t, 1, len(endpoint.Definition.MiddlewareStack), "Middleware stack should contain exactly one middleware wrapper")
	assert.Equal(t, "mockMiddleware", endpoint.Definition.MiddlewareStack[0].ID, "Middleware ID should match the mock middleware wrapper")
	assert.NotNil(t, endpoint.Definition.MiddlewareStack[0].Middleware, "Middleware function should not be nil")

	handler := api.ApplyMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		endpoint.Definition.MiddlewareStack.Middlewares()...,
	)

	req, err := http.NewRequest("GET", "/test-endpoint", nil)
	assert.NoError(t, err, "Creating the request should not fail")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "Mock Middleware Executed", rr.Body.String(), "Expected response from the mock middleware")
}

// TestGenericEndpointDefinition tests the GenericEndpointDefinition function.
func TestGenericEndpointDefinition(t *testing.T) {
	specification := InputSpecification[MockValidatedInput]{
		URL:    "/test",
		Method: http.MethodGet,
		InputFactory: func() *MockValidatedInput {
			return new(MockValidatedInput)
		},
	}

	callback := func(
		http.ResponseWriter,
		*http.Request,
		*MockValidatedInput,
	) (*any, error) {
		return nil, nil
	}

	expectedErrors := []inputlogic.ExpectedError{}
	stackBuilder := new(MockStackBuilder)
	stackBuilder.On("MustAddMiddleware", mock.Anything).Return(stackBuilder)
	stackBuilder.On("Build").Return(middleware.Stack{})

	sendFn := func(
		input *MockValidatedInput,
		host string,
	) (*client.Response[MockValidatedInput, any], error) {
		return nil, nil
	}

	mockObjectPicker := &MockObjectPicker[MockValidatedInput]{}

	mockOutputHandler := new(MockOutputHandler)
	mockOutputHandler.On(
		"ProcessOutput",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	opts := inputlogic.Options[MockValidatedInput]{
		ObjectPicker:  mockObjectPicker,
		OutputHandler: mockOutputHandler,
	}

	endpoint := GenericEndpointDefinition(
		specification,
		callback,
		expectedErrors,
		stackBuilder,
		opts,
		sendFn,
	)

	assert.NotNil(t, endpoint, "Endpoint should not be nil")
	assert.Equal(t, "/test", endpoint.Definition.URL, "Endpoint URL should match specification")
	assert.Equal(t, http.MethodGet, endpoint.Definition.Method, "Endpoint Method should match specification")
	assert.NotNil(t, endpoint.Client, "Endpoint client should not be nil")
}

// TestGenericEndpointDefinition_WithMiddlewareWrapper tests the
// WithMiddlewareWrapper function with a middleware wrapper.
func TestGenericEndpointDefinition_WithMiddlewareWrapper(t *testing.T) {
	specification := InputSpecification[MockValidatedInput]{
		URL:    "/middleware-test",
		Method: http.MethodGet,
		InputFactory: func() *MockValidatedInput {
			return new(MockValidatedInput)
		},
	}

	stackBuilder := new(MockStackBuilder)
	stackBuilder.On("MustAddMiddleware", mock.Anything).Return(stackBuilder)
	stackBuilder.On("Build").Return(middleware.Stack{})

	sendFn := func(
		input *MockValidatedInput,
		host string,
	) (*client.Response[MockValidatedInput, any], error) {
		return nil, nil
	}

	mockObjectPicker := new(MockObjectPicker[MockValidatedInput])
	mockOutputHandler := new(MockOutputHandler)

	opts := inputlogic.Options[MockValidatedInput]{
		ObjectPicker:  mockObjectPicker,
		OutputHandler: mockOutputHandler,
	}

	mockMiddlewareWrapper := MockMiddlewareWrapper(
		func(
			w http.ResponseWriter,
			r *http.Request,
			input *MockValidatedInput,
		) (*string, error) {
			output := "mock output"
			return &output, nil
		},
		specification.InputFactory,
		[]inputlogic.ExpectedError{},
		opts,
	)

	endpoint := GenericEndpointDefinition(
		specification,
		func(
			http.ResponseWriter,
			*http.Request,
			*MockValidatedInput,
		) (*string, error) {
			return nil, nil
		},
		[]inputlogic.ExpectedError{},
		stackBuilder,
		opts,
		sendFn,
		WithMiddlewareWrapper[MockValidatedInput, string, any](
			mockMiddlewareWrapper,
		),
	)

	assert.NotNil(t, endpoint, "Endpoint should not be nil")
	assert.Equal(t, "/middleware-test", endpoint.Definition.URL, "Endpoint URL should match specification")
	assert.Equal(t, http.MethodGet, endpoint.Definition.Method, "Endpoint Method should match specification")
	assert.NotNil(t, endpoint.Client, "Endpoint client should not be nil")
	assert.Equal(t, 1, len(endpoint.Definition.MiddlewareStack), "Middleware stack should have exactly one middleware")

	// Simulate a request and verify that the custom middleware wrapper is called
	handler := api.ApplyMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}),
		endpoint.Definition.MiddlewareStack.Middlewares()...,
	)

	req, err := http.NewRequest(http.MethodGet, "/middleware-test", nil)
	assert.NoError(t, err, "Creating the request should not fail")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// The actual middleware logic should be invoked here, so we expect the response from the mock middleware
	assert.Equal(t, "Success", rr.Body.String(), "Expected response from the custom middleware wrapper")
}

// TestGetEndpointDefinition tests the GetEndpointDefinition function.
func TestGetEndpointDefinition(t *testing.T) {
	specification := InputSpecification[MockParseableInput]{
		URL:    "/test",
		Method: http.MethodGet,
		InputFactory: func() *MockParseableInput {
			return new(MockParseableInput)
		},
	}

	stackBuilder := new(MockStackBuilder)
	stackBuilder.On("MustAddMiddleware", mock.Anything).Return(stackBuilder)
	stackBuilder.On("Build").Return(middleware.Stack(stackBuilder.middlewares))

	sendFn := func(
		input *MockParseableInput,
		host string,
	) (*client.Response[MockParseableInput, any], error) {
		return nil, nil
	}

	mockObjectPicker := new(MockObjectPicker[MockParseableInput])
	mockObjectPicker.On("PickObject", mock.Anything, mock.Anything, mock.Anything).
		Return(&MockParseableInput{}, nil)

	mockOutputHandler := new(MockOutputHandler)
	mockOutputHandler.On(
		"ProcessOutput",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	opts := inputlogic.Options[MockParseableInput]{
		ObjectPicker:  mockObjectPicker,
		OutputHandler: mockOutputHandler,
	}

	endpoint := GetEndpointDefinition(
		specification,
		MockGetServiceFunc,
		MockGetCountFunc,
		MockToGetEndpointOutput,
		[]inputlogic.ExpectedError{},
		stackBuilder,
		opts,
		sendFn,
	)

	assert.NotNil(t, endpoint, "Endpoint should not be nil")
	assert.Equal(t, "/test", endpoint.Definition.URL, "Endpoint URL should match specification")
	assert.Equal(t, http.MethodGet, endpoint.Definition.Method, "Endpoint Method should match specification")
	assert.NotNil(t, endpoint.Client, "Endpoint client should not be nil")

	// Simulate a request and verify that the callback is called
	handler := api.ApplyMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}),
		endpoint.Definition.MiddlewareStack.Middlewares()...,
	)

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err, "Creating the request should not fail")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "OK", rr.Body.String(), "Expected response from the callback in the middleware")

	stackBuilder.AssertExpectations(t)
	mockObjectPicker.AssertExpectations(t)
	mockOutputHandler.AssertExpectations(t)
}

// TestUpdateEndpointDefinition tests the UpdateEndpointDefinition function.
func TestUpdateEndpointDefinition(t *testing.T) {
	specification := InputSpecification[MockParseableUpdateInput]{
		URL:    "/update",
		Method: http.MethodPut,
		InputFactory: func() *MockParseableUpdateInput {
			return new(MockParseableUpdateInput)
		},
	}

	stackBuilder := new(MockStackBuilder)
	stackBuilder.On("MustAddMiddleware", mock.Anything).Return(stackBuilder)
	stackBuilder.On("Build").Return(middleware.Stack(stackBuilder.middlewares))

	sendFn := func(
		input *MockParseableUpdateInput,
		host string,
	) (*client.Response[MockParseableUpdateInput, any], error) {
		return nil, nil
	}

	mockObjectPicker := new(MockObjectPicker[MockParseableUpdateInput])
	mockObjectPicker.On(
		"PickObject",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).
		Return(&MockParseableUpdateInput{}, nil)

	mockOutputHandler := new(MockOutputHandler)
	mockOutputHandler.On(
		"ProcessOutput",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	opts := inputlogic.Options[MockParseableUpdateInput]{
		ObjectPicker:  mockObjectPicker,
		OutputHandler: mockOutputHandler,
	}

	mockUpdateServiceFunc := func(
		ctx context.Context,
		databaseSelectors []util.Selector,
		databaseUpdates []entity.Update,
	) (int64, error) {
		return 1, nil
	}

	endpoint := UpdateEndpointDefinition(
		specification,
		mockUpdateServiceFunc,
		MockToUpdateEndpointOutput,
		[]inputlogic.ExpectedError{},
		stackBuilder,
		opts,
		sendFn,
	)

	assert.NotNil(t, endpoint, "Endpoint should not be nil")
	assert.Equal(t, "/update", endpoint.Definition.URL, "Endpoint URL should match specification")
	assert.Equal(t, http.MethodPut, endpoint.Definition.Method, "Endpoint Method should match specification")
	assert.NotNil(t, endpoint.Client, "Endpoint client should not be nil")

	// Simulate a request and verify that the callback is called
	handler := api.ApplyMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}),
		endpoint.Definition.MiddlewareStack.Middlewares()...,
	)

	req, err := http.NewRequest(http.MethodPut, "/update", nil)
	assert.NoError(t, err, "Creating the request should not fail")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "OK", rr.Body.String(), "Expected response from the callback in the middleware")

	stackBuilder.AssertExpectations(t)
	mockObjectPicker.AssertExpectations(t)
	mockOutputHandler.AssertExpectations(t)
}

// TestDeleteEndpointDefinition tests the DeleteEndpointDefinition function.
func TestDeleteEndpointDefinition_WithCallbackExecution(t *testing.T) {
	specification := InputSpecification[MockParseableDeleteInput]{
		URL:    "/delete",
		Method: http.MethodDelete,
		InputFactory: func() *MockParseableDeleteInput {
			return new(MockParseableDeleteInput)
		},
	}

	stackBuilder := new(MockStackBuilder)
	stackBuilder.On("MustAddMiddleware", mock.Anything).Return(stackBuilder)
	stackBuilder.On("Build").Return(middleware.Stack(stackBuilder.middlewares))

	sendFn := func(
		input *MockParseableDeleteInput,
		host string,
	) (*client.Response[MockParseableDeleteInput, any], error) {
		return nil, nil
	}

	mockObjectPicker := new(MockObjectPicker[MockParseableDeleteInput])
	mockObjectPicker.On("PickObject", mock.Anything, mock.Anything, mock.Anything).
		Return(&MockParseableDeleteInput{}, nil)

	mockOutputHandler := new(MockOutputHandler)
	mockOutputHandler.On(
		"ProcessOutput",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	opts := inputlogic.Options[MockParseableDeleteInput]{
		ObjectPicker:  mockObjectPicker,
		OutputHandler: mockOutputHandler,
	}

	endpoint := DeleteEndpointDefinition(
		specification,
		MockDeleteServiceFunc,
		MockToDeleteEndpointOutput,
		[]inputlogic.ExpectedError{},
		stackBuilder,
		opts,
		sendFn,
	)

	assert.NotNil(t, endpoint, "Endpoint should not be nil")
	assert.Equal(t, "/delete", endpoint.Definition.URL, "Endpoint URL should match specification")
	assert.Equal(t, http.MethodDelete, endpoint.Definition.Method, "Endpoint Method should match specification")
	assert.NotNil(t, endpoint.Client, "Endpoint client should not be nil")

	handler := api.ApplyMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}),
		endpoint.Definition.MiddlewareStack.Middlewares()...,
	)

	req, err := http.NewRequest(http.MethodDelete, "/delete", nil)
	assert.NoError(t, err, "Creating the request should not fail")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "OK", rr.Body.String(), "Expected response from the callback in the middleware")

	stackBuilder.AssertExpectations(t)
	mockObjectPicker.AssertExpectations(t)
	mockOutputHandler.AssertExpectations(t)
}

// TestParseGetEndpointInput tests the ParseGetEndpointInput function with
// valid input.
func TestParseGetEndpointInput_ValidInput(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	orders := []order.Order{
		{
			Field:     "field1",
			Direction: order.DIRECTION_ASC,
		},
	}

	allowedOrderFields := []string{"field1"}
	inputPage := &page.Page{Offset: 0, Limit: 10}
	maxPageCount := 20
	getCount := true

	result, err := ParseGetEndpointInput(
		apiFields,
		selectors,
		orders,
		allowedOrderFields,
		inputPage,
		maxPageCount,
		getCount,
	)

	assert.NoError(t, err, "ParseGetEndpointInput should not return an error for valid input")
	assert.NotNil(t, result, "ParsedGetEndpointInput should not be nil")

	expectedDBOrders := []util.Order{
		{
			Table:     "table1",
			Field:     "column1",
			Direction: util.OrderAsc,
		},
	}
	assert.Equal(t, expectedDBOrders, result.Orders, "ParsedGetEndpointInput should have the correct translated orders")

	expectedDBSelectors := util.Selectors{
		{
			Table:     "table1",
			Field:     "column1",
			Predicate: util.EQUAL,
			Value:     "value1",
		},
	}
	assert.Equal(t, expectedDBSelectors, result.DatabaseSelectors, "ParsedGetEndpointInput should have the correct translated selectors")

	assert.Equal(t, inputPage, result.Page, "ParsedGetEndpointInput should have the correct page settings")
	assert.Equal(t, getCount, result.GetCount, "ParsedGetEndpointInput should have the correct GetCount value")
}

// TestParseGetEndpointInput_NilPage tests the ParseGetEndpointInput function
// with a nil Page.
func TestParseGetEndpointInput_NilPage(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{}
	orders := []order.Order{}
	allowedOrderFields := []string{"field1"}
	maxPageCount := 20
	getCount := true

	result, err := ParseGetEndpointInput(apiFields, selectors, orders, allowedOrderFields, nil, maxPageCount, getCount)

	assert.NoError(t, err, "ParseGetEndpointInput should not return an error for valid input with nil page")
	assert.NotNil(t, result, "ParsedGetEndpointInput should not be nil")
	assert.Equal(t, &page.Page{Offset: 0, Limit: maxPageCount}, result.Page, "ParsedGetEndpointInput should use default page settings when nil")
}

// TestParseGetEndpointInput_InvalidOrderField tests the ParseGetEndpointInput
// function with an invalid order field.
func TestParseGetEndpointInput_InvalidOrderField(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{}
	orders := []order.Order{
		{
			Field:     "invalid_field",
			Direction: order.DIRECTION_ASC,
		},
	}
	allowedOrderFields := []string{"field1"}
	maxPageCount := 20
	getCount := false

	result, err := ParseGetEndpointInput(apiFields, selectors, orders, allowedOrderFields, nil, maxPageCount, getCount)

	assert.Error(t, err, "ParseGetEndpointInput should return an error for invalid order field")
	assert.Nil(t, result, "ParsedGetEndpointInput should be nil for invalid order field")
}

// TestParseGetEndpointInput_InvalidSelectors tests the ParseGetEndpointInput
// function with an invalid selector.
func TestParseGetEndpointInput_InvalidSelectors(t *testing.T) {
	// Arrange: Define API fields mapping with a selector that doesn't exist in APIFields
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:     "non_existent_field",
			Predicate: predicate.EQUAL,
			Value:     "value1",
		},
	}
	orders := []order.Order{}
	allowedOrderFields := []string{"field1"}
	maxPageCount := 20
	getCount := false

	// Act: Call ParseGetEndpointInput with invalid selectors
	result, err := ParseGetEndpointInput(apiFields, selectors, orders, allowedOrderFields, nil, maxPageCount, getCount)

	// Assert: Validate that an error is returned
	assert.Error(t, err, "ParseGetEndpointInput should return an error for invalid selectors")
	assert.Nil(t, result, "ParsedGetEndpointInput should be nil for invalid selectors")
}

// TestParseGetEndpointInput_InvalidPage tests the ParseGetEndpointInput
// function with an invalid Page.
func TestParseGetEndpointInput_InvalidPage(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	orders := []order.Order{
		{
			Field:     "field1",
			Direction: order.DIRECTION_ASC,
		},
	}

	allowedOrderFields := []string{"field1"}
	invalidPage := &page.Page{Offset: -1, Limit: 10}
	maxPageCount := 5
	getCount := true

	_, err := ParseGetEndpointInput(
		apiFields,
		selectors,
		orders,
		allowedOrderFields,
		invalidPage,
		maxPageCount,
		getCount,
	)

	assert.Error(t, err, "ParseGetEndpointInput should return an error for invalid page settings")
}

// TestParseUpdateEndpointInput_ValidInput tests ParseUpdateEndpointInput with
// valid input.
func TestParseUpdateEndpointInput_ValidInput(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	updates := []update.Update{
		{
			Field: "field1",
			Value: "new_value",
		},
	}
	upsert := true

	result, err := ParseUpdateEndpointInput(
		apiFields,
		selectors,
		updates,
		upsert,
	)

	assert.NoError(t, err, "ParseUpdateEndpointInput should not return an error for valid input")
	assert.NotNil(t, result, "ParsedUpdateEndpointInput should not be nil")

	expectedSelectors := util.Selectors{
		{
			Table:     "table1",
			Field:     "column1",
			Predicate: util.EQUAL,
			Value:     "value1",
		},
	}
	assert.Equal(t, expectedSelectors, result.DatabaseSelectors, "ParsedUpdateEndpointInput should have the correct selectors")

	expectedUpdates := []entity.Update{
		{
			Field: "column1",
			Value: "new_value",
		},
	}
	assert.Equal(t, expectedUpdates, result.DatabaseUpdates, "ParsedUpdateEndpointInput should have the correct updates")

	assert.Equal(t, upsert, result.Upsert, "ParsedUpdateEndpointInput should have the correct upsert value")
}

// TestParseUpdateEndpointInput_InvalidSelector tests ParseUpdateEndpointInput
// with an invalid selector.
func TestParseUpdateEndpointInput_InvalidSelector(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:     "invalid_field",
			Predicate: predicate.EQUAL,
			Value:     "value1",
		},
	}
	updates := []update.Update{
		{
			Field: "field1",
			Value: "new_value",
		},
	}
	upsert := false

	result, err := ParseUpdateEndpointInput(
		apiFields,
		selectors,
		updates,
		upsert,
	)

	assert.Error(t, err, "ParseUpdateEndpointInput should return an error for an invalid selector")
	assert.Nil(t, result, "ParsedUpdateEndpointInput should be nil for an invalid selector")
}

// TestParseUpdateEndpointInput_NoSelectors tests ParseUpdateEndpointInput with
// no selectors.
func TestParseUpdateEndpointInput_NoSelectors(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{} // No selectors provided
	updates := []update.Update{
		{
			Field: "field1",
			Value: "new_value",
		},
	}
	upsert := false

	result, err := ParseUpdateEndpointInput(
		apiFields,
		selectors,
		updates,
		upsert,
	)

	assert.Error(t, err, "ParseUpdateEndpointInput should return an error when no selectors are provided")
	assert.Nil(t, result, "ParsedUpdateEndpointInput should be nil when no selectors are provided")
	assert.Equal(t, NeedAtLeastOneSelectorError, err, "Expected NeedAtLeastOneSelectorError")
}

// TestParseUpdateEndpointInput_NoUpdates tests ParseUpdateEndpointInput with
// no updates.
func TestParseUpdateEndpointInput_NoUpdates(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	updates := []update.Update{}
	upsert := false

	result, err := ParseUpdateEndpointInput(
		apiFields,
		selectors,
		updates,
		upsert,
	)
	assert.Error(t, err, "ParseUpdateEndpointInput should return an error when no updates are provided")
	assert.Nil(t, result, "ParsedUpdateEndpointInput should be nil when no updates are provided")
	assert.Equal(t, NeedAtLeastOneUpdateError, err, "Expected NeedAtLeastOneUpdateError")
}

// TestParseUpdateEndpointInput_InvalidUpdates tests ParseUpdateEndpointInput
// with an update that is not mappable.
func TestParseUpdateEndpointInput_InvalidUpdates(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	updates := []update.Update{
		{
			Field: "invalid_field", // Field does not exist in apiFields
			Value: "new_value",
		},
	}
	upsert := false

	result, err := ParseUpdateEndpointInput(
		apiFields,
		selectors,
		updates,
		upsert,
	)

	assert.Error(t, err, "ParseUpdateEndpointInput should return an error for an invalid update field")
	assert.Nil(t, result, "ParsedUpdateEndpointInput should be nil for an invalid update field")
}

// TestParseDeleteEndpointInput_ValidInput tests ParseDeleteEndpointInput with
// valid input.
func TestParseDeleteEndpointInput_ValidInput(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	orders := []order.Order{
		{
			Field:     "field1",
			Direction: order.DIRECTION_ASC,
		},
	}
	allowedOrderFields := []string{"field1"}
	limit := 10

	result, err := ParseDeleteEndpointInput(
		apiFields,
		selectors,
		orders,
		allowedOrderFields,
		limit,
	)

	assert.NoError(t, err, "ParseDeleteEndpointInput should not return an error for valid input")
	assert.NotNil(t, result, "ParsedDeleteEndpointInput should not be nil")

	expectedSelectors := util.Selectors{
		{
			Table:     "table1",
			Field:     "column1",
			Predicate: util.EQUAL,
			Value:     "value1",
		},
	}
	assert.Equal(t, expectedSelectors, result.DatabaseSelectors, "ParsedDeleteEndpointInput should have the correct selectors")

	expectedOrders := []util.Order{
		{
			Table:     "table1",
			Field:     "column1",
			Direction: util.OrderAsc,
		},
	}
	assert.Equal(t, expectedOrders, result.DeleteOpts.Orders, "ParsedDeleteEndpointInput should have the correct orders")

	assert.Equal(t, limit, result.DeleteOpts.Limit, "ParsedDeleteEndpointInput should have the correct limit")
}

// TestParseDeleteEndpointInput_InvalidSelector tests ParseDeleteEndpointInput
// with an invalid selector.
func TestParseDeleteEndpointInput_InvalidSelector(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:     "invalid_field",
			Predicate: predicate.EQUAL,
			Value:     "value1",
		},
	}
	orders := []order.Order{}
	allowedOrderFields := []string{"field1"}
	limit := 10

	result, err := ParseDeleteEndpointInput(
		apiFields,
		selectors,
		orders,
		allowedOrderFields,
		limit,
	)

	assert.Error(t, err, "ParseDeleteEndpointInput should return an error for an invalid selector")
	assert.Nil(t, result, "ParsedDeleteEndpointInput should be nil for an invalid selector")
}

// TestParseDeleteEndpointInput_NoSelectors tests ParseDeleteEndpointInput with
// no selectors.
func TestParseDeleteEndpointInput_NoSelectors(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{} // No selectors provided
	orders := []order.Order{}
	allowedOrderFields := []string{"field1"}
	limit := 10

	result, err := ParseDeleteEndpointInput(
		apiFields,
		selectors,
		orders,
		allowedOrderFields,
		limit,
	)

	assert.Error(t, err, "ParseDeleteEndpointInput should return an error when no selectors are provided")
	assert.Nil(t, result, "ParsedDeleteEndpointInput should be nil when no selectors are provided")
	assert.Equal(t, NeedAtLeastOneSelectorError, err, "Expected NeedAtLeastOneSelectorError")
}

// TestParseDeleteEndpointInput_InvalidOrderField tests ParseDeleteEndpointInput
// with an invalid order field.
func TestParseDeleteEndpointInput_InvalidOrderField(t *testing.T) {
	apiFields := APIFields{
		"field1": dbfield.DBField{Table: "table1", Column: "column1"},
	}
	selectors := []selector.Selector{
		{
			Field:             "field1",
			Predicate:         predicate.EQUAL,
			Value:             "value1",
			AllowedPredicates: []predicate.Predicate{predicate.EQUAL},
		},
	}
	orders := []order.Order{
		{
			Field:     "invalid_field", // Field that is not allowed
			Direction: order.DIRECTION_ASC,
		},
	}
	allowedOrderFields := []string{"field1"}
	limit := 10

	result, err := ParseDeleteEndpointInput(
		apiFields,
		selectors,
		orders,
		allowedOrderFields,
		limit,
	)

	assert.Error(t, err, "ParseDeleteEndpointInput should return an error for an invalid order field")
	assert.Nil(t, result, "ParsedDeleteEndpointInput should be nil for an invalid order field")
}
