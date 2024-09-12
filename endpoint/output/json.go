package output

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PakkaSys/fluidapi/core/api"
)

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

type Output[PayloadType any] struct {
	Payload PayloadType `json:"payload,omitempty"`
	Error   *api.Error  `json:"error,omitempty"`
}

func JSON(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	outputData any,
	outputError error,
	statusCode int,
) (*Output[any], error) {
	output := Output[any]{
		Payload: outputData,
		Error:   handleError(outputError),
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonData)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func handleError(outputError error) *api.Error {
	if outputError == nil {
		return nil
	}
	switch errType := outputError.(type) {
	case *api.Error:
		return errType
	default:
		return Error()
	}
}
