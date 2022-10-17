package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/trace"
)

// Recovery recovers from a panic within a handler.
func Recovery(globalLogger *zap.Logger) func(inner http.Handler) http.Handler {
	globalLogger.Debug("Recovery Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			defer func() {
				if obj := recover(); obj != nil {

					// Log the error and optionally the stacktrace
					fields := []zap.Field{
						zap.String("url", fmt.Sprint(r.URL)),
					}
					traceID := trace.FromContext(r.Context())
					if !traceID.IsNil() {
						fields = append(fields, zap.String("milmove_trace_id", traceID.String()))
					}
					if err, ok := obj.(error); ok {
						fields = append(fields, zap.Error(err))
						fields = append(fields, zap.String("stacktrace", string(debug.Stack())))
					} else {
						fields = append(fields, zap.Any("object", obj))
						fields = append(fields, zap.String("stacktrace", string(debug.Stack())))
					}
					logger.Error("http request panic", fields...)

					// Create a formatted server error
					jsonBody, _ := json.Marshal(struct {
						Title    string `json:"title"`
						Instance string `json:"instance"`
						Detail   string `json:"detail"`
					}{handlers.InternalServerErrMessage, traceID.String(), "An unexpected server error has occurred."})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					_, err := w.Write(jsonBody)
					if err != nil {
						logger.Error("Failed to write json error to the response body", zap.Error(err))
					}
				}
			}()
			inner.ServeHTTP(w, r)
		})
	}
}
