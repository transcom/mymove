package notifications

import (
	"context"
	"go.uber.org/zap"
)

// StubNotificationSender mocks an SES client for local usage
type StubNotificationSender NotificationSendingContext

// NewStubNotificationSender returns a new StubNotificationSender
func NewStubNotificationSender(domain string, logger *zap.Logger) StubNotificationSender {
	return StubNotificationSender{
		domain: domain,
		logger: logger,
	}
}

// SendNotification returns a dummy ID
func (m StubNotificationSender) SendNotification(ctx context.Context, notification notification) error {
	emails, err := notification.emails(ctx)
	if err != nil {
		return err
	}

	for _, email := range emails {
		rawMessage, err := formatRawEmailMessage(email, m.domain)
		if err != nil {
			return err
		}

		m.logger.Debug("Not sending this email",
			zap.String("destinations", email.recipientEmail),
			zap.String("raw message", string(rawMessage[:])))
	}

	return nil
}
