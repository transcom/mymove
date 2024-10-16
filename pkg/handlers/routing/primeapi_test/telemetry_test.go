package primeapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/telemetry/metrictest"
)

func (suite *PrimeAPISuite) TestTelemetry() {
	suite.Run("Authorized prime v1/move-task-orders", func() {
		// The NewAuthenticatedPrimeRequest method adds a header that,
		// if provided, is used by handlers.DevlocalClientCertMiddleware
		clientCert := factory.BuildClientCert(suite.DB(), nil, nil)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		targetURL := fmt.Sprintf("/prime/v1/move-task-orders/%s", move.ID.String())
		routePattern := "/prime/v1/move-task-orders/{moveID}"
		req := suite.NewAuthenticatedPrimeRequest("GET", targetURL, nil, clientCert)
		rr := httptest.NewRecorder()

		telemetryConfig := &telemetry.Config{
			Enabled:          true,
			Endpoint:         "memory",
			SamplingFraction: 1.0,
			CollectSeconds:   0,
			EnvironmentName:  "test",
		}

		shutdownFn, spanExporter, metricExporter := telemetry.Init(suite.Logger(), telemetryConfig)
		defer shutdownFn()

		routingConfig := suite.RoutingConfig()
		siteHandler := suite.SetupCustomSiteHandlerWithTelemetry(routingConfig, telemetryConfig)
		siteHandler.ServeHTTP(rr, req)
		suite.Equal(http.StatusOK, rr.Code)

		ctx := context.Background()
		tp := otel.GetTracerProvider()
		sdktp, ok := tp.(*sdktrace.TracerProvider)
		if !ok {
			suite.FailNow("Cannot cast tracer provider")
		}
		sdktp.ForceFlush(ctx)

		mp := otel.GetMeterProvider()
		mmp, ok := mp.(*sdkmetric.MeterProvider)
		if !ok {
			suite.FailNow("Cannot convert global metric provider to sdkmetric.MeterProvider")
		}
		// flush to export data
		suite.NoError(mmp.ForceFlush(ctx))

		mse, ok := spanExporter.(*tracetest.InMemoryExporter)
		suite.FatalTrue(ok)
		spans := mse.GetSpans()

		mme, ok := metricExporter.(*metrictest.InMemoryExporter)
		suite.FatalTrue(ok)
		metrics := mme.GetMetrics()

		suite.Equal(1, len(spans))
		suite.EqualServerName(spans[0].Name)
		suite.Contains(spans[0].Attributes, semconv.HTTPTarget(req.URL.String()))
		suite.Contains(spans[0].Attributes, semconv.HTTPRoute(routePattern))

		suite.Equal(1, len(metrics))

		// this is mostly duplicated from the
		// telemetry.OtelHTTPMiddlewareTrace
		//
		// but this is an full stack test, and we've seen some bugs
		// that crept in when running full stack because of changing
		// assumptions of middleware, etc
		metricExport := metrics[0]
		// 2 metric groups exported, one for our custom mymove/request and one
		// from contrib/otelhttp
		suite.Equal(2, len(metricExport.ScopeMetrics))

		var requestCountScope metricdata.ScopeMetrics
		var otelHTTPScope metricdata.ScopeMetrics

		// figure out which metric is which
		for i := range metricExport.ScopeMetrics {
			m := metricExport.ScopeMetrics[i]
			if m.Scope.Name == telemetry.RequestTelemetryName {
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
		suite.Contains(attrSlice, semconv.HTTPRoute(routePattern))
		suite.Contains(attrSlice, semconv.HTTPStatusCode(rr.Code))
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
		suite.Contains(attrSlice, semconv.HTTPRoute(routePattern))
		// we don't want the metric to have the target, as that has high
		// cardinality and will cause memory usage to explode
		suite.False(durationAgg.DataPoints[0].Attributes.HasValue(semconv.HTTPTargetKey))
	})
}
