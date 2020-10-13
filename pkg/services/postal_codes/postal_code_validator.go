package postalcode

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type postalCodeValidator struct {
	DB *pop.Connection
}

// make sure the zipcode is in the route zip_locations
// need the FetchRateAreaForZip5 (db, origin)
// FetchRegionForZip5 (db, destination)
func (v postalCodeValidator) ValidatePostalCode(postalCode string, postalCodeType services.PostalCodeType) (bool, error) {
	_, err := route.Zip5ToLatLong(postalCode)
	if err != nil {
		return false, err
	}

	switch postalCodeType {
	case services.Origin:
		if _, err := models.FetchRateAreaForZip5(v.DB, postalCode); err != nil {
			return false, err
		}
	case services.Destination:
		if _, err := models.FetchRegionForZip5(v.DB, postalCode); err != nil {
			return false, err
		}
	}
	return true, nil
}

// NewPostalCodeValidator is the public constructor for a `NewPostalCodeValidator`
// using Pop
func NewPostalCodeValidator(db *pop.Connection) services.PostalCodeValidator {
	return &postalCodeValidator{db}
}
