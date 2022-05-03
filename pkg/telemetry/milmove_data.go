package telemetry

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

type tableStatMetrics struct {
	liveTuples instrument.Int64ObservableGauge
	deadTuples instrument.Int64ObservableGauge
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

func registerTableLiveDeadCallback(appCtx appcontext.AppContext, meter metric.Meter) error {

	// this should be configurable
	const minDuration = time.Second * 60

	var lock sync.Mutex

	// do the query once at register time to get the list of tables
	// and create the gauges

	// lock prevents a race between batch observer and instrument registration.
	lock.Lock()
	defer lock.Unlock()

	allStats := []pgStatLiveDead{}

	err := appCtx.DB().RawQuery(liveDeadQuery).All(&allStats)
	if err != nil {
		appCtx.Logger().Error("Cannot get initial list of tables for table stats", zap.Error(err))
		return err
	}
	tableNameMetricMap := make(map[string]tableStatMetrics)
	tableAllInstruments := []instrument.Asynchronous{}

	for i := range allStats {
		tableName := allStats[i].TableName
		liveMetricName := tableName + ".live"
		liveMetricDesc := "The total number of live tuples in the " + tableName + " table"

		deadMetricName := tableName + ".dead"
		deadMetricDesc := "The total number of dead tuples in the " + tableName + " table"
		liveTuples, lerr := meter.Int64ObservableGauge(
			liveMetricName,
			instrument.WithDescription(liveMetricDesc),
		)
		if lerr != nil {
			appCtx.Logger().Error("Error creating live counter", zap.Any("liveCounter", liveMetricName), zap.Error(lerr))
			return lerr
		}
		deadTuples, derr := meter.Int64ObservableGauge(
			deadMetricName,
			instrument.WithDescription(deadMetricDesc),
		)
		if derr != nil {
			appCtx.Logger().Error("Error creating dead counter", zap.Any("deadCounter", deadMetricName), zap.Error(derr))
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
		func(ctx context.Context, observer metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()

			now := time.Now()

			// round to nearest second
			diff := now.Sub(lastStats).Round(time.Second)
			if diff < minDuration {
				appCtx.Logger().Warn("Skipping data telemetry update")
				return nil
			}

			allStats := []pgStatLiveDead{}
			aerr := appCtx.DB().RawQuery(liveDeadQuery).All(&allStats)
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

	meterProvider := global.MeterProvider()

	milmoveDataMeter := meterProvider.Meter("github.com/transcom/mymove/data",
		metric.WithInstrumentationVersion("0.4"))

	return registerTableLiveDeadCallback(appCtx, milmoveDataMeter)

}
