package server

import (
	"log"
	"net/http"

	"go.uber.org/zap"
)

// HTTPServer is
type HTTPServer struct {
	ListenAddress string
	HTTPHandler   http.Handler
	Logger        *zap.Logger
	Port          string
}

// ListenAndServe will create an HTTP listener
func (s HTTPServer) ListenAndServe() error {
	var serverFunc serverFunc
	var server *http.Server
	var standardLog *log.Logger
	var err error

	standardLog, err = stdLogAt(s.Logger)
	if err != nil {
		s.Logger.Error("failed to create a standard logger", zap.Error(err))
		return err
	}

	server = &http.Server{
		Addr:           addr(s.ListenAddress, s.Port),
		ErrorLog:       standardLog,
		Handler:        s.HTTPHandler,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderSize,
	}

	s.Logger.Info("start http server",
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
