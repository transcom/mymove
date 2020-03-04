package middleware

import (
	"context"
	"net/http"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/primerequester"
)

// PrimeRequester returns a Prime Requester (user) middleware.
func PrimeRequester(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pull requester info from context
			requester := PrimeRequesterFromContext(r.Context())
			if requester == nil {
				logger.Error("unauthorized user for ghc prime")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			// Store requester access to system along with the client cert used
			clientCert := authentication.ClientCertFromContext(r.Context())
			if clientCert == nil {
				logger.Error("unauthorized user for ghc prime")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			primerequester.SaveAccess(requester, clientCert)
			next.ServeHTTP(w, r)
		})
	}
}

// PrimeRequesterFromContext gets the Prime API requester field stored in the request.Context()
func PrimeRequesterFromContext(ctx context.Context) *string {
	if requester, ok := ctx.Value("requester").(*string); ok {
		return requester
	}
	return nil
}
