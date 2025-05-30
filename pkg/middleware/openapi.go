package middleware

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/audit"
	"github.com/transcom/mymove/pkg/telemetry"
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
				// See telemetry.NewOtelHTTPMiddleware for more explanation
				routePattern := telemetry.RoutePatternFromContext(r.Context())
				if routePattern != nil {
					*routePattern = matchedRoute.PathPattern
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
