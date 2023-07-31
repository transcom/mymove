package telemetry

import (
	"context"
	"net/http/httptest"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/transcom/mymove/pkg/telemetry/metrictest"
)

func (suite *TelemetrySuite) TestRequestStats() {
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
	suite.NoError(mmp.ForceFlush(ctx))

	mme, ok := metricExporter.(*metrictest.InMemoryExporter)
	suite.FatalTrue(ok)
	metrics := mme.GetMetrics()
	suite.Equal(1, len(metrics))

	metricData := metrics[0]
	suite.Equal(1, len(metricData.ScopeMetrics))
	// currently recording 1 request metrics: request count
	suite.Equal(1, len(metricData.ScopeMetrics[0].Metrics))
	suite.Equal("github.com/transcom/mymove/request",
		metricData.ScopeMetrics[0].Scope.Name)
}
