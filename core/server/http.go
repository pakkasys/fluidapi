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

	"github.com/PakkaSys/fluidapi/core/api"
)

type multiplexedEndpoints map[string]map[string]http.Handler

func HTTPServer(
	port int,
	httpEndpoints []api.Endpoint,
	loggerInfoFn func(r *http.Request) func(messages ...any),
	loggerErrorFn func(r *http.Request) func(messages ...any),
) error {
	mux := http.NewServeMux()
	endpoints := multiplexEndpoints(httpEndpoints, loggerErrorFn)

	for url := range endpoints {
		iterUrl := url
		log.Printf("Registering URL: %s %v", url, mapKeys(endpoints[iterUrl]))
		mux.Handle(iterUrl, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if handler, ok := endpoints[iterUrl][r.Method]; ok {
					handler.ServeHTTP(w, r)
					return
				}
				if loggerErrorFn != nil {
					loggerInfoFn(r)(fmt.Sprintf(
						"Method not allowed: %q, URL: %s",
						r.Method,
						iterUrl,
					))
				}
				http.Error(
					w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
			}),
		)
	}

	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			loggerInfoFn(r)(fmt.Sprintf("Not found: %s", r.URL))
			http.Error(
				w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)
		},
	))

	return startServer(mux, port)
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func startServer(mux *http.ServeMux, port int) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Listen for shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	// Block until a signal is received
	<-stopChan
	log.Printf("Shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	log.Printf("HTTP server shutdown")
	return nil
}

func multiplexEndpoints(
	httpEndpoints []api.Endpoint,
	loggerErrorFn func(r *http.Request) func(messages ...any),
) multiplexedEndpoints {
	endpoints := multiplexedEndpoints{}

	for i := range httpEndpoints {
		url := httpEndpoints[i].URL
		method := httpEndpoints[i].HTTPMethod

		if endpoints[url] == nil {
			endpoints[url] = make(map[string]http.Handler)
		}

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
				if loggerErrorFn == nil {
					fmt.Fprintf(
						os.Stderr,
						"Server Panic: %v, %v",
						err,
						stackTraceSlice(),
					)
				} else {
					loggerErrorFn(r)(
						"Server Panic",
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
