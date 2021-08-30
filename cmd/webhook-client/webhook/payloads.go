package webhook

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetWebhookNotificationPayload converts Webhook Notification model to Payload using the definition in support API
func GetWebhookNotificationPayload(model *models.WebhookNotification) *supportmessages.WebhookNotification {
	payload := supportmessages.WebhookNotification{
		EventKey:        model.EventKey,
		Object:          swag.String(model.Payload),
		ObjectID:        handlers.FmtUUIDPtr(model.ObjectID),
		MoveTaskOrderID: handlers.FmtUUIDPtr(model.MoveTaskOrderID),
	}

	if model.TraceID != nil {
		payload.TraceID = *handlers.FmtUUIDPtr(model.TraceID)
	}
	return &payload
}
