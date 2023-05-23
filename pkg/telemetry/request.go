package telemetry

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

type RequestTelemetry struct {
	requestCounter         metric.Int64Counter
	serverLatencyHistogram metric.Float64Histogram
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
	serverLatencyHistogram, err := requestMeter.Float64Histogram("http.server.duration",
		metric.WithDescription("Duration of request in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		logger.Error("Error registering latency histogram", zap.Error(err))
		return nil
	}

	return &RequestTelemetry{
		requestCounter:         requestCounter,
		serverLatencyHistogram: serverLatencyHistogram,
	}
}

var ourRequestAttributes = map[attribute.Key]bool{
	semconv.HTTPMethodKey:  true,
	semconv.HTTPSchemeKey:  true,
	semconv.HTTPFlavorKey:  true,
	semconv.HTTPTargetKey:  true,
	semconv.NetHostNameKey: true,
}

func (rt *RequestTelemetry) HandleRequest(r *http.Request, metrics httpsnoop.Metrics) {
	serverAttributes := httpconv.ServerRequest(r.Host, r)

	metricAttributes := []attribute.KeyValue{}
	for i := range serverAttributes {
		attr := serverAttributes[i]
		if ok := ourRequestAttributes[attr.Key]; ok {
			metricAttributes = append(metricAttributes, attr)
		}
	}

	routeStr := ""
	route := mux.CurrentRoute(r)
	if route != nil {
		var err error
		routeStr, err = route.GetPathTemplate()
		if err != nil {
			routeStr, err = route.GetPathRegexp()
			if err != nil {
				routeStr = ""
			}
		}
	}

	if routeStr != "" {
		metricAttributes = append(metricAttributes,
			semconv.HTTPRoute(routeStr))
	}

	metricAttributes = append(metricAttributes,
		semconv.HTTPStatusCode(metrics.Code))
	o := metric.WithAttributes(metricAttributes...)

	rt.requestCounter.Add(r.Context(), 1, o)
	rt.serverLatencyHistogram.Record(r.Context(), metrics.Duration.Seconds(), o)
}
