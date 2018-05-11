package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	size, err := rec.ResponseWriter.Write(b)
	rec.size += size
	return size, err
}

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
func LogRequestMiddleware(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		var protocol string
		start := time.Now()

		if r.TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}

		rec := responseRecorder{w, 200, 0}
		inner.ServeHTTP(&rec, r)
		zap.L().Info("Request",
			zap.String("accepted-language", r.Header.Get("accepted-language")),
			zap.Int64("content-length", r.ContentLength),
			zap.Float64("duration-ms", float64(time.Since(start))/float64(time.Millisecond)),
			zap.String("host", r.Host),
			zap.String("method", r.Method),
			zap.String("protocol", protocol),
			zap.String("protocol-version", r.Proto),
			zap.String("referer", r.Header.Get("referer")),
			zap.Int("resp-size-bytes", rec.size),
			zap.Int("resp-status", rec.status),
			zap.String("source", r.RemoteAddr),
			zap.String("url", r.URL.String()),
			zap.String("user-agent", r.UserAgent()),
			zap.String("x-amzn-trace-id", r.Header.Get("x-amzn-trace-id")),
			zap.String("x-forwarded-for", r.Header.Get("x-forwarded-for")),
			zap.String("x-forwarded-host", r.Header.Get("x-forwarded-host")),
			zap.String("x-forwarded-proto", r.Header.Get("x-forwarded-proto")),
		)

	}
	return http.HandlerFunc(mw)
}
