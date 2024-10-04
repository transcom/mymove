package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type UsPostRegion struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UsprZipID string    `db:"uspr_zip_id" json:"uspr_zip_id"`
	Zip3      string    `db:"zip3" json:"zip3"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (d UsPostRegion) TableName() string {
	return "us_post_region"
}
