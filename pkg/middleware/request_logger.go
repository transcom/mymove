package middleware

import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/gofrs/uuid"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/logging"
)

// RequestLogger returns a middleware that logs requests.
func RequestLogger(globalLogger *zap.Logger) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			logger := logging.FromContext(ctx)

			fields := []zap.Field{
				zap.String("accepted-language", r.Header.Get("accepted-language")),
				zap.Int64("content-length", r.ContentLength),
				zap.String("host", r.Host),
				zap.String("method", r.Method),
				zap.String("protocol-version", r.Proto),
				zap.String("referer", r.Header.Get("referer")),
				zap.String("source", r.RemoteAddr),
				zap.String("url", r.URL.String()),
				zap.String("user-agent", r.UserAgent()),
			}

			if r.TLS == nil {
				fields = append(fields, zap.String("protocol", "http"))
			} else {
				fields = append(fields, zap.String("protocol", "https"))
			}

			// Append x- headers, e.g., x-forwarded-for.
			for name, values := range r.Header {
				if nameLowerCase := strings.ToLower(name); strings.HasPrefix(nameLowerCase, "x-") {
					if len(values) > 0 {
						fields = append(fields, zap.String(nameLowerCase, values[0]))
					}
				}
			}

			// Log the number of headers, which can be used for finding abnormal requests
			fields = append(fields, zap.Int("headers", len(r.Header)))

			if session := auth.SessionFromContext(ctx); session != nil {
				if session.UserID != uuid.Nil {
					fields = append(fields, zap.String("user-id", session.UserID.String()))
				}
				if session.IsServiceMember() {
					fields = append(fields, zap.String("service-member-id", session.ServiceMemberID.String()))
				}
				if session.IsOfficeUser() {
					fields = append(fields, zap.String("office-user-id", session.OfficeUserID.String()))
				}
			} else if clientCert := authentication.ClientCertFromContext(ctx); clientCert != nil {
				fields = append(fields, zap.String("client-cert-id", clientCert.ID.String()))
			}

			metrics := httpsnoop.CaptureMetrics(inner, w, r)

			fields = append(fields, []zap.Field{
				zap.Duration("duration", metrics.Duration),
				zap.Int64("resp-size-bytes", metrics.Written),
				zap.Int("resp-status", metrics.Code),
			}...)

			logger.Info("Request", fields...)

		})
	}
}
