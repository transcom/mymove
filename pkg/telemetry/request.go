package telemetry

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.uber.org/zap"
)

type RequestTelemetry struct {
	requestCounter metric.Int64Counter
}

const (
	RequestTelemetryName    = "github.com/transcom/mymove/request"
	RequestTelemetryVersion = "0.1"
)

// NewRequestTelemetry provides a way for the request logger to
// provide stats. If we want accurate request counts with dimensions,
// this seems to be the best way to do it
//
// # According to the cloudwatch concepts documentation
//
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/cloudwatch_concepts.html
//
//	CloudWatch treats each unique combination of dimensions as a
//	separate metric, even if the metrics have the same metric name.
//
// This means if we try to use the "Sample count" statistic in cloud
// watch, it will count across all dimensions. It doesn't seem
// possible to get the count without dimensions.
//
// Increment a request count with the same dimensions
func NewRequestTelemetry(logger *zap.Logger) *RequestTelemetry {
	meterProvider := otel.GetMeterProvider()

	requestMeter := meterProvider.Meter(RequestTelemetryName,
		metric.WithInstrumentationVersion(RequestTelemetryVersion))

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

var allowedHTTPRequestAttributes = map[attribute.Key]bool{
	semconv.HTTPMethodKey:     true,
	semconv.HTTPRouteKey:      true,
	semconv.HTTPSchemeKey:     true,
	semconv.HTTPStatusCodeKey: true,
	semconv.NetHostNameKey:    true,
}

func allowedHTTPRequestAttributeFilter(kv attribute.KeyValue) bool {
	_, ok := allowedHTTPRequestAttributes[kv.Key]
	return ok
}

func (rt *RequestTelemetry) IncrementRequestCount(r *http.Request, routePattern string, statusCode int) {

	serverAttributes := httpconv.ServerRequest(r.Host, r)
	metricAttributes := []attribute.KeyValue{}

	for i := range serverAttributes {
		attr := serverAttributes[i]
		if allowedHTTPRequestAttributeFilter(attr) {
			metricAttributes = append(metricAttributes, attr)
		}
	}

	if routePattern != "" {
		metricAttributes = append(metricAttributes, semconv.HTTPRoute(routePattern))
	}

	if statusCode > 0 {
		metricAttributes = append(metricAttributes,
			semconv.HTTPStatusCode(statusCode))
	}
	o := metric.WithAttributes(metricAttributes...)

	rt.requestCounter.Add(r.Context(), 1, o)
}
