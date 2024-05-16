package telemetry

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/transcom/mymove/pkg/telemetry/metrictest"
)

// nolint:staticcheck
func (suite *TelemetrySuite) TestInitConfigDisabled() {
	config := &Config{
		Enabled: false,
	}
	Init(suite.Logger(), config)

	suite.Equal(trace.NewNoopTracerProvider(), otel.GetTracerProvider())
	suite.Equal(noop.NewMeterProvider(), otel.GetMeterProvider())
}

func (suite *TelemetrySuite) TestInitConfigStdoutTrace() {
	config := &Config{
		Enabled:          true,
		Endpoint:         "stdout",
		SamplingFraction: 1,
		EnvironmentName:  "test",
	}
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()

	r, fakeStdout, err := os.Pipe()
	suite.NoError(err)
	os.Stdout = fakeStdout

	shutdownFn, _, _ := Init(suite.Logger(), config)
	defer shutdownFn()

	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("test_tracer", trace.WithSchemaURL("url"),
		trace.WithInstrumentationVersion("1.0"))
	ctx := context.Background()
	testSpanName := "test_span"
	ctx, span := tracer.Start(ctx, testSpanName,
		trace.WithAttributes(attribute.KeyValue{
			Key:   "key",
			Value: attribute.StringValue("Value"),
		}))
	now := time.Now()
	span.End(trace.WithTimestamp(now))
	sdktp, ok := tp.(*sdktrace.TracerProvider)
	if !ok {
		suite.FailNow("Cannot cast tracer provider")
	}
	sdktp.ForceFlush(ctx)
	suite.NoError(fakeStdout.Close())
	bytes, err := io.ReadAll(r)
	suite.NoError(err)
	suite.NoError(r.Close())
	var spanData map[string]interface{}
	suite.NoError(json.Unmarshal(bytes, &spanData))
	suite.Equal(spanData["Name"], testSpanName)
}

func (suite *TelemetrySuite) TestInitConfigStdoutMetric() {
	config := &Config{
		Enabled:          true,
		Endpoint:         "stdout",
		SamplingFraction: 1,
		CollectSeconds:   1,
		EnvironmentName:  "test",
	}
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()

	r, fakeStdout, err := os.Pipe()
	suite.NoError(err)
	os.Stdout = fakeStdout

	shutdownFn, _, _ := Init(suite.Logger(), config)
	defer shutdownFn()

	mp := otel.GetMeterProvider()
	meter := mp.Meter("test_meter", metric.WithSchemaURL("url"))
	counter, err := meter.Int64Counter("test_counter")
	suite.NoError(err)
	ctx := context.Background()
	counter.Add(ctx, 1)
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
	suite.Equal(1, len(metricData.ScopeMetrics[0].Metrics))
	suite.Equal("test_counter", metricData.ScopeMetrics[0].Metrics[0].Name)
}

func (suite *TelemetrySuite) TestInitConfigMemoryTrace() {
	config := &Config{
		Enabled:          true,
		Endpoint:         "memory",
		SamplingFraction: 1,
		EnvironmentName:  "test",
	}

	shutdownFn, spanExporter, _ := Init(suite.Logger(), config)
	defer shutdownFn()

	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("test_tracer", trace.WithSchemaURL("url"),
		trace.WithInstrumentationVersion("1.0"))
	ctx := context.Background()
	testSpanName := "test_span"
	ctx, span := tracer.Start(ctx, testSpanName,
		trace.WithAttributes(attribute.KeyValue{
			Key:   "key",
			Value: attribute.StringValue("Value"),
		}))
	now := time.Now()
	span.End(trace.WithTimestamp(now))
	sdktp, ok := tp.(*sdktrace.TracerProvider)
	if !ok {
		suite.FailNow("Cannot cast tracer provider")
	}
	sdktp.ForceFlush(ctx)
	mse, ok := spanExporter.(*tracetest.InMemoryExporter)
	suite.FatalTrue(ok)
	spans := mse.GetSpans()
	suite.Equal(1, len(spans))
	suite.Equal(testSpanName, spans[0].Name)
}

func (suite *TelemetrySuite) TestInitConfigMemoryMetric() {
	config := &Config{
		Enabled:          true,
		Endpoint:         "memory",
		SamplingFraction: 1,
		CollectSeconds:   1,
		EnvironmentName:  "test",
	}

	shutdownFn, _, metricExporter := Init(suite.Logger(), config)
	defer shutdownFn()

	mp := otel.GetMeterProvider()
	meter := mp.Meter("test_meter", metric.WithSchemaURL("url"))
	counter, err := meter.Int64Counter("test_counter")
	suite.NoError(err)
	ctx := context.Background()
	counter.Add(ctx, 1)
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
	suite.Equal(1, len(metricData.ScopeMetrics[0].Metrics))
	suite.Equal("test_counter", metricData.ScopeMetrics[0].Metrics[0].Name)
}
