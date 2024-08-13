package mobilehomeshipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type mobileHomeShipmentValidator interface {
	// Validate checks the newMobileHomeShipment for adherence to business rules. The
	// oldMobileHomeShipment parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newMobileHomeShipment models.MobileHome, oldMobileHomeShipment *models.MobileHome, shipment *models.MTOShipment) error
}

// validateMobileHomeShipment checks a Mobile Home shipment against a passed-in set of business rule checks
func validateMobileHomeShipment(
	appCtx appcontext.AppContext,
	newMobileHomeShipment models.MobileHome,
	oldMobileHomeShipment *models.MobileHome,
	shipment *models.MTOShipment,
	checks ...mobileHomeShipmentValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newMobileHomeShipment, oldMobileHomeShipment, shipment); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = apperror.NewInvalidInputError(newMobileHomeShipment.ID, nil, verrs, "Invalid input found while validating the Mobile Home shipment.")
	}

	return result
}

// mobileHomeShipmentValidatorFunc is an adapter type for converting a function into an implementation of mobileHomeShipmentValidator
type mobileHomeShipmentValidatorFunc func(appcontext.AppContext, models.MobileHome, *models.MobileHome, *models.MTOShipment) error

// Validate fulfills the mobileHomeShipmentValidator interface
func (fn mobileHomeShipmentValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.MobileHome, older *models.MobileHome, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}

func mergeMobileHomeShipment(newMobileHomeShipment models.MobileHome, oldMobileHomeShipment *models.MobileHome) (*models.MobileHome, error) {
	var err error

	if oldMobileHomeShipment == nil {
		return &newMobileHomeShipment, nil
	}

	mobileHomeShipment := *oldMobileHomeShipment

	if newMobileHomeShipment.Year != nil {
		mobileHomeShipment.Year = newMobileHomeShipment.Year
	}
	if newMobileHomeShipment.Make != nil {
		mobileHomeShipment.Make = newMobileHomeShipment.Make
	}
	if newMobileHomeShipment.Model != nil {
		mobileHomeShipment.Model = newMobileHomeShipment.Model
	}
	if newMobileHomeShipment.LengthInInches != nil {
		mobileHomeShipment.LengthInInches = newMobileHomeShipment.LengthInInches
	}
	if newMobileHomeShipment.WidthInInches != nil {
		mobileHomeShipment.WidthInInches = newMobileHomeShipment.WidthInInches
	}
	if newMobileHomeShipment.HeightInInches != nil {
		mobileHomeShipment.HeightInInches = newMobileHomeShipment.HeightInInches
	}

	return &mobileHomeShipment, err
}