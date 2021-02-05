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
func (o *webhookSubscriptionCreator) CreateWebhookSubscription(subscription *models.WebhookSubscription, contractorIDFilter []services.QueryFilter) (*models.WebhookSubscription, *validate.Errors, error) {
	// Use FetchOne to see if we have an Contractor that matches the provided id (a.k.a subscriber id)
	var contractor models.Contractor
	fetchErr := o.builder.FetchOne(&contractor, contractorIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	var verrs *validate.Errors
	var err error

	txErr := o.db.Transaction(func(connection *pop.Connection) error {

		verrs, err = o.builder.CreateOne(subscription)
		if verrs != nil || err != nil {
			return err
		}

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	return subscription, nil, nil
}

// NewWebhookSubscriptionCreator returns a new admin user creator builder
func NewWebhookSubscriptionCreator(db *pop.Connection, builder webhookSubscriptionQueryBuilder) services.WebhookSubscriptionCreator {
	return &webhookSubscriptionCreator{db, builder}
}
