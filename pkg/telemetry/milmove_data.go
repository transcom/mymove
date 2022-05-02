package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

const milmoveDataTelemetryVersion = "0.1"

type tableStatGauge struct {
	liveGauge asyncint64.Gauge
	deadGauge asyncint64.Gauge
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
`

func registerTableLiveDeadCallback(appCtx appcontext.AppContext, meter metric.Meter) {

	// do the query once at register time to get the list of tables
	// and create the gauges

	allStats := []pgStatLiveDead{}
	err := appCtx.DB().RawQuery(liveDeadQuery).All(&allStats)
	if err != nil {
		appCtx.Logger().Fatal("Cannot get initial list of tables for table stats", zap.Error(err))
		return
	}
	tableNameGaugeMap := make(map[string]tableStatGauge)
	tableAllInstruments := []instrument.Asynchronous{}

	for i := range allStats {
		tableName := allStats[i].TableName
		liveGaugeName := "milmovedata." + tableName + ".live"
		liveGaugeDesc := "The total number of live tuples in the " + tableName + " table"

		deadGaugeName := "milmovedata." + tableName + ".dead"
		deadGaugeDesc := "The total number of dead tuples in the " + tableName + " table"
		liveGauge, lerr := meter.AsyncInt64().Gauge(liveGaugeName, instrument.WithDescription(liveGaugeDesc))
		if lerr != nil {
			appCtx.Logger().Error("Error creating live guage", zap.Any("liveGauge", liveGaugeName), zap.Error(lerr))
			continue
		}
		deadGauge, derr := meter.AsyncInt64().Gauge(deadGaugeName, instrument.WithDescription(deadGaugeDesc))
		if derr != nil {
			appCtx.Logger().Error("Error creating dead guage", zap.Any("deadGauge", deadGaugeName), zap.Error(derr))
			continue
		}
		tableNameGaugeMap[tableName] = tableStatGauge{
			liveGauge: liveGauge,
			deadGauge: deadGauge,
		}
		tableAllInstruments = append(tableAllInstruments, liveGauge)
		tableAllInstruments = append(tableAllInstruments, deadGauge)
	}

	err = meter.RegisterCallback(tableAllInstruments,
		func(ctx context.Context) {

			allStats := []pgStatLiveDead{}
			aerr := appCtx.DB().RawQuery(liveDeadQuery).All(&allStats)
			if aerr != nil {
				appCtx.Logger().Fatal("Cannot get live/dead stats", zap.Error(aerr))
				return
			}

			for i := range allStats {
				tableName := allStats[i].TableName
				tableNameGaugeMap[tableName].liveGauge.Observe(ctx, allStats[i].LiveTupleCount)
				tableNameGaugeMap[tableName].deadGauge.Observe(ctx, allStats[i].DeadTupleCount)
			}

		})

	if err != nil {
		appCtx.Logger().Fatal("Failed to register live/dead callback", zap.Error(err))
	}

}

func RegisterMilmoveDataObserver(appCtx appcontext.AppContext, config *Config) {
	if !config.Enabled {
		return
	}

	meterProvider := global.MeterProvider()

	milmoveDataMeter := meterProvider.Meter("github.com/transcom/mymove/data",
		metric.WithInstrumentationVersion(milmoveDataTelemetryVersion))

	registerTableLiveDeadCallback(appCtx, milmoveDataMeter)

}
