package server

import (
	"crypto/tls"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// TLSServer is
type TLSServer struct {
	CertPEMBlock  []byte
	KeyPEMBlock   []byte
	ListenAddress string
	HTTPHandler   http.Handler
	Logger        *zap.Logger
	Port          string
	TLSCerts      []TLSCert
}

// ListenAndServeTLS will create a
func (s TLSServer) ListenAndServeTLS() error {
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
		// TODO add logging message
		s.Logger.Error("failed to create ", zap.Error(err))
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: tlsCerts,
		NextProtos:   supportedProtocols,
	}

	server = &http.Server{
		Addr:           addr(s.ListenAddress, s.Port),
		ErrorLog:       standardLog,
		Handler:        s.HTTPHandler,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderSize,
		TLSConfig:      tlsConfig,
	}

	s.Logger.Info("start tls server",
		zap.Duration("idle-timeout", server.IdleTimeout),
		zap.Any("listen-address", s.ListenAddress),
		zap.Int("max-header-bytes", server.MaxHeaderBytes),
		zap.String("port", s.Port),
	)

	serverFunc = func(httpServer *http.Server) error {
		tlsListener, err := tls.Listen("tcp", server.Addr, tlsConfig)
		if err != nil {
			return err
		}
		defer tlsListener.Close()
		return server.Serve(tlsListener)
	}

	return serverFunc(server)
}
