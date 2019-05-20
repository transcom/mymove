package middleware

import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/gofrs/uuid"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
)

// RequestLogger returns a middleware that logs requests.
// The middleware trys to use the logger from the context.
// If the request context has no handler, then falls back to the server logger.
func RequestLogger(logger Logger) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

			if session := auth.SessionFromRequestContext(r); session != nil {
				if session.UserID != uuid.Nil {
					fields = append(fields, zap.String("user-id", session.UserID.String()))
				}
				if session.IsServiceMember() {
					fields = append(fields, zap.String("service-member-id", session.ServiceMemberID.String()))
				}
				if session.IsOfficeUser() {
					fields = append(fields, zap.String("office-user-id", session.OfficeUserID.String()))
				}
				if session.IsTspUser() {
					fields = append(fields, zap.String("tsp-user-id", session.TspUserID.String()))
				}
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
