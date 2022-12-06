package authentication

import (
	"net/http"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/audit"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/user"
)

func AddAuditUserToRequestContextMiddleware(appCtx appcontext.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			newAppCtx := appcontext.NewAppContextFromContext(r.Context(), appCtx)
			reqCtx := r.Context()
			if newAppCtx.Session() != nil {
				auditUser, err := user.NewUserFetcher(query.NewQueryBuilder()).FetchUser(newAppCtx, []services.QueryFilter{query.NewQueryFilter("id", "=", newAppCtx.Session().UserID)})
				if err != nil {
					newAppCtx.Logger().Info("No user attached to the session")
					http.Error(w, http.StatusText(500), http.StatusInternalServerError)
					return
				}
				reqCtx = audit.WithAuditUser(r.Context(), auditUser)
			}

			next.ServeHTTP(w, r.WithContext(reqCtx))
		}
		return http.HandlerFunc(mw)
	}
}
