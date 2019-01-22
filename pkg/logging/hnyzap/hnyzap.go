package hnyzap

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/honeycombio/beeline-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const fieldPrefix string = "logging"

// Logger is an embedded zap.Logger to extend zap logs into Honeycomb.
type Logger struct {
	*zap.Logger
}

// ZapFieldToHoneycombField converts a zap.Field in to a type supported by
// Honeycomb (bool, number, string). Unsupported types will have a string value
// of "unsupported field type"
func ZapFieldToHoneycombField(zapField zap.Field) interface{} {
	switch zapField.Type {
	case zapcore.BoolType:
		val := false
		if zapField.Integer >= 1 {
			val = true
		}
		return val
	case zapcore.Float32Type:
		return math.Float32frombits(uint32(zapField.Integer))
	case zapcore.Float64Type:
		return math.Float64frombits(uint64(zapField.Integer))
	case zapcore.Int32Type:
		return int32(zapField.Integer)
	case zapcore.Int64Type:
		return zapField.Integer
	case zapcore.StringType:
		return zapField.String
	case zapcore.Uint32Type:
		return uint32(zapField.Integer)
	case zapcore.Uint64Type:
		return uint64(zapField.Integer)
	case zapcore.ErrorType:
		return zapField.Interface.(error).Error()
	default:
		return "unsupported field type"
	}
}

// LogToHoneycombSpan translates zap.Fields into fields supported by Honeycomb's
// tracing service. Honeycomb currently support bool, number, and string
// types. zap.Field types not supported by Honeycomb will still be logged to
// zap, but will be sent to Honeycomb with a string "unsupported field type".
func LogToHoneycombSpan(ctx context.Context, level string, msg string, fields ...zap.Field) {
	_, span := beeline.StartSpan(ctx, level)
	defer span.Send()

	span.AddField(fmt.Sprintf("%s.level", fieldPrefix), strings.ToLower(level))
	span.AddField(fmt.Sprintf("%s.msg", fieldPrefix), msg)
	for _, zapField := range fields {
		fieldKey := fmt.Sprintf("%s.%s", fieldPrefix, zapField.Key)
		fieldValue := ZapFieldToHoneycombField(zapField)
		span.AddField(fieldKey, fieldValue)
	}
}

// TraceDebug logs a message at DebugLevel to a span within a Honeycomb trace as well as the configured zap logger.
func (l *Logger) TraceDebug(ctx context.Context, msg string, fields ...zap.Field) {
	LogToHoneycombSpan(ctx, "Debug", msg, fields...)
	skipLogger := l.Logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Debug(msg, fields...)
}

// TraceInfo logs a message at InfoLevel to a span within a Honeycomb trace as well as the configured zap logger.
func (l *Logger) TraceInfo(ctx context.Context, msg string, fields ...zap.Field) {
	LogToHoneycombSpan(ctx, "Info", msg, fields...)
	skipLogger := l.Logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Info(msg, fields...)
}

// TraceWarn logs a message at WarnLevel to a span within a Honeycomb trace as well as the configured zap logger.
func (l *Logger) TraceWarn(ctx context.Context, msg string, fields ...zap.Field) {
	LogToHoneycombSpan(ctx, "Warn", msg, fields...)
	skipLogger := l.Logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Warn(msg, fields...)
}

// TraceError logs a message at ErrorLevel to a span within a Honeycomb trace as well as the configured zap logger.
func (l *Logger) TraceError(ctx context.Context, msg string, fields ...zap.Field) {
	LogToHoneycombSpan(ctx, "Error", msg, fields...)
	skipLogger := l.Logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Error(msg, fields...)
}

// TraceFatal logs a message at FatalLevel to a span within a Honeycomb trace as well as the configured zap logger.
func (l *Logger) TraceFatal(ctx context.Context, msg string, fields ...zap.Field) {
	LogToHoneycombSpan(ctx, "Fatal", msg, fields...)
	skipLogger := l.Logger.WithOptions(zap.AddCallerSkip(1))
	skipLogger.Fatal(msg, fields...)
}
