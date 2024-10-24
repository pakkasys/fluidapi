package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/pakkasys/fluidapi/core/api"
)

type multiplexedEndpoints map[string]map[string]http.Handler

// IServer represents an HTTP server.
type IServer interface {
	ListenAndServe() error              // Start the server
	Shutdown(ctx context.Context) error // Stop the server
}

// HTTPServer sets up an HTTP server with the specified port and endpoints,
// using optional logging functions for requests and errors. If no custom server
// options are provided, it creates a default http.Server. The server listens
// for OS interrupt signals to gracefully shut down.
//
//   - server: Server implementation to use.
func HTTPServer(server IServer) error {
	return startServer(make(chan os.Signal, 1), server)
}

type LoggerFn func(r *http.Request) func(messages ...any)

// HTTPServer returns the default HTTP server implementation.
//
//   - port: Port for the HTTP server.
//   - httpEndpoints: Endpoints to register.
//   - loggerInfoFn: Function to log informational messages.
//   - loggerErrorFn: Function to log error messages.
func DefaultHTTPServer(
	port int,
	httpEndpoints []api.Endpoint,
	loggerInfoFn LoggerFn,
	loggerErrorFn LoggerFn,
) IServer {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: setupMux(httpEndpoints, loggerInfoFn, loggerErrorFn),
	}
}

func startServer(stopChan chan os.Signal, server IServer) error {
	// Listen for shutdown signals
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Capture the error from ListenAndServe
	errChan := make(chan error, 1)

	go func() {
		log.Printf("Starting HTTP server")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting HTTP server: %v", err)
			errChan <- err
			stopChan <- os.Interrupt
		} else {
			errChan <- nil
		}
	}()

	// Wait for a signal to shut down
	<-stopChan
	log.Printf("Shutting down HTTP server")

	// Give the server some time to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	log.Printf("HTTP server shutdown")
	return <-errChan
}

func setupMux(
	httpEndpoints []api.Endpoint,
	loggerInfoFn func(r *http.Request) func(messages ...any),
	loggerErrorFn func(r *http.Request) func(messages ...any),
) *http.ServeMux {
	mux := http.NewServeMux()
	endpoints := multiplexEndpoints(httpEndpoints, loggerErrorFn)

	for url := range endpoints {
		log.Printf("Registering URL: %s %v", url, mapKeys(endpoints[url]))
		iterUrl := url
		mux.Handle(
			iterUrl,
			createEndpointHandler(
				endpoints[iterUrl],
				loggerInfoFn,
				loggerErrorFn,
			),
		)
	}

	mux.Handle("/", createNotFoundHandler(loggerInfoFn))

	return mux
}

func createEndpointHandler(
	endpoints map[string]http.Handler,
	loggerInfoFn func(r *http.Request) func(messages ...any),
	loggerErrorFn func(r *http.Request) func(messages ...any),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := endpoints[r.Method]; ok {
			handler.ServeHTTP(w, r)
			return
		}
		if loggerErrorFn != nil {
			loggerInfoFn(r)(fmt.Sprintf(
				"Method not allowed: %s (%v)",
				r.URL,
				r.Method),
			)
		}
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func createNotFoundHandler(
	loggerInfoFn func(r *http.Request) func(messages ...any),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggerInfoFn(r)(fmt.Sprintf("Not found: %s (%v)", r.URL, r.Method))
		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func multiplexEndpoints(
	httpEndpoints []api.Endpoint,
	loggerErrorFn func(r *http.Request) func(messages ...any),
) multiplexedEndpoints {
	endpoints := multiplexedEndpoints{}
	for i := range httpEndpoints {
		url := httpEndpoints[i].URL
		method := httpEndpoints[i].Method
		if endpoints[url] == nil {
			endpoints[url] = make(map[string]http.Handler)
		}
		// Include panic handler with other middlewares
		endpoints[url][method] = serverPanicHandler(
			api.ApplyMiddlewares(
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
					},
				),
				httpEndpoints[i].Middlewares...,
			),
			loggerErrorFn,
		)
	}
	return endpoints
}

func serverPanicHandler(
	next http.Handler,
	loggerErrorFn func(r *http.Request) func(messages ...any),
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if loggerErrorFn != nil {
					loggerErrorFn(r)(
						"Server panic",
						fmt.Sprintf("%v, %v", err, stackTraceSlice()),
					)
				}
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func stackTraceSlice() []string {
	var stackTrace []string
	var skip int
	for {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		// Get the function name and format entry.
		fn := runtime.FuncForPC(pc)
		entry := fmt.Sprintf("%s:%d %s", file, line, fn.Name())
		stackTrace = append(stackTrace, entry)

		skip++
	}
	return stackTrace
}
