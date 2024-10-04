package util

import (
	"bytes"
	"io"
	"net/http"
)

// RequestWrapper wraps an http.Request, capturing its information for
// inspection and modification.
type RequestWrapper struct {
	*http.Request        // Embedded *http.Request
	BodyContent   []byte // Captured body
}

// NewRequestWrapper creates a new instance of RequestWrapper and captures the
// body content.
func NewRequestWrapper(r *http.Request) (*RequestWrapper, error) {
	// Capture the body content so that it can be read multiple times
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return &RequestWrapper{
		Request:     r,
		BodyContent: bodyBytes,
	}, nil
}

// GetBodyContent provides a way to read the captured body.
func (rw *RequestWrapper) GetBodyContent() []byte {
	return rw.BodyContent
}

// GetBody provides a way to read the captured body.
func (rw *RequestWrapper) GetBody() (func() (io.ReadCloser, error), error) {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(rw.BodyContent)), nil
	}, nil
}
