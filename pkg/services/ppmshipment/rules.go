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

// helper function to check if the secondary address is empty, but the tertiary is not
func isPPMShipmentAddressCreateSequenceValid(ppmShipmentToCheck models.PPMShipment) bool {
	bothPickupAddressesEmpty := (models.IsAddressEmpty(ppmShipmentToCheck.SecondaryPickupAddress) && models.IsAddressEmpty(ppmShipmentToCheck.TertiaryPickupAddress))
	bothDestinationAddressesEmpty := (models.IsAddressEmpty(ppmShipmentToCheck.SecondaryDestinationAddress) && models.IsAddressEmpty(ppmShipmentToCheck.TertiaryDestinationAddress))
	bothPickupAddressesNotEmpty := !bothPickupAddressesEmpty
	bothDestinationAddressesNotEmpty := !bothDestinationAddressesEmpty
	hasNoSecondaryHasTertiaryPickup := (models.IsAddressEmpty(ppmShipmentToCheck.SecondaryPickupAddress) && !models.IsAddressEmpty(ppmShipmentToCheck.TertiaryPickupAddress))
	hasNoSecondaryHasTertiaryDestination := (models.IsAddressEmpty(ppmShipmentToCheck.SecondaryDestinationAddress) && !models.IsAddressEmpty(ppmShipmentToCheck.TertiaryDestinationAddress))

	// need an explicit case to capture when both are empty or not empty
	if (bothPickupAddressesEmpty && bothDestinationAddressesEmpty) || (bothPickupAddressesNotEmpty && bothDestinationAddressesNotEmpty) {
		return true
	}
	if hasNoSecondaryHasTertiaryPickup || hasNoSecondaryHasTertiaryDestination {
		return false
	}
	return true
}

/* Checks if a valid address sequence is being maintained. This will return false if the tertiary address is being updated while the secondary address remains empty
*
 */
func isPPMAddressUpdateSequenceValid(shipmentToUpdateWith *models.PPMShipment, currentShipment *models.PPMShipment) bool {
	// if the incoming model has both fields, then we know the model will be updated with both secondary and tertiary addresses. therefore return true
	if !models.IsAddressEmpty(shipmentToUpdateWith.SecondaryPickupAddress) && !models.IsAddressEmpty(shipmentToUpdateWith.TertiaryPickupAddress) {
		return true
	}
	if !models.IsAddressEmpty(shipmentToUpdateWith.SecondaryDestinationAddress) && !models.IsAddressEmpty(shipmentToUpdateWith.TertiaryDestinationAddress) {
		return true
	}
	if currentShipment.SecondaryPickupAddress == nil && shipmentToUpdateWith.TertiaryPickupAddress != nil {
		return !hasTertiaryWithNoSecondaryAddress(currentShipment.SecondaryPickupAddress, shipmentToUpdateWith.TertiaryPickupAddress)
	}
	if currentShipment.SecondaryDestinationAddress == nil && shipmentToUpdateWith.TertiaryDestinationAddress != nil {
		return !hasTertiaryWithNoSecondaryAddress(currentShipment.SecondaryDestinationAddress, shipmentToUpdateWith.TertiaryDestinationAddress)
	}
	// no addresses are being updated, so correct address sequence is maintained, return true
	return true
}

// helper function to check if the secondary address is empty, but the tertiary is not
func hasTertiaryWithNoSecondaryAddress(secondaryAddress *models.Address, tertiaryAddress *models.Address) bool {
	return (models.IsAddressEmpty(secondaryAddress) && !models.IsAddressEmpty(tertiaryAddress))
}

func checkPPMShipmentSequenceValidForCreate() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(appCtx appcontext.AppContext, newer models.PPMShipment, _ *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		squenceIsValid := isPPMShipmentAddressCreateSequenceValid(newer)
		if !squenceIsValid {
			verrs.Add("error validating ppm shipment", "PPM Shipment cannot have a tertiary address without a secondary address present")
			return verrs
		}
		return nil
	})
}

func checkPPMShipmentSequenceValidForUpdate() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(appCtx appcontext.AppContext, newer models.PPMShipment, older *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		sequenceIsValid := isPPMAddressUpdateSequenceValid(&newer, older)
		if !sequenceIsValid {
			verrs.Add("error validating ppm shipment", "PPM Shipment cannot have a tertiary address without a secondary address present")
			return verrs
		}
		return nil
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, _ *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		// Check that we have something in the pickupPostalCode field:
		if newPPMShipment.PickupAddressID == nil || newPPMShipment.PickupAddressID == models.UUIDPointer(uuid.Nil) {
			verrs.Add("pickupAddressID", "cannot be nil or empty")
		}

		// Check that we have something in the destinationPostalCode field:
		if newPPMShipment.DestinationAddressID == nil || newPPMShipment.PickupAddressID == models.UUIDPointer(uuid.Nil) {
			verrs.Add("destinationAddressID", "cannot be nil or empty")
		}

		// Check that we have something in the expectedDepartureDate field:
		if newPPMShipment.ExpectedDepartureDate.IsZero() {
			verrs.Add("expectedDepartureDate", "cannot be a zero value")
		}

		if newPPMShipment.SITExpected == nil {
			verrs.Add("sitExpected", "cannot be nil")
		}

		return verrs
	})
}

// checkAdvanceAmountRequested()  checks that the advance fields are updated appropriately
func checkAdvanceAmountRequested() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, _ *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if newPPMShipment.HasRequestedAdvance == nil || !*newPPMShipment.HasRequestedAdvance {
			if newPPMShipment.AdvanceAmountRequested != nil {
				verrs.Add("advanceAmountRequested", "Advance amount requested must be nil because of the value of the field indicating if an advance was requested")
			}
		} else {

			if newPPMShipment.AdvanceAmountRequested == nil {
				verrs.Add("advanceAmountRequested", "An advance amount is required")
			} else if float64(*newPPMShipment.AdvanceAmountRequested) < float64(0) {
				verrs.Add("advanceAmountRequested", "Advance amount requested cannot be negative.")
			} else if float64(*newPPMShipment.AdvanceAmountRequested) > math.Floor(float64(*newPPMShipment.EstimatedIncentive)*0.6) {
				verrs.Add("advanceAmountRequested", "Advance amount requested can not be greater than 60% of the estimated incentive")
			}
		}

		return verrs
	})
}

// checkEstimatedWeight() checks that the weight estimate is available to the PPM Estimator
func checkEstimatedWeight() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, _ *models.PPMShipment, _ *models.MTOShipment) error {
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
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, _ *models.PPMShipment, _ *models.MTOShipment) error {
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
