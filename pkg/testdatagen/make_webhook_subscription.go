package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeWebhookSubscription creates a single WebhookSubscription.
func MakeWebhookSubscription(db *pop.Connection, assertions Assertions) models.WebhookSubscription {
	subscriber := assertions.Contractor
	if isZeroUUID(subscriber.ID) {
		subscriber = MakeContractor(db, assertions)
	}

	callbackURL := assertions.WebhookSubscription.CallbackURL
	if callbackURL == "" {
		callbackURL = DefaultWebhookSubscriptionCallbackURL
	}

	eventKey := assertions.WebhookSubscription.EventKey
	if eventKey == "" {
		eventKey = "PaymentRequest.Update"
	}

	webhookSubscription := models.WebhookSubscription{
		Subscriber:   subscriber,
		SubscriberID: subscriber.ID,
		EventKey:     eventKey,
		CallbackURL:  callbackURL,
	}

	mergeModels(&webhookSubscription, assertions.WebhookSubscription)

	mustCreate(db, &webhookSubscription)

	return webhookSubscription
}