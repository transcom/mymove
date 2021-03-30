package services

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
)

// WebhookSubscriptionFetcher is the service object interface for FetchWebhookSubscription
//go:generate mockery -name WebhookSubscriptionFetcher
type WebhookSubscriptionFetcher interface {
	FetchWebhookSubscription(filters []QueryFilter) (models.WebhookSubscription, error)
}

// WebhookSubscriptionCreator is the exported interface for creating an admin user
//go:generate mockery -name WebhookSubscriptionCreator
type WebhookSubscriptionCreator interface {
	CreateWebhookSubscription(subscription *models.WebhookSubscription, subscriberIDFilter []QueryFilter) (*models.WebhookSubscription, *validate.Errors, error)
}

//WebhookSubscriptionUpdater is the service object interface for UpdateWebhookSubscription
//go:generate mockery -name WebhookSubscriptionUpdater
type WebhookSubscriptionUpdater interface {
	UpdateWebhookSubscription(webhooksubscription *models.WebhookSubscription, severity *int64, eTag *string) (*models.WebhookSubscription, error)
}
