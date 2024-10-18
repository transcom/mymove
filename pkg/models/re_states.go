package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type State struct {
	ID        uuid.UUID `json:"id" db:"id"`
	State     string    `json:"state" db:"state"`
	StateName string    `json:"state_name" db:"state_name"`
	IsOconus  bool      `json:"is_oconus" db:"is_oconus"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (s State) TableName() string {
	return "re_states"
}
