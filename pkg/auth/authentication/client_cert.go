package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gobuffalo/pop"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

type authClientCertKey string

const clientCertContextKey authClientCertKey = "clientCert"

// SetClientCertInRequestContext returns a copy of the request's Context() with the client certificate data
func SetClientCertInRequestContext(r *http.Request, clientCert *models.ClientCert) context.Context {
	return context.WithValue(r.Context(), clientCertContextKey, clientCert)
}

// ClientCertFromRequestContext gets the reference to the ClientCert stored in the request.Context()
func ClientCertFromRequestContext(r *http.Request) *models.ClientCert {
	if clientCert, ok := r.Context().Value(clientCertContextKey).(*models.ClientCert); ok {
		return clientCert
	}
	return nil
}

// ClientCertMiddleware enforces that the incoming request includes a known client certificate, and stores the fetched permissions in the session
func ClientCertMiddleware(logger *zap.Logger, db *pop.Connection) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := beeline.StartSpan(r.Context(), "ClientCertMiddleware")
			defer span.Send()

			if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
				logger.Info("Unauthenticated")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			// get DER hash
			hash := sha256.Sum256(r.TLS.PeerCertificates[0].Raw)
			hashString := hex.EncodeToString(hash[:])

			clientCert, err := models.FetchClientCert(db, hashString)
			if err != nil {
				// This is not a known client certificate at all
				logger.Info("Unknown / unregistered client certificate")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			ctx = SetClientCertInRequestContext(r, clientCert)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		return http.HandlerFunc(mw)
	}
}
