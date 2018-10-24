package route

import "github.com/transcom/mymove/pkg/di"

// AddProviders registers all the dependency providers for the route package
func AddProviders(c *di.Container) {
	c.MustProvide(NewHEREPlanner)
}
