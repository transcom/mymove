package telemetry

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

// When deployed in production, we have multiple deployed servers
// running, but we only need one of them doing the queries to the
// database. As an optimization, maybe we can have only one of them
// running.
//
// We don't want to try to use locking as this telemetry code doesn't
// hold a database connection open
//
// The hacky approach is to look at the current connections in the
// pg_stat_activity table and have the server with the
// min(client_addr) do the work.
//
// This should also allow for a new server to take over if another
// server is restarted.
//
// If this doesn't work, the worst thing is the same stats will be
// reported from multiple clients, which is ok.
const currentClientAddrQuery = `
SELECT client_addr FROM pg_stat_activity WHERE pid = pg_backend_pid()
`
const isMinClientAddrQuery = `
SELECT min(client_addr) = ? FROM pg_stat_activity
`

type tableStatMetrics struct {
	liveTuples metric.Int64ObservableGauge
	deadTuples metric.Int64ObservableGauge
}

type pgStatLiveDead struct {
	TableName      string `db:"relname"`
	LiveTupleCount int64  `db:"n_live_tup"`
	DeadTupleCount int64  `db:"n_dead_tup"`

	//
	// should probably track these too in the future
	//
	// SeqScanCount     int64  `db:"seq_scan"`
	// IdxScanCount     int64  `db:"idx_scan"`
	// InsertTupleCount int64  `db:"n_tup_ins"`
	// UpdateTupleCount int64  `db:"n_tup_upd"`
	// DeleteTupleCount int64  `db:"n_tup_del"`
}

// postgres exposes the stats from a table
// see https://dba.stackexchange.com/a/193239 for more details
const liveDeadQuery = `
SELECT relname, n_live_tup, n_dead_tup
FROM pg_stat_user_tables
ORDER BY relname
`

const metricPrefix = "mymove.data."

func registerTableLiveDeadCallback(appCtx appcontext.AppContext, meter metric.Meter, config *Config) error {
	minDuration := time.Duration(config.CollectSeconds * int(time.Second))

	var lock sync.Mutex

	// do the query once at register time to get the list of tables
	// and create the gauges

	// lock prevents a race between batch observer and instrument registration.
	lock.Lock()
	defer lock.Unlock()

	var currentIP string
	err := appCtx.DB().RawQuery(currentClientAddrQuery).First(&currentIP)
	if err != nil {
		appCtx.Logger().Error("Cannot get current client ip", zap.Error(err))
		return err
	}

	// do the query once at register time to get the list of tables
	// and create the gauges. All servers do this once on startup

	allStats := []pgStatLiveDead{}

	err = appCtx.DB().RawQuery(liveDeadQuery).All(&allStats)
	if err != nil {
		appCtx.Logger().Error("Cannot get initial list of tables for table stats", zap.Error(err))
		return err
	}
	tableNameMetricMap := make(map[string]tableStatMetrics)
	tableAllInstruments := []metric.Observable{}

	for i := range allStats {
		tableName := allStats[i].TableName
		liveMetricName := metricPrefix + tableName + ".live"
		liveMetricDesc := "The total number of live tuples in the " + tableName + " table"

		deadMetricName := metricPrefix + tableName + ".dead"
		deadMetricDesc := "The total number of dead tuples in the " + tableName + " table"
		liveTuples, lerr := meter.Int64ObservableGauge(
			liveMetricName,
			metric.WithDescription(liveMetricDesc),
		)
		if lerr != nil {
			appCtx.Logger().Error("Error creating live gauge", zap.Any("liveGauge", liveMetricName), zap.Error(lerr))
			return lerr
		}
		deadTuples, derr := meter.Int64ObservableGauge(
			deadMetricName,
			metric.WithDescription(deadMetricDesc),
		)
		if derr != nil {
			appCtx.Logger().Error("Error creating dead gauge", zap.Any("deadGauge", deadMetricName), zap.Error(derr))
			return derr
		}
		tableNameMetricMap[tableName] = tableStatMetrics{
			liveTuples: liveTuples,
			deadTuples: deadTuples,
		}
		tableAllInstruments = append(tableAllInstruments, liveTuples)
		tableAllInstruments = append(tableAllInstruments, deadTuples)
	}

	lastStats := time.Now()
	_, err = meter.RegisterCallback(
		func(_ context.Context, observer metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()

			now := time.Now()

			// round to nearest second
			diff := now.Sub(lastStats).Round(time.Second)
			if diff < minDuration {
				appCtx.Logger().Warn("Skipping data telemetry update")
				return nil
			}

			var isMinClient bool
			aerr := appCtx.DB().RawQuery(isMinClientAddrQuery,
				currentIP).First(&isMinClient)

			if aerr != nil {
				appCtx.Logger().Fatal("Cannot get isMinClientAddr", zap.Error(aerr))
				return aerr
			}

			if !isMinClient {
				appCtx.Logger().Warn("This server min is not min client addr: skipping data telemetry update",
					zap.String("currentIp", currentIP))
				return nil
			}

			allStats := []pgStatLiveDead{}
			aerr = appCtx.DB().RawQuery(liveDeadQuery).All(&allStats)
			if aerr != nil {
				appCtx.Logger().Fatal("Cannot get live/dead stats", zap.Error(aerr))
				return aerr
			}

			for i := range allStats {
				tableName := allStats[i].TableName
				observer.ObserveInt64(tableNameMetricMap[tableName].liveTuples, allStats[i].LiveTupleCount)
				observer.ObserveInt64(tableNameMetricMap[tableName].deadTuples, allStats[i].DeadTupleCount)
			}

			lastStats = now
			return nil
		}, tableAllInstruments...)

	return err

}

func RegisterMilmoveDataObserver(appCtx appcontext.AppContext, config *Config) error {
	if !config.Enabled {
		return nil
	}

	meterProvider := otel.GetMeterProvider()

	milmoveDataMeter := meterProvider.Meter("github.com/transcom/mymove/data",
		metric.WithInstrumentationVersion("0.4"))

	return registerTableLiveDeadCallback(appCtx, milmoveDataMeter, config)

}
