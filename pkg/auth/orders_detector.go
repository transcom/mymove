package auth

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// OrdersDetectorMiddleware only allows ordersHost through
func OrdersDetectorMiddleware(logger *zap.Logger, ordersHostname string) func(next http.Handler) http.Handler {
	logger.Info("Creating orders detector", zap.String("ordersHost", ordersHostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Host, ":")
			if !strings.EqualFold(parts[0], ordersHostname) {
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
