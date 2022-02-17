package ppmshipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type ppmShipmentValidator interface {
	// Validate checks the newPPMShipment for adherence to business rules. The
	// oldPPMShipment parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, shipment *models.MTOShipment) error
}

// validatePPMShipment checks a PPM shipment against a passed-in set of business rule checks
func validatePPMShipment(
	appCtx appcontext.AppContext,
	newPPMShipment models.PPMShipment,
	oldPPMShipment *models.PPMShipment,
	shipment *models.MTOShipment,
	checks ...ppmShipmentValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newPPMShipment, oldPPMShipment, shipment); err != nil {
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
		result = apperror.NewInvalidInputError(newPPMShipment.ID, nil, verrs, "Invalid input found while validating the PPM shipment.")
	}

	return result
}

// ppmShipmentValidatorFunc is an adapter type for converting a function into an implementation of ppmShipmentValidator
type ppmShipmentValidatorFunc func(appcontext.AppContext, models.PPMShipment, *models.PPMShipment, *models.MTOShipment) error

// Validate fulfills the ppmShipmentValidator interface
func (fn ppmShipmentValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.PPMShipment, older *models.PPMShipment, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}
