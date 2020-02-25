package models

import (
	"time"

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
	ServiceMember    ServiceMember     `belongs_to:"service_member"`
	SESMessageID     string            `db:"ses_message_id"`
	NotificationType NotificationTypes `db:"notification_type"`
	CreatedAt        time.Time         `db:"created_at"`
}

// Notifications is a slice of notification structs
type Notifications []Notification
