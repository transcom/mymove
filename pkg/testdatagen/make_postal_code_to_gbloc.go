package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakePostalCodeToGBLOC creates a single PostalCodeToGBLOC and associated data.
func MakePostalCodeToGBLOC(db *pop.Connection, postalCode string, gbloc string) models.PostalCodeToGBLOC {
	postalCodeToGBLOC := models.PostalCodeToGBLOC{
		ID:         uuid.Must(uuid.NewV4()),
		PostalCode: postalCode,
		GBLOC:      gbloc,
	}

	// There's no reason to stub this model because it is only used in queries
	mustCreate(db, &postalCodeToGBLOC, false)

	return postalCodeToGBLOC
}
