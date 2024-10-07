package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockInputParser is a mock for the InternalParser interface.
type MockInputParser struct {
	mock.Mock
}

func (m *MockInputParser) ParseInput(
	method string,
	input any,
) (*ParsedInput, error) {
	args := m.Called(method, input)
	return args.Get(0).(*ParsedInput), args.Error(1)
}

// MockSender is a mock for the InternalSender interface.
type MockSender struct {
	mock.Mock
}

func (m *MockSender) ProcessAndSend(
	host string,
	url string,
	method string,
	inputData *RequestData,
) (*SendResult[any], error) {
	args := m.Called(
		host,
		url,
		method,
		inputData,
	)

	// Safely check for nil before type assertion
	var httpResp *http.Response
	if args.Get(0) != nil {
		httpResp = args.Get(0).(*http.Response)
	}

	var outputResp any
	if args.Get(1) != nil {
		outputResp = args.Get(1)
	}

	return &SendResult[any]{
		Response: httpResp,
		Output:   &outputResp,
	}, args.Error(2)
}

// TestSend_Success tests the Send function with a successful response.
func TestSend_Success(t *testing.T) {
	type Input struct{}
	// Initialize the mocks
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define mock behavior for ParseInput
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		&ParsedInput{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Cookies:       []http.Cookie{},
			URLParameters: map[string]any{},
			Body:          map[string]any{"k": "v"},
		},
		nil,
	)

	// Define mock behavior for ProcessAndSend
	mockSender.On(
		"ProcessAndSend",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		&http.Response{StatusCode: http.StatusOK},
		"output",
		nil,
	)

	// Call Send function
	input := Input{}

	mockURLEncoder := &MockURLEncoder{}

	resp, err := Send(
		&input,
		"/test-url",
		"localhost",
		"GET",
		mockURLEncoder,
		HandlerOpts[any]{
			InputParser: mockParser.ParseInput,
			Sender:      mockSender.ProcessAndSend,
		},
	)

	// Assert no errors
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Assert correct HTTP status code
	if resp.Response.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.Response.StatusCode)
	}

	// Verify that mock methods were called as expected
	mockParser.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

// TestSend_ParseInputError tests the Send function with a parse input error.
func TestSend_ParseInputError(t *testing.T) {
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define mock behavior for ParseInput to return an error.
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		(*ParsedInput)(nil),
		errors.New("parse error"),
	)

	mockURLEncoder := &MockURLEncoder{}

	// Call Send function
	input := struct{}{}
	_, err := Send(
		&input,
		"/test-url",
		"localhost",
		"GET",
		mockURLEncoder,
		HandlerOpts[any]{
			InputParser: mockParser.ParseInput,
			Sender:      mockSender.ProcessAndSend,
		},
	)

	// Assert that an error is returned
	if err == nil || err.Error() != "parse error" {
		t.Errorf("expected 'parse error', got %v", err)
	}

	mockParser.AssertExpectations(t)
}

// TestSend_ProcessAndSendError tests the Send function with an error from
// ProcessAndSend call.
func TestSend_ProcessAndSendError(t *testing.T) {
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define mock behavior for ParseInput
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		&ParsedInput{
			Headers:       map[string]string{},
			Cookies:       []http.Cookie{},
			URLParameters: map[string]any{},
			Body:          map[string]any{},
		},
		nil,
	)

	// Define mock behavior for ProcessAndSend to return an error
	mockSender.On(
		"ProcessAndSend",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		(*http.Response)(nil),
		nil,
		errors.New("send error"),
	)

	mockURLEncoder := &MockURLEncoder{}

	// Call Send function
	input := struct{}{}
	_, err := Send(
		&input,
		"/test-url",
		"localhost",
		"GET",
		mockURLEncoder,
		HandlerOpts[any]{
			InputParser: mockParser.ParseInput,
			Sender:      mockSender.ProcessAndSend,
		},
	)

	// Assert that an error is returned
	if err == nil || err.Error() != "send error" {
		t.Errorf("expected 'send error', got %v", err)
	}

	mockParser.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

// TestDetermineSendOpt_CustomOptions tests the determineSendOpt function with
// custom options.
func TestDetermineSendOpt_CustomOptions(t *testing.T) {
	// Create mock objects for custom options
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define custom options
	customOpts := HandlerOpts[any]{
		InputParser: mockParser.ParseInput,
		Sender:      mockSender.ProcessAndSend,
	}

	mockURLEncoder := &MockURLEncoder{}

	// Call determineSendOpt with custom options
	opts := determineSendOpt([]HandlerOpts[any]{customOpts}, mockURLEncoder)

	// Assert that input parsers are the same by calling them both
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		&ParsedInput{},
		nil,
	).Twice()
	if _, err := opts.InputParser("GET", mock.Anything); err != nil {
		t.Errorf("error parsing input: %v", err)
	}
	if _, err := mockParser.ParseInput("GET", mock.Anything); err != nil {
		t.Errorf("error parsing input: %v", err)
	}

	// Assert that senders are the same by calling them both
	mockSender.On(
		"ProcessAndSend",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		(*http.Response)(nil),
		nil,
		nil,
	).Twice()
	_, err := opts.Sender(
		"http://localhost",
		"/test-url",
		"GET",
		&RequestData{},
	)
	if err != nil {
		t.Errorf("error sending: %v", err)
	}
	_, err = mockSender.ProcessAndSend(
		"http://localhost",
		"/test-url",
		"GET",
		&RequestData{},
	)
	if err != nil {
		t.Errorf("error sending: %v", err)
	}
}

// TestDetermineSendOpt_NoOptions tests the determineSendOpt function with no
// provider options.
func TestDetermineSendOpt_NoOptions(t *testing.T) {
	// Create a mock server using httptest
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock response for the request
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{\"message\": \"success\" }"))
		}),
	)
	defer mockServer.Close()

	mockURLEncoder := &MockURLEncoder{}

	// Test with no options provided (should use default)
	opts := determineSendOpt[any](nil, mockURLEncoder)

	// Assert that default InputParser is not nil
	if opts.InputParser == nil {
		t.Error("expected InputParser, got nil")
	}

	// Assert that default Sender is not nil
	if opts.Sender == nil {
		t.Error("expected Sender, got nil")
	}

	// Test that the default Sender function works with the mock server
	_, err := opts.Sender(
		mockServer.URL, // Use the mock server's URL
		"/test-url",
		"GET",
		&RequestData{},
	)
	if err != nil {
		t.Errorf("error sending with default Sender: %v", err)
	}
}
