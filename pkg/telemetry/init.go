package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config defines the config necessary to enable telemetry gathering
type Config struct {
	Enabled          bool
	Endpoint         string
	UseXrayID        bool
	SamplingFraction float64
	CollectSeconds   int
	ReadEvents       bool
	WriteEvents      bool
}

const (
	defaultCollectSeconds = 30
)

// Init currently the target for distributed tracing / opentelemetry is
// local development environments, but this may change in the future
// to include hosted/deployed environments
func Init(logger *zap.Logger, config *Config) (shutdown func()) {
	ctx := context.Background()
	shutdown = func() {}

	logger.Info("Configuring tracing", zap.Any("TelemetryConfig", config))
	if !config.Enabled {
		tp := trace.NewNoopTracerProvider()
		otel.SetTracerProvider(tp)
		global.SetMeterProvider(metric.NewNoopMeterProvider())
		logger.Info("opentelemetry not enabled")
		return shutdown
	}

	var spanExporter sdktrace.SpanExporter
	var metricExporter export.Exporter

	var err error

	switch config.Endpoint {
	case "stdout":
		spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			logger.Error("unable to create otel stdout span exporter", zap.Error(err))
			break
		}
		metricExporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			logger.Error("unable to create otel stdout metric exporter", zap.Error(err))
			break
		}
	default:
		spanClient := otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
		)
		spanExporter, err = otlptrace.New(ctx, spanClient)
		if err != nil {
			logger.Error("failed to create otel trace exporter", zap.Error(err))
		}
		metricClient := otlpmetricgrpc.NewClient(
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpoint(config.Endpoint),
		)
		// use MetricExportKindSelector to prevent memory leak?
		// https://github.com/open-telemetry/opentelemetry-go/issues/2225#issuecomment-915517182
		metricExporter, err = otlpmetric.New(ctx, metricClient,
			otlpmetric.WithMetricAggregationTemporalitySelector(aggregation.DeltaTemporalitySelector()))
		if err != nil {
			logger.Error("failed to create otel metric exporter", zap.Error(err))
		}

	}
	// Create a tracer provider that processes spans using a
	// batch-span-processor.
	bsp := sdktrace.NewBatchSpanProcessor(spanExporter)

	sampler := sdktrace.TraceIDRatioBased(config.SamplingFraction)
	var idGenerator sdktrace.IDGenerator
	if config.UseXrayID {
		idGenerator = xray.NewIDGenerator()
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("milmove"))),
		sdktrace.WithSampler(sampler),
		sdktrace.WithIDGenerator(idGenerator),
		sdktrace.WithSpanProcessor(bsp),
	)
	// Instantiate a new ECS resource detector
	ecsResourceDetector := ecs.NewResourceDetector()
	ecsResource, err := ecsResourceDetector.Detect(ctx)
	if err != nil {
		logger.Error("failed to create ECS resource detector", zap.Error(err))
	}

	// Create pusher for metrics that runs in the background and pushes
	// metrics periodically.
	collectSeconds := config.CollectSeconds
	if collectSeconds == 0 {
		collectSeconds = defaultCollectSeconds
	}
	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			metricExporter,
		),
		controller.WithResource(ecsResource),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(time.Duration(collectSeconds)*time.Second),
	)
	err = pusher.Start(ctx)
	if err != nil {
		logger.Error("failed to start the metric controller", zap.Error(err))
	}

	logger.Info("emit tracing to local opentelemetry collector at " + config.Endpoint)
	shutdown = func() {
		if err = spanExporter.Shutdown(ctx); err != nil {
			logger.Error("shutdown problems with tracing exporter", zap.Error(err))
		}
		if err = pusher.Stop(ctx); err != nil {
			logger.Error("shutdown problems with metrics pusher", zap.Error(err))
		}
	}

	otel.SetTracerProvider(tp)
	global.SetMeterProvider(pusher)
	if config.UseXrayID {
		propagation.NewCompositeTextMapPropagator(
			xray.Propagator{},
			propagation.Baggage{},
			propagation.TraceContext{},
		)
	} else {
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.Baggage{},
				propagation.TraceContext{},
			),
		)
	}

	return shutdown
}
