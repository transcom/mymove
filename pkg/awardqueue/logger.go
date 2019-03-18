package awardqueue

import (
	"context"

	"go.uber.org/zap"
)

// Logger is an interface that describes the logging requirements of this package.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	TraceInfo(ctx context.Context, msg string, fields ...zap.Field)
	TraceError(ctx context.Context, msg string, fields ...zap.Field)
}
