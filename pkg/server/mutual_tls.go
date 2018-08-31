package server

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// MutualTLSServer is
type MutualTLSServer struct {
	CACertPEMBlock []byte
	HTTPHandler    http.Handler
	ListenAddress  string
	Logger         *zap.Logger
	Port           string
	TLSCerts       []TLSCert
}

// ListenAndServeMutualTLS will create a
func (s MutualTLSServer) ListenAndServeMutualTLS() error {
	var serverFunc serverFunc
	var server *http.Server
	var standardLog *log.Logger
	var err error

	standardLog, err = stdLogAt(s.Logger)
	if err != nil {
		s.Logger.Error("failed to create a standard logger", zap.Error(err))
		return err
	}

	tlsCerts, err := ParseTLSCert(s.TLSCerts)
	if err != nil {
		//TODO add logging message
		s.Logger.Error("failed to create ", zap.Error(err))
		return err
	}

	caCerts := x509.NewCertPool()
	ok := caCerts.AppendCertsFromPEM(s.CACertPEMBlock)
	if !ok {
		s.Logger.Fatal("failed to append client certificate authority", zap.Error(err))
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caCerts,
		Certificates: tlsCerts,
		NextProtos:   supportedProtocols,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server = &http.Server{
		Addr:           addr(s.ListenAddress, s.Port),
		ErrorLog:       standardLog,
		Handler:        s.HTTPHandler,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderSize,
		TLSConfig:      tlsConfig,
	}

	s.Logger.Info("start mutual tls server",
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
