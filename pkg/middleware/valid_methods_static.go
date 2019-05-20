package middleware

import (
	"net/http"
)

// ValidMethodsStatic only lets GET AND HEAD requests for static resources.
func ValidMethodsStatic(logger Logger) func(inner http.Handler) http.Handler {
	logger.Debug("ValidMethodsStatic Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			inner.ServeHTTP(w, r)
			return
		})
	}
}
