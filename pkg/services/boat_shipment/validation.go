package boatshipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type boatShipmentValidator interface {
	// Validate checks the newBoatShipment for adherence to business rules. The
	// oldBoatShipment parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newBoatShipment models.BoatShipment, oldBoatShipment *models.BoatShipment, shipment *models.MTOShipment) error
}

// validateBoatShipment checks a Boat shipment against a passed-in set of business rule checks
func validateBoatShipment(
	appCtx appcontext.AppContext,
	newBoatShipment models.BoatShipment,
	oldBoatShipment *models.BoatShipment,
	shipment *models.MTOShipment,
	checks ...boatShipmentValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newBoatShipment, oldBoatShipment, shipment); err != nil {
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
		result = apperror.NewInvalidInputError(newBoatShipment.ID, nil, verrs, "Invalid input found while validating the Boat shipment.")
	}

	return result
}

// boatShipmentValidatorFunc is an adapter type for converting a function into an implementation of boatShipmentValidator
type boatShipmentValidatorFunc func(appcontext.AppContext, models.BoatShipment, *models.BoatShipment, *models.MTOShipment) error

// Validate fulfills the boatShipmentValidator interface
func (fn boatShipmentValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.BoatShipment, older *models.BoatShipment, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}

func mergeBoatShipment(newBoatShipment models.BoatShipment, oldBoatShipment *models.BoatShipment) (*models.BoatShipment, error) {
	var err error

	if oldBoatShipment == nil {
		return &newBoatShipment, nil
	}

	boatShipment := *oldBoatShipment

	if newBoatShipment.Type == models.BoatShipmentTypeHaulAway || newBoatShipment.Type == models.BoatShipmentTypeTowAway {
		boatShipment.Type = newBoatShipment.Type
	}
	if newBoatShipment.Year != nil {
		boatShipment.Year = newBoatShipment.Year
	}
	if newBoatShipment.Make != nil {
		boatShipment.Make = newBoatShipment.Make
	}
	if newBoatShipment.Model != nil {
		boatShipment.Model = newBoatShipment.Model
	}
	if newBoatShipment.LengthInInches != nil {
		boatShipment.LengthInInches = newBoatShipment.LengthInInches
	}
	if newBoatShipment.WidthInInches != nil {
		boatShipment.WidthInInches = newBoatShipment.WidthInInches
	}
	if newBoatShipment.HeightInInches != nil {
		boatShipment.HeightInInches = newBoatShipment.HeightInInches
	}
	boatShipment.IsRoadworthy = newBoatShipment.IsRoadworthy
	if newBoatShipment.HasTrailer != nil {
		boatShipment.HasTrailer = newBoatShipment.HasTrailer
		if !*boatShipment.HasTrailer {
			boatShipment.IsRoadworthy = nil
		}
	}

	return &boatShipment, err
}
