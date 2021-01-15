package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/trace"
)

// Recovery recovers from a panic within a handler.
func Recovery(logger Logger) func(inner http.Handler) http.Handler {
	logger.Debug("Recovery Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if obj := recover(); obj != nil {

					// Log the error and optionally the stacktrace
					fields := []zap.Field{
						zap.String("url", fmt.Sprint(r.URL)),
					}
					traceID := trace.FromContext(r.Context())
					if traceID != "" {
						fields = append(fields, zap.String("milmove_trace_id", traceID))
					}
					if err, ok := obj.(error); ok {
						fields = append(fields, zap.Error(err))
					} else {
						fields = append(fields, zap.Any("object", obj))
						fields = append(fields, zap.String("stacktrace", fmt.Sprintf("%s", debug.Stack())))
					}
					logger.Error("http request panic", fields...)

					// Create a formatted server error
					jsonBody, _ := json.Marshal(struct {
						Title    string `json:"title"`
						Instance string `json:"instance"`
						Detail   string `json:"detail"`
					}{handlers.InternalServerErrMessage, traceID, "An unexpected server error has occurred."})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					//RA Summary: gosec - errcheck - Unchecked return value
					//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
					//RA: Function with unchecked return value in the file is used to write a response to the client
					//RA: Due to the nature of writing to the client being its sole function, any unexpected states and conditons would be inherently handled by the httpResponseWriter
					//RA Developer Status: Mitigated
					//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
					//RA Validator: jneuner@mitre.org
					//RA Modified Severity:
					w.Write(jsonBody) // nolint:errcheck
				}
			}()
			inner.ServeHTTP(w, r)
		})
	}
}
