package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

const (
	// WebhookSubscriptionStatusActive is the active status for Webhook Subscription
	WebhookSubscriptionStatusActive WebhookSubscriptionStatus = "ACTIVE"
	// WebhookSubscriptionStatusDisabled is the disabled status for Webhook Subscription
	WebhookSubscriptionStatusDisabled WebhookSubscriptionStatus = "DISABLED"
)

// WebhookSubscriptionStatus is a type representing the webhook subscription status type - string
type WebhookSubscriptionStatus string

// A WebhookSubscription represents a webhook subscription
type WebhookSubscription struct {
	ID           uuid.UUID                 `db:"id"`
	Subscriber   Contractor                `belongs_to:"contractors:"`
	SubscriberID uuid.UUID                 `db:"subscriber_id"`
	Status       WebhookSubscriptionStatus `db:"status"`
	EventKey     string                    `db:"event_key"`
	CallbackURL  string                    `db:"callback_url"`
	CreatedAt    time.Time                 `db:"created_at"`
	UpdatedAt    time.Time                 `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (wS *WebhookSubscription) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: wS.SubscriberID, Name: "SubscriberID"},
		&validators.StringIsPresent{Field: wS.EventKey, Name: "EventKey"},
		&validators.StringIsPresent{Field: wS.CallbackURL, Name: "CallbackURL"},
		&validators.StringInclusion{Field: string(wS.Status), Name: "Status", List: []string{
			string(WebhookSubscriptionStatusActive),
			string(WebhookSubscriptionStatusDisabled),
		}},
	), nil
}

// TableName overrides the table name used by Pop.
func (wS *WebhookSubscription) TableName() string {
	return "webhook_subscriptions"
}