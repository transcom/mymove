package ppmshipment

import (
	"math"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkShipmentType checks if the associated mtoShipment has the appropriate shipmentType
func checkShipmentType() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.PPMShipment, _ *models.PPMShipment, mtoShipment *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if mtoShipment.ShipmentType != models.MTOShipmentTypePPM {
			verrs.Add("ShipmentType", "ShipmentType must be of type "+string(models.MTOShipmentTypePPM))
		}
		return verrs
	})
}

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldPPMShipment == nil {
			if newPPMShipment.ShipmentID == uuid.Nil {
				verrs.Add("ShipmentID", "Shipment ID is required")
			}
		} else {
			if newPPMShipment.ShipmentID != uuid.Nil && newPPMShipment.ShipmentID != oldPPMShipment.ShipmentID {
				verrs.Add("ShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkPPMShipmentID checks that the user can't change the PPMShipment ID
func checkPPMShipmentID() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldPPMShipment == nil {
			if newPPMShipment.ID != uuid.Nil {
				verrs.Add("ID", "cannot manually set a new PPM Shipment's UUID")
			}
		} else {
			if newPPMShipment.ID != oldPPMShipment.ID {
				verrs.Add("ID", "ID can not be updated once it is set")
			}
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		// Check that we have something in the expectedDepartureDate field:
		if newPPMShipment.ExpectedDepartureDate.IsZero() {
			verrs.Add("expectedDepartureDate", "cannot be a zero value")
		}

		// Check that we have something in the pickupPostalCode field:
		if newPPMShipment.PickupPostalCode == "" {
			verrs.Add("pickupPostalCode", "cannot be nil or empty")
		}

		// Check that we have something in the destinationPostalCode field:
		if newPPMShipment.DestinationPostalCode == "" {
			verrs.Add("destinationPostalCode", "cannot be nil or empty")
		}

		if newPPMShipment.SitExpected == nil {
			verrs.Add("sitExpected", "cannot be nil")
		}

		return verrs
	})
}

// check Advance checks that the advance fields are updated appropriately
func checkAdvance() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if newPPMShipment.Advance == nil {
			if newPPMShipment.AdvanceRequested == nil || !*newPPMShipment.AdvanceRequested {
				return verrs
			}
		}

		if newPPMShipment.Advance != nil {
			if newPPMShipment.AdvanceRequested == nil || !*newPPMShipment.AdvanceRequested {
				verrs.Add("advance", "Advance must be nil because of the advance requested value")
				return verrs
			}
		}

		if *newPPMShipment.AdvanceRequested && newPPMShipment.Advance == nil {
			verrs.Add("advance", "An advance amount is required")
			return verrs
		}

		if float64(*newPPMShipment.Advance) > math.Floor(float64(*newPPMShipment.EstimatedIncentive)*0.6) {
			verrs.Add("advance", "Advance can not be greater than 60% of the estimated incentive")
		}

		if float64(*newPPMShipment.Advance) < float64(1) {
			verrs.Add("advance", "Advance can not be value less than 1")
		}

		return verrs
	})
}

// checkEstimatedWeight() checks that the weight estimate is availabele to the PPM Estimator
func checkEstimatedWeight() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		// Check that we have something in the estimatedWeight field.
		if newPPMShipment.EstimatedWeight == nil {
			verrs.Add("estimatedWeight", "cannot be empty")
		}

		return verrs
	})
}
