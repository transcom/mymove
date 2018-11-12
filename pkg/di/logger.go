package di

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger is the DI provider for constructing a new zap.Logger
func NewLogger(cfg *viper.Viper) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if cfg.GetString("env") != "development" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.GetBool("debug-logging") {
		debug := zap.NewAtomicLevel()
		debug.SetLevel(zap.DebugLevel)
		loggerConfig.Level = debug
	}
	return loggerConfig.Build()
}
