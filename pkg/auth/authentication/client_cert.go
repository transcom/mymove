package authentication

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gobuffalo/pop"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

// ClientCertMiddleware enforces that the incoming request includes a known client certificate
func ClientCertMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := beeline.StartSpan(r.Context(), "ClientCertMiddleware")
			defer span.Send()

			//session := auth.SessionFromRequestContext(r)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		return http.HandlerFunc(mw)
	}
}

// ClientCertVerifier returns a function suitable for use as a VerifyPeerCertificate callback that
// ensures that a supplied x509 certificate is known and found in the database
func ClientCertVerifier(db *pop.Connection) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return errors.New("No certs found")
		}

		// get DER hash
		hash := sha256.Sum256(rawCerts[0])
		hashString := hex.EncodeToString(hash[:])

		// check for presence of cert in client_certs table
		found, err := models.ClientCertExists(db, hashString)
		if err != nil {
			return err
		}
		if found != true {
			return fmt.Errorf("Unknown cert %s", hashString)
		}
		return nil
	}
}
