package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/trace"
)

const (
	// TelemetryEnabledFlag is the Trace Enable Flag
	TelemetryEnabledFlag string = "trace-enabled"
	// TelemetryEndpointFlag configures the endpoint used for open
	// telemetry tracing
	TelemetryEndpointFlag string = "trace-endpoint"
	// TelemetryUseXrayIDFlag enables using AWS Xray Trace IDs for open telemetry
	TelemetryUseXrayIDFlag string = "trace-use-xray-id"
	// TelemetrySamplingFractionFlag configures the percent of traces to sample
	TelemetrySamplingFractionFlag string = "trace-sampling-fraction"
	// TelemetryCollectSecondsFlag configures the metric collection
	// period in seconds
	TelemetryCollectSecondsFlag string = "trace-collect-seconds"
)

// InitTelemetryFlags initializes the open telemetry flags
func InitTelemetryFlags(flag *pflag.FlagSet) {
	flag.Bool(TelemetryEnabledFlag, false, "Is open telemetry tracing enabled")
	flag.String(TelemetryEndpointFlag, "stdout", "open telemetry tracing endpoint")
	flag.Bool(TelemetryUseXrayIDFlag, false, "Using AWS Xray Trace IDs")
	flag.Float32(TelemetrySamplingFractionFlag, 0.5, "Percent of traces to sample")
	flag.Int(TelemetryCollectSecondsFlag, 30, "Metric collection period in seconds")
}

// CheckTelemetry validates the telemetry config
func CheckTelemetry(v *viper.Viper) (*trace.TelemetryConfig, error) {
	config := &trace.TelemetryConfig{}
	config.Enabled = v.GetBool(TelemetryEnabledFlag)
	if !config.Enabled {
		return config, nil
	}
	config.Endpoint = v.GetString(TelemetryEndpointFlag)
	config.UseXrayID = v.GetBool(TelemetryUseXrayIDFlag)
	config.SamplingFraction = v.GetFloat64(TelemetrySamplingFractionFlag)
	if config.SamplingFraction < 0 || config.SamplingFraction > 1 {
		return nil, fmt.Errorf("%s must be between 0 and 1", TelemetrySamplingFractionFlag)
	}
	config.CollectSeconds = v.GetInt(TelemetryCollectSecondsFlag)
	return config, nil
}
