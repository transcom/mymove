package logging

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config configures a Zap logger based on the environment string and debugLevel
func Config(env string, loggingLevel string) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if env != "development" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch strings.ToLower(loggingLevel) {
	case "info":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "debug":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return loggerConfig.Build()
}
