package sitextension

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"

	"fmt"

	"github.com/gofrs/uuid"
)

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitExtension models.SITExtension,  _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if sitExtension.MTOShipmentID == uuid.Nil {
			verrs.Add("MTOShipmentID", "Shipment ID is required")
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitExtension models.SITExtension, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		var sitStatus models.SITExtensionStatus
		var sitExtensionReason models.SITExtensionRequestReason
		var approvedDays = sitExtension.ApprovedDays
		var decisionDate = sitExtension.DecisionDate


		sitStatus = sitExtension.Status
		sitExtensionReason = sitExtension.RequestReason

		// Check that we have something in the ApprovedDays field:
		if approvedDays == nil {
			verrs.Add("approvedDays", "cannot be empty")
		}

		// Check that we have something in the DecisionDate field:
		if decisionDate == nil {
			verrs.Add("decisionDate", "cannot be empty")
		}

		// Check that we have something in the Status field:
		if sitStatus == "" {
			verrs.Add("Status", "cannot be blank")
		}

		// Check that we have something in the RequestReason field:
		if sitExtensionReason == "" {
			verrs.Add("RequestReason", "cannot be blank")
		}

		return verrs
	})
}

//checks that the shipment associated with the reweigh is available to Prime
func checkPrimeAvailability(checker services.MoveTaskOrderChecker) sitExtensionValidator {
	return sitExtensionValidatorFunc(func(appCtx appcontext.AppContext, sitExtension models.SITExtension,  shipment *models.MTOShipment) error {
		if shipment == nil {
			return services.NewNotFoundError(sitExtension.ID, "while looking for Prime-available Shipment")
		}

		isAvailable, err := checker.MTOAvailableToPrime(appCtx, shipment.MoveTaskOrderID)
		if !isAvailable || err != nil {
			return services.NewNotFoundError(
				sitExtension.ID, fmt.Sprintf("while looking for Prime-available Shipment with id: %s", shipment.ID))
		}
		return nil
	})
}
