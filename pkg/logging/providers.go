package logging

import "github.com/transcom/mymove/pkg/dependencies"

func AddProviders(c *dependencies.Container) {
	c.MustProvide(NewLogRequestMiddleware)
}
