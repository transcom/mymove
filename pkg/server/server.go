package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	idleTimeout   = 120 * time.Second // 2 minutes
	maxHeaderSize = 1 * 1000 * 1000   // 1 Megabyte
)

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
	Logger       Logger
	HTTPHandler  http.Handler
	ClientAuth   tls.ClientAuthType
	Certificates []tls.Certificate
	ClientCAs    *x509.CertPool // CaCertPool
}

// NamedServer wraps *http.Server to override the definition of ListenAndServeTLS, but bypasses some restrictions.
type NamedServer struct {
	*http.Server
	Name string
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
	defer listener.Close()
	return s.Serve(listener)
}

// CreateNamedServer returns a no-tls, tls, or mutual-tls Server based on the input given and an error, if any.
func CreateNamedServer(input *CreateNamedServerInput) (*NamedServer, error) {

	address := fmt.Sprintf("%s:%d", input.Host, input.Port)

	var tlsConfig *tls.Config
	if len(input.Certificates) > 0 {

		if input.ClientAuth == tls.VerifyClientCertIfGiven || input.ClientAuth == tls.RequireAndVerifyClientCert {
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

		// Map certificates with the CommonName / DNSNames to support
		// Server Name Indication (SNI). In other words this will tell
		// the TLS listener to sever the appropriate certificate matching
		// the requested hostname.
		tlsConfig.BuildNameToCertificate()
	}

	srv := &NamedServer{
		Name: input.Name,
		Server: &http.Server{
			Addr:           address,
			ErrorLog:       newStandardLogger(input.Logger),
			Handler:        input.HTTPHandler,
			IdleTimeout:    idleTimeout,
			MaxHeaderBytes: maxHeaderSize,
			TLSConfig:      tlsConfig,
		},
	}
	return srv, nil

}
