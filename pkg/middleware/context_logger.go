package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/trace"
)

// ContextLogger returns a handler that injects the logger into the request context.
func ContextLogger(field string, original Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if id := trace.FromContext(ctx); len(id) > 0 {
				next.ServeHTTP(w, r.WithContext(logging.NewContext(ctx, original.With(zap.String(field, id)))))
			} else {
				next.ServeHTTP(w, r.WithContext(logging.NewContext(ctx, original)))
			}
		})
	}
}
