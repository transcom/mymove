package telemetry

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func (suite *TelemetrySuite) TestDBStatsObserver() {
	// use stdout metric to see what is reported

	config := &Config{
		Enabled:          true,
		Endpoint:         "stdout",
		SamplingFraction: 1,
		CollectSeconds:   0,
		EnvironmentName:  "test",
	}
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()

	r, fakeStdout, err := os.Pipe()
	suite.NoError(err)
	os.Stdout = fakeStdout

	shutdownFn := Init(suite.Logger(), config)
	defer shutdownFn()

	err = RegisterDBStatsObserver(suite.AppContextForTest(), config)
	suite.Assert().NoError(err)
	mp := otel.GetMeterProvider()
	ctx := context.Background()
	mmp, ok := mp.(*sdkmetric.MeterProvider)
	if !ok {
		suite.FailNow("Cannot convert global metric provider to sdkmetric.MeterProvider")
	}
	// flush to export data
	suite.NoError(mmp.Shutdown(ctx))

	suite.NoError(fakeStdout.Close())
	bytes, err := io.ReadAll(r)
	suite.NoError(err)
	suite.NoError(r.Close())
	var metricData metricdata.ResourceMetrics
	err = json.Unmarshal(bytes, &metricData)
	// this will always fail because of how otel defines private
	// interfaces for aggregations
	suite.NotNil(err)
	suite.Equal(1, len(metricData.ScopeMetrics))
	// currently recording 6 db metrics
	suite.Equal(6, len(metricData.ScopeMetrics[0].Metrics))
}
