package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	idleTimeout   = 120 * time.Second // 2 minutes
	maxHeaderSize = 1 * 1000 * 1000   // 1 Megabyte
)

// ErrMissingCACert represents an error caused by server config that requires
// certificate verification, but is missing a CA certificate
var ErrMissingCACert = errors.New("missing required CA certificate")

// ErrUnparseableCACert represents an error cause by a misconfigured CA certificate
// that was unable to be parsed.
var ErrUnparseableCACert = errors.New("unable to parse CA certificate")

type serverFunc func(server *http.Server) error

// Server represents an http or https listening server. HTTPS listeners support
// requiring client authentication with a provided CA.
type Server struct {
	CaCertPool     *x509.CertPool
	ClientAuthType tls.ClientAuthType
	HTTPHandler    http.Handler
	ListenAddress  string
	Logger         *zap.Logger
	Port           int
	TLSCerts       []tls.Certificate
}

// addr generates an address:port string to be used in defining an http.Server
func addr(listenAddress string, port int) string {
	return fmt.Sprintf("%s:%d", listenAddress, port)
}

// stdLogError creates a *log.logger based off an existing zap.Logger instance.
// Some libraries call log.logger directly, which isn't structured as JSON. This method
// Will reformat log calls as zap.Error logs.
func stdLogError(logger *zap.Logger) (*log.Logger, error) {
	standardLog, err := zap.NewStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		return nil, err

	}
	return standardLog, nil
}

// serverConfig generates a *http.Server with a structured error logger.
func (s Server) serverConfig(tlsConfig *tls.Config) (*http.Server, error) {
	// By detault http.Server will use the standard logging library which isn't
	// structured JSON. This will pass zap.Logger with log level error
	standardLog, err := stdLogError(s.Logger)
	if err != nil {
		s.Logger.Error("failed to create an error logger", zap.Error(err))
		return nil, errors.Wrap(err, "Faile")
	}

	serverConfig := &http.Server{
		Addr:           addr(s.ListenAddress, s.Port),
		ErrorLog:       standardLog,
		Handler:        s.HTTPHandler,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderSize,
		TLSConfig:      tlsConfig,
	}
	return serverConfig, err
}

// tlsConfig generates a new *tls.Config based on Mozilla's recommendations and returns an error, if any.
func (s Server) tlsConfig() (*tls.Config, error) {

	// Load client Certificate Authority (CA) if we are requiring client
	// cert authentication.
	if s.ClientAuthType == tls.VerifyClientCertIfGiven ||
		s.ClientAuthType == tls.RequireAndVerifyClientCert {
		if s.CaCertPool == nil || len(s.CaCertPool.Subjects()) == 0 {
			return nil, ErrMissingCACert
		}
	}

	// Follow Mozilla's "modern" server side TLS recommendations
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
	// https://statics.tls.security.mozilla.org/server-side-tls-conf-4.0.json
	// This configuration is compatible with Firefox 27, Chrome 30, IE 11 on
	// Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8
	tlsConfig := &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		Certificates: s.TLSCerts,
		ClientAuth:   s.ClientAuthType,
		ClientCAs:    s.CaCertPool,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion:               tls.VersionTLS12,
		NextProtos:               []string{"h2"},
		PreferServerCipherSuites: true,
	}

	// Map certificates with the CommonName / DNSNames to support
	// Server Name Indication (SNI). In other words this will tell
	// the TLS listener to sever the appropriate certificate matching
	// the requested hostname.
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

// ListenAndServeTLS returns a TLS Listener function for serving HTTPS requests
func (s Server) ListenAndServeTLS() error {
	var serverFunc serverFunc
	var server *http.Server
	var tlsConfig *tls.Config
	var err error

	tlsConfig, err = s.tlsConfig()
	if err != nil {
		s.Logger.Error("failed to generate a TLS config", zap.Error(err))
		return err
	}

	server, err = s.serverConfig(tlsConfig)
	if err != nil {
		s.Logger.Error("failed to generate a TLS server config", zap.Error(err))
		return err
	}

	s.Logger.Info("start https listener",
		zap.Duration("idle-timeout", server.IdleTimeout),
		zap.Any("listen-address", s.ListenAddress),
		zap.Int("max-header-bytes", server.MaxHeaderBytes),
		zap.Int("port", s.Port),
	)

	serverFunc = func(httpServer *http.Server) error {
		tlsListener, err := tls.Listen("tcp",
			server.Addr,
			tlsConfig)
		if err != nil {
			return err
		}
		defer tlsListener.Close()
		return server.Serve(tlsListener)
	}

	return serverFunc(server)
}

// ListenAndServe returns an HTTP ListenAndServe function for serving HTTP requests
func (s Server) ListenAndServe() error {
	var serverFunc serverFunc
	var server *http.Server
	var tlsConfig *tls.Config
	var err error

	server, err = s.serverConfig(tlsConfig)
	if err != nil {
		s.Logger.Error("failed to generate a server config", zap.Error(err))
		return err
	}

	s.Logger.Info("start http listener",
		zap.Duration("idle-timeout", server.IdleTimeout),
		zap.Any("listen-address", s.ListenAddress),
		zap.Int("max-header-bytes", server.MaxHeaderBytes),
		zap.Int("port", s.Port),
	)

	serverFunc = func(httpServer *http.Server) error {
		return httpServer.ListenAndServe()
	}

	return serverFunc(server)
}
