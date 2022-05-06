package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeWebhookSubscription creates a single WebhookSubscription.
func MakeWebhookSubscription(db *pop.Connection, assertions Assertions) models.WebhookSubscription {
	subscriber := assertions.Contractor
	if isZeroUUID(subscriber.ID) {
		subscriber = MakeContractor(db, assertions)
	}

	webhookSubscription := models.WebhookSubscription{
		Subscriber:   subscriber,
		SubscriberID: subscriber.ID,
		Status:       models.WebhookSubscriptionStatusActive,
		EventKey:     "PaymentRequest.Update",
		CallbackURL:  "/my/callback/url",
	}

	mergeModels(&webhookSubscription, assertions.WebhookSubscription)

	mustCreate(db, &webhookSubscription, assertions.Stub)

	return webhookSubscription
}

// MakeDefaultWebhookSubscription makes a WebhookSubscription with default values
func MakeDefaultWebhookSubscription(db *pop.Connection) models.WebhookSubscription {
	return MakeWebhookSubscription(db, Assertions{})
}
