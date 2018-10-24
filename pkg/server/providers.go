package server

import (
	"github.com/transcom/mymove/pkg/di"
)

// AddProviders adds all the DI providers from this package to the Container
func AddProviders(c *di.Container) {
	c.MustProvide(NewLogRequestMiddleware)
	c.MustProvide(NewSessionCookieMiddleware)
	c.MustProvide(NewAppDetectorMiddleware)
}
