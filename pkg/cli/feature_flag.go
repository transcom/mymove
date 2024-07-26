package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FeatureFlagServerURLFlag is the URL to the feature flag server
	FeatureFlagServerURLFlag string = "feature-flag-server-url"
	// FeatureFlagAPITokenFlag is the api token
	FeatureFlagAPITokenFlag string = "feature-flag-api-token"
	// FeatureFlagDODIDValidation
	FeatureFlagDODIDUnique string = "feature-flag-dodid-unique"
)

// InitFeatureFlags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.String(FeatureFlagServerURLFlag, "", "The endpoint of the feature flag server")
	flag.String(FeatureFlagAPITokenFlag, "", "The api token for the feature flag server")
	flag.Bool(FeatureFlagDODIDUnique, false, "The feature flag that determines if DODIDs need to be unique")
}

// CheckFeatureFlag validates the URL
func CheckFeatureFlag(_ *viper.Viper) error {
	// Right now, we have no mandatory checks as we can allow a server
	// URL without a token
	return nil
}

type FeatureFlagConfig struct {
	Namespace string
	URL       string
	Token     string
}

func GetFliptFetcherConfig(v *viper.Viper) FeatureFlagConfig {
	return FeatureFlagConfig{
		Namespace: v.GetString(EnvironmentFlag),
		URL:       v.GetString(FeatureFlagServerURLFlag),
		Token:     v.GetString(FeatureFlagAPITokenFlag),
	}
}
