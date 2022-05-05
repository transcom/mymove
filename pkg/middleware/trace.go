package middleware

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/trace"
)

// The traceid needs to be formatted like so
// https://docs.aws.amazon.com/xray/latest/devguide/xray-api-sendingdata.html
//
func awsXrayIDFromBytes(data []byte) (string, error) {
	if 16 != len(data) {
		return "",
			fmt.Errorf("AWS XRay ID must be exactly 16 bytes long, got %d bytes",
				len(data))
	}
	time := hex.EncodeToString(data[0:4])
	tid := hex.EncodeToString(data[4:])
	return fmt.Sprintf("1-%s-%s", time, tid), nil
}

const traceHeader = "X-MILMOVE-TRACE-ID"

// Trace returns a trace middleware that injects a unique trace id into every request.
func Trace(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			span := sdktrace.SpanFromContext(r.Context())
			var id uuid.UUID
			var xrayID string
			ctx := r.Context()
			if span.SpanContext().HasTraceID() {
				traceID := span.SpanContext().TraceID().String()
				// Now try to reformat the span's traceId for the
				// milmove_trace_id and AWS Xray
				bytes, err := hex.DecodeString(traceID)
				if err != nil {
					logger.Warn("Cannot hex decode span traceid", zap.Error(err))
				} else {
					id, err = uuid.FromBytes(bytes)
					if err != nil {
						logger.Warn("Cannot create uuid from span traceid", zap.Error(err))
					} else {
						// If we have a span and a traceid that can be
						// converted to an AWS X-Ray ID, include that
						// in the request context so that when we
						// create the logger in the ContextLogger
						// middleware, it can include the AWS X-Ray ID
						// in all logs
						xrayID, err = awsXrayIDFromBytes(id.Bytes())
						if err == nil {
							// add the xray as an attribute for easier correlation
							span.SetAttributes(attribute.String("transcom.milmove.xray.id",
								xrayID))
							ctx = trace.AwsXrayNewContext(ctx, xrayID)
						} else {
							logger.Warn("Cannot create AWS XRay ID from span traceid",
								zap.Error(err))
						}
					}
				}

			}

			// if we don't have an id, maybe tracing isn't enabled or
			// maybe the span traceid was in a bogus format
			if id.IsNil() {
				var err error
				id, err = uuid.NewV4()
				if err != nil {
					logger.Error(errors.Wrap(err, "error creating trace id").Error())
					next.ServeHTTP(w, r)
					return
				}
			}

			// Also insert as a key, value pair in the http request
			// context
			// This is needed for the ContextLogger middleware
			ctx = trace.NewContext(ctx, id)

			strID := id.String()

			// decorate the span with the milmove trace formatted as uuid
			span.SetAttributes(attribute.String("transcom.milmove.trace.uuid", strID))

			// Let a caller see what the traceID is
			w.Header().Add(traceHeader, strID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
