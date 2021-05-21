package middleware

import (
	"net/http"

	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/trace"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/trace"
)

const traceHeader = "X-MILMOVE-TRACE-ID"

// Trace returns a trace middleware that injects a unique trace id into every request.
func Trace(logger Logger, handlerContext interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.NewV4()
			if err != nil {
				logger.Error(errors.Wrap(err, "error creating trace id").Error())
				next.ServeHTTP(w, r)
				return
			}

			// Set traceID in the handlerContext
			context, ok := handlerContext.(*handlers.HandlerContext)
			if ok {
				(*context).SetTraceID(id)
			}

			strID := id.String()

			// decorate the span with the id
			sdktrace.SpanFromContext(r.Context()).SetAttributes(attribute.String(traceHeader, strID))

			// Let a caller see what the traceID is
			w.Header().Add(traceHeader, strID)

			// Also insert as a key, value pair in the http request context
			next.ServeHTTP(w, r.WithContext(trace.NewContext(r.Context(), strID)))
		})
	}
}
