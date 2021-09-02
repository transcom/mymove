package auth

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
)

// HostnameDetectorMiddleware only allows the given hostname through
func HostnameDetectorMiddleware(globalLogger *zap.Logger, hostname string) func(next http.Handler) http.Handler {
	globalLogger.Info("Creating hostname detector", zap.String("hostname", hostname))
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logging.FromContext(ctx)
			parts := strings.Split(r.Host, ":")
			if !strings.EqualFold(parts[0], hostname) {
				logger.Error("Bad hostname", zap.String("hostname", r.Host))
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}
