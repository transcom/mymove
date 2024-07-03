package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	Name         string
	Host         string
	Port         int
	Logger       *zap.Logger
	HTTPHandler  http.Handler
	ClientAuth   tls.ClientAuthType
	Certificates []tls.Certificate
	ClientCAs    *x509.CertPool // CaCertPool
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
			fmt.Println(fmt.Errorf("failed to close listener due to %w", closeErr))
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
		}
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
