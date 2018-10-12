package services

import (
	dep "github.com/transcom/mymove/pkg/dependencies"
	"github.com/transcom/mymove/pkg/services/user"
)

// AddProviders adds all the DI providers from the services package
func AddProviders(c *dep.Container) {
	c.MustProvide(user.NewFetchServiceMemberService)
}
