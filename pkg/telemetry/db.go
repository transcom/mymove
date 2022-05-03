package telemetry

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"

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
func RegisterDBStatsObserver(appCtx appcontext.AppContext, config *Config) error {
	// this should be configurable
	const minDuration = time.Second * 15

	if !config.Enabled {
		return nil
	}

	meterProvider := global.MeterProvider()

	dbMeter := meterProvider.Meter("github.com/transcom/mymove/db",
		metric.WithInstrumentationVersion(dbTelemetryVersion))

	// lock prevents a race between batch observer and instrument registration.
	var lock sync.Mutex

	lock.Lock()
	defer lock.Unlock()

	poolInUse, err := dbMeter.Int64ObservableUpDownCounter(dbPoolInUseName, instrument.WithDescription(dbPoolInUseDesc))
	if err != nil {
		return err
	}

	poolIdle, err := dbMeter.Int64ObservableUpDownCounter(dbPoolIdleName, instrument.WithDescription(dbPoolIdleDesc))
	if err != nil {
		return err
	}

	dbWait, err := dbMeter.Int64ObservableUpDownCounter(dbWaitDurationName, instrument.WithDescription(dbWaitDurationDesc))
	if err != nil {
		return err
	}

	_, err = dbMeter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				observer.ObserveInt64(poolInUse, int64(dbStats.InUse))
				observer.ObserveInt64(poolIdle, int64(dbStats.Idle))
				observer.ObserveInt64(dbWait, int64(dbStats.WaitDuration))
			}
			return dberr
		}, poolInUse, poolIdle, dbWait)

	if err != nil {
		return err
	}

	return nil
}
