package middleware

import (
	"net/http"
	"strings"

	"github.com/transcom/mymove/pkg/appcontext"
	apiversion "github.com/transcom/mymove/pkg/handlers/routing/api_version"
)

// APIVersionFlagger returns a handler that sets the API version flag in the request context.
func APIVersionFlagger(appCtx appcontext.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			apiVersion := apiversion.NoneSpecified
			if strings.Contains(r.URL.Path, "/prime/v1/") {
				apiVersion = apiversion.PrimeVersion1
			}
			if strings.Contains(r.URL.Path, "/prime/v2/") {
				apiVersion = apiversion.PrimeVersion2
			}
			context := apiversion.WithAPIVersion(r.Context(), apiVersion)
			next.ServeHTTP(w, r.WithContext(context))
		}
		return http.HandlerFunc(mw)
	}
}
