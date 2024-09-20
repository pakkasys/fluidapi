package api

// Endpoint represents an API endpoint.
type Endpoint struct {
	URL         string
	Method      string
	Middlewares []Middleware
}
