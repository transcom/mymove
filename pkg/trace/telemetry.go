package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TelemetryConfig defines the config necessary to enable telemetry gathering
type TelemetryConfig struct {
	Enabled          bool
	Endpoint         string
	UseXrayID        bool
	SamplingFraction float64
	CollectSeconds   int
}

// ConfigureTelemetry currently the target for distributed tracing / opentelemetry is
// local development environments, but this may change in the future
// to include hosted/deployed environments
func ConfigureTelemetry(logger *zap.Logger, config TelemetryConfig) (shutdown func()) {
	ctx := context.Background()
	shutdown = func() {}

	if !config.Enabled {
		tp := trace.NewNoopTracerProvider()
		otel.SetTracerProvider(tp)
		global.SetMeterProvider(metric.NoopMeterProvider{})
		logger.Info("opentelemetry not enabled")
		return shutdown
	}

	var exporter interface{}
	var err error

	switch config.Endpoint {
	case "stdout":
		exporter, err = stdout.NewExporter(stdout.WithPrettyPrint())
		if err != nil {
			logger.Error("unable to create otel stdout exporter", zap.Error(err))
			break
		}
	default:
		driver := otlpgrpc.NewDriver(
			otlpgrpc.WithInsecure(),
			otlpgrpc.WithEndpoint(config.Endpoint),
		)
		exporter, err = otlp.NewExporter(ctx, driver)
		if err != nil {
			logger.Error("failed to create otel exporter", zap.Error(err))
		}
	}
	spanExporter := exporter.(sdktrace.SpanExporter)
	// Create a tracer provider that processes spans using a
	// batch-span-processor.
	bsp := sdktrace.NewBatchSpanProcessor(spanExporter)

	sampler := sdktrace.TraceIDRatioBased(config.SamplingFraction)
	var idGenerator sdktrace.IDGenerator = nil
	if config.UseXrayID {
		idGenerator = xray.NewIDGenerator()
	}
	tp := sdktrace.NewTracerProvider(
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
		collectSeconds = 5
	}
	metricExporter := exporter.(exportmetric.Exporter)
	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
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
	global.SetMeterProvider(pusher.MeterProvider())
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
