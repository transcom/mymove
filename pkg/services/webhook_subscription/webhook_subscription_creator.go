package webhooksubscription

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type webhookSubscriptionCreator struct {
	builder webhookSubscriptionQueryBuilder
}

// CreateWebhookSubscription creates admin user
func (o *webhookSubscriptionCreator) CreateWebhookSubscription(appCtx appcontext.AppContext, subscription *models.WebhookSubscription) (*models.WebhookSubscription, *validate.Errors, error) {
	return o.createWebhookSubscription(appCtx, subscription, checkSubscriberExists(o.builder))
}

// createWebhookSubscription creates admin user
func (o *webhookSubscriptionCreator) createWebhookSubscription(appCtx appcontext.AppContext, subscription *models.WebhookSubscription, checks ...webhookSubscriptionValidator) (*models.WebhookSubscription, *validate.Errors, error) {
	var verrs *validate.Errors
	var err error

	e := validateWebhookSubscription(appCtx, *subscription, checks...)
	if e != nil {
		return nil, nil, e
	}
	verrs, err = o.builder.CreateOne(appCtx, subscription)
	if verrs != nil && verrs.HasAny() {
		return nil, verrs, nil
	} else if err != nil {
		return nil, verrs, apperror.NewQueryError("unknown", err, "")
	}

	return subscription, nil, nil
}

// NewWebhookSubscriptionCreator returns a new admin user creator builder
func NewWebhookSubscriptionCreator(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionCreator {
	return &webhookSubscriptionCreator{builder}
}
