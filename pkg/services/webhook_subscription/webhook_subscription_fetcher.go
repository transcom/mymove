package webhooksubscription

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type webhookSubscriptionQueryBuilder interface {
	FetchOne(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCfg appconfig.AppConfig, model interface{}) (*validate.Errors, error)
	UpdateOne(appCfg appconfig.AppConfig, model interface{}, eTag *string) (*validate.Errors, error)
}

type webhookSubscriptionFetcher struct {
	builder webhookSubscriptionQueryBuilder
}

// FetchWebhookSubscription fetches a webhookSubscription given a slice of filters
func (o *webhookSubscriptionFetcher) FetchWebhookSubscription(appCfg appconfig.AppConfig, filters []services.QueryFilter) (models.WebhookSubscription, error) {
	var webhookSubscription models.WebhookSubscription
	error := o.builder.FetchOne(appCfg, &webhookSubscription, filters)
	return webhookSubscription, error
}

// NewWebhookSubscriptionFetcher return an implementation of the WebhookSubscriptionFetcher interface
func NewWebhookSubscriptionFetcher(builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionFetcher {
	return &webhookSubscriptionFetcher{builder}
}
