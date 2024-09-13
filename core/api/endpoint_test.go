package api

import (
	"testing"
)

// TestNewEndpoint tests the creation of a new Endpoint.
func TestNewEndpoint(t *testing.T) {
	middlewares := []Middleware{}
	endpoint := NewEndpoint("/test", "GET", middlewares)

	if endpoint.URL != "/test" {
		t.Errorf("expected URL to be /test, got %s", endpoint.URL)
	}

	if endpoint.HTTPMethod != "GET" {
		t.Errorf("expected HTTP method to be GET, got %s", endpoint.HTTPMethod)
	}

	if len(endpoint.Middlewares) != 0 {
		t.Errorf("expected middlewares length to be 0, got %d", len(endpoint.Middlewares))
	}
}
