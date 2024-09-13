package client

import (
	"net/http"
)

// ParsedInput represents the parsed input data for a request.
type ParsedInput struct {
	Headers       map[string]string // HTTP headers to include in the request.
	Cookies       []http.Cookie     // Cookies to include in the request.
	URLParameters map[string]any    // URL parameters for the request.
	Body          map[string]any    // Request body data.
}

// RequestData represents the data ready for a request.
type RequestData struct {
	Headers       map[string]string // HTTP headers to include in the request.
	Cookies       []http.Cookie     // Cookies to include in the request.
	URLParameters map[string]any    // URL parameters for the request.
	Body          any               // Request body data.
}

// Response represents the response from a client request, including the HTTP
// response, input data, and output data.
type Response[Input any, Output any] struct {
	Response *http.Response // The HTTP response object.
	Input    *Input         // The original input data for the request.
	Output   *Output        // The output data of the API response.
}

// SendOptions is used to configure data for sending requests.
type InputDataBuilder interface {
	WithHeaders(headers map[string]string) InputDataBuilder
	WithCookies(cookies []http.Cookie) InputDataBuilder
	WithURLParameters(params map[string]any) InputDataBuilder
	WithBody(body any) InputDataBuilder
	Build() *RequestData
}
