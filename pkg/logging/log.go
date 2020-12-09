package logging

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// ZapConfig defines the configurable parameters you can pass in when
// setting up the logger. See examples where we call logging.Config
type ZapConfig struct {
	Environment      string
	LoggingLevel     string
	StacktraceLength int
}

// ZapConfigOption is the type for the possible options you can pass in
// to logging.Config
type ZapConfigOption func(*ZapConfig)

// WithEnvironment provides an option to pass in the environment
func WithEnvironment(environment string) ZapConfigOption {
	return func(c *ZapConfig) {
		c.Environment = environment
	}
}

// WithLoggingLevel provides an option to pass in the logging level
func WithLoggingLevel(level string) ZapConfigOption {
	return func(c *ZapConfig) {
		c.LoggingLevel = level
	}
}

// WithStacktraceLength provides an option to pass in the stack trace length
func WithStacktraceLength(length int) ZapConfigOption {
	return func(c *ZapConfig) {
		c.StacktraceLength = length
	}
}

// Config configures a Zap logger based on the environment string and debugLevel
func Config(opts ...ZapConfigOption) (*zap.Logger, error) {
	const (
		defaultEnvironment      = "development"
		defaultLoggingLevel     = "info"
		defaultStacktraceLength = 6
	)

	config := &ZapConfig{
		Environment:      defaultEnvironment,
		LoggingLevel:     defaultLoggingLevel,
		StacktraceLength: defaultStacktraceLength,
	}

	for _, opt := range opts {
		opt(config)
	}

	var loggerConfig zap.Config

	registerCustomZapEncoders(config.StacktraceLength)

	devEncoderConfig := zap.NewDevelopmentEncoderConfig()
	devEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	devConfig := zap.NewDevelopmentConfig()
	devConfig.EncoderConfig = devEncoderConfig

	if config.Environment != "development" {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.Encoding = "filtered-json"
	} else {
		loggerConfig = devConfig
		loggerConfig.Encoding = "filtered-console"
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch strings.ToLower(config.LoggingLevel) {
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

type filteredConsoleEncoder struct {
	*zapcore.EncoderConfig
	zapcore.Encoder
	consoleEncoder   zapcore.Encoder
	stacktraceLength int
}

type filteredJSONEncoder struct {
	*zapcore.EncoderConfig
	zapcore.Encoder
	jsonEncoder      zapcore.Encoder
	stacktraceLength int
}

// Clone implements the zapcore Encoder interface
func (fce *filteredConsoleEncoder) Clone() zapcore.Encoder {
	return &filteredConsoleEncoder{
		EncoderConfig:    fce.EncoderConfig,
		Encoder:          fce.Encoder.Clone(),
		consoleEncoder:   fce.consoleEncoder.Clone(),
		stacktraceLength: fce.stacktraceLength,
	}
}

// Clone implements the zapcore Encoder interface
func (fje *filteredJSONEncoder) Clone() zapcore.Encoder {
	return &filteredJSONEncoder{
		EncoderConfig:    fje.EncoderConfig,
		Encoder:          fje.Encoder.Clone(),
		jsonEncoder:      fje.jsonEncoder.Clone(),
		stacktraceLength: fje.stacktraceLength,
	}
}

func (fce *filteredConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	ent = filteredAndLimitedStackTrace(ent, fce.stacktraceLength)
	return fce.consoleEncoder.EncodeEntry(ent, fields)
}

func (fje *filteredJSONEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	ent = filteredAndLimitedStackTrace(ent, fje.stacktraceLength)
	return fje.jsonEncoder.EncodeEntry(ent, fields)
}

// Filter the stack trace to only return lines from the mymove codebase
// And limit the number of lines based on the STACKTRACE_LENGTH env var
// The ability to limit the stacktrace length is a STIG requirement.
func filteredAndLimitedStackTrace(ent zapcore.Entry, stacktraceLength int) zapcore.Entry {
	if ent.Stack == "" {
		return ent
	}

	var matchingLines []string
	searchTerm := "mymove/pkg"

	stacktraceLines := strings.Split(ent.Stack, "\n")

	for _, line := range stacktraceLines {
		if strings.Contains(line, searchTerm) {
			matchingLines = append(matchingLines, line)
			if len(matchingLines) == stacktraceLength {
				break
			}
		}
	}
	ent.Stack = strings.Join(matchingLines, "\n")

	return ent
}

func registerCustomZapEncoders(stacktraceLength int) {
	zap.RegisterEncoder("filtered-console", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		fce := filteredConsoleEncoder{
			EncoderConfig:    &cfg,
			Encoder:          zapcore.NewConsoleEncoder(cfg),
			consoleEncoder:   zapcore.NewConsoleEncoder(cfg),
			stacktraceLength: stacktraceLength,
		}

		return &fce, nil
	})

	zap.RegisterEncoder("filtered-json", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		fje := filteredJSONEncoder{
			EncoderConfig:    &cfg,
			Encoder:          zapcore.NewJSONEncoder(cfg),
			jsonEncoder:      zapcore.NewJSONEncoder(cfg),
			stacktraceLength: stacktraceLength,
		}

		return &fje, nil
	})
}
