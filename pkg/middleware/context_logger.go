package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/trace"
)

// ContextLogger returns a handler that injects the logger into the request context.
func ContextLogger(field string, original *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logs := original.With(zap.String("host", r.Host))
			if id := trace.FromContext(ctx); !id.IsNil() {
				logs = logs.With(zap.String(field, id.String()))
			}
			ctx = logging.NewContext(ctx, logs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
