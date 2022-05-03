package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"

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

	meterProvider := global.MeterProvider()

	dbMeter := meterProvider.Meter("github.com/transcom/mymove/db",
		metric.WithInstrumentationVersion(dbTelemetryVersion))

	poolInUse, _ := dbMeter.AsyncInt64().UpDownCounter(dbPoolInUseName, instrument.WithDescription(dbPoolInUseDesc))
	poolInUseInst := []instrument.Asynchronous{poolInUse}

	err := dbMeter.RegisterCallback(poolInUseInst,
		func(ctx context.Context) {
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				poolInUse.Observe(ctx, int64(dbStats.InUse))
			}
		})

	if err != nil {
		appCtx.Logger().Fatal("Failed to start db instrumentation", zap.Error(err))
	}

	poolIdle, _ := dbMeter.AsyncInt64().UpDownCounter(dbPoolIdleName, instrument.WithDescription(dbPoolIdleDesc))
	poolIdleInst := []instrument.Asynchronous{poolIdle}

	err = dbMeter.RegisterCallback(poolIdleInst,
		func(ctx context.Context) {
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				poolIdle.Observe(ctx, int64(dbStats.Idle))
			}
		})

	if err != nil {
		appCtx.Logger().Fatal("Failed to start db instrumentation", zap.Error(err))
	}

	dbWait, _ := dbMeter.AsyncInt64().UpDownCounter(dbWaitDurationName, instrument.WithDescription(dbWaitDurationDesc))
	dbWaitInst := []instrument.Asynchronous{dbWait}

	err = dbMeter.RegisterCallback(dbWaitInst,
		func(ctx context.Context) {
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				dbWait.Observe(ctx, int64(dbStats.WaitDuration))
			}
		})

	if err != nil {
		appCtx.Logger().Fatal("Failed to start db instrumentation", zap.Error(err))
	}

}
