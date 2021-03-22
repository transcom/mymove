package payloads

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// WebhookSubscriptionPayload converts a webhook subscription model to a payload
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
		CreatedAt:    strfmt.DateTime(sub.CreatedAt),
		UpdatedAt:    strfmt.DateTime(sub.UpdatedAt),
		ETag:         etag.GenerateEtag(sub.UpdatedAt),
	}
}
