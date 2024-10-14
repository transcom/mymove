package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	"github.com/transcom/mymove/pkg/telemetry/metrictest"
)

const (
	fakeURL          = "/fake/url"
	fakeRoutePattern = "/fake/{pattern}"
	fakeStatusCode   = 201
)

func (suite *TelemetrySuite) runOtelHTTPMiddleware(samplingFraction float64) (tracetest.SpanStubs, []metricdata.ResourceMetrics) {
	config := &Config{
		Enabled:          true,
		Endpoint:         "memory",
		SamplingFraction: samplingFraction,
		CollectSeconds:   0,
		EnvironmentName:  "test",
	}

	shutdownFn, spanExporter, metricExporter := Init(suite.Logger(), config)
	defer shutdownFn()

	mw := NewOtelHTTPMiddleware(config, "server_name", suite.Logger())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fakeURL, nil)

	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		routePattern := RoutePatternFromContext(r.Context())
		if routePattern != nil {
			*routePattern = fakeRoutePattern
		}
		statusCode := StatusCodeFromContext(r.Context())
		if statusCode != nil {
			*statusCode = fakeStatusCode
		}
	})

	// run the middleware with the request
	mw(next).ServeHTTP(rr, req)

	ctx := context.Background()
	tp := otel.GetTracerProvider()
	sdktp, ok := tp.(*sdktrace.TracerProvider)
	if !ok {
		suite.FailNow("Cannot cast tracer provider")
	}
	sdktp.ForceFlush(ctx)

	mse, ok := spanExporter.(*tracetest.InMemoryExporter)
	suite.FatalTrue(ok)
	spans := mse.GetSpans()

	mp := otel.GetMeterProvider()
	mmp, ok := mp.(*sdkmetric.MeterProvider)
	if !ok {
		suite.FailNow("Cannot convert global metric provider to sdkmetric.MeterProvider")
	}
	// flush to export data
	suite.NoError(mmp.ForceFlush(ctx))

	mme, ok := metricExporter.(*metrictest.InMemoryExporter)
	suite.FatalTrue(ok)
	metrics := mme.GetMetrics()

	return spans, metrics
}

func (suite *TelemetrySuite) TestOtelHTTPMiddlewareTrace() {
	// set sampling to 1.0 to ensure the span is exported
	spans, _ := suite.runOtelHTTPMiddleware(1.0)

	suite.Equal(1, len(spans))
	suite.Contains(spans[0].Attributes, semconv.HTTPTarget(fakeURL))
	suite.Contains(spans[0].Attributes, semconv.HTTPRoute(fakeRoutePattern))
}

func (suite *TelemetrySuite) TestOtelHTTPMiddlewareMetrics() {
	// set sampling to 0 to turn off tracing
	_, metrics := suite.runOtelHTTPMiddleware(0.0)

	// one metric export
	suite.Equal(1, len(metrics))
	metricExport := metrics[0]
	// 2 metric groups exported, one for our custom mymove/request and one
	// from contrib/otelhttp
	suite.Equal(2, len(metricExport.ScopeMetrics))

	var requestCountScope metricdata.ScopeMetrics
	var otelHTTPScope metricdata.ScopeMetrics

	// figure out which metric is which
	for i := range metricExport.ScopeMetrics {
		m := metricExport.ScopeMetrics[i]
		if m.Scope.Name == RequestTelemetryName {
			requestCountScope = m
		} else {
			otelHTTPScope = m
		}
	}

	suite.Equal(1, len(requestCountScope.Metrics))
	requestCountAgg := requestCountScope.Metrics[0]
	suite.Equal("http.server.request_count", requestCountAgg.Name)
	sumAgg, ok := requestCountAgg.Data.(metricdata.Sum[int64])
	suite.True(ok)
	// a single request count data point
	suite.Equal(1, len(sumAgg.DataPoints), sumAgg.DataPoints)
	attrSlice := sumAgg.DataPoints[0].Attributes.ToSlice()
	// does the metric have the custom attributes we set up?
	suite.Contains(attrSlice, semconv.HTTPRoute(fakeRoutePattern))
	suite.Contains(attrSlice, semconv.HTTPStatusCode(fakeStatusCode))
	// we don't want the metric to have the target, as that has high
	// cardinality and will cause memory usage to explode
	suite.False(sumAgg.DataPoints[0].Attributes.HasValue(semconv.HTTPTargetKey))

	expectedOtelMetrics := map[string]bool{
		"http.server.duration":      true,
		"http.server.request.size":  true,
		"http.server.response.size": true,
	}
	otelHTTPMetricNames := make(map[string]bool)
	for i := range otelHTTPScope.Metrics {
		otelHTTPMetricNames[otelHTTPScope.Metrics[i].Name] = true
	}
	suite.Equal(expectedOtelMetrics, otelHTTPMetricNames)

	// pick a single metric to make sure it has the custom attribute
	var durationAgg metricdata.Histogram[float64]
	for i := range otelHTTPScope.Metrics {
		if otelHTTPScope.Metrics[i].Name == "http.server.duration" {
			durationAgg, ok = otelHTTPScope.Metrics[i].Data.(metricdata.Histogram[float64])
			suite.FatalTrue(ok)
		}
	}
	suite.NotNil(durationAgg)
	suite.NotEmpty(durationAgg.DataPoints)
	attrSlice = durationAgg.DataPoints[0].Attributes.ToSlice()
	// does the metric have the custom attributes we set up?
	suite.Contains(attrSlice, semconv.HTTPRoute(fakeRoutePattern))
	// we don't want the metric to have the target, as that has high
	// cardinality and will cause memory usage to explode
	suite.False(durationAgg.DataPoints[0].Attributes.HasValue(semconv.HTTPTargetKey))
}
