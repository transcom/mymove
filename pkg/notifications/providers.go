package notifications

import (
	"github.com/transcom/mymove/pkg/dependencies"
	"go.uber.org/zap"
)

// NewNotificationSender provides either a local NotificationSender or one using SES depending on
// whether or not the SESNotificationConfig is supplied
func NewNotificationSender(cfg *SESNotificationConfig, l *zap.Logger) (NotificationSender, error) {
	if cfg == nil {
		return NewStubNotificationSender(l), nil
	}
	return NewSESNotificationSender(cfg, l)
}

// AddProviders adds the dependency providers supplied by the notifications module
func AddProviders(c *dependencies.Container) {
	c.MustProvide(NewNotificationSender)
}
