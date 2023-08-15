package server

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/ocsp"
)

const (
	idleTimeout       = 120 * time.Second // 2 minutes
	readHeaderTimeout = 60 * time.Second  // 1 minute
	maxHeaderSize     = 1 * 1000 * 1000   // 1 Megabyte
)

// the contextKey is typed so as not to conflict between similar keys from different pkgs
type contextKey string

var namedServerContextKey = contextKey("named_server")

// ErrMissingCACert represents an error caused by server config that requires
// certificate verification, but is missing a CA certificate
var ErrMissingCACert = errors.New("missing required CA certificate")

var cipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

var curvePreferences = []tls.CurveID{
	tls.CurveP256,
	tls.X25519,
}

// CreateNamedServerInput contains the input for the CreateServer function.
type CreateNamedServerInput struct {
	Name                  string
	Host                  string
	Port                  int
	Logger                *zap.Logger
	HTTPHandler           http.Handler
	ClientAuth            tls.ClientAuthType
	Certificates          []tls.Certificate
	ClientCAs             *x509.CertPool // CaCertPool
	VerifyPeerCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
}

// NamedServer wraps *http.Server to override the definition of ListenAndServeTLS, but bypasses some restrictions.
type NamedServer struct {
	*http.Server
	Name               string
	IsServerReady      bool
	IsServerReadyMutex sync.Mutex
}

// Port returns the port the server binds to.  Returns -1 if any error.
func (s *NamedServer) Port() int {
	if !strings.Contains(s.Addr, ":") {
		return -1
	}
	port, err := strconv.Atoi(strings.SplitN(s.Addr, ":", 2)[1])
	if err != nil {
		return -1
	}
	return port
}

// ListenAndServeTLS is similar to (*http.Server).ListenAndServeTLS, but bypasses some restrictions.
func (s *NamedServer) ListenAndServeTLS() error {
	listener, err := tls.Listen("tcp", s.Addr, s.TLSConfig)
	if err != nil {
		return err
	}
	s.IsServerReadyMutex.Lock()
	s.IsServerReady = true
	s.IsServerReadyMutex.Unlock()
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			fmt.Println(fmt.Errorf("Failed to close listener due to %w", closeErr))
		}
	}()
	return s.Serve(listener)
}

// IsReady returns if a server is ready
func (s *NamedServer) IsReady() bool {
	s.IsServerReadyMutex.Lock()
	defer s.IsServerReadyMutex.Unlock()
	return s.IsServerReady
}

// WaitUntilReady waits until the server is ready
func (s *NamedServer) WaitUntilReady() {
	times := 0
	// Wait for server to be ready
	for !s.IsReady() && times < 4 {
		times++
		time.Sleep(500 * time.Millisecond)
	}
}

// func fetchCRL(url string) (*x509.RevocationList, error) {
// 	httpResponse, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer httpResponse.Body.Close()

// 	if httpResponse.StatusCode >= 300 {
// 		return nil, errors.New("failed to retrieve CRL")
// 	}

// 	body, err := io.ReadAll(httpResponse.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return x509.ParseRevocationList(body)
// }

// // Request CRL response from server.
// // Returns error if the server can't get the certificate
// // func getCRLResponse(fetch storage.FileStorer, v *viper.Viper, crlFile string, clientCert *x509.Certificate, issuerCert *x509.Certificate) error {
// func getCRLResponse(clientCert *x509.Certificate, issuerCert *x509.Certificate) error {
// 	for _, url := range clientCert.CRLDistributionPoints {
// 		// TODO: Skip LDAP

// 		//x509.ParseRevocationList is not a direct ASN.1 representation, so leaves the option to add more detailed information
// 		parseCRL, err := fetchCRL(url)
// 		if err != nil {
// 			return err
// 		}

// 		//Parsed CRL against the issuer certificate
// 		err = parseCRL.CheckSignatureFrom(issuerCert)
// 		if err != nil {
// 			return err
// 		}

