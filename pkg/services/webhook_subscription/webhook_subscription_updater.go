package webhooksubscription

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type webhookSubscriptionUpdater struct {
	builder webhookSubscriptionQueryBuilder
}

// UpdateWebhookSubscription updates a webhookSubscription
func (o *webhookSubscriptionUpdater) UpdateWebhookSubscription(webhooksubscription *models.WebhookSubscription) (*models.WebhookSubscription, error) {
	webhookSubscriptionID := uuid.FromStringOrNil(webhooksubscription.ID.String())
	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscriptionID)}

	// Find the existing web subscription to update
	var foundWebhookSubscription models.WebhookSubscription
	err := o.builder.FetchOne(&foundWebhookSubscription, queryFilters)
	if err != nil {
		return nil, err
	}

	// Update webhook subscription new status for Active
	if webhooksubscription.Status != "" {
		foundWebhookSubscription.Status = webhooksubscription.Status
	}

	if webhooksubscription.SubscriberID != uuid.Nil {
		foundWebhookSubscription.SubscriberID = webhooksubscription.SubscriberID
	}

	if webhooksubscription.EventKey != "" {
		foundWebhookSubscription.EventKey = webhooksubscription.EventKey
	}

	if webhooksubscription.Severity != -1 {
		foundWebhookSubscription.Severity = webhooksubscription.Severity
	}

	if webhooksubscription.CallbackURL != "" {
		foundWebhookSubscription.CallbackURL = webhooksubscription.CallbackURL
	}

	verrs, err := o.builder.UpdateOne(&foundWebhookSubscription, nil)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(webhookSubscriptionID, err, verrs, "")
	}

	// First check to see if there is an error on the type and return a precondition fail error, if not return the error
	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(webhookSubscriptionID, err)
		default:
			return nil, err
		}
	}
	// return *webhooksubscription, nil
	return &foundWebhookSubscription, nil
}

// NewWebhookSubscriptionUpdater returns an instance of the WebhookSubscriptionUpdater interface
func NewWebhookSubscriptionUpdater(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionUpdater {
	return &webhookSubscriptionUpdater{builder}
}
