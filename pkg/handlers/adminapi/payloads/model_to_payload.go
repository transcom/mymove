package payloads

import (
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// WebhookSubscriptionPayload converts a webhook subscription payload to a model
func WebhookSubscriptionPayload(sub models.WebhookSubscription) *adminmessages.WebhookSubscription {
	severity := int64(sub.Severity)
	status := adminmessages.WebhookSubscriptionStatus(sub.Status)

	return &adminmessages.WebhookSubscription{
		ID:           *handlers.FmtUUID(sub.ID),
		SubscriberID: handlers.FmtUUID(sub.SubscriberID),
		CallbackURL:  &sub.CallbackURL,
		Severity:     &severity,
		EventKey:     &sub.EventKey,
		Status:       &status,
		CreatedAt:    *handlers.FmtDateTime(sub.CreatedAt),
		UpdatedAt:    *handlers.FmtDateTime(sub.UpdatedAt),
	}
}
