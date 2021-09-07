package middleware

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// OpenAPIWithContext descripes an API that implements go-openapi Context method
type OpenAPIWithContext interface {
	Context() *middleware.Context
}

type openApipathPatternContextKey string

var pathPatternContextKey = openApipathPatternContextKey("pathPattern")

// OpenAPITracing instruments a context with the path pattern for the
// request. This is useful for recording metrics
func OpenAPITracing(api OpenAPIWithContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			matchedRoute, _, found := api.Context().RouteInfo(r)
			if found {
				pctx := context.WithValue(r.Context(), pathPatternContextKey,
					matchedRoute.PathPattern)
				next.ServeHTTP(w, r.WithContext(pctx))
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(mw)
	}
}

// PathPatternFromContext retrieves the pattern from the Context
func PathPatternFromContext(ctx context.Context) string {
	pathPattern, ok := ctx.Value(pathPatternContextKey).(string)
	if !ok {
		return ""
	}
	return pathPattern
}
