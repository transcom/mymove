package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/route"
)

func checkUBShipmentOCONUSRequirement() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			isShipmentOCONUS := models.IsShipmentOCONUS(*newer)
			if isShipmentOCONUS != nil && !*isShipmentOCONUS {
				verrs.Add("UB shipment error", "At least one address for a UB shipment must be OCONUS")
			}
		}
		return verrs
	})
}

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
			return apperror.NewNotFoundError(newer.ID, "not available to prime for mtoShipment")
		}
		return nil
	})
}

func checkReweighAllowed() validator {
	return validatorFunc(func(_ appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		if newer.Status != models.MTOShipmentStatusApproved && newer.Status != models.MTOShipmentStatusApprovalsRequested && newer.Status != models.MTOShipmentStatusDiversionRequested {
			return apperror.NewConflictError(newer.ID, fmt.Sprintf("Can only reweigh a shipment that is Approved, Approvals Requested, or Diversion Requested. The shipment's current status is %s", newer.Status))
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
				if isTIO || isTOO || isServiceCounselor {
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
			case models.MTOShipmentStatusApprovalsRequested:
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

		if appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
			if move.Status != models.MoveStatusDRAFT && move.Status != models.MoveStatusNeedsServiceCounseling {
				return apperror.NewForbiddenError("Service Counselor: A shipment can only be deleted if the move is in 'Draft' or 'NeedsServiceCounseling' status")
			}
		}

		if appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
			if older.Status == models.MTOShipmentStatusApproved || older.Status == models.MTOShipmentStatusApprovalsRequested {
				return apperror.NewForbiddenError("TOO: A shipment cannot be deleted if it's in Approved or Approvals Requested status")
			}
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

// helper function to check if the secondary address is empty, but the tertiary is not
func isMTOShipmentAddressCreateSequenceValid(mtoShipmentToCheck models.MTOShipment) bool {
	bothPickupAddressesEmpty := (models.IsAddressEmpty(mtoShipmentToCheck.SecondaryPickupAddress) && models.IsAddressEmpty(mtoShipmentToCheck.TertiaryPickupAddress))
	bothDestinationAddressesEmpty := (models.IsAddressEmpty(mtoShipmentToCheck.SecondaryDeliveryAddress) && models.IsAddressEmpty(mtoShipmentToCheck.TertiaryDeliveryAddress))
	bothPickupAddressesNotEmpty := !bothPickupAddressesEmpty
	bothDestinationAddressesNotEmpty := !bothDestinationAddressesEmpty
	hasNoSecondaryHasTertiaryPickup := (models.IsAddressEmpty(mtoShipmentToCheck.SecondaryPickupAddress) && !models.IsAddressEmpty(mtoShipmentToCheck.TertiaryPickupAddress))
	hasNoSecondaryHasTertiaryDestination := (models.IsAddressEmpty(mtoShipmentToCheck.SecondaryDeliveryAddress) && !models.IsAddressEmpty(mtoShipmentToCheck.TertiaryDeliveryAddress))

	// need an explicit case to capture when both are empty or not empty
	if ((bothPickupAddressesEmpty && bothDestinationAddressesEmpty) || (bothPickupAddressesNotEmpty && bothDestinationAddressesNotEmpty)) && !(hasNoSecondaryHasTertiaryPickup || hasNoSecondaryHasTertiaryDestination) {
		return true
	}
	if hasNoSecondaryHasTertiaryPickup || hasNoSecondaryHasTertiaryDestination {
		return false
	}
	return true
}

// helper function to check if the secondary address is empty, but the tertiary is not
func hasTertiaryWithNoSecondaryAddress(secondaryAddress *models.Address, tertiaryAddress *models.Address) bool {
	return (models.IsAddressEmpty(secondaryAddress) && !models.IsAddressEmpty(tertiaryAddress))
}

/* Checks if a valid address sequence is being maintained. This will return false if the tertiary address is being updated while the secondary address remains empty
*
 */
func isMTOAddressUpdateSequenceValid(shipmentToUpdateWith *models.MTOShipment, currentShipment *models.MTOShipment) bool {
	// if the incoming model has both fields, then we know the model will be updated with both secondary and tertiary addresses. therefore return true
	if !models.IsAddressEmpty(shipmentToUpdateWith.SecondaryPickupAddress) && !models.IsAddressEmpty(shipmentToUpdateWith.TertiaryPickupAddress) {
		return true
	}
	if !models.IsAddressEmpty(shipmentToUpdateWith.SecondaryDeliveryAddress) && !models.IsAddressEmpty(shipmentToUpdateWith.TertiaryDeliveryAddress) {
		return true
	}
	if currentShipment.SecondaryPickupAddress == nil && shipmentToUpdateWith.TertiaryPickupAddress != nil {
		return !hasTertiaryWithNoSecondaryAddress(currentShipment.SecondaryPickupAddress, shipmentToUpdateWith.TertiaryPickupAddress)
	}
	if currentShipment.SecondaryDeliveryAddress == nil && shipmentToUpdateWith.TertiaryDeliveryAddress != nil {
		return !hasTertiaryWithNoSecondaryAddress(currentShipment.SecondaryDeliveryAddress, shipmentToUpdateWith.TertiaryDeliveryAddress)
	}
	// no addresses are being updated, so correct address sequence is maintained, return true
	return true
}

func MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer != nil && older != nil {
			squenceIsValid := isMTOAddressUpdateSequenceValid(newer, older)
			if !squenceIsValid {
				verrs.Add("error validating mto shipment", "Shipment cannot have a third address without a second address present")
				return apperror.NewInvalidInputError(newer.ID, nil, verrs, "Invalid input found while validating the MTO shipment")
			}
		}
		return nil
	})
}
func MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer != nil {
			squenceIsValid := isMTOShipmentAddressCreateSequenceValid(*newer)
			if !squenceIsValid {
				verrs.Add("error validating mto shipment", "Shipment cannot have a third address without a second address present")
				return apperror.NewInvalidInputError(newer.ID, nil, verrs, "Invalid input found while validating the MTO shipment")
			}
		}
		return nil
	})
}

