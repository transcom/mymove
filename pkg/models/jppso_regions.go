package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
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

func FetchJppsoRegionByCode(db *pop.Connection, code string) (*JppsoRegions, error) {
	var jppsoRegions JppsoRegions
	err := db.Q().
		Where("jppso_regions.code = ?", code).
		First(&jppsoRegions)
	if err != nil {
		return nil, err
	}
	return &jppsoRegions, nil
}
