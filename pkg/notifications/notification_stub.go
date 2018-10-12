package notifications

import (
	"go.uber.org/zap"
)

// StubNotificationSender mocks an SES client for local usage
type StubNotificationSender SESNotificationSender

// NewStubNotificationSender returns a new StubNotificationSender
func NewStubNotificationSender(logger *zap.Logger) StubNotificationSender {
	return StubNotificationSender{
		logger: logger,
	}
}

// SendNotification returns a dummy ID
func (m StubNotificationSender) SendNotification(notification notification) error {
	emails, err := notification.emails()
	if err != nil {
		return err
	}

	for _, email := range emails {
		rawMessage, err := formatRawEmailMessage(email)
		if err != nil {
			return err
		}

		m.logger.Info("Not sending this email",
			zap.String("destinations", email.recipientEmail),
			zap.String("raw message", string(rawMessage[:])))
	}

	return nil
}
