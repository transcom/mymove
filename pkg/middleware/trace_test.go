package middleware

import (
	"net/http"
	"net/http/httptest"

	"github.com/gofrs/uuid"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *testSuite) TestTraceWithSpan() {
	telemetryConfig := &telemetry.Config{
		UseXrayID: true,
	}
	mw := Trace(telemetryConfig)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", testURL, nil)

	spanTraceID := uuid.Must(uuid.NewV4())

	// set up a test handler to extract the expected info from the
	// request context
	var traceID uuid.UUID
	var xrayID string
	var span oteltrace.Span
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		traceID = trace.FromContext(r.Context())
		span = oteltrace.SpanFromContext(r.Context())
		xrayID = trace.AwsXrayFromContext(r.Context())
	})

	// fake a span to test setting the aws xray id
	sc := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		SpanID:  [8]byte{1},
		TraceID: oteltrace.TraceID(spanTraceID.Bytes()),
	})
	req = req.WithContext(oteltrace.ContextWithSpanContext(req.Context(), sc))

	// run the middleware
	suite.do(mw, next, rr, req)
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	suite.Equal(spanTraceID, traceID)
	suite.True(span.SpanContext().HasTraceID())
	suite.NotEqual("", xrayID)
}

func (suite *testSuite) TestTraceWithoutSpan() {
	telemetryConfig := &telemetry.Config{
		UseXrayID: true,
	}
	mw := Trace(telemetryConfig)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", testURL, nil)

	// set up a test handler to extract the expected info from the
	// request context
	var traceID uuid.UUID
	var xrayID string
	var span oteltrace.Span
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		traceID = trace.FromContext(r.Context())
		span = oteltrace.SpanFromContext(r.Context())
		xrayID = trace.AwsXrayFromContext(r.Context())
	})

	// run the middleware
	suite.do(mw, next, rr, req)
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code

	// with no span, a traceID should still be set
	suite.False(traceID.IsNil())

	// but no span and no xrayID
	suite.False(span.SpanContext().HasTraceID())
	suite.Equal("", xrayID)
}

func (suite *testSuite) TestXRayIDFromBytes() {
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	id, err := awsXrayIDFromBytes(data)
	suite.NoError(err)
	suite.Equal("1-00010203-0405060708090a0b0c0d0e0f", id)
}
