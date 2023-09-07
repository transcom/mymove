package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

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
