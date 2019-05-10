package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

// Recovery recovers from a panic within a handler.
func Recovery(logger Logger) func(inner http.Handler) http.Handler {
	logger.Debug("Recovery Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if obj := recover(); obj != nil {
					w.WriteHeader(http.StatusInternalServerError) // don't write error to body, since body might have already been written.
					fields := []zap.Field{
						zap.String("url", fmt.Sprint(r.URL)),
					}
					if err, ok := obj.(error); ok {
						fields = append(fields, zap.Error(err))
					} else {
						fields = append(fields, zap.Any("object", obj))
						zap.String("stacktrace", fmt.Sprintf("%s", debug.Stack()))
					}
					logger.Error("http request panic", fields...)
				}
			}()
			inner.ServeHTTP(w, r)
		})
	}
}
