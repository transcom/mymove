package authentication

import (
	"github.com/transcom/mymove/pkg/di"
)

// AddProviders adds all the dependency providers to the DI container
func AddProviders(c *di.Container) {
	c.MustProvide(NewUserAuthMiddleware)
	c.MustProvide(NewLoginGovProvider)
	c.MustProvide(NewCallbackHandler)
	c.MustProvide(NewLogoutHandler)
	c.MustProvide(NewRedirectHandler)
	c.MustProvide(NewAssignUserHandler)
	c.MustProvide(NewCreateUserHandler)
	c.MustProvide(NewUserListHandler)
	c.MustProvide(NewAuthContext)
}
