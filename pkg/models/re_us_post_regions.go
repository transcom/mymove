package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReUsPostRegions struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UsprZipID string    `json:"uspr_zip_id" db:"uspr_zip_id"`
	StateId   uuid.UUID `json:"state_id" db:"state_id"`
	Zip3      string    `json:"zip3" db:"zip3"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (d ReUsPostRegions) TableName() string {
	return "re_us_post_regions"
}
