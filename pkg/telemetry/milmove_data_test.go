package telemetry

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func (suite *TelemetrySuite) TestMilmoveDataObserver() {
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

	// use a tempfile because the data is larger than os.Pipe can handle
	fakeStdout, err := os.CreateTemp("", "milmove_data_test")
	suite.NoError(err)
	os.Stdout = fakeStdout

	defer func() {
		suite.NoError(os.Remove(fakeStdout.Name()))
	}()

	r, err := os.Open(fakeStdout.Name())
	suite.NoError(err)

	shutdownFn := Init(suite.Logger(), config)
	defer shutdownFn()

	err = RegisterMilmoveDataObserver(suite.AppContextForTest(), config)
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
