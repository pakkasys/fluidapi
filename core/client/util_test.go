package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// errorReader is a helper type to simulate a reader that always errors.
type errorReader struct{}

// Read always returns an error to simulate a faulty reader.
func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}

type TestStruct struct {
	Name  string `json:"name" source:"url"`
	Age   int    `json:"age" source:"body"`
	Token string `json:"token" source:"headers"`
	Auth  string `json:"auth" source:"cookies"`
}

// MockOutput is a mock output struct.
type MockOutput struct {
	Message string `json:"message"`
}

// MockRoundTripper simulates errors during HTTP requests.
type MockRoundTripper struct {
	Err error
}

func (m *MockRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	return nil, m.Err
}

// TestProcessAndSend tests the processAndSend function.
func TestProcessAndSend(t *testing.T) {
	// Create a mock server to simulate an API endpoint
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response for POST method
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			out := MockOutput{Message: "success"}
			err := json.NewEncoder(w).Encode(out)
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
			return
		}
		// Mock response for GET method
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			out := MockOutput{Message: "success"}
			err := json.NewEncoder(w).Encode(out)
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
			return
		}
	}))
	defer mockServer.Close()

	// Create a mock client using httptest
	mockClient := &http.Client{}

	// Case 1: Valid POST request with a body
	input := &RequestData{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    map[string]any{"key": "value"},
	}

	result, err := processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusOK, result.Response.StatusCode)
	assert.Equal(t, "success", result.Output.Message)

	// Case 2: Invalid GET request with a body
	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodGet,
		input,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "body cannot be set for GET requests")

	// Case 3: Error in marshalBody
	unmarshalableBody := make(chan int) // Channels cannot be marshaled
	input = &RequestData{
		Body: unmarshalableBody,
	}
	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type: chan int")

	// Case 4: Error in constructURL
	input = &RequestData{
		URLParameters: map[string]any{
			"invalid_param": make(chan int), // Channels cannot be encoded in URL
		},
	}
	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
	)
	assert.NotNil(t, err)
	assert.Contains(
		t,
		err.Error(),
		"value type not supported by URL encoding: chan",
	)

	// Case 5: Error in createRequest
	input = &RequestData{}
	_, err = processAndSend[MockOutput](
		mockClient,
		"http://[::1]:namedport",
		"/test",
		http.MethodGet,
		input,
	)
	assert.NotNil(t, err)

	// Case 6: Error in client.Do
	mockClient = &http.Client{
		Transport: &MockRoundTripper{Err: errors.New("network error")},
	}

	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "network error")

	// Case 7: Error in responseToPayload
	mockServerInvalidJSON := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("invalid json"))
			assert.Nil(t, err)
		}),
	)
	defer mockServerInvalidJSON.Close()

	_, err = processAndSend[MockOutput](
		&http.Client{},
		mockServerInvalidJSON.URL,
		"/test",
		http.MethodPost,
		&RequestData{},
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "JSON unmarshal error")
}

