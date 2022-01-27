package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

type ActivityLog struct {
	ID           uuid.UUID `db:"id"`
	ActivityUser string    `db:"activity_user"`
	Source       string    `db:"source"`
	Entity       string    `db:"entity"`
	LogEventType string    `db:"log_event_type"`
	LogData      string    `db:"log_data"`
	MoveID       string    `db:"move_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// CreateActivityLog looks for an office user with a specific email
func CreateActivityLog(tx *pop.Connection) error {
	activityLog := ActivityLog{
		ActivityUser: "kenneth_chow",
		Source:       "shipment_create",
		Entity:       "shipment",
		LogEventType: "Create",
		MoveID:       "ABCD-1234",
		LogData:      `{"test": "test"}`,
	}

	err := tx.Create(&activityLog)
	if err != nil {
		return err
	}
	return nil
}
