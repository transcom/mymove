package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type State struct {
	ID        uuid.UUID `json:"id" db:"id" rw:"r"`
	State     string    `json:"state" db:"state" rw:"r"`
	StateName string    `json:"state_name" db:"state_name" rw:"r"`
	IsOconus  bool      `json:"is_oconus" db:"is_oconus" rw:"r"`
	CreatedAt time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (s State) TableName() string {
	return "re_states"
}
