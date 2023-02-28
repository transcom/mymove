package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/telemetry"
)

type LoggingTransport struct {
	f func(req *http.Request) (*http.Response, error)
}

func (lt *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return lt.f(req)
}

// NewClientTelemetryHandler creates a handler for receiving client
// telemetry and forwarding it to the aws otel collector
func NewClientTelemetryHandler(appCtx appcontext.AppContext, telemetryConfig *telemetry.Config) (http.Handler, error) {
	telemetryURL, err := url.Parse(telemetryConfig.HTTPEndpoint)
	if err != nil {
		appCtx.Logger().Error("Cannot create client collector handler",
			zap.String("httpEndpoint", telemetryConfig.HTTPEndpoint),
			zap.Error(err),
		)
		return nil, err
	}

	director := func(req *http.Request) {
		req.URL = telemetryURL
		req.RequestURI = telemetryURL.Path
		if req.RequestURI == "" {
			req.RequestURI = "/"
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	defaultTransport := http.DefaultTransport
	transport := &LoggingTransport{
		f: func(req *http.Request) (*http.Response, error) {
			return defaultTransport.RoundTrip(req)
		},
	}
	reverseProxy := httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
	rHandler := func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)
	}

	return http.HandlerFunc(rHandler), nil
}

type ClientLoggerStats struct {
	DroppedLogsCount int `json:"droppedLogsCount"`
	FailedSendCount  int `json:"failedSendCount"`
	FailedTimerCount int `json:"failedTimerCount"`
}

type ClientLogMessage []interface{}

type ClientLogEntry struct {
	Level string           `json:"level"`
	Args  ClientLogMessage `json:"args"`
}

type ClientLogUpload struct {
	App         string            `json:"app"`
	LoggerStats ClientLoggerStats `json:"loggerStats"`
	LogEntries  []ClientLogEntry  `json:"logEntries"`
}

// NewClientLogHandler creates a handler for receiving client logs
func NewClientLogHandler(appCtx appcontext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			appCtx.Logger().Error("client logs handler error", zap.Error(err))
			return
		}
		var logUpload ClientLogUpload
		err = json.Unmarshal(data, &logUpload)
		if err != nil {
			appCtx.Logger().Error("Error unmarshalling ClientLogUpload", zap.Error(err))
			// really no need to tell the client the server had
			// problems in this case, it can't do anything about it
			return
		}
		// use the appCtx logger and not the one associated with the
		// request so that client logs have only the attributes
		// configured here
		clientLogger := appCtx.Logger().With(
			zap.String("app", logUpload.App),
		)

		if sessionID := auth.SessionIDFromContext(r.Context()); sessionID != "" {
			clientLogger = clientLogger.With(zap.String("session_id", sessionID))
		}

		clientLogger.Info("client log upload stats",
			zap.String("source", "client_stats"),
			zap.Int("logEntryCount", len(logUpload.LogEntries)),
			zap.Int("droppedLogCount", logUpload.LoggerStats.DroppedLogsCount),
			zap.Int("failedSendCount", logUpload.LoggerStats.FailedSendCount),
			zap.Int("failedTimerCount", logUpload.LoggerStats.FailedTimerCount),
		)

		for i := range logUpload.LogEntries {
			logEntry := logUpload.LogEntries[i]
			clientLogger.Info("client log entry",
				zap.String("source", "client_log_entry"),
				zap.String("logLevel", logEntry.Level),
				zap.Any("args", logEntry.Args),
			)
		}
	}
}