// TestCreateRequest tests the createRequest function.
func TestCreateRequest(t *testing.T) {
	// Case 1: Valid request with headers and cookies
	method := http.MethodPost
	fullURL := "http://localhost/test"
	body := bytes.NewReader([]byte(`{"key": "value"}`))
	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token"}
	cookies := []http.Cookie{
		{Name: "session_id", Value: "12345"},
	}

	req, err := createRequest(method, fullURL, body, headers, cookies)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, fullURL, req.URL.String())

	// Check headers
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))

	// Check cookies
	cookie, err := req.Cookie("session_id")
	assert.Nil(t, err)
	assert.Equal(t, "12345", cookie.Value)

	// Check body
	bodyBytes, err := io.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"key": "value"}`, string(bodyBytes))

	// Case 2: Request with no body
	req, err = createRequest(http.MethodGet, fullURL, nil, headers, cookies)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, fullURL, req.URL.String())

	// Case 3: Invalid URL
	invalidURL := "http://[::1]:namedport"
	req, err = createRequest(http.MethodGet, invalidURL, nil, headers, cookies)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

// TestMarshalBody_NilBody tests marshalBody with a nil body.
func TestMarshalBody_NilBody(t *testing.T) {
	// Call marshalBody with nil
	reader, err := marshalBody(nil)

	// Verify no error is returned
	assert.Nil(t, err, "expected no error for nil body")

	// Verify the reader is not nil and empty
	assert.NotNil(t, reader, "expected non-nil reader for nil body")
	bodyBytes, _ := io.ReadAll(reader)
	assert.Equal(t, []byte{}, bodyBytes, "expected empty body bytes")
}

// TestMarshalBody_ValidBody tests marshalBody with a valid body.
func TestMarshalBody_ValidBody(t *testing.T) {
	// Define a valid body
	body := map[string]string{"key": "value"}

	// Call marshalBody with a valid body
	reader, err := marshalBody(body)

	// Verify no error is returned
	assert.Nil(t, err, "expected no error for valid body")

	// Verify the reader is not nil and contains the correct JSON
	assert.NotNil(t, reader, "expected non-nil reader for valid body")
	bodyBytes, _ := io.ReadAll(reader)
	expectedBytes, _ := json.Marshal(body)
	assert.Equal(t, expectedBytes, bodyBytes, "unexpected body bytes")
}

// TestMarshalBody_UnmarshalableBody tests marshalBody with an unmarshalable
// body.
func TestMarshalBody_UnmarshalableBody(t *testing.T) {
	// Define an unmarshalable body (channel type cannot be marshaled)
	body := make(chan int)

	// Call marshalBody with an unmarshalable body
	reader, err := marshalBody(body)

	// Verify that an error is returned
	assert.NotNil(t, err, "expected error for unmarshalable body")

	// Verify the reader is nil
	assert.Nil(t, reader, "expected nil reader for unmarshalable body")
}

// TestConstructURL_NoParams tests constructing a URL without any URL
// parameters.
func TestConstructURL_NoParams(t *testing.T) {
	host := "http://localhost"
	path := "/api/v1/resource"
	urlParams := map[string]any{}

	// Call constructURL
	result, err := constructURL(host, path, urlParams)

	// Verify that no error is returned
	assert.Nil(t, err, "unexpected error")

	// Verify that the URL is correct
	assert.Equal(t, host+path, *result, "unexpected URL")
}

// TestConstructURL_WithParams tests constructing a URL with URL parameters.
func TestConstructURL_WithParams(t *testing.T) {
	host := "http://localhost"
	path := "/api/v1/resource"
	urlParams := map[string]any{
		"param1": "value1",
		"param2": "value2",
	}

	// Call constructURL
	result, err := constructURL(host, path, urlParams)

	// Verify that no error is returned
	assert.Nil(t, err, "unexpected error")

	// Verify that the URL is correct
	expected := host + path + "?param1=value1&param2=value2"
	assert.Equal(t, expected, *result, "unexpected URL")
}

// // TestConstructURL_WithNilParamValue tests error handling when constructing
// // an URL.
// func TestConstructURL_WithError(t *testing.T) {
// 	host := "http://localhost"
// 	path := "/api/v1/resource"
// 	urlParams := map[string]any{
// 		"param1": nil,
// 	}
// 	// Call constructURL
// 	result, err := constructURL(host, path, urlParams)

// 	// Verify that an error is returned
// 	assert.NotNil(t, err, "expected error")
// 	assert.Equal(t, host+path, *result, "unexpected URL")
// }

// TestResponseToPayload tests the responseToPayload function.
func TestResponseToPayload(t *testing.T) {
	// Define a mock JSON payload
	mockPayload := map[string]string{"key": "value"}
	mockBody, _ := json.Marshal(mockPayload)

	// Create a mock HTTP response
	response := &http.Response{
		Body: io.NopCloser(bytes.NewBuffer(mockBody)),
	}

	var output map[string]string

	// Call responseToPayload
	result, err := responseToPayload(response, &output)

	// Check if there is no error
	assert.Nil(t, err, "expected no error converting response to payload")

	// Verify that the payload is correctly unmarshalled
	assert.Equal(t, mockPayload, *result, "unexpected payload")
}

func TestResponseToPayload_ReadAllError(t *testing.T) {
	// Create a mock HTTP response with an error
	response := &http.Response{
		Body: &errorReader{},
	}

	var output map[string]string

	// Call responseToPayload
	result, err := responseToPayload(response, &output)

	// Check if there is an error
	assert.NotNil(t, err, "expected an error reading response body")

	// Verify that the result is nil due to the error
	assert.Nil(t, result, "expected nil result due to read error")
}

// TestResponseToPayload_UnmarshalError tests the responseToPayload function
// with an unmarshalling error.
func TestResponseToPayload_UnmarshalError(t *testing.T) {
	// Create a mock HTTP response with invalid JSON
	invalidJSON := "invalid json"
	response := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(invalidJSON)),
	}

	var output map[string]string

	// Call responseToPayload
	result, err := responseToPayload(response, &output)

	// Check if there is an error
	assert.NotNil(t, err, "expected an error unmarshalling invalid JSON")

	// Verify that the result is nil due to the error
	assert.Nil(t, result, "expected nil result due to unmarshalling error")
}

// TestToURLParamString_Map tests toURLParamString with a valid map input.
func TestToURLParamString_Map(t *testing.T) {
	input := map[string]any{"key1": "value1", "key2": "value2"}

	// Call toURLParamString with a valid map
	result, err := toURLParamString(input)

	// Verify no error is returned
	assert.Nil(t, err, "expected no error for valid map input")

	// Verify the result contains the correct URL parameters
	expectedResult1 := "key1=value1&key2=value2"
	expectedResult2 := "key2=value2&key1=value1" // Order is not guaranteed

	assert.True(t, *result == expectedResult1 || *result == expectedResult2,
		"unexpected URL parameter string for map input: %s", result)
}

// TestToURLParamString_EmptyMap tests toURLParamString with an empty map.
func TestToURLParamString_EmptyMap(t *testing.T) {
	input := map[string]any{}

	// Call toURLParamString with an empty map
	result, err := toURLParamString(input)

	// Verify no error is returned
	assert.Nil(t, err, "expected no error for empty map input")

	// Verify the result is an empty string
	assert.Equal(t, "", *result, "expected empty string for empty map input")
}

// TestToURLParamString_NilMap tests toURLParamString with a nil map.
func TestToURLParamString_NilMap(t *testing.T) {
	// Call toURLParamString with a nil map
	_, err := toURLParamString(nil)

	// Verify an error is returned
	assert.NotNil(t, err, "expected error for nil map input")
	assert.Equal(t, "input map is nil", err.Error(), "unexpected error")
}

// Test parseInput function
func TestParseInput(t *testing.T) {
	// Test nil input
	_, err := parseInput(http.MethodGet, nil)
	assert.NotNil(t, err, "expected error for nil input")
	assert.Contains(t, err.Error(), "parsed input is nil")

	// Test valid input with different HTTP methods
	input := &TestStruct{Name: "John", Age: 30, Token: "abc123"}
	result, err := parseInput(http.MethodGet, input)
	assert.Nil(t, err, "expected no error for valid input")
	assert.Equal(t, map[string]any{"name": "John"}, result.URLParameters, "unexpected URL parameters")
	assert.Equal(t, map[string]any{"age": int64(30)}, result.Body, "unexpected body")
	assert.Equal(t, map[string]string{"token": "abc123"}, result.Headers, "unexpected headers")

	result, err = parseInput(http.MethodPost, input)
	assert.Nil(t, err, "expected no error for valid input")
	assert.Equal(t, map[string]any{"name": "John"}, result.URLParameters, "unexpected URL parameters")
	assert.Equal(t, map[string]any{"age": int64(30)}, result.Body, "unexpected body")
	assert.Equal(t, map[string]string{"token": "abc123"}, result.Headers, "unexpected headers")

	// Test invalid source tag
	type InvalidStruct struct {
		Field string `json:"field" source:"invalid"`
	}
	_, err = parseInput(http.MethodGet, &InvalidStruct{Field: "test"})
	assert.NotNil(t, err, "expected error for invalid source tag")
	assert.Contains(t, err.Error(), "invalid source tag")
}

// Test determineDefaultPlacement function
func TestDetermineDefaultPlacement(t *testing.T) {
	assert.Equal(t, requestURL, determineDefaultPlacement(http.MethodGet))
	assert.Equal(t, requestBody, determineDefaultPlacement(http.MethodPost))
	assert.Equal(t, requestBody, determineDefaultPlacement("UNKNOWN_METHOD"))
}

// Test processField function
func TestProcessField(t *testing.T) {
	headers := make(map[string]string)
	cookies := make([]http.Cookie, 0)
	urlParameters := make(map[string]any)
	body := make(map[string]any)

	input := TestStruct{Name: "John", Age: 30, Token: "abc123"}
	inputVal := reflect.ValueOf(&input).Elem()
	inputType := inputVal.Type()

	// Case 1: Test valid field processing with explicit source tag
	field := inputVal.FieldByName("Name")
	fieldInfo, _ := inputType.FieldByName("Name")

	_, err := processField(
		field,
		fieldInfo,
		requestURL,
		headers,
		cookies,
		urlParameters,
		body,
	)
	assert.Nil(t, err, "expected no error")
	assert.Equal(
		t,
		map[string]any{"name": "John"},
		urlParameters,
		"unexpected URL parameters",
	)

	// Case 2: Test field processing with default placement (no source tag)
	type DefaultPlacementStruct struct {
		FieldWithoutSourceTag string `json:"default_field"`
	}

	input2 := DefaultPlacementStruct{FieldWithoutSourceTag: "test value"}
	inputVal2 := reflect.ValueOf(&input2).Elem()
	inputType2 := inputVal2.Type()

	field2 := inputVal2.FieldByName("FieldWithoutSourceTag")
	fieldInfo2, _ := inputType2.FieldByName("FieldWithoutSourceTag")

	_, err = processField(
		field2,
		fieldInfo2,
		requestBody,
		headers,
		cookies,
		urlParameters,
		body,
	)
	assert.Nil(t, err, "expected no error")
	assert.Equal(
		t,
		map[string]any{"default_field": "test value"},
		body,
		"unexpected body content for default placement",
	)
}

// Test determineFieldName function
func TestDetermineFieldName(t *testing.T) {
	assert.Equal(t, "name", determineFieldName("name", "Field"))
	assert.Equal(t, "Field", determineFieldName("", "Field"))
}

// Test extractFieldValue function
func TestExtractFieldValue(t *testing.T) {
	val := reflect.ValueOf(true)
	assert.Equal(t, true, extractFieldValue(val))

	val = reflect.ValueOf("test")
	assert.Equal(t, "test", extractFieldValue(val))

	val = reflect.ValueOf(123)
	assert.Equal(t, int64(123), extractFieldValue(val))

	val = reflect.ValueOf(uint(123))
	assert.Equal(t, uint64(123), extractFieldValue(val))

	val = reflect.ValueOf(123.45)
	assert.Equal(t, 123.45, extractFieldValue(val))

	val = reflect.ValueOf(struct{}{})
	assert.Equal(t, val.Interface(), extractFieldValue(val))
}

// Test placeFieldValue function
func TestPlaceFieldValue(t *testing.T) {
	headers := make(map[string]string)
	cookies := make([]http.Cookie, 0)
	urlParameters := make(map[string]any)
	body := make(map[string]any)

	placements := []string{
		requestURL,
		requestBody,
		requestHeaders,
		requestCookies,
	}
	for _, placement := range placements {
		var err error
		cookies, err = placeFieldValue(
			placement,
			"key",
			"value",
			headers,
			cookies,
			urlParameters,
			body,
		)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"key": "value"}, urlParameters)
	}

	_, err := placeFieldValue(
		"invalid",
		"key",
		"value",
		headers,
		cookies,
		urlParameters,
		body,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid source tag")
}
