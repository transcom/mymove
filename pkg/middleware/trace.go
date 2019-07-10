package middleware

import (
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/trace"
)

// Trace returns a trace middleware that injects a unique trace id into every request.
func Trace(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.NewV4()
			if err != nil {
				logger.Error(errors.Wrap(err, "error creating trace id").Error())
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r.WithContext(trace.NewContext(r.Context(), id.String())))
			}
		})
	}
}
