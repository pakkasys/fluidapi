package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/stretchr/testify/assert"
)

// Mock HTTP server that we can control
type MockServer struct {
	ListenAndServeFunc func() error
	ShutdownFunc       func(ctx context.Context) error
}

func (m *MockServer) ListenAndServe() error {
	if m.ListenAndServeFunc != nil {
		return m.ListenAndServeFunc()
	}
	return nil
}

func (m *MockServer) Shutdown(ctx context.Context) error {
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc(ctx)
	}
	return nil
}

func startTestHTTPServer(
	port int,
	httpEndpoints []api.Endpoint,
) (shutdownFunc func()) {
	mockLogger := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			fmt.Println(messages...)
		}
	}

	// Start the HTTP server in a goroutine
	go func() {
		err := HTTPServer(DefaultHTTPServer(
			port,
			httpEndpoints,
			mockLogger,
			mockLogger,
		))
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Return a function to gracefully shutdown the server
	return func() {
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)  // Send interrupt signal to current process
		time.Sleep(1 * time.Second) // Wait for graceful shutdown
	}
}

// Test case to simulate a graceful shutdown
func TestHTTPServer_Lifecycle(t *testing.T) {
	httpEndpoints := []api.Endpoint{{URL: "/test", HTTPMethod: "GET"}}

	shutdown := startTestHTTPServer(8080, httpEndpoints)
	defer shutdown()

	// Simulate sending a shutdown signal
	time.Sleep(500 * time.Millisecond)
	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(os.Interrupt)

	// Give some time for graceful shutdown
	time.Sleep(1 * time.Second)
}

// Test case for calling an unregistered (not found) endpoint
func TestHTTPServer_CallNotFoundEndpoint(t *testing.T) {
	// Define endpoints (no /notfound registered)
	httpEndpoints := []api.Endpoint{
		{
			URL:        "/test",
			HTTPMethod: "GET",
			Middlewares: []api.Middleware{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
							_, err := w.Write([]byte("Hello, World!"))
							assert.NoError(t, err)
						},
					)
				},
			},
		},
	}

	// Start the HTTP server
	shutdown := startTestHTTPServer(8080, httpEndpoints)
	defer shutdown()

	// Make a request to an unregistered endpoint
	resp, err := http.Get("http://localhost:8080/notfound")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert the response
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Test case for calling an endpoint with the wrong method
func TestHTTPServer_CallEndpointWithWrongMethod(t *testing.T) {
	// Define endpoints
	httpEndpoints := []api.Endpoint{
		{
			URL:        "/test",
			HTTPMethod: "GET",
			Middlewares: []api.Middleware{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
							_, err := w.Write([]byte("Hello, World!"))
							assert.NoError(t, err)
						},
					)
				},
			},
		},
	}

	// Start the HTTP server
	shutdown := startTestHTTPServer(8080, httpEndpoints)
	defer shutdown()

	// Make a request to the registered endpoint with the wrong method (POST)
	resp, err := http.Post("http://localhost:8080/test", "application/json", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert the response
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

// Test case for calling a registered endpoint
func TestHTTPServer_CallRegisteredEndpoint(t *testing.T) {
	mockLogger := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			fmt.Println(messages...)
		}
	}

	// Define endpoints
	httpEndpoints := []api.Endpoint{
		{
			URL:        "/test",
			HTTPMethod: "GET",
			Middlewares: []api.Middleware{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
							_, err := w.Write([]byte("Hello, World!"))
							assert.NoError(t, err)
						},
					)
				},
			},
		},
	}

	// Create a mock server
	mockServer := &MockServer{
		ListenAndServeFunc: func() error { return nil },
		ShutdownFunc:       func(ctx context.Context) error { return nil },
	}

	// Create a test HTTP server
	go func() {
		_ = HTTPServer(mockServer)
	}()

	// Create a test request to the registered endpoint
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Setup mux and test handler
	mux := setupMux(httpEndpoints, mockLogger, mockLogger)
	mux.ServeHTTP(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, World!", rr.Body.String())
}

// Test case for DefaultHTTPServer function
func TestDefaultHTTPServer(t *testing.T) {
	mockLogger := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			fmt.Println(messages...)
		}
	}

	// Define test parameters
	port := 8080
	httpEndpoints := []api.Endpoint{
		{
			URL:        "/test",
			HTTPMethod: "GET",
			Middlewares: []api.Middleware{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(
						func(
							w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
							_, err := w.Write([]byte("Hello, World!"))
							assert.NoError(t, err)
						},
					)
				},
			},
		},
	}

	// Call DefaultHTTPServer to create the server
	server := DefaultHTTPServer(port, httpEndpoints, mockLogger, mockLogger)

	// Assert that the server is of type *http.Server
	httpServer, ok := server.(*http.Server)
	assert.True(t, ok, "Expected server to be of type *http.Server")
	assert.NotNil(t, httpServer, "Expected non-nil http.Server instance")

	// Check the server address (port)
	expectedAddr := fmt.Sprintf(":%d", port)
	assert.Equal(
		t,
		expectedAddr,
		httpServer.Addr,
		"Expected server to listen on the correct port",
	)

	// Check that the server handler is set correctly
	assert.NotNil(
		t,
		httpServer.Handler,
		"Expected server handler to be non-nil",
	)
}

