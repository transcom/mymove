package notifications

import (
	"github.com/gofrs/uuid"

	"go.uber.org/zap"
)

// StubNotificationSender mocks an SES client for local usage
type StubNotificationSender NotificationSendingContext

// NewStubNotificationSender returns a new StubNotificationSender
func NewStubNotificationSender(domain string, logger Logger) StubNotificationSender {
	return StubNotificationSender{
		domain: domain,
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
		rawMessage, err := formatRawEmailMessage(email, m.domain)
		if err != nil {
			return err
		}
		if email.onSuccess != nil {
			id, _ := uuid.NewV4()
			err := email.onSuccess(id.String())
			if err != nil {
				m.logger.Error("email.onSuccess error", zap.Error(err))
			}
		}

		m.logger.Debug("Not sending this email",
			zap.String("destinations", email.recipientEmail),
			zap.String("raw message", string(rawMessage[:])))
	}

	return nil
}
