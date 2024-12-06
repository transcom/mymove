package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type GblocAors struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	JppsoRegionID       uuid.UUID `json:"jppso_regions_id" db:"jppso_regions_id"`
	OconusRateAreaID    uuid.UUID `json:"oconus_rate_area_id" db:"oconus_rate_area_id"`
	DepartmentIndicator *string   `json:"department_indicator" db:"department_indicator"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (c GblocAors) TableName() string {
	return "gbloc_aors"
}
