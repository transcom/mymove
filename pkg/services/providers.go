package services

import (
	"github.com/transcom/mymove/pkg/services/user"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// AddProviders adds all the DI providers from the services package
func AddProviders(c *dig.Container) {
	err := c.Provide(user.NewFetchServiceMemberService)
	if err != nil {
		c.Invoke(func(l *zap.Logger) {
			l.Fatal("services.AddProviders(NewFetchServiceMemberService", zap.Error(err))
		})
	}
}
