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

		if newPPMShipment.SITExpected == nil {
			verrs.Add("sitExpected", "cannot be nil")
		}

		return verrs
	})
}

// checkAdvanceAmountRequested()  checks that the advance fields are updated appropriately
func checkAdvanceAmountRequested() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if newPPMShipment.HasRequestedAdvance == nil || !*newPPMShipment.HasRequestedAdvance {
			if newPPMShipment.AdvanceAmountRequested != nil {
				verrs.Add("advanceAmountRequested", "Advance amount requested must be nil because of the value of the field indicating if an advance was requested")
			}
		} else {

			if newPPMShipment.AdvanceAmountRequested == nil {
				verrs.Add("advanceAmountRequested", "An advance amount is required")
			} else if float64(*newPPMShipment.AdvanceAmountRequested) < float64(100) {
				verrs.Add("advanceAmountRequested", "Advance amount requested cannot be a value less than $1")
			} else if float64(*newPPMShipment.AdvanceAmountRequested) > math.Floor(float64(*newPPMShipment.EstimatedIncentive)*0.6) {
				verrs.Add("advanceAmountRequested", "Advance amount requested can not be greater than 60% of the estimated incentive")
			}
		}

		return verrs
	})
}

// checkEstimatedWeight() checks that the weight estimate is available to the PPM Estimator
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

// checkSITRequiredFields() checks that if SIT is expected that the other dependent fields are all unset or all valid
func checkSITRequiredFields() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if newPPMShipment.SITExpected == nil || !*newPPMShipment.SITExpected {
			return verrs
		}

		// If the customer is selecting SITExpected then we will be missing the details only the services counselor can
		// provide later
		if newPPMShipment.SITLocation == nil &&
			newPPMShipment.SITEstimatedWeight == nil &&
			newPPMShipment.SITEstimatedEntryDate == nil &&
			newPPMShipment.SITEstimatedDepartureDate == nil {
			return verrs
		}

		if newPPMShipment.SITLocation == nil {
			verrs.Add("sitLocation", "cannot be empty")
		}

		if newPPMShipment.SITEstimatedWeight == nil {
			verrs.Add("sitEstimatedWeight", "cannot be empty")
		}

		if newPPMShipment.SITEstimatedEntryDate == nil || newPPMShipment.SITEstimatedEntryDate.IsZero() {
			verrs.Add("sitEstimatedEntryDate", "cannot be empty")
		}

		if newPPMShipment.SITEstimatedDepartureDate == nil || newPPMShipment.SITEstimatedDepartureDate.IsZero() {
			verrs.Add("sitEstimatedDepartureDate", "cannot be empty")
		} else if newPPMShipment.SITEstimatedEntryDate != nil && newPPMShipment.SITEstimatedDepartureDate.Before(*newPPMShipment.SITEstimatedEntryDate) {
			verrs.Add("sitEstimatedDepartureDate", "cannot come before SIT Estimated Entry Date")
		}

		return verrs
	})
}
