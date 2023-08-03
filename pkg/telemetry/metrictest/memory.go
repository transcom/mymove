package metrictest

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// ensure InMemoryExporter implements Exporter interface
var _ metric.Exporter = &InMemoryExporter{}

type InMemoryExporter struct {
	mu sync.Mutex
	rm []metricdata.ResourceMetrics
}

func NewInMemoryExporter() *InMemoryExporter {
	return new(InMemoryExporter)
}

func (i *InMemoryExporter) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return metric.DefaultTemporalitySelector(k)
}

func (i *InMemoryExporter) Aggregation(k metric.InstrumentKind) aggregation.Aggregation {
	return metric.DefaultAggregationSelector(k)
}

func (i *InMemoryExporter) Export(_ context.Context, data *metricdata.ResourceMetrics) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.rm = append(i.rm, *data)
	return nil
}

func (i *InMemoryExporter) ForceFlush(ctx context.Context) error {
	// exporter holds no state, nothing to flush.
	return ctx.Err()
}

func (i *InMemoryExporter) Shutdown(_ context.Context) error {
	i.Reset()
	return nil
}

func (i *InMemoryExporter) Reset() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.rm = nil
}

func (i *InMemoryExporter) GetMetrics() []metricdata.ResourceMetrics {
	i.mu.Lock()
	defer i.mu.Unlock()
	ret := make([]metricdata.ResourceMetrics, len(i.rm))
	copy(ret, i.rm)
	return ret
}
