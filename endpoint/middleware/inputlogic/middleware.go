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

// IErrorHandler represents an interface for handling errors.
type IErrorHandler interface {
	Handle(
		handleError error,
		expectedErrors []ExpectedError,
	) (int, *api.Error[any])
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
	logger func(*http.Request) ILogger,
) api.Middleware {
	if objectPicker == nil {
		panic("object picker cannot be nil")
	}
	if outputHandler == nil {
		panic("output handler cannot be nil")
	}

	allExpectedErrors := slices.Concat(internalExpectedErrors, expectedErrors)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handle the input and execute the callback.
			out, callbackError := handleInputAndRunCallback(
				w,
				r,
				callback,
				*inputFactory(),
				objectPicker,
				logger,
			)
			statusCode := http.StatusOK

			// Handle any callback errors.
			var outError error
			if callbackError != nil {
				statusCode, outError = ErrorHandler{}.Handle(
					callbackError,
					allExpectedErrors,
				)
				if logger != nil {
					logger(r).Tracef(
						"Error handler, status code: %d, callback error: %s, output error: %s",
						statusCode,
						callbackError,
						outError,
					)
				}
			}

			err := outputHandler.ProcessOutput(w, r, out, outError, statusCode)
			if err != nil {
				if logger != nil {
					logger(r).Errorf(
						"Error processing output: %s",
						err,
					)
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// handleInputAndRunCallback handles the input by picking the object from the
// request, validating it, and executing the callback if validation succeeds.
func handleInputAndRunCallback[Input ValidatedInput, Output any](
	w http.ResponseWriter,
	r *http.Request,
	callback Callback[Input, Output],
	inputObject Input,
	objectPicker IObjectPicker[Input],
	logger func(*http.Request) ILogger,
) (*Output, error) {
	input, err := objectPicker.PickObject(r, w, inputObject)
	if err != nil {
		return nil, err
	}
	if logger != nil {
		logger(r).Tracef("Picked object", input)
	}

	// Validate the input.
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
