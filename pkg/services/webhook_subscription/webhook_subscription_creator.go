package webhooksubscription

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type webhookSubscriptionCreator struct {
	db      *pop.Connection
	builder webhookSubscriptionQueryBuilder
}

// CreateWebhookSubscription creates admin user
func (o *webhookSubscriptionCreator) CreateWebhookSubscription(subscription *models.WebhookSubscription, subscriberIDFilter []services.QueryFilter) (*models.WebhookSubscription, *validate.Errors, error) {
	var contractor models.Contractor
	var verrs *validate.Errors
	var err error

	// check to see if subscriber exists
	fetchErr := o.builder.FetchOne(&contractor, subscriberIDFilter)
	if fetchErr != nil {
		return nil, nil, services.NewNotFoundError(subscription.SubscriberID, "while looking for SubscriberID")
	}

	verrs, err = o.builder.CreateOne(subscription)
	if verrs != nil && verrs.HasAny() {
		return nil, verrs, nil
	} else if err != nil {
		return nil, verrs, services.NewQueryError("unknown", err, "")
	}

	return subscription, nil, nil
}

// NewWebhookSubscriptionCreator returns a new admin user creator builder
func NewWebhookSubscriptionCreator(db *pop.Connection, builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionCreator {
	return &webhookSubscriptionCreator{db, builder}
}
