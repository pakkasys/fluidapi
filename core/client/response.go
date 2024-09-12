package client

import (
	"net/http"

	"github.com/PakkaSys/fluidapi/core/api"
	"github.com/PakkaSys/fluidapi/endpoint/output"
)

type Response[Input any, Output any] struct {
	Response *http.Response
	Input    *Input
	Output   *output.Output[Output]
}

// APIPayload returns the payload of the API response from the output if there
// are no API errors. Returns nil otherwise.
func (r *Response[Input, Output]) APIPayload() *Output {
	if r.APIError() != nil {
		return nil
	}
	if r.Output != nil {
		return &r.Output.Payload
	}
	return nil
}

// APIError returns the error of the API response from the output.
func (r *Response[Input, Output]) APIError() *api.Error {
	if r.Output == nil {
		return nil
	}
	return r.Output.Error
}
