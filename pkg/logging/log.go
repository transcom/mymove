package logging

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/server"
)

// Config contains the environment name and debug logging flag for configuring zap.Logging
type Config struct {
	DebugLogging bool
}

// NewLogger is the DI provider for constructing a new zap.Logger
func NewLogger(config *Config, env *server.LocalEnvConfig) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if env.Environment != "development" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if config.DebugLogging {
		debug := zap.NewAtomicLevel()
		debug.SetLevel(zap.DebugLevel)
		loggerConfig.Level = debug
	}
	return loggerConfig.Build()
}

// LogRequestMiddleware is a type marker for DI
type LogRequestMiddleware func(http.Handler) http.Handler

// NewLogRequestMiddleware generates an piece of MW which logs HTTP/HTTPS request using the passed in logger
func NewLogRequestMiddleware(l *zap.Logger) LogRequestMiddleware {
	return func(inner http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			var protocol, officeUserID, serviceMemberID, userID string

			if r.TLS == nil {
				protocol = "http"
			} else {
				protocol = "https"
			}

			session := auth.SessionFromRequestContext(r)
			if session.UserID != uuid.Nil {
				userID = session.UserID.String()
			}
			if session.IsServiceMember() {
				serviceMemberID = session.ServiceMemberID.String()
			}

			if session.IsOfficeUser() {
				officeUserID = session.OfficeUserID.String()
			}

			metrics := httpsnoop.CaptureMetrics(inner, w, r)
			l.Info("Request",
				zap.String("accepted-language", r.Header.Get("accepted-language")),
				zap.Int64("content-length", r.ContentLength),
				zap.Duration("duration", metrics.Duration),
				zap.String("host", r.Host),
				zap.String("method", r.Method),
				zap.String("office-user-id", officeUserID),
				zap.String("protocol", protocol),
				zap.String("protocol-version", r.Proto),
				zap.String("referer", r.Header.Get("referer")),
				zap.Int64("resp-size-bytes", metrics.Written),
				zap.Int("resp-status", metrics.Code),
				zap.String("service-member-id", serviceMemberID),
				zap.String("source", r.RemoteAddr),
				zap.String("url", r.URL.String()),
				zap.String("user-agent", r.UserAgent()),
				zap.String("user-id", userID),
				zap.String("x-amzn-trace-id", r.Header.Get("x-amzn-trace-id")),
				zap.String("x-forwarded-for", r.Header.Get("x-forwarded-for")),
				zap.String("x-forwarded-host", r.Header.Get("x-forwarded-host")),
				zap.String("x-forwarded-proto", r.Header.Get("x-forwarded-proto")),
			)

		}
		return http.HandlerFunc(mw)
	}
}