// Test case to simulate a normal server operation and shutdown
func TestStartServer_NormalOperation(t *testing.T) {
	stopChan := make(chan os.Signal, 1)

	// Mock server to simulate normal operation
	mockServer := &MockServer{
		ListenAndServeFunc: func() error {
			time.Sleep(500 * time.Millisecond)
			return http.ErrServerClosed
		},
		ShutdownFunc: func(ctx context.Context) error {
			return nil
		},
	}

	go func() {
		err := startServer(stopChan, mockServer)
		assert.NoError(t, err)
	}()

	// Simulate sending a shutdown signal
	time.Sleep(200 * time.Millisecond)
	stopChan <- os.Interrupt

	// Give some time for graceful shutdown
	time.Sleep(1 * time.Second)
}

// Test case to cover the "Error starting HTTP server" path
func TestStartServer_ServerStartError(t *testing.T) {
	stopChan := make(chan os.Signal, 1)

	// Mock server to simulate an error during startup
	mockServer := &MockServer{
		ListenAndServeFunc: func() error {
			return errors.New("mock server start error")
		},
		ShutdownFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Capture log output
	log.SetFlags(0)
	var buf bytes.Buffer
	log.SetOutput(&buf)

	go func() {
		err := startServer(stopChan, mockServer)
		assert.EqualError(t, err, "mock server start error")
	}()

	// Give some time for graceful shutdown
	time.Sleep(1 * time.Second)
}

// Test case to cover the "Error during shutdown" path
func TestStartServer_ServerShutdownError(t *testing.T) {
	stopChan := make(chan os.Signal, 1)

	// Mock server to simulate normal startup but error during shutdown
	mockServer := &MockServer{
		ListenAndServeFunc: func() error {
			time.Sleep(500 * time.Millisecond) // Simulate running
			return http.ErrServerClosed        // Simulate graceful shutdown
		},
		ShutdownFunc: func(ctx context.Context) error {
			return errors.New("mock server shutdown error") // Simulate shutdown error
		},
	}

	// Start the server in a goroutine
	go func() {
		err := startServer(stopChan, mockServer)
		assert.Error(t, err) // We expect an error here due to simulated shutdown error
	}()

	// Simulate sending a shutdown signal
	time.Sleep(200 * time.Millisecond)
	stopChan <- os.Interrupt

	// Give some time for graceful shutdown
	time.Sleep(1 * time.Second)
}
func TestSetupMux(t *testing.T) {
	// Mock logger functions to capture log messages
	var infoLogMessages, errorLogMessages []string

	loggerInfoFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			infoLogMessages = append(infoLogMessages, fmt.Sprint(messages...))
		}
	}

	loggerErrorFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			errorLogMessages = append(errorLogMessages, fmt.Sprint(messages...))
		}
	}

	// Define test middleware and endpoints
	getMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	}
	postMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte("Created"))
			assert.Nil(t, err)
		})
	}
	endpoints := []api.Endpoint{
		{
			URL:         "/test",
			HTTPMethod:  "GET",
			Middlewares: []api.Middleware{getMiddleware},
		},
		{
			URL:         "/test",
			HTTPMethod:  "POST",
			Middlewares: []api.Middleware{postMiddleware},
		},
	}

	// Setup Mux
	mux := setupMux(endpoints, loggerInfoFn, loggerErrorFn)

	// Case 1: Test registered endpoint with allowed GET method
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Test-Middleware"))

	// Case 2: Test registered endpoint with allowed POST method
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "Created", recorder.Body.String())

	// Case 3: Test registered endpoint with not allowed method (PUT to /test)
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("PUT", "/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	assert.Contains(t, infoLogMessages[0], "Method not allowed")
}

// TestCreateEndpointHandler tests the createEndpointHandler function.
func TestCreateEndpointHandler(t *testing.T) {
	mockLogger := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			fmt.Println(messages...)
		}
	}

	// Mock handlers for different HTTP methods
	getHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "GET method")
		},
	)
	postHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "POST method")
		},
	)

	// Create a map of endpoints
	endpoints := map[string]http.Handler{
		"GET":  getHandler,
		"POST": postHandler,
	}

	// Create an instance of the handler
	handler := createEndpointHandler(endpoints, "/test", mockLogger, mockLogger)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{"Good GET", "GET", http.StatusOK, "GET method\n"},
		{"Good POST", "POST", http.StatusOK, "POST method\n"},
		{"Bad PUT", "PUT", http.StatusMethodNotAllowed, "Method Not Allowed\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req := httptest.NewRequest(tt.method, "/test", nil)
			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check the response body
			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					body, tt.expectedBody)
			}
		})
	}
}

