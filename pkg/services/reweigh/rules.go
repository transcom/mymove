package reweigh

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() reweighValidator {
	return reweighValidatorFunc(func(_ appcontext.AppContext, newReweigh models.Reweigh, oldReweigh *models.Reweigh, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldReweigh == nil {
			if newReweigh.ShipmentID == uuid.Nil {
				verrs.Add("ShipmentID", "Shipment ID is required")
			}
		} else {
			if newReweigh.ShipmentID != uuid.Nil && newReweigh.ShipmentID != oldReweigh.ShipmentID {
				verrs.Add("ShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkReweighID checks that the user can't change the reweigh ID
func checkReweighID() reweighValidator {
	return reweighValidatorFunc(func(_ appcontext.AppContext, newReweigh models.Reweigh, oldReweigh *models.Reweigh, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldReweigh == nil {
			if newReweigh.ID != uuid.Nil {
				verrs.Add("ID", "cannot manually set a new reweigh's UUID")
			}
		} else {
			if newReweigh.ID != oldReweigh.ID {
				return apperror.NewImplementationError(
					fmt.Sprintf("the newReweigh ID (%s) must match oldReweigh ID (%s).", newReweigh.ID, oldReweigh.ID),
				)
			}
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() reweighValidator {
	return reweighValidatorFunc(func(_ appcontext.AppContext, reweigh models.Reweigh, oldReweigh *models.Reweigh, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		var requestedAt time.Time
		var requestedBy models.ReweighRequester

		// Set any pre-existing values as the baseline:
		if oldReweigh != nil {
			requestedAt = oldReweigh.RequestedAt
			requestedBy = oldReweigh.RequestedBy
		}

		// Override pre-existing values with anything sent in for the update/create:
		if !reweigh.RequestedAt.IsZero() {
			requestedAt = reweigh.RequestedAt
		}
		if reweigh.RequestedBy != "" {
			requestedBy = reweigh.RequestedBy
		}

		// Check that we have something in the RequestedAt field:
		if requestedAt.IsZero() {
			verrs.Add("requestedAt", "cannot be empty")
		}

		// Check that we have something in the RequestedBy field:
		if requestedBy == "" {
			verrs.Add("requestedBy", "cannot be blank")
		}

		return verrs
	})
}

//checks that the shipment associated with the reweigh is available to Prime
func checkPrimeAvailability(checker services.MoveTaskOrderChecker) reweighValidator {
	return reweighValidatorFunc(func(appCtx appcontext.AppContext, newReweigh models.Reweigh, oldReweigh *models.Reweigh, shipment *models.MTOShipment) error {
		if shipment == nil {
			return apperror.NewNotFoundError(newReweigh.ID, "while looking for Prime-available Shipment")
		}

		if shipment.UsesExternalVendor {
			return apperror.NewNotFoundError(
				newReweigh.ID, fmt.Sprintf("while looking for Prime-available Shipment with id: %s", shipment.ID))
		}

		isAvailable, err := checker.MTOAvailableToPrime(appCtx, shipment.MoveTaskOrderID)
		if !isAvailable || err != nil {
			return apperror.NewNotFoundError(
				newReweigh.ID, fmt.Sprintf("while looking for Prime-available Shipment with id: %s", shipment.ID))
		}
		return nil
	})
}
