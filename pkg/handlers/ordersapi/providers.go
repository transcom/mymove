package ordersapi

import "github.com/transcom/mymove/pkg/di"

// AddProviders adds all the DI providers for this package
func AddProviders(c *di.Container) {
	c.MustProvide(NewOrdersAPIHandler)
}
