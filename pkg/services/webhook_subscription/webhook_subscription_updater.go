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

// UpdateWebhookSubscription updates a webhook subscription
// It uses the id in the passed in model to find the subscription.
// For the severity field, it uses severity from the parameter not the model. If nil, severity will not be updated.
// For all other fields, it uses the values found in the model.
func (o *webhookSubscriptionUpdater) UpdateWebhookSubscription(requestedUpdate *models.WebhookSubscription, severity *int64, eTag *string) (*models.WebhookSubscription, error) {
	webhookSubscriptionID := uuid.FromStringOrNil(requestedUpdate.ID.String())
	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscriptionID)}

	// Find the existing web subscription to update
	var foundSub models.WebhookSubscription
	err := o.builder.FetchOne(&foundSub, queryFilters)
	if err != nil {
		return nil, err
	}

	// Update webhook subscription new status for Active
	if requestedUpdate.Status != "" {
		foundSub.Status = requestedUpdate.Status
	}

	if requestedUpdate.SubscriberID != uuid.Nil {
		foundSub.SubscriberID = requestedUpdate.SubscriberID
	}

	if requestedUpdate.EventKey != "" {
		foundSub.EventKey = requestedUpdate.EventKey
	}

	if severity != nil {
		foundSub.Severity = int(*severity)
	}

	if requestedUpdate.CallbackURL != "" {
		foundSub.CallbackURL = requestedUpdate.CallbackURL
	}

	verrs, err := o.builder.UpdateOne(&foundSub, eTag)

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
	return &foundSub, nil
}

// NewWebhookSubscriptionUpdater returns an instance of the WebhookSubscriptionUpdater interface
func NewWebhookSubscriptionUpdater(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionUpdater {
	return &webhookSubscriptionUpdater{builder}
}
