package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/PakkaSys/fluidapi/core/api"
)

const (
	CORSMiddlewareID = "cors"

	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowCredentials = "Access-Control-Allow-Credentials"

	originHeader = "Origin"
)

var (
	corsAllowHeaders = []string{"Content-Type"}
)

func CORSMiddlewareWrapper(
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
) *api.MiddlewareWrapper {
	return api.NewMiddlewareWrapperBuilder().
		ID(CORSMiddlewareID).
		Middleware(
			CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders),
		).
		Build()
}

func CORSMiddleware(
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
) api.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(originHeader)
			if slices.Contains(allowedOrigins, origin) {
				w.Header().Set(headerAllowOrigin, origin)
			}

			w.Header().Set(
				headerAllowMethods,
				strings.Join(allowedMethods, ","),
			)

			w.Header().Set(
				headerAllowHeaders,
				strings.Join(
					slices.Concat(corsAllowHeaders, allowedHeaders), ",",
				),
			)

			w.Header().Set(headerAllowCredentials, "true")

			next.ServeHTTP(w, r)
		})
	}
}
