package telemetry

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

type RequestTelemetry struct {
	requestCounter metric.Int64Counter
}

const requestTelemetryVersion = "0.1"

func NewRequestTelemetry(logger *zap.Logger) *RequestTelemetry {
	meterProvider := global.MeterProvider()

	requestMeter := meterProvider.Meter("github.com/transcom/mymove/request",
		metric.WithInstrumentationVersion(requestTelemetryVersion))

	requestCounter, err := requestMeter.Int64Counter("http.server.request_count",
		metric.WithDescription("Count of http requests"),
	)

	if err != nil {
		logger.Error("Error registering request counter", zap.Error(err))
		return nil
	}

	return &RequestTelemetry{
		requestCounter: requestCounter,
	}
}

func (rt *RequestTelemetry) CountRequest(r *http.Request, metrics httpsnoop.Metrics) {
	metricAttributes := httpconv.ServerRequest(r.Host, r)
	metricAttributes = append(metricAttributes,
		semconv.HTTPStatusCode(metrics.Code))
	o := metric.WithAttributes(metricAttributes...)
	rt.requestCounter.Add(r.Context(), 1, o)
}
