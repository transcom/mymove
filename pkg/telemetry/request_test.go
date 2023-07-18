package telemetry

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func (suite *TelemetrySuite) TestRequestStats() {
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

	rt := NewRequestTelemetry(suite.Logger())
	suite.NotNil(rt)

	req := httptest.NewRequest("GET", "http://test.example.com/foobad", nil)
	routePattern := "/foobad"
	statusCode := 200
	rt.IncrementRequestCount(req, routePattern, statusCode)

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
	// currently recording 1 request metrics: request count
	suite.Equal(1, len(metricData.ScopeMetrics[0].Metrics))
	suite.Equal("github.com/transcom/mymove/request",
		metricData.ScopeMetrics[0].Scope.Name)
}
