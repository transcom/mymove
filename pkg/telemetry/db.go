package telemetry

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/stats"
)

const dbTelemetryVersion = "0.1"

// RegisterDBStatsObserver creates a custom metric that is updated
// automatically using an observer
func RegisterDBStatsObserver(appCtx appcontext.AppContext, config *Config) error {
	minDuration := time.Duration(config.CollectSeconds * int(time.Second))

	if !config.Enabled {
		return nil
	}

	meterProvider := otel.GetMeterProvider()

	dbMeter := meterProvider.Meter("github.com/transcom/mymove/db",
		metric.WithInstrumentationVersion(dbTelemetryVersion))

	// lock prevents a race between batch observer and instrument registration.
	var lock sync.Mutex

	lock.Lock()
	defer lock.Unlock()

	// Ideally the "go.opentelemetry.io/otel/semconv" package would
	// provide standard names for the stats below, but they don't seem to
	// See https://github.com/open-telemetry/opentelemetry-go/blob/main/semconv/v1.7.0/trace.go

	// https://pkg.go.dev/database/sql#DB.Stats
	poolInUse, err := dbMeter.Int64ObservableUpDownCounter(
		"db.pool.inuse",
		metric.WithDescription("The number of connections currently in use"))
	if err != nil {
		return err
	}

	// https://pkg.go.dev/database/sql#DB.Stats
	poolIdle, err := dbMeter.Int64ObservableUpDownCounter(
		"db.pool.idle",
		metric.WithDescription("The number of idle connections"))
	if err != nil {
		return err
	}

	// https://pkg.go.dev/database/sql#DB.Stats
	dbWait, err := dbMeter.Int64ObservableUpDownCounter(
		"db.waitduration",
		metric.WithUnit("ms"),
		metric.WithDescription("Milliseconds blocked waiting for a new connection"))
	if err != nil {
		return err
	}

	maxIdleClosed, err := dbMeter.Int64ObservableUpDownCounter(
		"db.maxidleclosed",
		metric.WithDescription("The total number of connections closed due to SetMaxIdleConns"),
	)
	if err != nil {
		return err
	}

	maxIdleTimeClosed, err := dbMeter.Int64ObservableUpDownCounter(
		"db.maxidletimeclosed",
		metric.WithDescription("The total number of connections closed due to SetConnMaxIdleTime"),
	)
	if err != nil {
		return err
	}

	maxLifetimeClosed, err := dbMeter.Int64ObservableUpDownCounter(
		"db.maxlifetimeclosed",
		metric.WithDescription("The total number of connections closed due to SetConnMaxLifetime"),
	)
	if err != nil {
		return err
	}

	lastStats := time.Now()
	_, err = dbMeter.RegisterCallback(
		func(_ context.Context, observer metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()

			now := time.Now()

			// round to nearest second
			diff := now.Sub(lastStats).Round(time.Second)
			if diff < minDuration {
				appCtx.Logger().Warn("Skipping db telemetry update")
				return nil
			}
			dbStats, dberr := stats.DBStats(appCtx)
			if dberr == nil {
				observer.ObserveInt64(poolInUse, int64(dbStats.InUse))
				observer.ObserveInt64(poolIdle, int64(dbStats.Idle))
				observer.ObserveInt64(dbWait, int64(dbStats.WaitDuration.Milliseconds()))
				observer.ObserveInt64(maxIdleClosed, dbStats.MaxIdleClosed)
				observer.ObserveInt64(maxIdleTimeClosed, dbStats.MaxIdleTimeClosed)
				observer.ObserveInt64(maxLifetimeClosed, dbStats.MaxLifetimeClosed)
				lastStats = now
			}
			return dberr
		}, poolInUse, poolIdle, dbWait, maxIdleClosed, maxIdleTimeClosed, maxLifetimeClosed)

	if err != nil {
		return err
	}

	return nil
}
