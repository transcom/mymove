package ppmshipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
	if float64(*newPPMShipment.Advance) > float64(*newPPMShipment.EstimatedIncentive)*0.6 {
		result = apperror.NewInvalidInputError(newPPMShipment.ID, nil, verrs, "Advance can not be greater than 60% of the estimated incentive")
	}

	if float64(*newPPMShipment.Advance) < float64(1) {
		result = apperror.NewInvalidInputError(newPPMShipment.ID, nil, verrs, "Advance can not be  value less than 1")
	}

	if *newPPMShipment.AdvanceRequested && *newPPMShipment.Advance == 0 {
		result = apperror.NewInvalidInputError(newPPMShipment.ID, nil, verrs, "An advance amount is required")
	}
	return result
}

// ppmShipmentValidatorFunc is an adapter type for converting a function into an implementation of ppmShipmentValidator
type ppmShipmentValidatorFunc func(appcontext.AppContext, models.PPMShipment, *models.PPMShipment, *models.MTOShipment) error

// Validate fulfills the ppmShipmentValidator interface
func (fn ppmShipmentValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.PPMShipment, older *models.PPMShipment, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}

func mergePPMShipment(newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment) *models.PPMShipment {
	if oldPPMShipment == nil {
		return &newPPMShipment
	}

	ppmShipment := *oldPPMShipment

	ppmShipment.ActualMoveDate = services.SetOptionalDateTimeField(newPPMShipment.ActualMoveDate, ppmShipment.ActualMoveDate)

	ppmShipment.SecondaryPickupPostalCode = services.SetOptionalStringField(newPPMShipment.SecondaryPickupPostalCode, ppmShipment.SecondaryPickupPostalCode)
	ppmShipment.SecondaryDestinationPostalCode = services.SetOptionalStringField(newPPMShipment.SecondaryDestinationPostalCode, ppmShipment.SecondaryDestinationPostalCode)
	ppmShipment.HasProGear = services.SetNoNilOptionalBoolField(newPPMShipment.HasProGear, ppmShipment.HasProGear)
	ppmShipment.EstimatedWeight = services.SetNoNilOptionalPoundField(newPPMShipment.EstimatedWeight, ppmShipment.EstimatedWeight)
	ppmShipment.NetWeight = services.SetNoNilOptionalPoundField(newPPMShipment.NetWeight, ppmShipment.NetWeight)
	ppmShipment.ProGearWeight = services.SetNoNilOptionalPoundField(newPPMShipment.ProGearWeight, ppmShipment.ProGearWeight)
	ppmShipment.SpouseProGearWeight = services.SetNoNilOptionalPoundField(newPPMShipment.SpouseProGearWeight, ppmShipment.SpouseProGearWeight)
	ppmShipment.EstimatedIncentive = services.SetNoNNilOptionalInt32Field(newPPMShipment.EstimatedIncentive, ppmShipment.EstimatedIncentive)
	ppmShipment.Advance = services.SetNoNilOptionalCentField(newPPMShipment.Advance, ppmShipment.Advance)
	ppmShipment.AdvanceRequested = services.SetNoNilOptionalBoolField(newPPMShipment.AdvanceRequested, ppmShipment.AdvanceRequested)

	if newPPMShipment.Advance != nil {
		ppmShipment.Advance = newPPMShipment.Advance
	}

	if !newPPMShipment.ExpectedDepartureDate.IsZero() {
		ppmShipment.ExpectedDepartureDate = newPPMShipment.ExpectedDepartureDate
	}

	if newPPMShipment.PickupPostalCode != "" {
		ppmShipment.PickupPostalCode = newPPMShipment.PickupPostalCode
	}
	if newPPMShipment.DestinationPostalCode != "" {
		ppmShipment.DestinationPostalCode = newPPMShipment.DestinationPostalCode
	}

	// TODO: handle sitExpected

	return &ppmShipment
}
