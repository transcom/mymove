package trace

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.uber.org/zap"
)

// RegisterRuntimeObserver creates a custom metric that is updated
// automatically using an observer
func RegisterRuntimeObserver(logger *zap.Logger, config TelemetryConfig) {
	if !config.Enabled {
		return
	}
	collectSeconds := config.CollectSeconds
	if collectSeconds == 0 {
		collectSeconds = defaultCollectSeconds
	}

	if err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Duration(collectSeconds) * time.Second),
	); err != nil {
		logger.Fatal("failed to start runtime instrumentation:", zap.Error(err))
	}

}
