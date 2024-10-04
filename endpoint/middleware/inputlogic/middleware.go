package inputlogic

import (
	"net/http"
	"slices"

	"github.com/pakkasys/fluidapi/core/api"
)

// MiddlewareID is a constant used to identify the middleware within the system.
const MiddlewareID = "inputlogic"

// Internal server expected errors for validation failures.
var internalExpectedErrors []ExpectedError = []ExpectedError{
	{
		ID:         ValidationError.ID,
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
}

// Callback represents the function signature used by the middleware to process
// requests. Input is the type of input to the callback, and Output is the type
// of the output.
type Callback[Input any, Output any] func(
	w http.ResponseWriter,
	r *http.Request,
	i *Input,
) (*Output, error)

// ValidatedInput is an interface that should be implemented by input types that
// can be validated.
type ValidatedInput interface {
	Validate() []FieldError
}

// IObjectPicker represents an interface for picking objects from an HTTP
// request.
type IObjectPicker[T any] interface {
	PickObject(r *http.Request, w http.ResponseWriter, obj T) (*T, error)
}

type ILogger interface {
	Tracef(message string, params ...any)
	Errorf(message string, params ...any)
}

// IOutputHandler represents an interface for processing and sending output
// based on the request context.
type IOutputHandler interface {
	ProcessOutput(
		w http.ResponseWriter,
		r *http.Request,
		out any,
		outError error,
		statusCode int,
	) error
}

// Options represents options that can be configured for the middleware.
// It includes an object picker, output handler, and logging functions.
type Options[Input any] struct {
	// Used to extract the input object from the request.
	ObjectPicker IObjectPicker[Input]
	// Handles output processing.
	OutputHandler IOutputHandler
	// Gets an instance of the logger.
	Logger func(*http.Request) ILogger
}

// MiddlewareWrapper wraps the callback and creates a MiddlewareWrapper
// instance that can be used as a middleware handler.
//
// Input is the type of input to the callback.
// Output is the type of the output to be returned.
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
			opts.OutputHandler,
			opts.Logger,
		),
		Inputs: []any{*inputFactory()},
	}
}

// Middleware constructs a new middleware that handles input validation,
// error handling, and request-response management.
//
// Input is the type of input to the callback.
// Output is the type of output expected from the callback.
//
//   - callback: The function that handles the request.
//   - inputFactory: A function that returns a pointer to the input object.
//   - expectedErrors: A list of expected errors that can be handled by the
//     middleware.
//   - objectPicker: An object picker that can be used to extract the input
//     object from the request.
//   - outputHandler: The handler that processes and sends the output to the
//     client.
//   - traceLoggerFn: A function that can be used to log trace messages.
//   - errorLoggerFn: A function that can be used to log error messages.
func Middleware[Input ValidatedInput, Output any](
	callback Callback[Input, Output],
	inputFactory func() *Input,
	expectedErrors []ExpectedError,
	objectPicker IObjectPicker[Input],
	outputHandler IOutputHandler,
	loggerFn func(*http.Request) ILogger,
) api.Middleware {
	if objectPicker == nil {
		panic("object picker cannot be nil")
	}
	if outputHandler == nil {
		panic("output handler cannot be nil")
	}

	allErrors := slices.Concat(internalExpectedErrors, expectedErrors)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			input, err := handleInput(
				w,
				r,
				*inputFactory(),
				objectPicker,
				loggerFn,
			)
			if err != nil {
				handleError(w, r, err, outputHandler, allErrors, loggerFn)
				return
			}

			out, err := callback(w, r, input)
			if err != nil {
				handleError(w, r, err, outputHandler, allErrors, loggerFn)
				return
			}

			handleOutput(
				w,
				r,
				out,
				nil,
				http.StatusOK,
				outputHandler,
				loggerFn,
			)

			next.ServeHTTP(w, r)
		})
	}
}

func handleError(
	w http.ResponseWriter,
	r *http.Request,
	handleError error,
	outputHandler IOutputHandler,
	expectedErrors []ExpectedError,
	loggerFn func(*http.Request) ILogger,
) {
	statusCode, outError := ErrorHandler{}.Handle(handleError, expectedErrors)

	if loggerFn != nil {
		loggerFn(r).Tracef(
			"Error handler, status code: %d, callback error: %s, output error: %s",
			statusCode,
			handleError,
			outError,
		)
	}

	handleOutput(w, r, nil, outError, statusCode, outputHandler, loggerFn)
}

func handleInput[Input ValidatedInput](
	w http.ResponseWriter,
	r *http.Request,
	inputObject Input,
	objectPicker IObjectPicker[Input],
	loggerFn func(*http.Request) ILogger,
) (*Input, error) {
	input, err := objectPicker.PickObject(r, w, inputObject)
	if err != nil {
		return nil, err
	}
	if loggerFn != nil {
		loggerFn(r).Tracef("Picked object: %v", *input)
	}

	validationErrors := (*input).Validate()
	if len(validationErrors) > 0 {
		return nil, ValidationError.WithData(ValidationErrorData{
			Errors: validationErrors,
		})
	}

	return input, nil
}

func handleOutput(
	w http.ResponseWriter,
	r *http.Request,
	output any,
	outputError error,
	statusCode int,
	outputHandler IOutputHandler,
	loggerFn func(*http.Request) ILogger,
) {
	err := outputHandler.ProcessOutput(w, r, output, outputError, statusCode)

	if err != nil {
		if loggerFn != nil {
			loggerFn(r).Errorf(
				"Error processing output: %s",
				err,
			)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
