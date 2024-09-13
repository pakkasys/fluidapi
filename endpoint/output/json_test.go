package output

// // TestAPIPayload tests the APIPayload function of the Response struct.
// func TestAPIPayload(t *testing.T) {
// 	output:="test_payload"
// 	resp := &Response[string]{
// 		Response: &http.Response{},
// 		Output: &output,
// 	}

// 	// Case where APIError is nil
// 	payload := resp.APIPayload()
// 	assert.NotNil(t, payload, "expected non-nil payload")
// 	assert.Equal(t, "test_payload", payload, "unexpected payload")

// 	// Case where APIError is not nil
// 	resp.Output.Error = &api.Error{ID: "test_error"}
// 	payload = resp.APIPayload()
// 	assert.Nil(t, payload, "expected nil payload when there is an API error")

// 	// Case where error and output are nil
// 	resp.Output.Error = nil
// 	resp.Output = nil
// 	payload = resp.APIPayload()
// 	assert.Nil(t, payload, "expected nil payload when payload is nil")
// }

// // TestAPIError tests the APIError function of the Response struct.
// func TestAPIError(t *testing.T) {
// 	resp := &Response{}

// 	// Case where Output is nil
// 	assert.Nil(t, resp.APIError(), "expected nil error when Output is nil")

// 	// Case where Output is not nil and has an error
// 	resp.Output = &output.Output[any]{Error: &api.Error{ID: "test_error"}}
// 	err := resp.APIError()
// 	assert.NotNil(t, err, "expected non-nil error")
// 	assert.Equal(t, "test_error", err.ID, "unexpected error ID")
// }

// // TestHandleSendError tests the HandleSendError function with error.
// func TestHandleSendError_Error(t *testing.T) {
// 	mockResponse := &Response{
// 		Response: &http.Response{StatusCode: http.StatusOK},
// 		Output:   nil,
// 	}
// 	err := errors.New("test error")

// 	resp, returnedErr := HandleSendError(mockResponse, err)

// 	// Assert that the same error is returned
// 	if returnedErr != err {
// 		t.Errorf("expected error %v, got %v", err, returnedErr)
// 	}

// 	// Assert that the response is the same as provided
// 	if resp != mockResponse {
// 		t.Errorf("expected response %v, got %v", mockResponse, resp)
// 	}
// }

// // TestHandleSendError_APIError tests the HandleSendError function with API
// // error.
// func TestHandleSendError_APIError(t *testing.T) {
// 	apiErr := &api.Error{ID: "api_error"}
// 	mockResponse := &Response{
// 		Response: &http.Response{StatusCode: http.StatusOK},
// 		Output:   &output.Output[any]{Error: apiErr},
// 	}

// 	resp, returnedErr := HandleSendError(mockResponse, nil)

// 	// Assert that the APIError is returned
// 	if returnedErr != apiErr {
// 		t.Errorf("expected APIError %v, got %v", apiErr, returnedErr)
// 	}

// 	// Assert that the response is the same as provided
// 	if resp != mockResponse {
// 		t.Errorf("expected response %v, got %v", mockResponse, resp)
// 	}
// }

// // TestHandleSendError_NoError tests the HandleSendError function with no error.
// func TestHandleSendError_NoError(t *testing.T) {
// 	mockResponse := &Response{
// 		Response: &http.Response{StatusCode: http.StatusOK},
// 		Output:   &output.Output[any]{Error: nil},
// 	}

// 	resp, returnedErr := HandleSendError(mockResponse, nil)

// 	// Assert that no error is returned
// 	if returnedErr != nil {
// 		t.Errorf("expected no error, got %v", returnedErr)
// 	}

// 	// Assert that the response is the same as provided
// 	if resp != mockResponse {
// 		t.Errorf("expected response %v, got %v", mockResponse, resp)
// 	}
// }
