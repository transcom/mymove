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

// Logger is a wrapped zap.Logger to extend zap logs into Honeycomb.
type Logger struct {
	*zap.Logger
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
		switch zapField.Type {
		case zapcore.BoolType:
			val := false
			if zapField.Integer >= 1 {
				val = true
			}
			span.AddField(fieldKey, val)
		case zapcore.Float32Type:
			span.AddField(fieldKey, math.Float32frombits(uint32(zapField.Integer)))
		case zapcore.Float64Type:
			span.AddField(fieldKey, math.Float64frombits(uint64(zapField.Integer)))
		case zapcore.Int64Type:
			span.AddField(fieldKey, zapField.Integer)
		case zapcore.Int32Type:
			span.AddField(fieldKey, int32(zapField.Integer))
		case zapcore.StringType:
			span.AddField(fieldKey, zapField.String)
		case zapcore.Uint64Type:
			span.AddField(fieldKey, uint64(zapField.Integer))
		case zapcore.Uint32Type:
			span.AddField(fieldKey, uint32(zapField.Integer))
		case zapcore.ErrorType:
			span.AddField(fieldKey, zapField.Interface.(error).Error())
		default:
			span.AddField(fieldKey, "unsupported field type")
		}
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
