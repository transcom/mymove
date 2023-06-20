package featureflag

import (
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services"
)

func NewFeatureFlagFetcher(config cli.FeatureFlagConfig) (services.FeatureFlagFetcher, error) {
	if config.URL != "" {
		return NewFliptFetcher(config)
	}

	return NewEnvFetcher(config)
}
