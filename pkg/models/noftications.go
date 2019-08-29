package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type NotificationTypes string

const (
	// MoveReviewedEmail captures enum value "move_reviewed_email"
	MoveReviewedEmail NotificationTypes = "MOVE_REVIEWED_EMAIL"
)

type Notification struct {
	ID               uuid.UUID         `db:"id"`
	ServiceMemberID  uuid.UUID         `db:"service_member_id"`
	SESMessageID     string            `db:"ses_message_id"`
	NotificationType NotificationTypes `db:"notification_type"`
	CreatedAt        time.Time         `db:"created_at"`
}
