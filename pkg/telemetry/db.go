package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.uber.org/zap"

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

const dbTelemetryVersion = "0.1"

// RegisterDBStatsObserver creates a custom metric that is updated
// automatically using an observer
func RegisterDBStatsObserver(appCtx appcontext.AppContext, config *Config) {
	if !config.Enabled {
		return
	}

	meterProvider := global.GetMeterProvider()

	dbMeter := meterProvider.Meter("github.com/transcom/mymove/db",
		metric.WithInstrumentationVersion(dbTelemetryVersion))

	_, err := dbMeter.NewInt64UpDownCounterObserver(
		dbPoolInUseName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				result.Observe(int64(dbStats.InUse),
					attribute.String(dbPoolInUseName, dbPoolInUseDesc))
			}
		},
		metric.WithDescription(dbPoolInUseDesc))
	if err != nil {
		appCtx.Logger().Fatal("Failed to start db instrumentation", zap.Error(err))
	}

	_, err = dbMeter.NewInt64UpDownCounterObserver(
		dbPoolIdleName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				result.Observe(int64(dbStats.Idle),
					attribute.String(dbPoolIdleName, dbPoolIdleDesc))
			}
		},
		metric.WithDescription(dbPoolIdleDesc))
	if err != nil {
		appCtx.Logger().Fatal("Failed to start db instrumentation", zap.Error(err))
	}

	_, err = dbMeter.NewInt64CounterObserver(
		dbWaitDurationName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				result.Observe(int64(dbStats.WaitDuration),
					attribute.String(dbWaitDurationName, dbWaitDurationDesc))
			}
		},
		metric.WithDescription(dbWaitDurationDesc))
	if err != nil {
		appCtx.Logger().Fatal("Failed to start db instrumentation", zap.Error(err))
	}

}
