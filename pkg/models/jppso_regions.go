package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type JppsoRegions struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (c JppsoRegions) TableName() string {
	return "jppso_regions"
}
