package webhooksubscription

import (
	"database/sql"

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
	// logger := h.LoggerFromRequest(params.HTTPRequest)
	var foundWebhookSubscription models.WebhookSubscription

	// Find the existing web subscription to update
	err := o.builder.FetchOne(&foundWebhookSubscription, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, err
		default:
			return nil, err
		}
	}

	// Update webhook subscription new status for Active
	if &webhooksubscription.Status != nil {
		foundWebhookSubscription.Status = webhooksubscription.Status
	}

	if &webhooksubscription.SubscriberID != nil {
		foundWebhookSubscription.SubscriberID = webhooksubscription.SubscriberID
	}

	if &webhooksubscription.EventKey != nil {
		foundWebhookSubscription.EventKey = webhooksubscription.EventKey
	}

	if &webhooksubscription.Severity != nil {
		foundWebhookSubscription.Severity = webhooksubscription.Severity
	}

	if &webhooksubscription.CallbackURL != nil {
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
