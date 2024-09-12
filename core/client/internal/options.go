package internal

import "net/http"

type SendOptions struct {
	Host   string
	URL    string
	Method string
}

func NewSendOptions(
	host string,
	url string,
	method string,
) *SendOptions {
	return &SendOptions{
		Host:   host,
		URL:    url,
		Method: method,
	}
}

type SendInputData struct {
	Headers       map[string]string
	Cookies       []http.Cookie
	URLParameters map[string]any
	Body          any
}

type SendInputDataBuilder struct {
	SendInputData
}

func NewSendInputDataBuilder() *SendInputDataBuilder {
	return &SendInputDataBuilder{}
}

func (b *SendInputDataBuilder) Build() *SendInputData {
	return &b.SendInputData
}

func (b *SendInputDataBuilder) WithBody(body any) *SendInputDataBuilder {
	b.Body = body
	return b
}

func (b *SendInputDataBuilder) WithHeaders(
	headers map[string]string,
) *SendInputDataBuilder {
	b.Headers = headers
	return b
}

func (b *SendInputDataBuilder) WithCookies(
	cookies []http.Cookie,
) *SendInputDataBuilder {
	b.Cookies = cookies
	return b
}

func (b *SendInputDataBuilder) WithURLParameters(
	urlParameters map[string]any,
) *SendInputDataBuilder {
	b.URLParameters = urlParameters
	return b
}
