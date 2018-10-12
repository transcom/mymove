package route

import "github.com/transcom/mymove/pkg/dependencies"

// AddProviders registers all the dependency providers for the route package
func AddProviders(c *dependencies.Container) {
	c.MustProvide(NewHEREPlanner)
}
