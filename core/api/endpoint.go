package api

type Endpoint struct {
	URL         string
	HTTPMethod  string
	Middlewares []Middleware
}

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
