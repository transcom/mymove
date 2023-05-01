package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildWebhookNotification(db *pop.Connection, customs []Customization, traits []Trait) models.WebhookNotification {
	customs = setupCustomizations(customs, traits)

	// Find webhook notification customization and extract custom webhook notification
	var cNotification models.WebhookNotification
	if result := findValidCustomization(customs, WebhookNotification); result != nil {
		cNotification = result.Model.(models.WebhookNotification)
		if result.LinkOnly {
			return cNotification
		}
	}

	move := BuildMove(db, customs, traits)

	// Create a default notification
	traceID := uuid.Must(uuid.NewV4())
	notification := models.WebhookNotification{
		EventKey:        "Payment.Create",
		MoveTaskOrderID: &move.ID,
		Payload:         "{\"message\":\"This is a default Payment.Create notification.\"}",
		Status:          models.WebhookNotificationPending,
		TraceID:         &traceID,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&notification, cNotification)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &notification)
	}

	return notification

}
