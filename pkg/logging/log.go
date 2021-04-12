package logging

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

type stackTracer interface {
	StackTrace() errors.StackTrace
}

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
	config := &ZapConfig{}

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
	modifiedFields := filterErrorFields(fields, fce.stacktraceLength)
	ent = filteredAndLimitedStackTrace(ent, fce.stacktraceLength)
	return fce.Encoder.EncodeEntry(ent, modifiedFields)
}

func (fje *filteredJSONEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	modifiedFields := filterErrorFields(fields, fje.stacktraceLength)
	ent = filteredAndLimitedStackTrace(ent, fje.stacktraceLength)
	return fje.Encoder.EncodeEntry(ent, modifiedFields)
}

func filterErrorFields(fields []zapcore.Field, lineLimit int) []zapcore.Field {
	var modifiedFields []zapcore.Field

	for _, field := range fields {
		if field.Type == zapcore.ErrorType {
			fieldError := field.Interface.(error)
			stacktraceError, ok := fieldError.(stackTracer) // error implements the stackTracer interface

			// error may be of type errorString without a stacktrace so will
			// only create error message key value
			if !ok {
				modifiedFields = append(modifiedFields, field)
				continue
			}

			stacktrace := stacktraceError.StackTrace()

			filteredStacktrace := filterStacktraceFrames(stacktrace, lineLimit)

			// converts the error field to a string field with the same key (defaults to "error" if not named)
			// TODO: how are slices of errors handled?
			modifiedFields = append(modifiedFields, zap.String(field.Key, fieldError.Error()))
			// preserve the stacktrace in the keyVerbose field
			modifiedFields = append(modifiedFields, zap.String(field.Key+"Verbose", fmt.Sprintf("%+v", filteredStacktrace)))
		} else {
			modifiedFields = append(modifiedFields, field)
		}
	}

	return modifiedFields
}

func filterStacktraceFrames(frames []errors.Frame, lineLimit int) []errors.Frame {
	if frames == nil || (len(frames)*2) <= lineLimit {
		return frames
	}

	var filteredFrames []errors.Frame
	for _, frame := range frames {
		// %+v will return the package function and a filename path with line
		// number seperated by a newline
		if strings.Contains(fmt.Sprintf("%+v", frame), "mymove") {
			filteredFrames = append(filteredFrames, frame)

			// a frame is a pair of 2 lines
			if len(filteredFrames)*2 == lineLimit {
				break
			}
		}
	}

	return filteredFrames
}

// Filter the stack trace to only return lines from the mymove codebase
// And limit the number of lines based on the STACKTRACE_LENGTH env var
// The ability to limit the stacktrace length is a STIG requirement.
func filteredAndLimitedStackTrace(ent zapcore.Entry, stacktraceLength int) zapcore.Entry {
	if ent.Stack == "" {
		return ent
	}

	stacktraceLines := strings.Split(ent.Stack, "\n")

	// We don't need to filter if the stacktrace is beneath than the limit
	if len(stacktraceLines) <= stacktraceLength {
		return ent
	}

	var matchingLines []string
	searchTerm := "mymove"

	for _, line := range stacktraceLines {
		if strings.Contains(line, searchTerm) {
			matchingLines = append(matchingLines, line)
			if len(matchingLines) >= stacktraceLength {
				break
			}
		}
	}
	ent.Stack = strings.Join(matchingLines, "\n")

	return ent
}

func registerCustomZapEncoders(stacktraceLength int) {
	_ = zap.RegisterEncoder("filtered-console", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		fce := filteredConsoleEncoder{
			EncoderConfig:    &cfg,
			Encoder:          zapcore.NewConsoleEncoder(cfg),
			consoleEncoder:   zapcore.NewConsoleEncoder(cfg),
			stacktraceLength: stacktraceLength,
		}

		return &fce, nil
	})

	_ = zap.RegisterEncoder("filtered-json", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		fje := filteredJSONEncoder{
			EncoderConfig:    &cfg,
			Encoder:          zapcore.NewJSONEncoder(cfg),
			jsonEncoder:      zapcore.NewJSONEncoder(cfg),
			stacktraceLength: stacktraceLength,
		}

		return &fje, nil
	})
}
