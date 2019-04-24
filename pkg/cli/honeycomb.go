package cli

import (
	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// ServiceNameFlag is the Service Name Flag
	ServiceNameFlag string = "service-name"
	// HoneycombEnabledFlag is the Honeycomb Enabled Flag
	HoneycombEnabledFlag string = "honeycomb-enabled"
	// HoneycombAPIHostFlag is the Honeycomb API Host Flag
	HoneycombAPIHostFlag string = "honeycomb-api-host"
	// HoneycombAPIKeyFlag is the Honeycomb API Key Flag
	HoneycombAPIKeyFlag string = "honeycomb-api-key"
	// HoneycombDatasetFlag is the Honeycomb Dataset Flag
	HoneycombDatasetFlag string = "honeycomb-dataset"
	// HoneycombDebugFlag is the Honeycomb Debug Flag
	HoneycombDebugFlag string = "honeycomb-debug"
)

// InitHoneycombFlags initializes Honeycomb command line flags
func InitHoneycombFlags(flag *pflag.FlagSet) {
	flag.String(ServiceNameFlag, "app", "The service name identifies the application for instrumentation.")
	flag.Bool(HoneycombEnabledFlag, false, "Honeycomb enabled")
	flag.String(HoneycombAPIHostFlag, "https://api.honeycomb.io/", "API Host for Honeycomb")
	flag.String(HoneycombAPIKeyFlag, "", "API Key for Honeycomb")
	flag.String(HoneycombDatasetFlag, "", "Dataset for Honeycomb")
	flag.Bool(HoneycombDebugFlag, false, "Debug honeycomb using stdout.")
}

// InitHoneycomb initilizes the honeycomb service
func InitHoneycomb(v *viper.Viper, logger Logger) bool {

	honeycombEnabled := v.GetBool(HoneycombEnabledFlag)
	honeycombAPIHost := v.GetString(HoneycombAPIHostFlag)
	honeycombAPIKey := v.GetString(HoneycombAPIKeyFlag)
	honeycombDataset := v.GetString(HoneycombDatasetFlag)
	honeycombDebug := v.GetBool(HoneycombDebugFlag)
	honeycombServiceName := v.GetString(ServiceNameFlag)

	if honeycombEnabled && len(honeycombAPIKey) > 0 && len(honeycombDataset) > 0 && len(honeycombServiceName) > 0 {
		logger.Debug("Honeycomb Integration enabled",
			zap.String(HoneycombAPIHostFlag, honeycombAPIHost),
			zap.String(HoneycombDatasetFlag, honeycombDataset))
		beeline.Init(beeline.Config{
			APIHost:     honeycombAPIHost,
			WriteKey:    honeycombAPIKey,
			Dataset:     honeycombDataset,
			Debug:       honeycombDebug,
			ServiceName: honeycombServiceName,
		})
		return true
	}

	logger.Debug("Honeycomb Integration disabled")
	return false
}

// CheckHoneycomb validates Honeycomb command line flags
func CheckHoneycomb(v *viper.Viper) error {
	if serviceName := v.GetString(ServiceNameFlag); len(serviceName) == 0 {
		return errors.Errorf("Must provide service name for honeycomb")
	}

	return nil
}
