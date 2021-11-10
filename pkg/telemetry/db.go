package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/stats"
)

// Ideally the "go.opentelemetry.io/otel/semconv" package would
// provide standard names for the stats below, but they don't seem to
// See https://github.com/open-telemetry/opentelemetry-go/blob/main/semconv/v1.7.0/trace.go

// https://pkg.go.dev/database/sql#DB.Stats InUse
const dbPoolInUseName = "db.pool.inuse"
const dbPoolInUseDesc = "The number of connections currently in use"

// https://pkg.go.dev/database/sql#DB.Stats Idle
const dbPoolIdleName = "db.pool.idle"
const dbPoolIdleDesc = "The number of connections currently in use"

// https://pkg.go.dev/database/sql#DB.Stats WaitDuration
const dbWaitDurationName = "db.waitduration"
const dbWaitDurationDesc = "The total time blocked waiting for a new connection"

// RegisterDBStatsObserver creates a custom metric that is updated
// automatically using an observer
func RegisterDBStatsObserver(appCtx appcontext.AppContext, config *Config) {
	if !config.Enabled {
		return
	}

	metric.Must(meter).NewInt64UpDownCounterObserver(
		dbPoolInUseName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, err := stats.DBStats(appCtx)
			if err == nil {
				result.Observe(int64(dbStats.InUse),
					attribute.String(dbPoolInUseName, dbPoolInUseDesc))
			}
		},
		metric.WithDescription(dbPoolInUseDesc))

	metric.Must(meter).NewInt64UpDownCounterObserver(
		dbPoolIdleName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, err := stats.DBStats(appCtx)
			if err == nil {
				result.Observe(int64(dbStats.Idle),
					attribute.String(dbPoolIdleName, dbPoolIdleDesc))
			}
		},
		metric.WithDescription(dbPoolInUseDesc))

	metric.Must(meter).NewInt64UpDownCounterObserver(
		dbWaitDurationName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, err := stats.DBStats(appCtx)
			if err == nil {
				result.Observe(int64(dbStats.WaitDuration),
					attribute.String(dbPoolInUseName, dbPoolInUseDesc))
			}
		},
		metric.WithDescription(dbWaitDurationDesc))

}
