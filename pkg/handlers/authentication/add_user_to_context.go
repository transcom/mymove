package authentication

import (
	"net/http"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/audit"
)

func AddAuditUserIDToRequestContextMiddleware(appCtx appcontext.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			newAppCtx := appcontext.NewAppContextFromContext(r.Context(), appCtx)
			reqCtx := r.Context()
			if newAppCtx.Session() != nil {
				reqCtx = audit.WithAuditUserID(r.Context(), newAppCtx.Session().UserID)
			}
			next.ServeHTTP(w, r.WithContext(reqCtx))
		}
		return http.HandlerFunc(mw)
	}
}
