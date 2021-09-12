package middleware

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// OpenAPIWithContext descripes an API that implements go-openapi Context method
type OpenAPIWithContext interface {
	Context() *middleware.Context
}

// OpenAPITracing instruments a context with the path pattern for the
// request. This is useful for recording metrics
func OpenAPITracing(api OpenAPIWithContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			matchedRoute, _, found := api.Context().RouteInfo(r)
			if found {
				span := trace.SpanFromContext(r.Context())
				span.SetAttributes(semconv.HTTPTargetKey.String(matchedRoute.PathPattern))
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}
