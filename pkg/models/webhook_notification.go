package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// This file is auto-generated with soda g model table_name. It is not generated automatically
// though, so changes will persist.

// WebhookNotificationStatus represents the possible statuses for a mto shipment
type WebhookNotificationStatus string

const (
	// WebhookNotificationPending is the pending status type for a WebhookNotification
	WebhookNotificationPending WebhookNotificationStatus = "PENDING"
	// WebhookNotificationSent is the sent status type for a WebhookNotification
	WebhookNotificationSent WebhookNotificationStatus = "SENT"
	// WebhookNotificationSkipped is the skipped status type for a WebhookNotification
	WebhookNotificationSkipped WebhookNotificationStatus = "SKIPPED"
	// WebhookNotificationFailing is the failing status type for a WebhookNotification
	// - indicates the send has failed at least once but we are still retrying it
	WebhookNotificationFailing WebhookNotificationStatus = "FAILING"
	// WebhookNotificationFailed is the failed status type for a WebhookNotification
	WebhookNotificationFailed WebhookNotificationStatus = "FAILED"
)

// WebhookNotification is used by pop to map your webhook_notifications database table to your go code.
type WebhookNotification struct {
	ID               uuid.UUID                 `db:"id"`
	EventKey         string                    `db:"event_key"`
	TraceID          *uuid.UUID                `db:"trace_id"`
	MoveTaskOrderID  *uuid.UUID                `db:"move_id"`
	MoveTaskOrder    Move                      `belongs_to:"moves"`
	ObjectID         *uuid.UUID                `db:"object_id"`
	Payload          string                    `db:"payload"`
	Status           WebhookNotificationStatus `db:"status"`
	CreatedAt        time.Time                 `db:"created_at"`
	UpdatedAt        time.Time                 `db:"updated_at"`
	FirstAttemptedAt *time.Time                `db:"first_attempted_at"`
}

// String is not required by pop and may be deleted
func (w WebhookNotification) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

// WebhookNotifications is not required by pop and may be deleted
type WebhookNotifications []WebhookNotification

// String is not required by pop and may be deleted
func (w WebhookNotifications) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (w *WebhookNotification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&OptionalRegexMatch{Field: &w.EventKey, Name: "EventKey", Expr: `\w+\.\w+`, Message: "Eventkey should be in Subject.Action format."},
		&OptionalUUIDIsPresent{Field: w.TraceID, Name: "TraceID"},
		&validators.StringIsPresent{Field: w.Payload, Name: "Payload"},
		&validators.StringInclusion{Field: string(w.Status), Name: "Status", List: []string{
			string(WebhookNotificationPending),
			string(WebhookNotificationSent),
			string(WebhookNotificationSkipped),
			string(WebhookNotificationFailing),
			string(WebhookNotificationFailed),
		}},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (w *WebhookNotification) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (w *WebhookNotification) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
