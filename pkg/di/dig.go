package di

import (
	"go.uber.org/dig"
	"go.uber.org/zap"
	"log"
)

// Config contains the environment name and debug logging flag for configuring zap.Logging
type Config struct {
	Environment  string
	DebugLogging bool
}

// Container wraps dig.Container so we can add MustProvide and MustInvoke wrapper methods
type Container struct {
	dig.Container
}

// NewContainer constructs a dependency injection Container.
// configProvider should, minimally, provide *Config
func NewContainer(configProvider interface{}, opts ...dig.ProvideOption) *Container {

	dc := dig.New()
	if err := dc.Provide(configProvider, opts...); err != nil {
		log.Fatal("Provide(configProvider)", zap.Error(err))
	}
	// Set up logger so we can invoke MustProvide & MustInvoke
	if err := dc.Provide(NewLogger); err != nil {
		log.Fatal("Provide(loggingProvider", zap.Error(err))
	}
	return &Container{*dc}
}

// MustProvide wraps dig.Container.Provide in a fatal error check. Used for required initialization
func (c *Container) MustProvide(constructor interface{}, opts ...dig.ProvideOption) {
	if err := c.Provide(constructor, opts...); err != nil {
		c.Invoke(func(l *zap.Logger) {
			l.Fatal("MustProvide", zap.Any("Constructor", constructor), zap.Error(err))
		})
	}
}

// MustInvoke wraps dig.Container.Invoke in a fatal error check. Used for required initialization
func (c *Container) MustInvoke(function interface{}, opts ...dig.InvokeOption) {
	if err := c.Invoke(function, opts...); err != nil {
		c.Invoke(func(l *zap.Logger) {
			l.Fatal("MustInvoke", zap.Any("Function", function), zap.Error(err))
		})
	}
}
