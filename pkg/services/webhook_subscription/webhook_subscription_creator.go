package webhooksubscription

import (
	"database/sql"

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
func (o *webhookSubscriptionCreator) CreateWebhookSubscription(appCtx appcontext.AppContext, subscription *models.WebhookSubscription, subscriberIDFilter []services.QueryFilter) (*models.WebhookSubscription, *validate.Errors, error) {
	var contractor models.Contractor
	var verrs *validate.Errors
	var err error

	// check to see if subscriber exists
	fetchErr := o.builder.FetchOne(appCtx, &contractor, subscriberIDFilter)
	if fetchErr != nil {
		switch fetchErr {
		case sql.ErrNoRows:
			return nil, nil, apperror.NewNotFoundError(subscription.SubscriberID, "while looking for SubscriberID")
		default:
			return nil, nil, apperror.NewQueryError("Contractor", fetchErr, "")
		}
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
