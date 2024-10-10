package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReStates struct {
	ID         uuid.UUID `json:"id" db:"id"`
	State      string    `json:"state" db:"state"`
	State_Name string    `json:"state_name" db:"state_ame"`
	IsConus    bool      `json:"is_conus" db:"is_conus"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (s ReStates) TableName() string {
	return "re_states"
}
