package webhooksubscription

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type webhookSubscriptionQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type webhookSubscriptionFetcher struct {
	builder webhookSubscriptionQueryBuilder
}

// FetchUser fetches an  user given a slice of filters
func (o *webhookSubscriptionFetcher) FetchWebhookSubscription(filters []services.QueryFilter) (*models.WebhookSubscription, error) {
	var webhookSubscription models.WebhookSubscription
	error := o.builder.FetchOne(&webhookSubscription, filters)
	return &webhookSubscription, error
}

// NewWebhookSubscriptionFetcher return an implementation of the WebhookSubscriptionFetcher interface
func NewWebhookSubscriptionFetcher(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionFetcher {
	return &webhookSubscriptionFetcher{builder}
}
