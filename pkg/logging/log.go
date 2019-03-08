package logging

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/auth"
)

// Config configures a Zap logger based on the environment string and debugLevel
func Config(env string, debugLogging bool) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if env != "development" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if debugLogging {
		debug := zap.NewAtomicLevel()
		debug.SetLevel(zap.DebugLevel)
		loggerConfig.Level = debug
	}

	return loggerConfig.Build()
}

// LogRequestMiddleware generates an HTTP/HTTPS request logs using Zap
func LogRequestMiddleware(gitBranch string, gitCommit string) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			var protocol, tspUserID, officeUserID, serviceMemberID, userID string

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

			if session.IsTspUser() {
				tspUserID = session.TspUserID.String()
			}

			metrics := httpsnoop.CaptureMetrics(inner, w, r)
			zap.L().Info("Request",
				zap.String("git-branch", gitBranch),
				zap.String("git-commit", gitCommit),
				zap.String("accepted-language", r.Header.Get("accepted-language")),
				zap.Int64("content-length", r.ContentLength),
				zap.Duration("duration", metrics.Duration),
				zap.String("host", r.Host),
				zap.String("method", r.Method),
				zap.String("office-user-id", officeUserID),
				zap.String("tsp-user-id", tspUserID),
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
