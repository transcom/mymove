package services

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// WebhookSubscriptionFetcher is the service object interface for FetchWebhookSubscription
//go:generate mockery --name WebhookSubscriptionFetcher --disable-version-string
type WebhookSubscriptionFetcher interface {
	FetchWebhookSubscription(appCfg appconfig.AppConfig, filters []QueryFilter) (models.WebhookSubscription, error)
}

// WebhookSubscriptionCreator is the exported interface for creating an admin user
//go:generate mockery --name WebhookSubscriptionCreator --disable-version-string
type WebhookSubscriptionCreator interface {
	CreateWebhookSubscription(appCfg appconfig.AppConfig, subscription *models.WebhookSubscription, subscriberIDFilter []QueryFilter) (*models.WebhookSubscription, *validate.Errors, error)
}

//WebhookSubscriptionUpdater is the service object interface for UpdateWebhookSubscription
//go:generate mockery --name WebhookSubscriptionUpdater --disable-version-string
type WebhookSubscriptionUpdater interface {
	UpdateWebhookSubscription(appCfg appconfig.AppConfig, webhooksubscription *models.WebhookSubscription, severity *int64, eTag *string) (*models.WebhookSubscription, error)
}
