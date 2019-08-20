package middleware

import (
	"net/http"
)

// NoCache sets the Cache-Cotnrol header so this route is never cached by clients.
func NoCache(logger Logger) func(inner http.Handler) http.Handler {
	logger.Debug("NoCache Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			inner.ServeHTTP(w, r)
		})
	}
}
