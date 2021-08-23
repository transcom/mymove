package webhooksubscription

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type webhookSubscriptionQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type webhookSubscriptionFetcher struct {
	builder webhookSubscriptionQueryBuilder
}

// FetchWebhookSubscription fetches a webhookSubscription given a slice of filters
func (o *webhookSubscriptionFetcher) FetchWebhookSubscription(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.WebhookSubscription, error) {
	var webhookSubscription models.WebhookSubscription
	error := o.builder.FetchOne(appCtx, &webhookSubscription, filters)
	return webhookSubscription, error
}

// NewWebhookSubscriptionFetcher return an implementation of the WebhookSubscriptionFetcher interface
func NewWebhookSubscriptionFetcher(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionFetcher {
	return &webhookSubscriptionFetcher{builder}
}
