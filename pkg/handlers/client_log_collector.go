package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/telemetry"
)

type LoggingTransport struct {
	f func(req *http.Request) (*http.Response, error)
}

func (lt *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return lt.f(req)
}

// NewClientCollectorHandler creates a handler for receiving client
// telemetry and forwarding it to the aws otel collector
func NewClientCollectorHandler(appCtx appcontext.AppContext, telemetryConfig *telemetry.Config) (http.Handler, error) {
	appCtx.Logger().Info("DREW DEBUG client collector telemetry",
		zap.String("http_endpoint", telemetryConfig.HTTPEndpoint),
	)
	telemetryURL, err := url.Parse(telemetryConfig.HTTPEndpoint)
	if err != nil {
		appCtx.Logger().Error("Cannot create client collector handler",
			zap.String("httpEndpoint", telemetryConfig.HTTPEndpoint),
			zap.Error(err),
		)
		return nil, err
	}
	appCtx.Logger().Info("DREW DEBUG client collector telemetry",
		zap.String("telemetryURL", telemetryURL.String()),
	)

	director := func(req *http.Request) {
		rAppCtx := appcontext.NewAppContextFromContext(req.Context(), appCtx)
		rAppCtx.Logger().Info("DREW DEBUG original request",
			zap.Any("req.Header", req.Header),
			zap.Any("req.RequestURI", req.RequestURI),
			zap.Any("req.URL", req.URL),
		)
		req.URL = telemetryURL
		req.RequestURI = telemetryURL.Path
		if req.RequestURI == "" {
			req.RequestURI = "/"
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		rAppCtx.Logger().Info("DREW DEBUG new request",
			zap.Any("req.Header", req.Header),
			zap.Any("req.RequestURI", req.RequestURI),
			zap.Any("req.URL", req.URL),
		)
	}

	defaultTransport := http.DefaultTransport
	transport := &LoggingTransport{
		f: func(req *http.Request) (*http.Response, error) {
			rAppCtx := appcontext.NewAppContextFromContext(req.Context(), appCtx)
			resp, err := defaultTransport.RoundTrip(req)
			var status string
			var statusCode int
			if err != nil && resp != nil {
				statusCode = resp.StatusCode
				status = resp.Status
			}
			rAppCtx.Logger().Info("DREW DEBUG roundtrip",
				zap.Any("req.Header", req.Header),
				zap.Any("req.RequestURI", req.RequestURI),
				zap.Any("req.URL", req.URL),
				zap.Any("resp.Status", status),
				zap.Any("resp.StatusCode", statusCode),
				zap.Error(err),
			)
			if statusCode == http.StatusOK {
				return resp, err
			}
			// if the collector redirection failed for any reason, no
			// need to report that to the client, so fake a response
			data := `{"partialSuccess":{}}`
			buf := bytes.NewBuffer([]byte(data))
			resp = &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(buf),
			}
			return resp, nil
		},
	}
	reverseProxy := httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
	rHandler := func(w http.ResponseWriter, r *http.Request) {
		rAppCtx := appcontext.NewAppContextFromContext(r.Context(), appCtx)
		rAppCtx.Logger().Info("DREW DEBUG client collector")
		reverseProxy.ServeHTTP(w, r)
	}

	return http.HandlerFunc(rHandler), nil
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	data, err := io.ReadAll(r.Body)
	// 	if err != nil {
	// 		appCtx.Logger().Error("client logs handler error", zap.Error(err))
	// 		return
	// 	}
	// 	appCtx.Logger().Info("DREW DEBUG CLIENT LOGGING",
	// 		zap.String("data", string(data)))
	// }
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
		appCtx.Logger().Info("client log upload stats",
			zap.String("source", "client_stats"),
			zap.String("app", logUpload.App),
			zap.Int("logEntryCount", len(logUpload.LogEntries)),
			zap.Int("droppedLogCount", logUpload.LoggerStats.DroppedLogsCount),
			zap.Int("failedSendCount", logUpload.LoggerStats.FailedSendCount),
			zap.Int("failedTimerCount", logUpload.LoggerStats.FailedTimerCount),
		)

		for i := range logUpload.LogEntries {
			logEntry := logUpload.LogEntries[i]
			appCtx.Logger().Info("client log entry",
				zap.String("source", "client_log_entry"),
				zap.String("app", logUpload.App),
				zap.String("logLevel", logEntry.Level),
				zap.Any("args", logEntry.Args),
			)
		}
	}
}
