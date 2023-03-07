package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// PostalCodeToGBLOC is a mapping from Postal Codes to GBLOCs
type PostalCodeToGBLOC struct {
	ID         uuid.UUID `db:"id"`
	PostalCode string    `db:"postal_code"`
	GBLOC      string    `db:"gbloc"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (p PostalCodeToGBLOC) TableName() string {
	return "postal_code_to_gblocs"
}

type PostalCodeToGBLOCs []PostalCodeToGBLOC
