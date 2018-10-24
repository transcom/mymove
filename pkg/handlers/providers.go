package handlers

import (
	"github.com/transcom/mymove/pkg/di"
)

// AddProviders adds all the DI providers from the handlers package
func AddProviders(c *di.Container) {
	c.MustProvide(NewHandlerContext)
}
