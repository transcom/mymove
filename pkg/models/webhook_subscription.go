package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

const (
	// WebhookSubscriptionStatusActive is the active status for Webhook Subscription
	WebhookSubscriptionStatusActive WebhookSubscriptionStatus = "ACTIVE"
	// WebhookSubscriptionStatusDisabled is the disabled status for Webhook Subscription
	WebhookSubscriptionStatusDisabled WebhookSubscriptionStatus = "DISABLED"
	// WebhookSubscriptionStatusFailing is the failing status for Webhook Subscription
	// - it indicates that we have experienced notifications failing to be sent, but
	// have not disabled this subscription yet.
	WebhookSubscriptionStatusFailing WebhookSubscriptionStatus = "FAILING"
)

// WebhookSubscriptionStatus is a type representing the webhook subscription status type - string
type WebhookSubscriptionStatus string

// A WebhookSubscription represents a webhook subscription
type WebhookSubscription struct {
	ID           uuid.UUID                 `db:"id"`
	Subscriber   Contractor                `belongs_to:"contractors" fk_id:"subscriber_id"`
	SubscriberID uuid.UUID                 `db:"subscriber_id"`
	Status       WebhookSubscriptionStatus `db:"status"`
	Severity     int                       `db:"severity"` // Zero indicates no severity value, 1 is highest
	EventKey     string                    `db:"event_key"`
	CallbackURL  string                    `db:"callback_url"`
	CreatedAt    time.Time                 `db:"created_at"`
	UpdatedAt    time.Time                 `db:"updated_at"`
}

// WebhookSubscriptions is an array of webhook subscriptions
type WebhookSubscriptions []WebhookSubscription

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (wS *WebhookSubscription) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: wS.SubscriberID, Name: "SubscriberID"},
		&validators.StringIsPresent{Field: wS.EventKey, Name: "EventKey"},
		&validators.StringIsPresent{Field: wS.CallbackURL, Name: "CallbackURL"},
		&validators.StringInclusion{Field: string(wS.Status), Name: "Status", List: []string{
			string(WebhookSubscriptionStatusActive),
			string(WebhookSubscriptionStatusDisabled),
			string(WebhookSubscriptionStatusFailing),
		}},
	), nil
}

// TableName overrides the table name used by Pop.
func (wS *WebhookSubscription) TableName() string {
	return "webhook_subscriptions"
}
