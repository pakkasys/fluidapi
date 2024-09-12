package util

import (
	"bytes"
	"io"
	"net/http"
)

// RequestWrapper wraps an http.Request, capturing its URL, headers, and body
// for inspection and potentially modification.
type RequestWrapper struct {
	*http.Request        // Embedding *http.Request
	BodyContent   []byte // Captured body content for potential reuse
}

// NewRequestWrapper creates a new instance of RequestWrapper, capturing the body.
func NewRequestWrapper(r *http.Request) (*RequestWrapper, error) {
	// Capture the body content without disrupting the original request.
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err // Return error if reading the body fails
	}
	// Restore the body to the original request to ensure it can be read again later.
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return &RequestWrapper{
		Request:     r,
		BodyContent: bodyBytes,
	}, nil
}

// GetBodyContent provides a way to read the captured body content multiple times.
func (rw *RequestWrapper) GetBodyContent() []byte {
	return rw.BodyContent
}

// Replicate the original GetBody method behavior for HTTP/2:
func (rw *RequestWrapper) GetBody() (func() (io.ReadCloser, error), error) {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(rw.BodyContent)), nil
	}, nil
}
