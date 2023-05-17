package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func checkStatus() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer.Status != "" && newer.Status != models.MTOShipmentStatusDraft && newer.Status != models.MTOShipmentStatusSubmitted {
			verrs.Add("status", "can only update status to DRAFT or SUBMITTED. use UpdateMTOShipmentStatus for other status updates")
		}
		return verrs
	})
}

func validatePrimeEstimatedWeightRecordedDate(estimatedWeightRecordedDate time.Time, scheduledPickupDate time.Time) error {
	recordedYear, recordedMonth, recordedDate := estimatedWeightRecordedDate.Date()
	scheduledYear, scheduledMonth, scheduledDate := scheduledPickupDate.Date()

	if estimatedWeightRecordedDate.Before(scheduledPickupDate) {
		return nil
	}

	if recordedYear == scheduledYear && recordedMonth == scheduledMonth && recordedDate == scheduledDate {
		return nil
	}

	return apperror.InvalidInputError{}
}

func checkAvailToPrime() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		var move models.Move
		availToPrime, err := appCtx.DB().Q().
			Join("mto_shipments", "moves.id = mto_shipments.move_id").
			Where("available_to_prime_at IS NOT NULL").
			Where("mto_shipments.id = ?", newer.ID).
			Where("show = TRUE").
			Where("uses_external_vendor = FALSE").
			Exists(&move)
		if err != nil {
			return apperror.NewQueryError("Move", err, "Unexpected error")
		}
		if !availToPrime {
			return apperror.NewNotFoundError(newer.ID, "for mtoShipment")
		}
		return nil
	})
}

func checkReweighAllowed() validator {
	return validatorFunc(func(_ appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		if newer.Status != models.MTOShipmentStatusApproved && newer.Status != models.MTOShipmentStatusDiversionRequested {
			return apperror.NewConflictError(newer.ID, fmt.Sprintf("Can only reweigh a shipment that is Approved or Diversion Requested. The shipment's current status is %s", newer.Status))
		}
		if newer.Reweigh.RequestedBy != "" {
			return apperror.NewConflictError(newer.ID, "Cannot request a reweigh on a shipment that already has one.")
		}
		return nil
	})
}

// Checks if an office user is able to update a shipment based on shipment status
func checkUpdateAllowed() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, _ *models.MTOShipment, older *models.MTOShipment) error {
		msg := fmt.Sprintf("%v is not updatable", older.ID)
		err := apperror.NewForbiddenError(msg)

		if appCtx.Session().IsOfficeApp() && appCtx.Session().IsOfficeUser() {
			isServiceCounselor := appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor)
			isTOO := appCtx.Session().Roles.HasRole(roles.RoleTypeTOO)
			isTIO := appCtx.Session().Roles.HasRole(roles.RoleTypeTIO)
			switch older.Status {
			case models.MTOShipmentStatusSubmitted:
				if isServiceCounselor || isTOO {
					return nil
				}
			case models.MTOShipmentStatusApproved:
				if isTIO || isTOO {
					return nil
				}
			case models.MTOShipmentStatusCancellationRequested:
				if isTOO {
					return nil
				}
			case models.MTOShipmentStatusCanceled:
				if isTOO {
					return nil
				}
			case models.MTOShipmentStatusDiversionRequested:
				if isTOO {
					return nil
				}
			default:
				return err
			}

			return err
		}

		return err
	})
}

// Checks if a shipment can be deleted
func checkDeleteAllowed() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, _ *models.MTOShipment, older *models.MTOShipment) error {
		move := older.MoveTaskOrder
		if move.Status != models.MoveStatusDRAFT && move.Status != models.MoveStatusNeedsServiceCounseling {
			return apperror.NewForbiddenError("A shipment can only be deleted if the move is in Draft or NeedsServiceCounseling")
		}

		return nil
	})
}

// Checks if a shipment can be deleted
func checkPrimeDeleteAllowed() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, _ *models.MTOShipment, older *models.MTOShipment) error {
		if older.MoveTaskOrder.AvailableToPrimeAt == nil {
			return apperror.NewNotFoundError(older.ID, "for mtoShipment")
		}
		if older.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewForbiddenError("Prime can only delete PPM shipments")
		}
		if older.PPMShipment != nil && older.PPMShipment.Status == models.PPMShipmentStatusWaitingOnCustomer {
			return apperror.NewForbiddenError(fmt.Sprintf("A PPM shipment with the status %v cannot be deleted", models.PPMShipmentStatusWaitingOnCustomer))
		}
		return nil
	})
}

