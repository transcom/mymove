package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const defaultLocalCollector = "localhost:55680"

// currently the target for distributed tracing / opentelemetry is
// local development environments, but this may change in the future
// to include hosted/deployed environments
func configureTracing(logger *zap.Logger) (shutdown func()) {
	ctx := context.Background()
	shutdown = func() {}
	tp := trace.NewNoopTracerProvider()

	switch os.Getenv("OTEL_CONFIG") {
	case "STDOUT":
		prov, _, err := stdout.InstallNewPipeline([]stdout.Option{stdout.WithPrettyPrint()}, nil)
		if err != nil {
			logger.Error("unable to create otel stdout exporter", zap.Error(err))
			break
		}
		tp = prov
	case "LOCAL_COLLECTOR":
		exp, prov, _, err := otlp.NewExportPipeline(
			ctx,
			otlpgrpc.NewDriver(
				otlpgrpc.WithInsecure(),
				otlpgrpc.WithEndpoint(defaultLocalCollector),
			),
		)
		if err != nil {
			logger.Error("unable to create otel otlp grpc exporter", zap.Error(err))
			break
		}
		logger.Info("emit tracing to local opentelemetry collector at " + defaultLocalCollector)
		tp = prov
		shutdown = func() {
			if err := exp.Shutdown(ctx); err != nil {
				logger.Error("shutdown problems with tracing exporter", zap.Error(err))
			}
		}
	default:
		logger.Info("opentelemetry not enabled")
	}

	// set the global configuration
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	return shutdown
}
