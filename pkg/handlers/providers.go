package handlers

import (
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// AddProviders adds all the DI providers from the handlers package
func AddProviders(c *dig.Container) {
	err := c.Provide(NewHandlerContext)
	if err != nil {
		c.Invoke(func(l *zap.Logger) {
			l.Fatal("pkg.handlers.AddProviders(NewHandlerContext)", zap.Error(err))
		})
	}
	err = c.Provide(internalapi.NewInternalAPIHandler)
	if err != nil {
		c.Invoke(func(l *zap.Logger) {
			l.Fatal("pkg.handlers.AddProviders(internalapi.NewInternalAPIHandler)", zap.Error(err))
		})
	}
}
