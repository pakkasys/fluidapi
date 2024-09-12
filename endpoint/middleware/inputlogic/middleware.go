package inputlogic

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/PakkaSys/fluidapi/core/api"
	endpointoutput "github.com/PakkaSys/fluidapi/endpoint/output"
	"github.com/PakkaSys/fluidapi/endpoint/util"
)

const MiddlewareID = "inputlogic"

var internalExpectedErrors []ExpectedError = []ExpectedError{
	{
		ErrorID:       util.INVALID_INPUT_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: false,
	},
	{
		ErrorID:       VALIDATION_ERROR_ID,
		StatusCode:    http.StatusBadRequest,
		DataIsVisible: true,
	},
}

type Callback[Input any, Output any] func(
	writer http.ResponseWriter,
	request *http.Request,
	input *Input,
) (*Output, error)

type IErrorHandler interface {
	Handle(
		handleError error,
		expectedErrors []ExpectedError,
	) (int, *api.Error)
}

type IObjectPicker[T any] interface {
	PickObject(r *http.Request, w http.ResponseWriter, obj T) (*T, error)
}

type IValidator interface {
	ValidateStruct(obj any) error
	GetErrorStrings(err error) []string
}

func MiddlewareWrapper[Input any, Output any](
	callback Callback[Input, Output],
	inputFactory func() *Input,
	expectedErrors []ExpectedError,
	errorHandler IErrorHandler,
	objectPicker IObjectPicker[Input],
	validator IValidator,
	traceLoggerFn func(r *http.Request) func(messages ...any),
	errorLoggerFn func(r *http.Request) func(messages ...any),
) *api.MiddlewareWrapper {
	return api.NewMiddlewareWrapperBuilder().
		ID(MiddlewareID).
		Middleware(Middleware(
			callback,
			inputFactory,
			expectedErrors,
			errorHandler,
			objectPicker,
			validator,
			traceLoggerFn,
			errorLoggerFn,
		)).
		Build()
}

func Middleware[Input any, Output any](
	callback Callback[Input, Output],
	inputFactory func() *Input,
	expectedErrors []ExpectedError,
	errorHandler IErrorHandler,
	objectPicker IObjectPicker[Input],
	validator IValidator,
	traceLoggerFn func(r *http.Request) func(messages ...any),
	errorLoggerFn func(r *http.Request) func(messages ...any),
) api.Middleware {
	if errorHandler == nil {
		errorHandler = &ErrorHandler{}
	}
	if objectPicker == nil {
		objectPicker = &util.ObjectPicker[Input]{}
	}
	if validator == nil {
		validator = util.NewValidation()
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
				validator,
				traceLoggerFn,
			)
			statusCode := http.StatusOK

			var outputError error
			if callbackError != nil {
				statusCode, outputError = errorHandler.Handle(
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

func handleInputAndRunCallback[Input any, Output any](
	w http.ResponseWriter,
	r *http.Request,
	callback Callback[Input, Output],
	inputObject Input,
	objectPicker IObjectPicker[Input],
	validator IValidator,
	traceLoggerFn func(r *http.Request) func(messages ...any),
) (*Output, error) {
	input, err := objectPicker.PickObject(r, w, inputObject)
	if err != nil {
		return nil, err
	}
	if traceLoggerFn != nil {
		traceLoggerFn(r)("Picked object", input)
	}

	if err := validator.ValidateStruct(input); err != nil {
		return nil, ValidationError(validator.GetErrorStrings(err))
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
