package telemetry

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/transcom/mymove/pkg/telemetry/metrictest"
)

func (suite *TelemetrySuite) TestMilmoveDataObserver() {
	// use memory metric to see what is reported

	config := &Config{
		Enabled:          true,
		Endpoint:         "memory",
		SamplingFraction: 1,
		CollectSeconds:   0,
		EnvironmentName:  "test",
	}

	shutdownFn, _, metricExporter := Init(suite.Logger(), config)
	defer shutdownFn()

	err := RegisterMilmoveDataObserver(suite.AppContextForTest(), config)
	suite.FatalNoError(err)
	mp := otel.GetMeterProvider()
	ctx := context.Background()
	mmp, ok := mp.(*sdkmetric.MeterProvider)
	if !ok {
		suite.FailNow("Cannot convert global metric provider to sdkmetric.MeterProvider")
	}
	// flush to export data
	suite.NoError(mmp.ForceFlush(ctx))

	mme, ok := metricExporter.(*metrictest.InMemoryExporter)
	suite.FatalTrue(ok)
	metrics := mme.GetMetrics()

	// one metric export
	suite.Equal(1, len(metrics))
	metricData := metrics[0]

	suite.Equal(1, len(metricData.ScopeMetrics))
	allMetrics := metricData.ScopeMetrics[0].Metrics
	suite.Greater(len(allMetrics), 100)

	foundUsersLive := false
	foundUsersDead := false
	for i := range allMetrics {
		if strings.Contains(allMetrics[i].Name, "users.live") {
			foundUsersLive = true
		}
		if strings.Contains(allMetrics[i].Name, "users.dead") {
			foundUsersDead = true
		}
	}

	suite.True(foundUsersLive, "Missing live users metric")
	suite.True(foundUsersDead, "Missing dead users metric")
}
