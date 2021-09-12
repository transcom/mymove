package trace

import (
	"context"

	"github.com/gobuffalo/pop/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/transcom/mymove/pkg/db/stats"
)

const dbPoolName = "dbPool"
const dbPoolDesc = "dbPool Description"
const dbWaitDurationName = "dbWaitDuration"
const dbWaitDurationDesc = "dbWaitDuration description"

// RegisterDBStatsObserver creates a custom metric that is updated
// automatically using an observer
func RegisterDBStatsObserver(c *pop.Connection, config TelemetryConfig) {
	if !config.Enabled {
		return
	}

	metric.Must(meter).NewInt64UpDownCounterObserver(
		dbPoolName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, err := stats.DBStats(c)
			if err == nil {
				result.Observe(int64(dbStats.InUse),
					attribute.String(dbPoolName, dbPoolDesc))
			}
		},
		metric.WithDescription(dbPoolDesc))

	metric.Must(meter).NewInt64UpDownCounterObserver(
		dbWaitDurationName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, err := stats.DBStats(c)
			if err == nil {
				result.Observe(int64(dbStats.WaitDuration),
					attribute.String(dbPoolName, dbPoolDesc))
			}
		},
		metric.WithDescription(dbWaitDurationDesc))

}
