package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	idleTimeout = 120 * time.Second
	// max request headers size is 1 mb
	maxHeaderSize = 1 * 1000 * 1000
)

var supportedProtocols = []string{"h2"}

type serverFunc func(server *http.Server) error

// Server is
type Server struct {
	CACertPEMBlock []byte
	ClientAuthType tls.ClientAuthType
	HTTPHandler    http.Handler
	ListenAddress  string
	Logger         *zap.Logger
	Port           string
	TLSCerts       []TLSCert
}

func addr(listenAddress, port string) string {
	return fmt.Sprintf("%s:%s", listenAddress, port)
}

func stdLogError(logger *zap.Logger) (*log.Logger, error) {
	standardLog, err := zap.NewStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		return nil, err
	}
	return standardLog, nil
}

func (s Server) serverConfig(tlsConfig *tls.Config) (*http.Server, error) {
	// By detault http.Server will use the standard logging library which isn't
	// structured JSON. This will pass zap.Logger with log level error
	standardLog, err := stdLogError(s.Logger)
	if err != nil {
		s.Logger.Error("failed to create a standard logger", zap.Error(err))
		return nil, err
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

func (s Server) tlsConfig() (*tls.Config, error) {
	var caCerts *x509.CertPool

	// Load client Certificate Authority (CA) we are requiring client cert
	// authentication
	if s.ClientAuthType == tls.VerifyClientCertIfGiven ||
		s.ClientAuthType == tls.RequireAndVerifyClientCert {
		caCerts = x509.NewCertPool()
		ok := caCerts.AppendCertsFromPEM(s.CACertPEMBlock)
		if !ok {
			s.Logger.Fatal("failed to append client certificate authority")
		}
	}

	tlsCerts, err := ParseTLSCert(s.TLSCerts)
	if err != nil {
		//TODO add logging message
		s.Logger.Error("failed to create ", zap.Error(err))
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caCerts,
		Certificates: tlsCerts,
		NextProtos:   supportedProtocols,
		ClientAuth:   s.ClientAuthType,
	}
	return tlsConfig, err
}

// ListenAndServeTLS will create a
func (s Server) ListenAndServeTLS() error {
	var serverFunc serverFunc
	var server *http.Server
	var tlsConfig *tls.Config
	var err error

	tlsConfig, err = s.tlsConfig()
	if err != nil {
		//TODO handle error message
		return err
	}

	server, err = s.serverConfig(tlsConfig)
	if err != nil {
		//TODO handle error message
		return err
	}

	s.Logger.Info("start https listener",
		zap.Duration("idle-timeout", server.IdleTimeout),
		zap.Any("listen-address", s.ListenAddress),
		zap.Int("max-header-bytes", server.MaxHeaderBytes),
		zap.String("port", s.Port),
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

// ListenAndServe will create an HTTP listener
func (s Server) ListenAndServe() error {
	var serverFunc serverFunc
	var server *http.Server
	var tlsConfig *tls.Config
	var err error

	server, err = s.serverConfig(tlsConfig)
	if err != nil {
		//TODO handle error message
		return err
	}

	s.Logger.Info("start http listener",
		zap.Duration("idle-timeout", server.IdleTimeout),
		zap.Any("listen-address", s.ListenAddress),
		zap.Int("max-header-bytes", server.MaxHeaderBytes),
		zap.String("port", s.Port),
	)

	serverFunc = func(httpServer *http.Server) error {
		return httpServer.ListenAndServe()
	}

	return serverFunc(server)
}
