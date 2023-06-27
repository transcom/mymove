package telemetry

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"github.com/transcom/mymove/pkg/appcontext"
)

// RegisterRuntimeObserver creates a custom metric that is updated
// automatically using an observer
func RegisterRuntimeObserver(_ appcontext.AppContext, config *Config) error {
	if !config.Enabled {
		return nil
	}
	collectSeconds := config.CollectSeconds
	if collectSeconds == 0 {
		collectSeconds = defaultCollectSeconds
	}

	return runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Duration(collectSeconds) * time.Second))
}
