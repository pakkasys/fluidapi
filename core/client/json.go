package client

import (
	"net/http"
)

type SendResult[Payload any] struct {
	Response *http.Response
	Output   *Payload
}

// HandlerOpts contains data handler options.
type HandlerOpts[Payload any] struct {
	InputParser func(method string, input any) (*ParsedInput, error)
	Sender      func(
		host string,
		url string,
		method string,
		inputData *RequestData,
	) (*SendResult[Payload], error)
}

// Send sends a request to the specified URL with the provided input, host, and
// HTTP method and returns a Response containing the input, output, and HTTP
// response.
//   - input: The request input data.
//   - url: The endpoint URL path.
//   - host: The host server to send the request to.
//   - method: The HTTP method (e.g., GET, POST).
//   - opts: Optional SendOpts. If not provided, the default SendOpts will be
//     used.
func Send[Input any, Output any](
	input *Input,
	url string,
	host string,
	method string,
	opts ...HandlerOpts[Output],
) (*Response[Input, Output], error) {
	useOpts := determineSendOpt(opts)

	parsedInput, err := useOpts.InputParser(method, input)
	if err != nil {
		return nil, err
	}

	requestData := RequestData{
		Headers:       parsedInput.Headers,
		Cookies:       parsedInput.Cookies,
		URLParameters: parsedInput.URLParameters,
	}
	if len(parsedInput.Body) != 0 {
		requestData.Body = parsedInput.Body
	}

	sendResult, err := useOpts.Sender(
		host,
		url,
		method,
		&requestData,
	)
	if err != nil {
		return nil, err
	}

	return &Response[Input, Output]{
		Response: sendResult.Response,
		Input:    input,
		Output:   sendResult.Output,
	}, nil
}

func determineSendOpt[Payload any](
	opts []HandlerOpts[Payload],
) *HandlerOpts[Payload] {
	if len(opts) == 0 {
		return &HandlerOpts[Payload]{
			InputParser: parseInput,
			Sender: func(
				host string,
				url string,
				method string,
				inputData *RequestData,
			) (*SendResult[Payload], error) {
				return processAndSend[Payload](
					&http.Client{},
					host,
					url,
					method,
					inputData,
				)
			},
		}
	} else {
		return &opts[0]
	}
}
