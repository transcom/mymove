package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeWebhookNotification creates a single webhook notification
func MakeWebhookNotification(db *pop.Connection, assertions Assertions) models.WebhookNotification {
	// Get the passed in Move object
	move := assertions.WebhookNotification.MoveTaskOrder
	// But if no id was set, create a Move object and pass in assertions
	if assertions.WebhookNotification.MoveTaskOrderID == nil ||
		isZeroUUID(*assertions.WebhookNotification.MoveTaskOrderID) {
		move = MakeMove(db, assertions)
	}

	// Create a default notification
	notification := models.WebhookNotification{
		EventKey:        "Payment.Create",
		MoveTaskOrderID: &move.ID,
		Payload:         swag.String("{\"message\":\"This is a default Payment.Create notification.\"}"),
		Status:          models.WebhookNotificationPending,
	}

	// Overwrite the defaults with values provided
	mergeModels(&notification, assertions.WebhookNotification)

	mustCreate(db, &notification, assertions.Stub)

	return notification
}

// MakeDefaultWebhookNotification makes an Notification with default values
func MakeDefaultWebhookNotification(db *pop.Connection) models.WebhookNotification {
	return MakeWebhookNotification(db, Assertions{})
}
