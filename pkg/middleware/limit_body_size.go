package middleware

import (
	"net/http"
)

// LimitBodySize is a middleware
func LimitBodySize(maxBodySize int64, logger Logger) func(inner http.Handler) http.Handler {
	logger.Debug("LimitBodySize Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
			inner.ServeHTTP(w, r)
		})
	}
}