// MTOShipmentHasValidRequestedPickupDate validates the RequestedPickupDate field of an MTOShipment.
// It ensures that non-PPM shipments have a non-nil, non-zero RequestedPickupDate.
// Additionally, it checks that the RequestedPickupDate is not in the past (i.e., it must be tomorrow or later)
// when the date is newly set or updated. Returns an error if validation fails, otherwise nil.
func MTOShipmentHasValidRequestedPickupDate() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer.IsPPMShipment() {
			return nil
		}

		newerShipmentType := models.MTOShipmentType("")
		if newer != nil && newer.ShipmentType != "" {
			newerShipmentType = newer.ShipmentType
		}
		newerHasDate := newer != nil && newer.RequestedPickupDate != nil && !newer.RequestedPickupDate.IsZero()
		olderHasDate := older != nil && older.RequestedPickupDate != nil && !older.RequestedPickupDate.IsZero()
		isPickupDateRequired := newerShipmentType == models.MTOShipmentTypeBoatHaulAway || newerShipmentType == models.MTOShipmentTypeBoatTowAway ||
			newerShipmentType == models.MTOShipmentTypeMobileHome || newerShipmentType == models.MTOShipmentTypeHHGOutOfNTS
		if !newerHasDate && !isPickupDateRequired {
			verrs.Add("error validating mto shipment", "RequestedPickupDate is required to create a shipment")
			return apperror.NewInvalidInputError(newer.ID, nil, validate.NewErrors(),
				fmt.Sprintf("RequestedPickupDate is required to create %s %s shipment", GetAorAnByShipmentType(newerShipmentType), newerShipmentType))
		}
		isDateUpdated := olderHasDate && newerHasDate && !newer.RequestedPickupDate.Equal(*older.RequestedPickupDate)
		if (newerHasDate && !olderHasDate) || isDateUpdated {
			today := time.Now().Truncate(24 * time.Hour) // Truncate to date only (midnight)
			requestedDate := newer.RequestedPickupDate.Truncate(24 * time.Hour)
			if !requestedDate.After(today) {
				verrs.Add("error validating mto shipment", "Requested pickup must be greater than or equal to tomorrow's date.")
				return apperror.NewInvalidInputError(newer.ID, nil, verrs, "RequestedPickupDate must be greater than or equal to tomorrow's date.")
			}
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
			if newer.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS {
				verrs.Add("primeEstimatedWeight", "cannot be updated for nts-release shipments, please contact the TOO directly to request updates to this field")
			}
			if older.PrimeEstimatedWeight != nil {
				verrs.Add("primeEstimatedWeight", "cannot be updated after initial estimation")
			}
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

		if newer.NTSRecordedWeight != nil {
			if newer.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS {
				verrs.Add("ntsRecordedWeight", "cannot be updated for nts-release shipments, please contact the TOO directly to request updates to this field")
			}
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
		case models.MTOShipmentTypeHHGIntoNTS:
			if older.StorageFacility == nil {
				// latestDestinationAddress is only used for calculating RDD.
				// We don't want to block an update because we're missing info to calculate RDD
				break
			}
			latestPickupAddress = older.PickupAddress
			latestDestinationAddress = &older.StorageFacility.Address
		case models.MTOShipmentTypeHHGOutOfNTS:
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

		if older.TertiaryPickupAddress != nil && newer.TertiaryPickupAddress != nil { // If both are populated, return error
			verrs.Add("tertiaryPickupAddress", "the tertiary pickup address already exists and cannot be updated with this endpoint")
		}
		if older.TertiaryDeliveryAddress != nil && newer.TertiaryDeliveryAddress != nil {
			verrs.Add("tertiaryDeliveryAddress", "the tertiary delivery address already exists and cannot be updated with this endpoint")
		}

		// If we have all the data, calculate RDD
		if latestSchedPickupDate != nil && (latestEstimatedWeight != nil || (older.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS &&
			older.NTSRecordedWeight != nil)) && latestPickupAddress != nil && latestDestinationAddress != nil && older.ShipmentType != models.MTOShipmentTypeUnaccompaniedBaggage {
			weight := latestEstimatedWeight
			if older.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS && older.NTSRecordedWeight != nil {
				weight = older.NTSRecordedWeight
			}
			requiredDeliveryDate, err := CalculateRequiredDeliveryDate(appCtx, planner, *latestPickupAddress,
				*latestDestinationAddress, *latestSchedPickupDate, weight.Int(), older.MarketCode, older.MoveTaskOrderID, older.ShipmentType)
			if err != nil {
				verrs.Add("requiredDeliveryDate", err.Error())
			}
			newer.RequiredDeliveryDate = requiredDeliveryDate
		}
		return verrs
	})
}

