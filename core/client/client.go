package client

import (
	"github.com/pakkasys/fluidapi/core/client/internal"
)

// Send sends the request and returns the response.
func Send[Input any, Output any](
	input *Input,
	url string,
	host string,
	method string,
) (*Response[Input, Output], error) {
	headers, cookies, urlParams, body, err := internal.ParseInput(method, input)
	if err != nil {
		return nil, err
	}

	inputDataBuilder := internal.NewSendInputDataBuilder().
		WithHeaders(headers).
		WithCookies(cookies).
		WithURLParameters(urlParams)

	if len(body) != 0 {
		inputDataBuilder.WithBody(body)
	}

	response, output, err := internal.ProcessAndSend[Output](
		internal.NewSendOptions(host, url, method),
		inputDataBuilder.Build(),
	)
	if err != nil {
		return nil, err
	}

	return &Response[Input, Output]{
		Response: response,
		Input:    input,
		Output:   output,
	}, nil
}

// HandleSendError returns the output and error if there is an error in the
// request. It will return either request or API error.
func HandleSendError[Input any, Output any](
	output *Response[Input, Output],
	err error,
) (*Response[Input, Output], error) {
	if err != nil {
		return output, err
	}
	if output.APIError() != nil {
		return output, output.APIError()
	}
	return output, err
}
