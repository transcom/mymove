package auth

import (
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/dependencies"
)

// AddProviders adds all the auth providers to the DI container
func AddProviders(c *dependencies.Container) {
	c.MustProvide(NewSessionCookieMiddleware)
	c.MustProvide(authentication.NewLoginGovProvider)
}
