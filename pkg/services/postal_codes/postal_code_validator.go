package postalcode

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type postalCodeValidator struct {
}

// make sure the zipcode is in the route zip_locations
// need the FetchRateAreaForZip5 (db, origin)
// FetchRegionForZip5 (db, destination)
func (v postalCodeValidator) ValidatePostalCode(appCtx appcontext.AppContext, postalCode string, postalCodeType services.PostalCodeType) (bool, error) {
	_, err := route.Zip5ToLatLong(postalCode)
	if err != nil {
		return false, err
	}

	switch postalCodeType {
	case services.Origin:
		if _, err := models.FetchRateAreaForZip5(appCtx.DB(), postalCode); err != nil {
			return false, err
		}
	case services.Destination:
		if _, err := models.FetchRegionForZip5(appCtx.DB(), postalCode); err != nil {
			return false, err
		}
	}
	return true, nil
}

// NewPostalCodeValidator is the public constructor for a `NewPostalCodeValidator`
// using Pop
func NewPostalCodeValidator() services.PostalCodeValidator {
	return &postalCodeValidator{}
}
