package sitextension

import (
	"fmt"
	"slices"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
)

// checkShipmentID checks that a shipmentID is not nil and returns a verification error if it is
func checkShipmentID() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitExtension models.SITDurationUpdate, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if sitExtension.MTOShipmentID == uuid.Nil {
			verrs.Add("MTOShipmentID", "Shipment ID is required")
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitExtension models.SITDurationUpdate, _ *models.MTOShipment) error {
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
	return sitExtensionValidatorFunc(func(appCtx appcontext.AppContext, sitExtension models.SITDurationUpdate, shipment *models.MTOShipment) error {
		id := sitExtension.ID
		shipmentID := shipment.ID
		//status := sitExtension.Status
		var emptySITExtensionArray []models.SITDurationUpdate
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

// checks that the shipment associated with the sit extension is available to Prime
func checkPrimeAvailability(checker services.MoveTaskOrderChecker) sitExtensionValidator {
	return sitExtensionValidatorFunc(func(appCtx appcontext.AppContext, sitExtension models.SITDurationUpdate, shipment *models.MTOShipment) error {
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

// checks that the total SIT duration for a shipment is not reduced below 1 day by a newly-approved SITDurationUpdate
// since SITDurationUpdate.approvedDays can be negative
func checkMinimumSITDuration() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitDurationUpdate models.SITDurationUpdate, shipment *models.MTOShipment) error {
		if sitDurationUpdate.ApprovedDays == nil {
			return apperror.NewInvalidInputError(sitDurationUpdate.ID, nil, nil, "missing sitDurationUpdate.ApprovedDays, can't calculate newSITDuration")
		}
		if shipment.SITDaysAllowance == nil {
			return apperror.NewInvalidInputError(shipment.ID, nil, nil, "missing shipment.SITDaysAllowance, can't calculate newSITDuration")
		}
		newSITDuration := int(*sitDurationUpdate.ApprovedDays) + int(*shipment.SITDaysAllowance)
		if newSITDuration < 1 {
			return apperror.NewInvalidInputError(sitDurationUpdate.ID, nil, nil, "can't reduce a SIT duration to less than one day")
		}
		return nil
	},
	)
}

func checkDepartureDate() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(appCtx appcontext.AppContext, _ models.SITDurationUpdate, shipment *models.MTOShipment) error {
		// The prime cannot create a SIT Extension if the current SIT departure date
		// is before or equal to the authorized end date.
		var endDate *time.Time
		var si *models.MTOServiceItem
		shipmentSITStatus := sitstatus.NewShipmentSITStatus()

		sitGroupings, err := shipmentSITStatus.RetrieveShipmentSIT(appCtx, *shipment)
		if err != nil {
			return err
		}
		sorted := sitstatus.SortShipmentSITs(sitGroupings, time.Now())

		// Check if any current SITs
		if sorted.CurrentSITs != nil {
			for _, serviceItem := range shipment.MTOServiceItems {
				if serviceItem.SITDepartureDate != nil {
					// Check if valid service SIT service item to get correct authorized end date.
					if slices.Contains(models.ValidOriginAdditionalDaySITReServiceCodes, serviceItem.ReService.Code) &&
						shipment.OriginSITAuthEndDate != nil && shipment.DestinationSITAuthEndDate == nil {
						si = &serviceItem
						endDate = shipment.OriginSITAuthEndDate
					} else if (slices.Contains(models.ValidDestinationAdditionalDaySITReServiceCodes, serviceItem.ReService.Code)) && shipment.DestinationSITAuthEndDate != nil {
						si = &serviceItem
						endDate = shipment.DestinationSITAuthEndDate
					}
				}
			}
		}

		format := "2006-01-02"
		if endDate != nil && si != nil {
			if si.SITDepartureDate.Before(*endDate) || si.SITDepartureDate.Equal(*endDate) {
				sitErr := fmt.Sprintf("\nSIT extension cannot be created: SIT departure date (%s) cannot be prior or equal to the SIT end date (%s)", si.SITDepartureDate.Format(format), endDate.Format(format))
				return apperror.NewConflictError(shipment.ID, sitErr)
			}
		}
		return nil
	},
	)
}