// TestCreateNotFoundHandler tests the createNotFoundHandler function.
func TestCreateNotFoundHandler(t *testing.T) {
	mockLogger := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			fmt.Println(messages...)
		}
	}

	// Create an instance of the not found handler
	handler := createNotFoundHandler(mockLogger)

	// Define the test request
	req := httptest.NewRequest("GET", "/non-existent", nil)
	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusNotFound,
		)
	}

	// Check the response body
	expectedBody := "Not Found\n"
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			body,
			expectedBody,
		)
	}
}

// TestMapKeys tests the mapKeys function.
func TestMapKeys(t *testing.T) {
	// Case 1: Empty map
	emptyMap := map[string]int{}
	result := mapKeys(emptyMap)
	assert.Equal(t, 0, len(result))

	// Case 2: Map with keys
	testMap := map[string]int{"one": 1, "two": 2, "three": 3}
	result = mapKeys(testMap)
	expectedKeys := []string{"one", "two", "three"}
	assert.ElementsMatch(t, expectedKeys, result)
}

// TestMultiplexEndpoints tests the multiplexEndpoints function.
func TestMultiplexEndpoints(t *testing.T) {
	loggerErrorFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			t.Log("Error:", messages) // Logging for verification
		}
	}

	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	}

	// Create dummy endpoints with a middleware
	endpoints := []api.Endpoint{
		{
			URL:        "/test",
			HTTPMethod: "GET",
			Middlewares: []api.Middleware{
				testMiddleware,
			},
		},
	}

	muxEndpoints := multiplexEndpoints(endpoints, loggerErrorFn)

	// Test the middleware is applied
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	muxEndpoints["/test"]["GET"].ServeHTTP(recorder, request)

	assert.Equal(t, "true", recorder.Header().Get("X-Test-Middleware"))
	assert.Equal(t, http.StatusOK, recorder.Code)
}

// TestServerPanicHandler tests the serverPanicHandler function.
func TestServerPanicHandler(t *testing.T) {
	var loggedMessages []string

	// Mock loggerErrorFn to capture log messages
	loggerErrorFn := func(r *http.Request) func(messages ...any) {
		return func(messages ...any) {
			loggedMessages = append(loggedMessages, fmt.Sprint(messages...))
		}
	}

	// Case 1: Test normal execution without panic
	normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		assert.Nil(t, err)
	})

	handler := serverPanicHandler(normalHandler, loggerErrorFn)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "OK", recorder.Body.String())
	assert.Empty(t, loggedMessages) // No panic, no log should be captured

	// Case 2: Test execution with panic
	loggedMessages = []string{} // Reset log messages

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	handler = serverPanicHandler(panicHandler, loggerErrorFn)

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), http.StatusText(http.StatusInternalServerError))
	assert.Len(t, loggedMessages, 1)
	assert.Contains(t, loggedMessages[0], "Server panic")
	assert.Contains(t, loggedMessages[0], "test panic")
}

// TestStackTraceSlice tests the stackTraceSlice function.
func TestStackTraceSlice(t *testing.T) {
	// Call stackTraceSlice to get the stack trace
	stackTrace := stackTraceSlice()

	// Case 1: Check that the stack trace is not empty
	assert.NotEmpty(t, stackTrace)

	// Case 2: Verify each entry format (e.g., "filename:line functionName")
	for _, entry := range stackTrace {
		parts := strings.Split(entry, " ")
		assert.Equal(
			t,
			2,
			len(parts),
			"Entry should have two parts: 'filename:line' and 'functionName'",
		)

		fileLine := parts[0]
		funcName := parts[1]

		// Check if fileLine is in the format "filename:line"
		fileLineParts := strings.Split(fileLine, ":")
		assert.GreaterOrEqual(
			t,
			len(fileLineParts),
			2,
			"file:line should have at least two parts")
		assert.NotEmpty(t, fileLineParts[0], "Filename should not be empty")
		assert.NotEmpty(t, fileLineParts[1], "Line number should not be empty")

		// Check if funcName is not empty
		assert.NotEmpty(t, funcName, "Function name should not be empty")
	}

	// Case 3: Verify that the current test function appears in the stack trace
	found := false
	for _, entry := range stackTrace {
		if strings.Contains(entry, "TestStackTraceSlice") {
			found = true
			break
		}
	}
	assert.True(
		t,
		found,
		"'TestStackTraceSlice' should be present in the stack trace",
	)

	// Case 4: Check deep stack trace (call from another function)
	checkDeepStackTrace(t)
}

func checkDeepStackTrace(t *testing.T) {
	stackTrace := stackTraceSlice()

	// Verify that function names appear in the stack trace
	foundCheckFunc := false
	foundTestFunc := false

	for _, entry := range stackTrace {
		if strings.Contains(entry, "checkDeepStackTrace") {
			foundCheckFunc = true
		}
		if strings.Contains(entry, "TestStackTraceSlice") {
			foundTestFunc = true
		}
	}

	assert.True(
		t,
		foundCheckFunc,
		"Function 'checkDeepStackTrace' should be present in the stack trace",
	)
	assert.True(
		t,
		foundTestFunc,
		"Function 'TestStackTraceSlice' should be present in the stack trace",
	)
}
