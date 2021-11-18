package telemetry

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

// RegisterRuntimeObserver creates a custom metric that is updated
// automatically using an observer
func RegisterRuntimeObserver(appCtx appcontext.AppContext, config *Config) {
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
		appCtx.Logger().Fatal("failed to start runtime instrumentation:", zap.Error(err))
	}

}
