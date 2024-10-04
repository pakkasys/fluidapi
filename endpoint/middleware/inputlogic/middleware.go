package inputlogic

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/pakkasys/fluidapi/core/api"
	endpointoutput "github.com/pakkasys/fluidapi/endpoint/output"
	"github.com/pakkasys/fluidapi/endpoint/util"
)

const MiddlewareID = "inputlogic"

var internalExpectedErrors []ExpectedError = []ExpectedError{
	{
		ErrorID:       util.InvalidInputError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
	{
		ErrorID:       ValidationError.ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
}

type Callback[Input any, Output any] func(
	wrappers http.ResponseWriter,
	r *http.Request,
	i *Input,
) (*Output, error)

type ValidatedInput interface {
	Validate() []FieldError
}

type IErrorHandler interface {
	Handle(
		handleError error,
		expectedErrors []ExpectedError,
	) (int, *api.Error[any])
}

type IObjectPicker[T any] interface {
	PickObject(r *http.Request, w http.ResponseWriter, obj T) (*T, error)
}

type Options[Input any] struct {
	ObjectPicker  IObjectPicker[Input]
	TraceLoggerFn func(r *http.Request) func(messages ...any)
	ErrorLoggerFn func(r *http.Request) func(messages ...any)
}

func MiddlewareWrapper[Input ValidatedInput, Output any](
	callback Callback[Input, Output],
	inputFactory func() *Input,
	expectedErrors []ExpectedError,
	opts Options[Input],
) *api.MiddlewareWrapper {
	return &api.MiddlewareWrapper{
		ID: MiddlewareID,
		Middleware: Middleware(
			callback,
			inputFactory,
			expectedErrors,
			opts.ObjectPicker,
			opts.TraceLoggerFn,
			opts.ErrorLoggerFn,
		),
		Inputs: []any{*inputFactory()},
	}
}

func Middleware[Input ValidatedInput, Output any](
	callback Callback[Input, Output],
	inputFactory func() *Input,
	expectedErrors []ExpectedError,
	objectPicker IObjectPicker[Input],
	traceLoggerFn func(r *http.Request) func(messages ...any),
	errorLoggerFn func(r *http.Request) func(messages ...any),
) api.Middleware {
	if objectPicker == nil {
		objectPicker = &util.ObjectPicker[Input]{}
	}
	allExpectedErrors := slices.Concat(internalExpectedErrors, expectedErrors)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			output, callbackError := handleInputAndRunCallback(
				w,
				r,
				callback,
				*inputFactory(),
				objectPicker,
				traceLoggerFn,
			)
			statusCode := http.StatusOK

			var outputError error
			if callbackError != nil {
				statusCode, outputError = ErrorHandler{}.Handle(
					callbackError,
					allExpectedErrors,
				)
				if traceLoggerFn != nil {
					traceLoggerFn(r)(fmt.Sprintf(
						"Error handler, status code: %d, callback error: %s, output error: %s",
						statusCode,
						callbackError,
						outputError,
					))
				}
			}

			err := processOutput(
				w,
				r,
				output,
				outputError,
				statusCode,
				traceLoggerFn,
				errorLoggerFn,
			)
			if err != nil {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func handleInputAndRunCallback[Input ValidatedInput, Output any](
	w http.ResponseWriter,
	r *http.Request,
	callback Callback[Input, Output],
	inputObject Input,
	objectPicker IObjectPicker[Input],
	traceLoggerFn func(r *http.Request) func(messages ...any),
) (*Output, error) {
	input, err := objectPicker.PickObject(r, w, inputObject)
	if err != nil {
		return nil, err
	}
	if traceLoggerFn != nil {
		traceLoggerFn(r)("Picked object", input)
	}

	validationErrors := (*input).Validate()
	if len(validationErrors) > 0 {
		validationError := ValidationError.WithData(
			ValidationErrorData{
				Errors: validationErrors,
			},
		)
		return nil, validationError
	}

	return callback(w, r, input)
}

func processOutput(
	w http.ResponseWriter,
	r *http.Request,
	output any,
	outputError error,
	statusCode int,
	traceLoggerFn func(r *http.Request) func(messages ...any),
	errorLoggerFn func(r *http.Request) func(messages ...any),
) error {
	clientOutput, err := endpointoutput.JSON(
		r.Context(),
		w,
		r,
		output,
		outputError,
		statusCode,
	)
	if err != nil {
		if errorLoggerFn != nil {
			errorLoggerFn(r)(fmt.Sprintf(
				"Error handling output JSON: %s",
				err,
			))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if traceLoggerFn != nil {
		traceLoggerFn(r)("Client output", *clientOutput)
	}

	return nil
}
