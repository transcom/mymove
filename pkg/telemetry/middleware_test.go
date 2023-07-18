package telemetry

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	fakeURL          = "/fake/url"
	fakeRoutePattern = "/fake/{pattern}"
)

func (suite *TelemetrySuite) runOtelHTTPMiddleware(samplingFraction float64, postRun func()) []byte {
	config := &Config{
		Enabled:          true,
		Endpoint:         "stdout",
		SamplingFraction: samplingFraction,
		CollectSeconds:   0,
		EnvironmentName:  "test",
	}
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()

	r, fakeStdout, err := os.Pipe()
	suite.NoError(err)
	os.Stdout = fakeStdout

	shutdownFn := Init(suite.Logger(), config)
	defer shutdownFn()

	mw := NewOtelHTTPMiddleware(config, "server_name", suite.Logger())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fakeURL, nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routePattern := RoutePatternFromContext(r.Context())
		if routePattern != nil {
			*routePattern = fakeRoutePattern
		}
	})

	// run the middleware with the request
	mw(next).ServeHTTP(rr, req)

	postRun()

	suite.NoError(fakeStdout.Close())
	bytes, err := io.ReadAll(r)
	suite.NoError(err)
	suite.NoError(r.Close())

	return bytes
}

func (suite *TelemetrySuite) TestOtelHTTPMiddlewareTrace() {
	// set sampling to 1.0 to ensure the span is exported
	bytes := suite.runOtelHTTPMiddleware(1.0, func() {
		tp := otel.GetTracerProvider()
		ttp, ok := tp.(*sdktrace.TracerProvider)
		if !ok {
			suite.FailNow("Cannot convert global tracer provider to sdktrace.TracerProvider")
		}

		// flush to export data
		ctx := context.Background()
		suite.NoError(ttp.Shutdown(ctx))
	})

	// unfortunately otel makes it pretty impossible to unmarshal the
	// json
	data := string(bytes)
	suite.True(strings.Contains(data, fakeURL), data)
	suite.True(strings.Contains(data, fakeRoutePattern), data)
}

func (suite *TelemetrySuite) TestOtelHTTPMiddlewareMetrics() {
	// set sampling to 0 to turn off tracing
	bytes := suite.runOtelHTTPMiddleware(0.0, func() {
		mp := otel.GetMeterProvider()
		mmp, ok := mp.(*sdkmetric.MeterProvider)
		if !ok {
			suite.FailNow("Cannot convert global metric provider to sdkmetric.MeterProvider")
		}

		// flush to export data
		ctx := context.Background()
		suite.NoError(mmp.Shutdown(ctx))

	})

	// unfortunately otel makes it pretty impossible to unmarshal the
	// json
	data := string(bytes)
	suite.False(strings.Contains(data, fakeURL), data)
	suite.True(strings.Contains(data, fakeRoutePattern), data)
}
