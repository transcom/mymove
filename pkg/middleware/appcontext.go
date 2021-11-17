package middleware

import (
	"net/http"

	"github.com/transcom/mymove/pkg/appcontext"
)

// AppContextMiddleware returns a handler that inject the AppContext into the request context
func AppContextMiddleware(appCtx appcontext.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			newCxt := appcontext.NewContext(ctx, appCtx)

			next.ServeHTTP(w, r.WithContext(newCxt))
		})
	}
}
