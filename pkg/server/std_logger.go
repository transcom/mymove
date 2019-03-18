package server

import (
	"bytes"
	"log"

	"go.uber.org/zap"
)

type stdLogger struct {
	Logger func(msg []byte)
}

func (l *stdLogger) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.Logger(p)
	return len(p), nil
}

func newStandardLogger(l Logger) *log.Logger {
	logger := l.WithOptions(zap.AddCallerSkip(4)).Error
	return log.New(
		&stdLogger{Logger: func(msg []byte) { logger(string(msg)) }},
		"",
		0,
	)
}
