package middleware

import (
	"net/http"

	sdktrace "go.opentelemetry.io/otel/trace"

	"github.com/go-openapi/runtime/middleware"
	"github.com/trussworks/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/transcom/mymove/pkg/audit"
	"github.com/transcom/mymove/pkg/logging"
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
				logger := logging.FromContext(r.Context())
				labeler, ok := otelhttp.LabelerFromContext(r.Context())
				if !ok {
					logger.Warn("Cannot get labeler from context")
				} else {
					labeler.Add(semconv.HTTPTargetKey.String(matchedRoute.PathPattern))
					sdktrace.SpanFromContext(r.Context()).SetAttributes(semconv.HTTPRouteKey.String(matchedRoute.PathPattern))
				}

				// save the swagger operationId
				eventNameCtx := r.WithContext(audit.WithEventName(r.Context(),
					matchedRoute.Operation.ID))
				next.ServeHTTP(w, eventNameCtx)
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(mw)
	}
}