// This func helps protect against accidental V2 intended endpoint creation hitting V1
// Eg: Somebody is testing the V2 endpoint for diversion false, divertedShipmentID true
// That example test would pass which we do not want
// ! Prime V1 rule
func protectV1Diversion() validator {
	var verrs *validate.Errors
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		// Ensure that if DivertedFromShipmentID is provided that we kick back an invalid input response
		if newer.DivertedFromShipmentID != nil {
			return apperror.NewInvalidInputError(newer.ID, nil, verrs, "The divertedFromShipmentId parameter is meant for the V2 endpoint. You are currently using the V1 endpoint.")
		}
		return nil
	})
}

// Checks if diversion parameters are valid
// ! This is a Prime V2 rule
func checkDiversionValid() validator {
	var verrs *validate.Errors
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		// Ensure that if diversion is true that DivertedFromShipmentID is provided
		if newer.Diversion && newer.DivertedFromShipmentID == nil {
			return apperror.NewInvalidInputError(newer.ID, nil, verrs, "The divertedFromShipmentId parameter must be provided if you're creating a new diversion")
		}
		// Ensure that diversion is true if a diverted from shipment ID parameter is passed in
		if !newer.Diversion && newer.DivertedFromShipmentID != nil {
			return apperror.NewInvalidInputError(newer.ID, nil, verrs, "The diversion parameter must be true if a DivertedFromShipmentID is provided")
		}
		// Ensure that the "DivertedFromShipmentID" exists if it is provided
		// Also ensure that if these conditions are met that the PrimeActualWeight has not been provided in the new diversion
		if newer.Diversion && newer.DivertedFromShipmentID != nil {
			exists, err := appCtx.DB().Q().
				Where("id = ?", *newer.DivertedFromShipmentID).
				Exists(&models.MTOShipment{})
			if err != nil {
				return apperror.NewQueryError("Move", err, "Unexpected error")
			}
			if !exists {
				return apperror.NewNotFoundError(newer.ID, "DivertedFromShipmentID shipment not found")
			}

			// Ensure that if an actual weight is provided in this shipment that we inform the user the endpoint is beign utilized incorrectly
			// The prime actual weight should be inherited from the parent if it is a diversion, not provided on creation
			if newer.PrimeActualWeight != nil {
				return apperror.NewInvalidInputError(newer.ID, nil, verrs, "The prime actual weight should not be provided inside of a newly created diversion. It will be automatically inherited by the parent. This rule does not apply for updating a shipment if that was your intention.")
			}
		}
		// Ensure that the diverted from ID is not equal to itself
		if newer.Diversion && newer.DivertedFromShipmentID == &newer.ID {
			return apperror.NewInvalidInputError(newer.ID, nil, verrs, "The DivertedFromShipmentID parameter can not be equal to the current shipment ID")
		}
		return nil
	})
}

// This func automatically sets the actual weight of the newly created shipment to be equal to the parent shipment's actual weight
func childDiversionPrimeWeightRule() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		// Ensure that if "DivertedFromShipmentID" exists, that we set the actual weight of the new shipment to be equal to the parent
		if newer.DivertedFromShipmentID != nil {
			var parentShipment models.MTOShipment
			err := appCtx.DB().Q().
				Where("id = ?", *newer.DivertedFromShipmentID).
				First(&parentShipment)
			if err != nil {
				return apperror.NewQueryError("Move", err, "Unexpected error")
			}
			if parentShipment.PrimeActualWeight == nil {
				return apperror.NewQueryError("Move", err, "Unexpected error with parent shipment actual weight being nil")
			}
			newer.PrimeActualWeight = parentShipment.PrimeActualWeight
		}
		return nil
	})
}

func GetAorAnByShipmentType(shipmentType models.MTOShipmentType) string {
	switch shipmentType {
	case models.MTOShipmentTypeHHG,
		models.MTOShipmentTypeHHGIntoNTS,
		models.MTOShipmentTypeHHGOutOfNTS,
		models.MTOShipmentTypeUnaccompaniedBaggage:
		return "an"
	default:
		return "a"
	}
}