// 		// Check that the revocation list can be trusted
// 		if parseCRL.NextUpdate.Before(time.Now()) {
// 			return fmt.Errorf("CRL expired")
// 		}

// 		// Check id cert shows up in Revoked List
// 		for _, revokedCertificate := range parseCRL.RevokedCertificates {
// 			fmt.Printf("Revoked certificate serial number: %s\n", revokedCertificate.SerialNumber.String())
// 			if revokedCertificate.SerialNumber.Cmp(clientCert.SerialNumber) == 0 {
// 				return fmt.Errorf("revoked certificate")
// 			}
// 		}
// 	}

// 	return nil
// }

// Request OCSP response from server.
// Returns error if the server can't get the certificate
func getOCSPResponse(logger *zap.Logger, ocspServer string, clientCert *x509.Certificate, issuerCert *x509.Certificate) error {

	logger.Info("Checking OCSP",
		zap.String("server", ocspServer),
		zap.String("clientCert", clientCert.Subject.CommonName),
		zap.String("issuerCert", issuerCert.Subject.CommonName),
	)

	var httpClient = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	ocspRequestOpts := &ocsp.RequestOptions{Hash: crypto.SHA1}

	// buffer contains the serialized request that will be sent to the server.
	buffer, err := ocsp.CreateRequest(clientCert, issuerCert, ocspRequestOpts)
	if err != nil {
		return err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, ocspServer, bytes.NewBuffer(buffer))
	if err != nil {
		return err
	}

	ocspURL, err := url.Parse(ocspServer)
	if err != nil {
		return err
	}
	httpRequest.Header.Add("Content-Type", "application/ocsp-request")
	httpRequest.Header.Add("Accept", "application/ocsp-response")
	httpRequest.Header.Add("host", ocspURL.Host)

	httpResponse, err := httpClient.Do(httpRequest)
	// This means if we cannot reach the OSCP server, the certificate
	// will be invalid. Is that what we want?
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}
	logger.Info("OSCP cert response", zap.String("body", string(body)))
	ocspResponse, err := ocsp.ParseResponseForCert(body, clientCert, issuerCert)
	if err != nil {
		return err
	}

	switch ocspResponse.Status {
	case ocsp.Good:
		logger.Info("Certificate status: Good. It is still valid",
			zap.String("subject", clientCert.Subject.CommonName),
		)
	case ocsp.Revoked:
		logger.Info("Certificate status: Revoked.",
			zap.String("subject", clientCert.Subject.CommonName),
		)
		return fmt.Errorf("Revoked Certificate")
	default:
		logger.Info("Certificate status: Unknown",
			zap.Any("status", ocspResponse.Status),
			zap.String("subject", clientCert.Subject.CommonName),
		)
		return nil
	}

	return nil
}

func ocspRevokedCertCheck(logger *zap.Logger, clientCertificate *x509.Certificate, verifiedChains [][]*x509.Certificate) error {

	if len(clientCertificate.OCSPServer) == 0 {
		return nil
	}

	ocspURL := clientCertificate.OCSPServer[0]

	issuer := clientCertificate.Issuer.String()
	for _, vchains := range verifiedChains {
		for _, vchain := range vchains {
			if issuer == vchain.Subject.String() {
				err := getOCSPResponse(logger, ocspURL, clientCertificate, vchain)
				if err != nil {
					logger.Error("ocsp response error", zap.Error(err))
					return err
				}
			}
		}
	}
	return nil
}

