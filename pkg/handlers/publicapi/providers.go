package publicapi

import "github.com/transcom/mymove/pkg/di"

// AddProviders adds the DI providers for this package
func AddProviders(c *di.Container) {
	c.MustProvide(NewPublicAPIHandler)
}
