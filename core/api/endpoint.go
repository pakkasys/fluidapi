package api

// Endpoint represents an API endpoint with a URL, HTTP method, and associated
// middlewares.
type Endpoint struct {
	URL         string
	HTTPMethod  string
	Middlewares []Middleware
}

// NewEndpoint creates a new Endpoint with the provided URL, HTTP method, and
// middlewares.
// - url: The URL path of the endpoint.
// - httpMethod: The HTTP method for the endpoint (e.g., GET, POST).
// - middlewares: A slice of middlewares to apply to the endpoint.
func NewEndpoint(
	url string,
	httpMethod string,
	middlewares []Middleware,
) *Endpoint {
	return &Endpoint{
		URL:         url,
		HTTPMethod:  httpMethod,
		Middlewares: middlewares,
	}
}
