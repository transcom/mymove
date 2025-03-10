package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// DutyLocationName represents an alternative name for a DutyLocation
type DutyLocationName struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	Name           string       `json:"name" db:"name"`
	DutyLocationID uuid.UUID    `json:"duty_location_id" db:"duty_location_id"`
	DutyLocation   DutyLocation `belongs_to:"duty_locations" fk_id:"duty_location_id"`
}

// TableName overrides the table name used by Pop.
func (d DutyLocationName) TableName() string {
	return "duty_location_names"
}

type DutyLocationNames []DutyLocationName
