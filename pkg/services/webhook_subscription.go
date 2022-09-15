package services

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// WebhookSubscriptionFetcher is the service object interface for FetchWebhookSubscription
//go:generate mockery --name WebhookSubscriptionFetcher --disable-version-string
type WebhookSubscriptionFetcher interface {
	FetchWebhookSubscription(appCtx appcontext.AppContext, filters []QueryFilter) (models.WebhookSubscription, error)
}

// WebhookSubscriptionCreator is the exported interface for creating an admin user
//go:generate mockery --name WebhookSubscriptionCreator --disable-version-string
type WebhookSubscriptionCreator interface {
	CreateWebhookSubscription(appCtx appcontext.AppContext, subscription *models.WebhookSubscription) (*models.WebhookSubscription, *validate.Errors, error)
}

//WebhookSubscriptionUpdater is the service object interface for UpdateWebhookSubscription
//go:generate mockery --name WebhookSubscriptionUpdater --disable-version-string
type WebhookSubscriptionUpdater interface {
	UpdateWebhookSubscription(appCtx appcontext.AppContext, webhooksubscription *models.WebhookSubscription, severity *int64, eTag *string) (*models.WebhookSubscription, error)
}
