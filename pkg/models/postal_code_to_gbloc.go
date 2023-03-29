package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

// Fetches the GBLOC for a specific Postal Code
func FetchGBLOCForPostalCode(db *pop.Connection, postalCode string) (PostalCodeToGBLOC, error) {
	var postalCodeToGBLOC PostalCodeToGBLOC
	err := db.Where("postal_code = $1", postalCode).First(&postalCodeToGBLOC)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return PostalCodeToGBLOC{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return PostalCodeToGBLOC{}, err
	}

	return postalCodeToGBLOC, nil
}
