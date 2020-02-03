package rateengine

import (
	"fmt"
	"go.uber.org/zap"
)

// Logger is an interface that describes the logging requirements of this package.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func AppendID(msg string, id string) string {
	if id != "" {
		return fmt.Sprintf(msg + ": " + id)
	}
	return msg
}
