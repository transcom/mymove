package server

import (
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

func addr(listenAddress, port string) string {
	return fmt.Sprintf("%s:%s", listenAddress, port)
}

func stdLogAt(logger *zap.Logger) (*log.Logger, error) {
	standardLog, err := zap.NewStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		return nil, err
	}
	return standardLog, nil
}
