package notifications

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

// StubNotificationSender mocks an SES client for local usage
type StubNotificationSender NotificationSendingContext

// NewStubNotificationSender returns a new StubNotificationSender
func NewStubNotificationSender(domain string) StubNotificationSender {
	return StubNotificationSender{
		domain: domain,
	}
}

// SendNotification returns a dummy ID
func (m StubNotificationSender) SendNotification(appCtx appcontext.AppContext, notification Notification) error {
	emails, err := notification.emails(appCtx)
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
				appCtx.Logger().Error("email.onSuccess error", zap.Error(err))
			}
		}

		appCtx.Logger().Debug("Not sending this email",
			zap.String("destinations", email.recipientEmail),
			zap.String("raw message", string(rawMessage[:])))
	}

	return nil
}
