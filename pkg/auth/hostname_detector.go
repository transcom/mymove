package auth

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// HostnameDetectorMiddleware only allows the given hostname through
func HostnameDetectorMiddleware(logger Logger, hostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating hostname detector", zap.String("hostname", hostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Host, ":")
			if !strings.EqualFold(parts[0], hostname) {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		return http.HandlerFunc(mw)
	}
}
