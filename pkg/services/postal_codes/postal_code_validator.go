package postalcode

import (
	"database/sql"
	"strconv"

	"github.com/benbjohnson/clock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type postalCodeValidator struct {
	clock clock.Clock
}

// NewPostalCodeValidator is the public constructor for a `NewPostalCodeValidator`
// using Pop
func NewPostalCodeValidator(clock clock.Clock) services.PostalCodeValidator {
	return &postalCodeValidator{
		clock: clock,
	}
}

// ValidatePostalCode will ensure that the zip code is found in several data sources so we avoid issues
// with pricing and such:
//   - zip5ToLatLongMap map (may be temporary until the HERE planner is removed)
//   - postal_code_to_gblocs table
//   - re_zip3s table (and re_zip5_rate_areas if a zip3 with multiple rate areas)
func (v postalCodeValidator) ValidatePostalCode(appCtx appcontext.AppContext, postalCode string) (bool, error) {
	// Get the zip5 and zip3 after verifying proper format.
	if len(postalCode) < 5 {
		return false, apperror.NewUnsupportedPostalCodeError(postalCode, "less than 5 characters")
	}
	zip5 := postalCode[:5]
	if _, err := strconv.Atoi(zip5); err != nil {
		return false, apperror.NewUnsupportedPostalCodeError(postalCode, "should only contain digits")
	}
	zip3 := zip5[:3]

	// Note: We don't appear to use the zip3ToLatLongMap currently, so not looking for a zip3 there.

	// Should be able to delete this validation if the HERE planner is eventually removed.
	_, err := route.Zip5ToLatLong(zip5)
	if err != nil {
		return false, err // Already returns an UnsupportedPostalCodeError
	}

	// Check that the postal code exists in the postal_code_to_gblocs table.
	exists, err := appCtx.DB().Where("postal_code = ?", zip5).Exists(&models.PostalCodeToGBLOC{})
	if err != nil {
		return false, err
	} else if !exists {
		return false, apperror.NewUnsupportedPostalCodeError(zip5, "not found in postal_code_to_gblocs")
	}

	// Get contract ID based on the current date.
	var reContract models.ReContract
	err = appCtx.DB().Q().
		InnerJoin("re_contract_years cy", "re_contracts.id = cy.contract_id").
		Where("? between cy.start_date and cy.end_date", v.clock.Now()).First(&reContract)
	if err == sql.ErrNoRows {
		return false, apperror.NewUnsupportedPostalCodeError(zip5, "could not find contract year")
	} else if err != nil {
		return false, err
	}

	// Make sure that the postal code exists in the re_zip3s table (and the re_zip5_rate_areas table if it's
	// a zip3 with multiple rate areas).
	var reZip3 models.ReZip3
	err = appCtx.DB().Q().
		Where("contract_id = ? and zip3 = ?", reContract.ID, zip3).First(&reZip3)
	if err == sql.ErrNoRows {
		return false, apperror.NewUnsupportedPostalCodeError(zip3, "not found in re_zip3s")
	} else if err != nil {
		return false, err
	}

	if reZip3.HasMultipleRateAreas {
		exists, err = appCtx.DB().Q().
			Where("contract_id = ? and zip5 = ?", reContract.ID, zip5).Exists(&models.ReZip5RateArea{})
		if err != nil {
			return false, err
		} else if !exists {
			return false, apperror.NewUnsupportedPostalCodeError(zip5, "not found in re_zip5_rate_areas")
		}
	}

	return true, nil
}
