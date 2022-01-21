package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// NewTable represents a new table
type LogEventType struct {
	ID        uuid.UUID `db:"id"`
	EventType string    `db:"event_type"`
	EventName *string   `db:"event_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
