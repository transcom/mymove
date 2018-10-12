package handlers

import (
	dep "github.com/transcom/mymove/pkg/dependencies"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
)

// AddProviders adds all the DI providers from the handlers package
func AddProviders(c *dep.Container) {
	c.MustProvide(NewHandlerContext)
	c.Provide(internalapi.NewInternalAPIHandler)
}
