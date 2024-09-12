package internal

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PakkaSys/fluidapi/endpoint/output"
)

func ProcessAndSend[Output any](
	options *SendOptions,
	input *SendInputData,
) (*http.Response, *output.Output[Output], error) {
	if input.Body != nil && options.Method == http.MethodGet {
		return nil, nil, fmt.Errorf("body cannot be set for GET requests")
	}

	bodyReader, err := MarshalBody(input.Body)
	if err != nil {
		return nil, nil, err
	}

	request, err := createRequest(
		options.Method,
		ConstructURL(options.Host, options.URL, input.URLParameters),
		bodyReader,
		input.Headers,
		input.Cookies,
	)
	if err != nil {
		return nil, nil, err
	}

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return response, nil, err
	}

	output, err := ResponseToPayload(
		response,
		&output.Output[Output]{},
	)
	if err != nil {
		return response, nil, err
	}

	return response, output, nil
}

func createRequest(
	method string,
	fullURL string,
	bodyReader io.Reader,
	headers map[string]string,
	cookies []http.Cookie,
) (*http.Request, error) {
	var body io.Reader
	if bodyReader != nil {
		body = bodyReader
	}

	req, err := http.NewRequest(
		method,
		fullURL,
		body,
	)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	return req, nil
}
