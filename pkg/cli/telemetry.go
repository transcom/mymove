package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/telemetry"
)

const (
	// TelemetryEnabledFlag is the Trace Enable Flag
	TelemetryEnabledFlag string = "telemetry-enabled"
	// TelemetryEndpointFlag configures the endpoint used for open
	// telemetry tracing
	TelemetryEndpointFlag string = "telemetry-endpoint"
	// TelemetryUseXrayIDFlag enables using AWS Xray Trace IDs for open telemetry
	TelemetryUseXrayIDFlag string = "telemetry-use-xray-id"
	// TelemetrySamplingFractionFlag configures the percent of traces to sample
	TelemetrySamplingFractionFlag string = "telemetry-sampling-fraction"
	// TelemetryCollectSecondsFlag configures the metric collection
	// period in seconds
	TelemetryCollectSecondsFlag string = "telemetry-collect-seconds"
	// TelemetryReadEventsEnabledFlag enables read events
	TelemetryReadEventsEnabledFlag string = "telemetry-read-events-enabled"
	// TelemetryWriteEventsEnabledFlag enables write events
	TelemetryWriteEventsEnabledFlag string = "telemetry-write-events-enabled"
)

// InitTelemetryFlags initializes the open telemetry flags
func InitTelemetryFlags(flag *pflag.FlagSet) {
	flag.Bool(TelemetryEnabledFlag, false, "Is open telemetry tracing enabled")
	flag.String(TelemetryEndpointFlag, "stdout", "open telemetry tracing endpoint")
	flag.Bool(TelemetryUseXrayIDFlag, false, "Using AWS Xray Trace IDs")
	flag.Float32(TelemetrySamplingFractionFlag, 0.5, "Percent of traces to sample")
	flag.Int(TelemetryCollectSecondsFlag, 30, "Metric collection period in seconds")
	flag.Bool(TelemetryReadEventsEnabledFlag, true, "Enable read event traces")
	flag.Bool(TelemetryWriteEventsEnabledFlag, true, "Enable write event traces")
}

// CheckTelemetry validates the telemetry config
func CheckTelemetry(v *viper.Viper) (*telemetry.Config, error) {
	config := &telemetry.Config{}
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
	config.ReadEvents = v.GetBool(TelemetryReadEventsEnabledFlag)
	config.WriteEvents = v.GetBool(TelemetryWriteEventsEnabledFlag)
	return config, nil
}
