package trace

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const heapMemoryName = "heapMemory"
const heapMemoryDesc = "heatMemory Desc"

// RegisterRuntimeObserver creates a custom metric that is updated
// automatically using an observer
func RegisterRuntimeObserver(config TelemetryConfig) {
	if !config.Enabled {
		return
	}
	metric.Must(meter).NewInt64ValueObserver(
		heapMemoryName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			result.Observe(int64(mem.HeapAlloc),
				attribute.String(heapMemoryName, heapMemoryDesc))
		},
		metric.WithDescription(heapMemoryDesc))

}