// NewRevokedCertCheck creates a callback function to validate
// certificates, using the provided logger
func NewRevokedCertCheck(logger *zap.Logger) func(_ [][]byte, verifiedChains [][]*x509.Certificate) error {

	// If this callback returns nil, then the handshake continues and
	// will not be aborted
	//
	// rawCerts contain chains of certificates in raw ASN.1 format
	// each raw certs starts with the leafCert and ends with a root
	// self-signed CA certificate
	//
	// verifiedChains have a certificate chain that verifies the
	// signature validity and ends with a trusted certificate in the
	// chain
	return func(_ [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(verifiedChains) == 0 || len(verifiedChains[0]) == 0 {
			return nil
		}
		clientCertificate := verifiedChains[0][0]

		err := ocspRevokedCertCheck(logger, clientCertificate, verifiedChains)
		if err != nil {
			return err
		}

		// err := ocspRevokedCertCheck(logger, clientCertificate, verifiedChains)
		// if err != nil {
		// 	return err
		// }

		return nil
	}
}

// CreateNamedServer returns a no-tls, tls, or mutual-tls Server based on the input given and an error, if any.
func CreateNamedServer(input *CreateNamedServerInput) (*NamedServer, error) {

	address := fmt.Sprintf("%s:%d", input.Host, input.Port)

	var tlsConfig *tls.Config
	if len(input.Certificates) > 0 {

		if input.ClientAuth == tls.VerifyClientCertIfGiven || input.ClientAuth == tls.RequireAndVerifyClientCert {
			// RA Summary: staticcheck - SA1019 - Using a deprecated function, variable, constant or field
			// RA: Linter is flagging: input.ClientCAs.Subjects is deprecated: if s was returned by SystemCertPool, Subjects will not include the system roots.
			// RA: Why code valuable: It allows us to ensure we error if missing expected client certs.
			// RA: Mitigation: The deprecation notes this is a problem when reading SystemCertPool, but we do not use this here and are building up our own cert pool instead.
			// RA Developer Status: Mitigated
			// RA Validator Status: Mitigated
			// RA Validator: leodis.f.scott.civ@mail.mil
			// RA Modified Severity: CAT III
			// nolint:staticcheck
			if input.ClientCAs == nil || len(input.ClientCAs.Subjects()) == 0 {
				return nil, ErrMissingCACert
			}
		}

		// Follow Mozilla's "modern" server side TLS recommendations
		// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
		// https://statics.tls.security.mozilla.org/server-side-tls-conf-4.0.json
		// This configuration is compatible with Firefox 27, Chrome 30, IE 11 on
		// Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8
		tlsConfig = &tls.Config{
			CipherSuites:             cipherSuites,
			Certificates:             input.Certificates,
			ClientAuth:               input.ClientAuth,
			ClientCAs:                input.ClientCAs,
			CurvePreferences:         curvePreferences,
			MinVersion:               tls.VersionTLS12,
			NextProtos:               []string{"h2"},
			PreferServerCipherSuites: true,
			VerifyPeerCertificate:    input.VerifyPeerCertificate,
		}
		//option 1: if devLocal flag to switch between APIs that use mtls connection and those that do not
		//if auth.ApplicationServername == "AdminServername" || "PrimeServername" || "OrdersServername" {
		//if input.VerifyPeerCertificate {
		//	tlsConfig.VerifyPeerCertificate = certRevokedCheck
		//}

		//}

		//Option 2: set flag when server starts up that can turn off or on to test locally.
	}

	// wrappedHandler includes the name of the server in the context
	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), namedServerContextKey, input.Name)
		input.HTTPHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	srv := &NamedServer{
		Name: input.Name,
		Server: &http.Server{
			Addr:              address,
			ErrorLog:          newStandardLogger(input.Logger),
			Handler:           wrappedHandler,
			IdleTimeout:       idleTimeout,
			MaxHeaderBytes:    maxHeaderSize,
			ReadHeaderTimeout: readHeaderTimeout,
			TLSConfig:         tlsConfig,
		},
	}
	return srv, nil

}

// NamedServerFromContext returns name name of the server that was previously added into the context, if any.
func NamedServerFromContext(ctx context.Context) string {
	name, ok := ctx.Value(namedServerContextKey).(string)
	if !ok {
		return ""
	}
	return name
}
