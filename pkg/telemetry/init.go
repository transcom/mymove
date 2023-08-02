package telemetry

import (
	"context"
	"os"
	"time"

	"github.com/go-logr/zapr"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/telemetry/metrictest"
)

// Config defines the config necessary to enable telemetry gathering
type Config struct {
	Enabled          bool
	Endpoint         string
	HTTPEndpoint     string
	UseXrayID        bool
	SamplingFraction float64
	CollectSeconds   int
	ReadEvents       bool
	WriteEvents      bool
	EnvironmentName  string
}

const (
	defaultCollectSeconds = 30
)

func configNoopTelemetry(logger *zap.Logger) {
	otel.SetTracerProvider(trace.NewNoopTracerProvider())
	otel.SetMeterProvider(noop.NewMeterProvider())
	logger.Info("opentelemetry not enabled")
}

// Init sets up open telemetry as specified by config. It returns a
// shutdown function, and also the span and metric exporters. The
// latter two are useful in testing, but would almost certainly be
// ignored in production
//
// If setting up telemetry fails, the failures will be logged, but it
// will not prevent the server from starting. Thus, no error is
// returned here
func Init(logger *zap.Logger, config *Config) (func(), sdktrace.SpanExporter, sdkmetric.Exporter) {

	ctx := context.Background()
	shutdown := func() {}

	logger.Info("Configuring open telemetry", zap.Any("TelemetryConfig", config))
	if !config.Enabled {
		configNoopTelemetry(logger)
		return shutdown, nil, nil
	}

	// convert our zap logger to the go-logr interface expected by
	// otel, but only log otel errors
	otel.SetLogger(zapr.NewLogger(logger))

	// explicitly set error handler
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error("opentelemetry error", zap.Error(err))
	}))

	var spanExporter sdktrace.SpanExporter
	var metricExporter sdkmetric.Exporter

	var err error

	// overloading Endpoint this way isn't great
	switch config.Endpoint {
	case "memory":
		spanExporter = tracetest.NewInMemoryExporter()
		metricExporter = metrictest.NewInMemoryExporter()
	case "stdout":
		// explictly call WithWriter so we can override os.Stdout in tests
		spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(os.Stdout))
		if err != nil {
			logger.Error("unable to create otel stdout span exporter", zap.Error(err))
			// setting up telemetry is not fatal to server startup
			return shutdown, nil, nil
		}
		// stdoutmetric now pretty prints by default
		metricExporter, err = stdoutmetric.New()
		if err != nil {
			// setting up telemetry is not fatal to server startup
			logger.Error("unable to create otel stdout metric exporter", zap.Error(err))
			shutdownErr := spanExporter.Shutdown(ctx)
			if shutdownErr != nil {
				logger.Error("unable to shutdown stdout span exporter", zap.Error(shutdownErr))
			}
			configNoopTelemetry(logger)
			return shutdown, nil, nil
		}
	default:
		spanClient := otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
		)
		spanExporter, err = otlptrace.New(ctx, spanClient)
		if err != nil {
			// setting up telemetry is not fatal to server startup
			logger.Error("failed to create otel trace exporter", zap.Error(err))
			return shutdown, nil, nil
		}
		metricExporter, err = otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpoint(config.Endpoint),
		)
		if err != nil {
			// setting up telemetry is not fatal to server startup
			logger.Error("failed to create otel metric client", zap.Error(err))
			return shutdown, nil, nil
		}

	}
	// Create a tracer provider that processes spans using a
	// batch-span-processor.
	bsp := sdktrace.NewBatchSpanProcessor(spanExporter, sdktrace.WithBatchTimeout(time.Duration(config.CollectSeconds*int(time.Second))))

	sampler := sdktrace.TraceIDRatioBased(config.SamplingFraction)
	resourceAttrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String("milmove-server"),
		semconv.DeploymentEnvironmentKey.String(config.EnvironmentName)}

	// Instantiate a new ECS resource detector
	ecsResourceDetector := ecs.NewResourceDetector()
	ecsResource, err := ecsResourceDetector.Detect(ctx)
	if err != nil {
		logger.Warn("failed to create ECS resource detector", zap.Error(err))
	} else {
		if ecsResource.Attributes() != nil {
			logger.Info("ECS resource for telemetry", zap.Any("attributes", ecsResource.Attributes()))
			resourceAttrs = append(resourceAttrs, ecsResource.Attributes()...)
		}
	}

	var idGenerator sdktrace.IDGenerator

	// we could consider automatically using xray if running in ECS,
	// but they are technically orthogonal
	if config.UseXrayID {
		idGenerator = xray.NewIDGenerator()
	}

	// only add a single sdktrace.WithResource option, as adding more
	// than one just overwrites earlier resources
	milmoveResource := resource.NewWithAttributes(semconv.SchemaURL, resourceAttrs...)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(milmoveResource),
		sdktrace.WithSampler(sampler),
		sdktrace.WithIDGenerator(idGenerator),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Create pusher for metrics that runs in the background and pushes
	// metrics periodically.
	collectSeconds := config.CollectSeconds
	if collectSeconds == 0 {
		collectSeconds = defaultCollectSeconds
	}
	pr := sdkmetric.NewPeriodicReader(metricExporter,
		sdkmetric.WithInterval(time.Duration(collectSeconds)*time.Second),
	)

	// create a view to filter otelhttp attributes; otherwise we have
	// a memory leak as otel tracks attributes with an infinite number
	// of values (e.g. user-agent)
	//
	// inspired by
	// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/3071#issuecomment-1419366500
	//

	otelhttpView := sdkmetric.NewView(
		sdkmetric.Instrument{
			Scope: instrumentation.Scope{
				// this constant is not exported by otelhttp
				Name: "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
			},
		},
		sdkmetric.Stream{
			AttributeFilter: allowedHTTPRequestAttributeFilter,
		},
	)
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(milmoveResource),
		sdkmetric.WithReader(pr),
		sdkmetric.WithView(otelhttpView),
	)

	logger.Info("emitting tracing to local opentelemetry collector at " + config.Endpoint)
	shutdown = func() {
		if err = spanExporter.Shutdown(ctx); err != nil {
			logger.Error("shutdown problems with tracing exporter", zap.Error(err))
		}
		if err = metricExporter.Shutdown(ctx); err != nil {
			logger.Error("shutdown problems with metrics pusher", zap.Error(err))
		}
		logger.Info("Shutting down telemetry")
	}

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	if config.UseXrayID {
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				xray.Propagator{},
				propagation.Baggage{},
				propagation.TraceContext{},
			),
		)
	} else {
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.Baggage{},
				propagation.TraceContext{},
			),
		)
	}

	return shutdown, spanExporter, metricExporter
}
