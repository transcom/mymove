package webhooksubscription

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type webhookSubscriptionCreator struct {
	builder webhookSubscriptionQueryBuilder
}

// CreateWebhookSubscription creates admin user
func (o *webhookSubscriptionCreator) CreateWebhookSubscription(appCfg appconfig.AppConfig, subscription *models.WebhookSubscription, subscriberIDFilter []services.QueryFilter) (*models.WebhookSubscription, *validate.Errors, error) {
	var contractor models.Contractor
	var verrs *validate.Errors
	var err error

	// check to see if subscriber exists
	fetchErr := o.builder.FetchOne(appCfg, &contractor, subscriberIDFilter)
	if fetchErr != nil {
		return nil, nil, services.NewNotFoundError(subscription.SubscriberID, "while looking for SubscriberID")
	}

	verrs, err = o.builder.CreateOne(appCfg, subscription)
	if verrs != nil && verrs.HasAny() {
		return nil, verrs, nil
	} else if err != nil {
		return nil, verrs, services.NewQueryError("unknown", err, "")
	}

	return subscription, nil, nil
}

// NewWebhookSubscriptionCreator returns a new admin user creator builder
func NewWebhookSubscriptionCreator(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionCreator {
	return &webhookSubscriptionCreator{builder}
}
