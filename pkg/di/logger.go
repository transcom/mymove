package di

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger is the DI provider for constructing a new zap.Logger
func NewLogger(cfg *Config) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if cfg.Environment != "development" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.DebugLogging {
		debug := zap.NewAtomicLevel()
		debug.SetLevel(zap.DebugLevel)
		loggerConfig.Level = debug
	}
	return loggerConfig.Build()
}
