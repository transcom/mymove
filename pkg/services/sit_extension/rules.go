package sitextension

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// checkShipmentID checks that a shipmentID is not nil and returns a verification error if it is
func checkShipmentID() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitExtension models.SITExtension, _ *models.MTOShipment) error {
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

		sitStatus := sitExtension.Status
		sitExtensionReason := sitExtension.RequestReason
		sitRequestedDays := sitExtension.RequestedDays

		// Check that we have something in the SIT RequestedDays field:
		if sitRequestedDays == 0 {
			verrs.Add("RequestedDays", "cannot be blank")
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

func checkSITExtensionPending() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(appCtx appcontext.AppContext, sitExtension models.SITExtension, shipment *models.MTOShipment) error {
		id := sitExtension.ID
		shipmentID := shipment.ID
		//status := sitExtension.Status
		var emptySITExtensionArray []models.SITExtension
		err := appCtx.DB().Where("status = ?", models.SITExtensionStatusPending).Where("mto_shipment_id = ?", shipmentID).All(&emptySITExtensionArray)
		// Prevent a new SIT extension request if a sit extension is pending
		if err != nil {
			return err
		}

		if len(emptySITExtensionArray) > 0 {
			return apperror.NewConflictError(id, "All SIT extensions must be approved or denied to review this new SIT extension")
		}
		return err
	})
}

//checks that the shipment associated with the reweigh is available to Prime
func checkPrimeAvailability(checker services.MoveTaskOrderChecker) sitExtensionValidator {
	return sitExtensionValidatorFunc(func(appCtx appcontext.AppContext, sitExtension models.SITExtension, shipment *models.MTOShipment) error {
		if shipment == nil {
			return apperror.NewNotFoundError(sitExtension.ID, "while looking for Prime-available Shipment")
		}

		isAvailable, err := checker.MTOAvailableToPrime(appCtx, shipment.MoveTaskOrderID)
		if !isAvailable || err != nil {
			return apperror.NewNotFoundError(
				sitExtension.ID, fmt.Sprintf("while looking for Prime-available Shipment with id: %s", shipment.ID))
		}
		return nil
	})
}
