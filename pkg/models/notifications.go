package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// NotificationTypes represents types of notifications
type NotificationTypes string

const (
	// MoveReviewedEmail captures enum value "MOVE_REVIEWED_EMAIL"
	MoveReviewedEmail NotificationTypes = "MOVE_REVIEWED_EMAIL"
	// MovePaymentReminderEmail captures enum value "MOVE_PAYMENT_REMINDER_EMAIL"
	MovePaymentReminderEmail NotificationTypes = "MOVE_PAYMENT_REMINDER_EMAIL"
)

// Notification represents an email sent to a service member
type Notification struct {
	ID               uuid.UUID         `db:"id"`
	ServiceMemberID  uuid.UUID         `db:"service_member_id"`
	ServiceMember    ServiceMember     `belongs_to:"service_member" fk_id:"service_member_id"`
	SESMessageID     string            `db:"ses_message_id"`
	NotificationType NotificationTypes `db:"notification_type"`
	CreatedAt        time.Time         `db:"created_at"`
}

// Notifications is a slice of notification structs
type Notifications []Notification

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (n *Notification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: n.ServiceMemberID, Name: "ServiceMemberID"},
	), nil
}