// This function checks Prime specific validations on the model
// It expects older to represent what's in the db and mtoShipment to represent the requested update
// It updates mtoShipment accordingly if there are dependent updates like requiredDeliveryDate
// On completion it either returns a list of errors or an updated MTOShipment that should be stored to the database.
func checkPrimeValidationsOnModel(planner route.Planner) validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error {
		verrs := validate.NewErrors()
		// Prime cannot edit the customer's requestedPickupDate
		if newer.RequestedPickupDate != nil {
			requestedPickupDate := newer.RequestedPickupDate
			// if !requestedPickupDate.Equal(*older.RequestedPickupDate) {
			// 	verrs.Add("requestedPickupDate", "must match what customer has requested")
			// }
			newer.RequestedPickupDate = requestedPickupDate
		}

		// Get the latest scheduled pickup date as it's needed to calculate the update range for PrimeEstimatedWeight
		// And the RDD
		latestSchedPickupDate := older.ScheduledPickupDate
		if newer.ScheduledPickupDate != nil {
			latestSchedPickupDate = newer.ScheduledPickupDate
		}

		// Prime can update the estimated weight once within a set period of time
		// If it's expired, they can no longer update it.
		latestEstimatedWeight := older.PrimeEstimatedWeight
		if newer.PrimeEstimatedWeight != nil {
			// if older.PrimeEstimatedWeight != nil {
			// 	verrs.Add("primeEstimatedWeight", "cannot be updated after initial estimation")
			// }
			// Validate if we are in the allowed period of time
			now := time.Now()
			if latestSchedPickupDate != nil {
				err := validatePrimeEstimatedWeightRecordedDate(now, *latestSchedPickupDate)
				if err != nil {
					verrs.Add("primeEstimatedWeight", "the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight")
					verrs.Add("primeEstimatedWeight", err.Error())
				}
			}
			// If they can update it, it will be the latestEstimatedWeight (needed for RDD calc)
			// And we also record the date at which it happened
			latestEstimatedWeight = newer.PrimeEstimatedWeight
			newer.PrimeEstimatedWeightRecordedDate = &now
		}

		// Prime cannot update or add agents with this endpoint, so this should always be empty
		if len(newer.MTOAgents) > 0 {
			if len(older.MTOAgents) < len(newer.MTOAgents) {
				verrs.Add("agents", "cannot add or update MTO agents to a shipment")
			}
		}

		// Prime can create a new address, but cannot update it.
		// So if address exists, return an error. But also set the local pointer to nil, so we don't recalculate requiredDeliveryDate
		var latestPickupAddress *models.Address
		var latestDestinationAddress *models.Address

		switch older.ShipmentType {
		case models.MTOShipmentTypeHHGIntoNTSDom:
			if older.StorageFacility == nil {
				// latestDestinationAddress is only used for calculating RDD.
				// We don't want to block an update because we're missing info to calculate RDD
				break
			}
			latestPickupAddress = older.PickupAddress
			latestDestinationAddress = &older.StorageFacility.Address
		case models.MTOShipmentTypeHHGOutOfNTSDom:
			if older.StorageFacility == nil {
				// latestPickupAddress is only used for calculating RDD.
				// We don't want to block an update because we're missing info to calculate RDD
				break
			}
			latestPickupAddress = &older.StorageFacility.Address
			latestDestinationAddress = older.DestinationAddress
		default:
			latestPickupAddress = older.PickupAddress
			latestDestinationAddress = older.DestinationAddress
		}
		// We also track the latestPickupAddress for the RDD calculation
		if older.PickupAddress != nil && newer.PickupAddress != nil { // If both are populated, return error
			verrs.Add("pickupAddress", "the pickup address already exists and cannot be updated with this endpoint")
		} else if newer.PickupAddress != nil { // If only the update has an address, that's the latest address
			latestPickupAddress = newer.PickupAddress
		}
		if older.DestinationAddress != nil && newer.DestinationAddress != nil {
			verrs.Add("destinationAddress", "the destination address already exists and cannot be updated with this endpoint")
		} else if newer.DestinationAddress != nil {
			latestDestinationAddress = newer.DestinationAddress
		}

		// For secondary addresses we do the same, but don't have to track the latest values for RDD
		if older.SecondaryPickupAddress != nil && newer.SecondaryPickupAddress != nil { // If both are populated, return error
			verrs.Add("secondaryPickupAddress", "the secondary pickup address already exists and cannot be updated with this endpoint")
		}
		if older.SecondaryDeliveryAddress != nil && newer.SecondaryDeliveryAddress != nil {
			verrs.Add("secondaryDeliveryAddress", "the secondary delivery address already exists and cannot be updated with this endpoint")
		}

		// If we have all the data, calculate RDD
		if latestSchedPickupDate != nil && (latestEstimatedWeight != nil || (older.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom &&
			older.NTSRecordedWeight != nil)) && latestPickupAddress != nil && latestDestinationAddress != nil {
			weight := latestEstimatedWeight
			if older.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && older.NTSRecordedWeight != nil {
				weight = older.NTSRecordedWeight
			}
			requiredDeliveryDate, err := CalculateRequiredDeliveryDate(appCtx, planner, *latestPickupAddress,
				*latestDestinationAddress, *latestSchedPickupDate, weight.Int())
			if err != nil {
				verrs.Add("requiredDeliveryDate", err.Error())
			}
			newer.RequiredDeliveryDate = requiredDeliveryDate
		}
		return verrs
	})
}
