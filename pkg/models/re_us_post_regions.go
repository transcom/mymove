package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type UsPostRegion struct {
	ID        uuid.UUID `json:"id" db:"id" rw:"r"`
	UsprZipID string    `json:"uspr_zip_id" db:"uspr_zip_id" rw:"r"`
	StateId   uuid.UUID `json:"state_id" db:"state_id" rw:"r"`
	State     State     `belongs_to:"re_states" fk_id:"state_id" rw:"r"`
	Zip3      string    `json:"zip3" db:"zip3" rw:"r"`
	IsPoBox   bool      `json:"is_po_box" db:"is_po_box" rw:"r"`
	CreatedAt time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (r UsPostRegion) TableName() string {
	return "re_us_post_regions"
}
