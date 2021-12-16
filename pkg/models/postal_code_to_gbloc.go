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

// PostalCodeToGBLOCs is not required by pop and may be deleted
type PostalCodeToGBLOCs []PostalCodeToGBLOC
