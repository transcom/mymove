package middleware

import (
	"net/http"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/trace"
)

const traceHeader = "X-MILMOVE-TRACE-ID"

// Trace returns a trace middleware that injects a unique trace id into every request.
func Trace(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			id, err := uuid.NewV4()
			if err != nil {
				logger.Error(errors.Wrap(err, "error creating trace id").Error())
				next.ServeHTTP(w, r)
				return
			}

			// Let a caller see what the traceID is
			w.Header().Add(traceHeader, id.String())

			// Also insert as a key, value pair in the http request context
			next.ServeHTTP(w, r.WithContext(trace.NewContext(r.Context(), id.String())))
		})
	}
}
