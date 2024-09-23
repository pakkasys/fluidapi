package util

import (
	"net/http"
)

// ResponseWrapper wraps an http.ResponseWriter, capturing headers and the
// status code written to it.
type ResponseWrapper struct {
	http.ResponseWriter
	Headers       http.Header
	StatusCode    int
	Body          []byte
	headerWritten bool
}

// NewResponseWrapper creates a new instance of ResponseWrapper.
func NewResponseWrapper(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{
		ResponseWriter: w,
		Headers:        make(http.Header),
		StatusCode:     http.StatusOK, // Default status code
		headerWritten:  false,
	}
}

// Header overrides the Header method of the http.ResponseWriter interface.
// It returns the captured headers without modifying the underlying
// ResponseWriter's headers.
func (rw *ResponseWrapper) Header() http.Header {
	return rw.Headers
}

// WriteHeader captures the status code to be written, delaying its execution.
func (rw *ResponseWrapper) WriteHeader(statusCode int) {
	if !rw.headerWritten { // Only write headers once
		rw.StatusCode = statusCode
		// Apply the captured headers to the underlying ResponseWriter.
		for key, values := range rw.Headers {
			for _, value := range values {
				rw.ResponseWriter.Header().Add(key, value)
			}
		}
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
	}
}

// Write writes the response body and ensures headers and status code are written.
func (rw *ResponseWrapper) Write(data []byte) (int, error) {
	// Ensure headers and status code are written before writing the body.
	if !rw.headerWritten {
		rw.WriteHeader(rw.StatusCode)
	}

	// Append the data to the response body buffer.
	rw.Body = append(rw.Body, data...)

	// Write the data to the underlying ResponseWriter.
	return rw.ResponseWriter.Write(data)
}
