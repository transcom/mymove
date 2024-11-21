package mtoserviceitem

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"golang.org/x/exp/slices"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdateMTOServiceItemBasicValidator is the key for generic validation on the MTO Service Item
const UpdateMTOServiceItemBasicValidator string = "UpdateMTOServiceItemBasicValidator"

// UpdateMTOServiceItemPrimeValidator is the key for validating the MTO Service Item for the Prime contractor
const UpdateMTOServiceItemPrimeValidator string = "UpdateMTOServiceItemPrimeValidator"

// UpdateMTOServiceItemValidators is the map connecting the constant keys to the correct validator
var UpdateMTOServiceItemValidators = map[string]updateMTOServiceItemValidator{
	UpdateMTOServiceItemBasicValidator: new(basicUpdateMTOServiceItemValidator),
	UpdateMTOServiceItemPrimeValidator: new(primeUpdateMTOServiceItemValidator),
}

// Manual updates to SIT Departure dates are allowed for these service items
var (
	OriginReServiceCodesAllowedForSITDepartureDateUpdate      = []models.ReServiceCode{models.ReServiceCodeDOPSIT, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT}
	DestinationReServiceCodesAllowedForSITDepartureDateUpdate = []models.ReServiceCode{models.ReServiceCodeDDDSIT, models.ReServiceCodeDDASIT}
	ReServiceCodesAllowedForSITDepartureDateUpdate            = append(DestinationReServiceCodesAllowedForSITDepartureDateUpdate, OriginReServiceCodesAllowedForSITDepartureDateUpdate...)
)

var allSITServiceItemsToCheck = []models.ReServiceCode{
	models.ReServiceCodeDDDSIT,
	models.ReServiceCodeDDASIT,
	models.ReServiceCodeDDFSIT,
	models.ReServiceCodeDDSFSC,
	models.ReServiceCodeDOPSIT,
	models.ReServiceCodeDOFSIT,
	models.ReServiceCodeDOASIT,
	models.ReServiceCodeDOSFSC,
}

var destSITServiceItems = []models.ReServiceCode{
	models.ReServiceCodeDDDSIT,
	models.ReServiceCodeDDASIT,
	models.ReServiceCodeDDFSIT,
	models.ReServiceCodeDDSFSC,
}

type updateMTOServiceItemValidator interface {
	validate(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error
}

// basicUpdateMTOServiceItemValidator is the type for validation that should happen no matter who uses this service object
type basicUpdateMTOServiceItemValidator struct{}

func (v *basicUpdateMTOServiceItemValidator) validate(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error {
	err := serviceItemData.checkLinkedIDs()
	if err != nil {
		return err
	}

	// Checks that SITDestinationOriginalAddress isn't added/or updated using the updater
	// Should only be set when approving a service item
	err = serviceItemData.checkSITDestinationOriginalAddress(appCtx)
	if err != nil {
		return err
	}

	err = serviceItemData.getVerrs()
	if err != nil {
		return err
	}

	return nil
}

// primeUpdateMTOServiceItemValidator is the type for validation that is just for updates from the Prime contractor
type primeUpdateMTOServiceItemValidator struct{}

func (v *primeUpdateMTOServiceItemValidator) validate(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error {
	// Checks that the MTO ID, Shipment ID, and ReService IDs haven't changed
	err := serviceItemData.checkLinkedIDs()
	if err != nil {
		return err
	}

	// Checks that the Service Item is indeed available to the Prime
	err = serviceItemData.checkPrimeAvailability(appCtx)
	if err != nil {
		return err
	}

	// Checks that none of the fields that the Prime cannot update have been changed
	err = serviceItemData.checkNonPrimeFields(appCtx)
	if err != nil {
		return err
	}

	// Checks that if the updated service item has customer contacts that
	// the time fields are in the expected military time
	for index := range serviceItemData.updatedServiceItem.CustomerContacts {
		customerContacts := &serviceItemData.updatedServiceItem.CustomerContacts[index]
		err = validateTimeMilitaryField(appCtx, customerContacts.TimeMilitary)
		if err != nil {
			return err
		}
	}

	// Checks that there aren't any pending payment requests for this service item
	err = serviceItemData.checkPaymentRequests(appCtx, serviceItemData)
	if err != nil {
		return err
	}

	// Checks that only SITDepartureDate is only updated for DDDSIT and DOPSIT objects
	err = serviceItemData.checkSITDeparture(appCtx)
	if err != nil {
		return err
	}

	// Checks that the SITDepartureDate
	// - is not later than the authorized end date
	err = serviceItemData.checkSITDepartureDate(appCtx)
	if err != nil {
		return err
	}

	// Checks that only SIT Entry Date occurs AFTER the FADD
	err = serviceItemData.checkSITEntryDateAndFADD(appCtx)
	if err != nil {
		return err
	}

	// Checks that SITDestinationOriginalAddress isn't added/or updated using the updater
	// Should only be set when approving a service item
	err = serviceItemData.checkSITDestinationOriginalAddress(appCtx)
	if err != nil {
		return err
	}

	// Checks that SITDestinationFinalAddress isn't updated through this endpoint
	err = serviceItemData.checkSITDestinationFinalAddress(appCtx)
	if err != nil {
		return err
	}

	// Checks that the Old MTO SIT Service Item has a REJECTED status. If not the update req is rejected
	err = serviceItemData.checkOldServiceItemStatus(appCtx, serviceItemData)
	if err != nil {
		return err
	}

	// Check to see if the updated service item is different than the old one
	err = serviceItemData.checkForSITItemChanges(serviceItemData)
	if err != nil {
		return err
	}

	// Check to see if the updated service item is different than the old one
	err = serviceItemData.checkReasonWasUpdatedOnRejectedSIT(appCtx)
	if err != nil {
		return err
	}

	// Gets any validation errors from the above checks
	err = serviceItemData.getVerrs()
	if err != nil {
		return err
	}

	return nil
}

// updateMTOServiceItemData represents the data needed to validate an update on an MTOServiceItem
type updateMTOServiceItemData struct {
	updatedServiceItem  models.MTOServiceItem
	oldServiceItem      models.MTOServiceItem
	availabilityChecker services.MoveTaskOrderChecker
	verrs               *validate.Errors
}

// Check to see if the updatedSIT service item is different than the old one
// Turns out creating a custom comparsion method using if-statements has better performance than using a library in go
func (v *updateMTOServiceItemData) checkForSITItemChanges(serviceItemData *updateMTOServiceItemData) error {

	oldServiceItem := serviceItemData.oldServiceItem

	// Check will only be executed for serviceItems with reservice codes in the serviceItemsToCheck array
	if slices.Contains(allSITServiceItemsToCheck, oldServiceItem.ReService.Code) {

		updatedServiceItem := serviceItemData.updatedServiceItem

		// Start checking for differences. If a difference is found return nil. No need to reject the request if there are changes.
		// For now only check fields that the prime can actually submit a change for

		if updatedServiceItem.ReService.Code.String() != "" && updatedServiceItem.ReService.Code != oldServiceItem.ReService.Code {
			return nil
		}

		if updatedServiceItem.SITDepartureDate != nil && oldServiceItem.SITDepartureDate != nil {
			if updatedServiceItem.SITDepartureDate.UTC() != oldServiceItem.SITDepartureDate.UTC() {
				return nil
			}
		} else if updatedServiceItem.SITDepartureDate != nil && oldServiceItem.SITDepartureDate == nil {
			return nil
		}

		if updatedServiceItem.SITDestinationFinalAddress != nil && updatedServiceItem.SITDestinationFinalAddress != oldServiceItem.SITDestinationFinalAddress {
			return nil
		}

		if updatedServiceItem.SITCustomerContacted != nil && updatedServiceItem.SITCustomerContacted != oldServiceItem.SITCustomerContacted {
			return nil
		}

		if updatedServiceItem.SITRequestedDelivery != nil && oldServiceItem.SITRequestedDelivery != nil {
			if updatedServiceItem.SITRequestedDelivery.UTC() != oldServiceItem.SITRequestedDelivery.UTC() {
				return nil
			}
		} else if updatedServiceItem.SITRequestedDelivery != nil && oldServiceItem.SITRequestedDelivery == nil {
			return nil
		}

		if updatedServiceItem.SITEntryDate != nil && updatedServiceItem.SITEntryDate.UTC() != oldServiceItem.SITEntryDate.UTC() {
			return nil
		}

		if updatedServiceItem.Reason != nil && *updatedServiceItem.Reason != *oldServiceItem.Reason {
			return nil
		}

		if updatedServiceItem.SITPostalCode != nil && *updatedServiceItem.SITPostalCode != *oldServiceItem.SITPostalCode {
			return nil
		}

		if updatedServiceItem.RequestedApprovalsRequestedStatus != nil && *updatedServiceItem.RequestedApprovalsRequestedStatus != *oldServiceItem.RequestedApprovalsRequestedStatus {
			return nil
		}

		// If execution made it this far no changes were detected. Reject the request.
		return apperror.NewConflictError(oldServiceItem.ID,
			"- To re-submit a SIT sevice item the new SIT service item must be different than the previous one.")

	}

	return nil
}

// checkLinkedIDs checks that the user didn't attempt to change the service item's move, shipment, or reService IDs
func (v *updateMTOServiceItemData) checkLinkedIDs() error {
	if v.updatedServiceItem.MoveTaskOrderID != uuid.Nil && v.updatedServiceItem.MoveTaskOrderID != v.oldServiceItem.MoveTaskOrderID {
		v.verrs.Add("moveTaskOrderID", "cannot be updated")
	}
	if v.updatedServiceItem.MTOShipmentID != nil && *v.updatedServiceItem.MTOShipmentID != *v.oldServiceItem.MTOShipmentID {
		v.verrs.Add("mtoShipmentID", "cannot be updated")
	}
	if v.updatedServiceItem.ReServiceID != uuid.Nil && v.updatedServiceItem.ReServiceID != v.oldServiceItem.ReServiceID {
		v.verrs.Add("reServiceID", "cannot be updated")
	}

	return nil
}

// checkOldServiceItemStatus checks that the old service item has a REJECTED status
func (v *updateMTOServiceItemData) checkOldServiceItemStatus(_ appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error {

	if slices.Contains(allSITServiceItemsToCheck, serviceItemData.oldServiceItem.ReService.Code) {
		if serviceItemData.oldServiceItem.Status == models.MTOServiceItemStatusRejected {
			return nil
		} else if serviceItemData.oldServiceItem.Status == models.MTOServiceItemStatusApproved {

			invalidFieldChange := false
			// Fields that are not allowed to change when status is approved

			if serviceItemData.updatedServiceItem.ReService.Code.String() != "" && serviceItemData.updatedServiceItem.ReService.Code.String() != serviceItemData.oldServiceItem.ReService.Code.String() {
				invalidFieldChange = true
			}

			if serviceItemData.updatedServiceItem.SITEntryDate != nil {
				invalidFieldChange = true
			}

			if serviceItemData.updatedServiceItem.Reason != nil {
				invalidFieldChange = true
			}

			if serviceItemData.updatedServiceItem.SITPostalCode != nil {
				invalidFieldChange = true
			}

			if serviceItemData.updatedServiceItem.RequestedApprovalsRequestedStatus != nil {
				invalidFieldChange = true
			}

			if invalidFieldChange {
				return apperror.NewConflictError(serviceItemData.oldServiceItem.ID,
					"- one or more fields is not allowed to be updated when the SIT service item has an approved status.")
			}

			// Fields allowed to changed when status is approved
			if serviceItemData.updatedServiceItem.SITDepartureDate != nil {
				serviceItemData.updatedServiceItem.Status = models.MTOServiceItemStatusApproved
				return nil
			}
			if serviceItemData.updatedServiceItem.SITRequestedDelivery != nil {
				serviceItemData.updatedServiceItem.Status = models.MTOServiceItemStatusApproved
				return nil
			}
			if serviceItemData.updatedServiceItem.SITCustomerContacted != nil {
				serviceItemData.updatedServiceItem.Status = models.MTOServiceItemStatusApproved
				return nil
			}

			return apperror.NewConflictError(serviceItemData.oldServiceItem.ID,
				"- unknown field or fields attempting to be updated.")
			//nolint:revive // This is intentionally returning an error
		} else {
			// Rejects the update if the original SIT does not have a REJECTED status
			return apperror.NewConflictError(serviceItemData.oldServiceItem.ID,
				"- this SIT service item cannot be updated because the status is not in an editable state.")
		}
	}

	return nil
}

// checkReasonWasUpdatedOnRejectedSIT checks that the reason provided is different from the previously rejected SIT service item
func (v *updateMTOServiceItemData) checkReasonWasUpdatedOnRejectedSIT(_ appcontext.AppContext) error {

	if slices.Contains(allSITServiceItemsToCheck, v.oldServiceItem.ReService.Code) {
		// if the SI is approved then we don't need to check this
		if v.oldServiceItem.Status == models.MTOServiceItemStatusApproved {
			return nil
		}
		if v.oldServiceItem.Status == models.MTOServiceItemStatusRejected && v.updatedServiceItem.Status == models.MTOServiceItemStatusSubmitted {
			if v.updatedServiceItem.Reason == nil {
				return apperror.NewConflictError(v.oldServiceItem.ID,
					"- you must provide a new reason when resubmitting a previously rejected SIT service item")
			}
			if *v.updatedServiceItem.Reason == "" {
				return apperror.NewConflictError(v.oldServiceItem.ID,
					"- reason cannot be empty when resubmitting a previously rejected SIT service item")
			}
			if v.updatedServiceItem.Reason != nil && *v.updatedServiceItem.Reason == *v.oldServiceItem.Reason {
				return apperror.NewConflictError(v.oldServiceItem.ID,
					"- please provide a new reason when resubmitting a previously rejected SIT service item")
			}
		}
	}

	return nil
}

// checkPrimeAvailability checks that the service item is connected to a Prime-available move
func (v *updateMTOServiceItemData) checkPrimeAvailability(appCtx appcontext.AppContext) error {
	isAvailable, err := v.availabilityChecker.MTOAvailableToPrime(appCtx, v.oldServiceItem.MoveTaskOrderID)

	if !isAvailable || err != nil {
		return apperror.NewNotFoundError(v.oldServiceItem.ID, "while looking for Prime-available MTOServiceItem")
	}

	return nil
}

// checkNonPrimeFields checks that no fields were modified that are not allowed to be updated by the Prime
func (v *updateMTOServiceItemData) checkNonPrimeFields(_ appcontext.AppContext) error {

	if v.updatedServiceItem.Status != "" && v.updatedServiceItem.Status != v.oldServiceItem.Status && (!slices.Contains(allSITServiceItemsToCheck, v.oldServiceItem.ReService.Code)) {
		v.verrs.Add("status", "cannot be updated")
	}

	if v.updatedServiceItem.RejectionReason != nil && v.updatedServiceItem.RejectionReason != v.oldServiceItem.RejectionReason {
		v.verrs.Add("rejectionReason", "cannot be updated")
	}

	if v.updatedServiceItem.ApprovedAt != nil && v.updatedServiceItem.ApprovedAt != v.oldServiceItem.ApprovedAt {
		v.verrs.Add("approvedAt", "cannot be updated")
	}

	if v.updatedServiceItem.RejectedAt != nil && v.updatedServiceItem.RejectedAt != v.oldServiceItem.RejectedAt {
		v.verrs.Add("rejectedAt", "cannot be updated")
	}

	return nil
}

// checkSITDeparture checks that the service item is a DDDSIT or DOPSIT if the user is trying to update the
// SITDepartureDate
func (v *updateMTOServiceItemData) checkSITDeparture(_ appcontext.AppContext) error {
	if v.updatedServiceItem.SITDepartureDate == nil || v.updatedServiceItem.SITDepartureDate == v.oldServiceItem.SITDepartureDate {
		return nil // the SITDepartureDate isn't being updated, so we're fine here
	}

	if slices.Contains(allSITServiceItemsToCheck, v.oldServiceItem.ReService.Code) {
		return nil // the service item is a SIT departure service or SIT Domestic origin 1st day SIT , so we're fine
	}

	return apperror.NewConflictError(v.updatedServiceItem.ID,
		fmt.Sprintf("- SIT Departure Date may only be manually updated for the following service items: %s, %s, %s, %s, %s, %s, %s, %s",
			models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDOPSIT,
			models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDSFSC, models.ReServiceCodeDOSFSC))
}

// checkSITEntryDateAndFADD checks that the SIT entry date occurs after the FADD for destination SIT
func (v *updateMTOServiceItemData) checkSITEntryDateAndFADD(_ appcontext.AppContext) error {
	if slices.Contains(destSITServiceItems, v.oldServiceItem.ReService.Code) {
		// if they are attempting to update the SIT entry date
		if v.updatedServiceItem.SITEntryDate != nil {
			var contacts models.MTOServiceItemCustomerContacts
			if v.updatedServiceItem.CustomerContacts != nil {
				contacts = v.updatedServiceItem.CustomerContacts
			} else if v.oldServiceItem.CustomerContacts != nil {
				contacts = v.oldServiceItem.CustomerContacts
			}

			for _, contact := range contacts {
				if !contact.FirstAvailableDeliveryDate.IsZero() &&
					v.updatedServiceItem.SITEntryDate.Before(contact.FirstAvailableDeliveryDate) {
					return apperror.NewUnprocessableEntityError(fmt.Sprintf("the SIT Entry Date (%s) cannot be before the First Available Delivery Date (%s)",
						v.updatedServiceItem.SITEntryDate.Format("2006-01-02"), contact.FirstAvailableDeliveryDate.Format("2006-01-02")))
				}
			}
		}
		return nil
	}
	return nil
}

// checkSITDepartureDate checks that the SITDepartureDate:
// - is not later than the authorized end date
func (v *updateMTOServiceItemData) checkSITDepartureDate(_ appcontext.AppContext) error {
	if v.updatedServiceItem.SITDepartureDate == nil || v.updatedServiceItem.SITDepartureDate == v.oldServiceItem.SITDepartureDate {
		return nil // the SITDepartureDate isn't being updated, so we're fine here
	}

	if v.updatedServiceItem.SITDepartureDate != nil {
		if v.oldServiceItem.MTOShipmentID == nil || *v.oldServiceItem.MTOShipmentID == uuid.Nil {
			// If we are updating the SIT departure date, there must be an attached MTO shipment to this service item
			// Otherwise, authorized end dates cannot be calculated or validated properly
			return apperror.NewInternalServerError(fmt.Sprintf("The requested service item updates for ID %s did not have an attached MTO Shipment, preventing proper lookup of the authorized end date. This occurs on the server not preloading necessary data.", v.updatedServiceItem.ID))
		}
		// Set the SIT entry date we are going to use for comparison
		if v.updatedServiceItem.SITEntryDate == nil && v.oldServiceItem.SITEntryDate == nil {
			return apperror.NewInternalServerError(fmt.Sprintf("The requested service item updates for ID %s did not have a SIT entry date attached.", v.updatedServiceItem.ID))
		}
		SITEntryDate := v.oldServiceItem.SITEntryDate
		if v.updatedServiceItem.SITEntryDate != nil {
			SITEntryDate = v.updatedServiceItem.SITEntryDate
		}
		// Check that departure date is not before the current entry date
		if v.updatedServiceItem.SITDepartureDate.Before(*SITEntryDate) {
			v.verrs.Add("SITDepartureDate", "SIT departure date cannot be set before the SIT entry date.")
		}
	}
	return nil
}

// checkSITDestinationOriginalAddress checks that SITDestinationOriginalAddress isn't being changed
func (v *updateMTOServiceItemData) checkSITDestinationOriginalAddress(_ appcontext.AppContext) error {
	if v.updatedServiceItem.SITDestinationOriginalAddress == nil {
		return nil // SITDestinationOriginalAddress isn't being updated, so we're fine here
	}

	if v.oldServiceItem.SITDestinationOriginalAddressID == nil {
		v.verrs.Add("SITDestinationOriginalAddress", "cannot be manually set")
		return apperror.NewInvalidInputError(v.updatedServiceItem.ID, nil, v.verrs,
			"SITDestinationOriginalAddress cannot be manually set")
	}

	if *v.oldServiceItem.SITDestinationOriginalAddressID != uuid.Nil &&
		v.updatedServiceItem.SITDestinationOriginalAddress != nil &&
		v.updatedServiceItem.SITDestinationOriginalAddress.ID != *v.oldServiceItem.SITDestinationOriginalAddressID {
		v.verrs.Add("SITDestinationOriginalAddress", "cannot be updated")
		return apperror.NewInvalidInputError(v.updatedServiceItem.ID, nil, v.verrs,
			"SITDestinationOriginalAddress cannot be updated")
	}

	return nil
}

// checkSITDestinationFinalAddress checks that SITDestinationFinalAddress isn't being changed
func (v *updateMTOServiceItemData) checkSITDestinationFinalAddress(_ appcontext.AppContext) error {
	if v.updatedServiceItem.SITDestinationFinalAddress == nil {
		return nil // SITDestinationFinalAddress isn't being updated, so we're fine here
	}

	if v.oldServiceItem.SITDestinationFinalAddressID == nil {
		return nil // the SITDestinationFinalAddress is being created, so we're fine here
	}

	reServiceCodesDestination := []models.ReServiceCode{models.ReServiceCodeDDDSIT, models.ReServiceCodeDDASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC}
	if slices.Contains(reServiceCodesDestination, v.oldServiceItem.ReService.Code) {
		v.verrs.Add("SITDestinationFinalAddress", "Update the shipment destination address to update the service item's SIT final destination address.")
		return apperror.NewInvalidInputError(v.updatedServiceItem.ID, nil, v.verrs,
			"Update the shipment destination address to update the service item's SIT final destination address.")
	}

	if *v.oldServiceItem.SITDestinationFinalAddressID != uuid.Nil &&
		v.updatedServiceItem.SITDestinationFinalAddress != nil &&
		v.updatedServiceItem.SITDestinationFinalAddress.ID != *v.oldServiceItem.SITDestinationFinalAddressID {
		v.verrs.Add("SITDestinationFinalAddress", "Update the shipment destination address to update the service item's SIT final destination address.")
		return apperror.NewInvalidInputError(v.updatedServiceItem.ID, nil, v.verrs,
			"Update the shipment destination address to update the service item's SIT final destination address.")
	}

	return nil
}

// checkPaymentRequests looks for any existing payment requests connected to this service item and returns a
// Conflict Error if any are found
func (v *updateMTOServiceItemData) checkPaymentRequests(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error {
	var paymentServiceItem models.PaymentServiceItem

	// Check what fields are being updated to allow this update
	allowUpdateBasedOnAllowableFieldChange := paymentRequestCheckAllowableFieldCheck(serviceItemData)

	err := appCtx.DB().Where("mto_service_item_id = $1", v.updatedServiceItem.ID).First(&paymentServiceItem)

	if err == nil && paymentServiceItem.ID != uuid.Nil && !allowUpdateBasedOnAllowableFieldChange {
		return apperror.NewConflictError(v.updatedServiceItem.ID,
			"- this service item has an existing payment request and can no longer be updated.")
	} else if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		return err
	}
	// NOTE the third error case is when there are no payment requests found, which is good in this case!

	return nil
}

// getVerrs looks for any validation errors and returns a formatted InvalidInputError if any are found.
// Should only be called after the other check methods have been called.
func (v *updateMTOServiceItemData) getVerrs() error {
	if v.verrs.HasAny() {
		return apperror.NewInvalidInputError(v.updatedServiceItem.ID, nil, v.verrs,
			"Invalid input found while validating the service item.")
	}

	return nil
}

// setNewMTOServiceItem compares updatedServiceItem and oldServiceItem and updates a new MTOServiceItem instance with
// all data (changed and unchanged) filled in. Does not return an error, data must be checked for validation before
// this step.
func (v *updateMTOServiceItemData) setNewMTOServiceItem() *models.MTOServiceItem {
	newMTOServiceItem := v.oldServiceItem

	if v.updatedServiceItem.Status != "" {
		newMTOServiceItem.Status = v.updatedServiceItem.Status
	}

	// If the updated RequestedApprovalsRequestedStatus param is not null/nil then update the new serviceItem
	if v.updatedServiceItem.RequestedApprovalsRequestedStatus != nil {
		newMTOServiceItem.RequestedApprovalsRequestedStatus = v.updatedServiceItem.RequestedApprovalsRequestedStatus
	}

	// Set string fields:
	newMTOServiceItem.Reason = services.SetOptionalStringField(v.updatedServiceItem.Reason, newMTOServiceItem.Reason)

	newMTOServiceItem.Description = services.SetOptionalStringField(
		v.updatedServiceItem.Description, newMTOServiceItem.Description)

	newMTOServiceItem.RejectionReason = services.SetOptionalStringField(
		v.updatedServiceItem.RejectionReason, newMTOServiceItem.RejectionReason)

	if v.updatedServiceItem.SITPostalCode != nil {
		newMTOServiceItem.SITPostalCode = services.SetOptionalStringField(
			v.updatedServiceItem.SITPostalCode, newMTOServiceItem.SITPostalCode)
	}

	newMTOServiceItem.SITPostalCode = services.SetOptionalStringField(
		v.updatedServiceItem.SITPostalCode, newMTOServiceItem.SITPostalCode)

	// TODO are we going to remove this field from the model at some point?
	newMTOServiceItem.PickupPostalCode = services.SetOptionalStringField(
		v.updatedServiceItem.PickupPostalCode, newMTOServiceItem.PickupPostalCode)

	// Set date fields:
	newMTOServiceItem.ApprovedAt = services.SetOptionalDateTimeField(v.updatedServiceItem.ApprovedAt, newMTOServiceItem.ApprovedAt)

	newMTOServiceItem.RejectedAt = services.SetOptionalDateTimeField(v.updatedServiceItem.RejectedAt, newMTOServiceItem.RejectedAt)

	if v.updatedServiceItem.SITEntryDate != nil {
		newMTOServiceItem.SITEntryDate = services.SetOptionalDateTimeField(
			v.updatedServiceItem.SITEntryDate, newMTOServiceItem.SITEntryDate)
	}

	newMTOServiceItem.SITEntryDate = services.SetOptionalDateTimeField(
		v.updatedServiceItem.SITEntryDate, newMTOServiceItem.SITEntryDate)

	if v.updatedServiceItem.SITDepartureDate != nil {
		newMTOServiceItem.SITDepartureDate = services.SetOptionalDateTimeField(
			v.updatedServiceItem.SITDepartureDate, newMTOServiceItem.SITDepartureDate)
	}

	newMTOServiceItem.SITCustomerContacted = services.SetOptionalDateTimeField(v.updatedServiceItem.SITCustomerContacted, newMTOServiceItem.SITCustomerContacted)
	newMTOServiceItem.SITRequestedDelivery = services.SetOptionalDateTimeField(v.updatedServiceItem.SITRequestedDelivery, newMTOServiceItem.SITRequestedDelivery)

	if v.updatedServiceItem.SITDestinationFinalAddress != nil {
		newMTOServiceItem.SITDestinationFinalAddress = v.updatedServiceItem.SITDestinationFinalAddress
		newMTOServiceItem.SITDestinationFinalAddressID = &v.updatedServiceItem.SITDestinationFinalAddress.ID
	}

	// Set customer contact fields
	newMTOServiceItem.CustomerContacts = v.setNewCustomerContacts()

	// Set weight fields:
	newMTOServiceItem.EstimatedWeight = services.SetOptionalPoundField(
		v.updatedServiceItem.EstimatedWeight, newMTOServiceItem.EstimatedWeight)

	newMTOServiceItem.ActualWeight = services.SetOptionalPoundField(
		v.updatedServiceItem.ActualWeight, newMTOServiceItem.ActualWeight)

	newMTOServiceItem.PricingEstimate = services.SetNoNilOptionalCentField(
		v.updatedServiceItem.PricingEstimate, newMTOServiceItem.PricingEstimate)

	return &newMTOServiceItem
}

func (v *updateMTOServiceItemData) setNewCustomerContacts() models.MTOServiceItemCustomerContacts {
	// If there are no updated customer contacts we will just use the old ones, it doesn't matter if there are no old ones.
	if len(v.updatedServiceItem.CustomerContacts) == 0 {
		return v.oldServiceItem.CustomerContacts
	}

	// If there are no old customer contacts we will just use the updated ones.
	if len(v.oldServiceItem.CustomerContacts) == 0 {
		return v.updatedServiceItem.CustomerContacts
	}

	var newCustomerContacts models.MTOServiceItemCustomerContacts

	// Iterate through the updated and the old customer contacts to see if they correspond to one another.
	for _, updatedCustomerContact := range v.updatedServiceItem.CustomerContacts {
		foundCorrespondingOldContact := false
		var newCustomerContact models.MTOServiceItemCustomerContact
		for _, oldCustomerContact := range v.oldServiceItem.CustomerContacts {
			// We use the type field to determine if the CustomerContacts correspond to each other
			// If they correspond we update the information on the old CustomerContact
			if updatedCustomerContact.Type == oldCustomerContact.Type {
				newCustomerContact = oldCustomerContact
				newCustomerContact.DateOfContact = updatedCustomerContact.DateOfContact
				newCustomerContact.TimeMilitary = updatedCustomerContact.TimeMilitary
				newCustomerContact.FirstAvailableDeliveryDate = updatedCustomerContact.FirstAvailableDeliveryDate
				foundCorrespondingOldContact = true
			}
		}
		// If there is no corresponding old CustomerContact we use the updated CustomerContact
		if !foundCorrespondingOldContact {
			newCustomerContact = updatedCustomerContact
		}
		newCustomerContacts = append(newCustomerContacts, newCustomerContact)
	}

	// We need to iterate once more through the old CustomerContacts
	// to find any that don't have a corresponding updated CustomerContact
	for _, oldCustomerContact := range v.oldServiceItem.CustomerContacts {
		foundCorrespondingUpdatedContact := false
		for _, updatedCustomerContact := range v.updatedServiceItem.CustomerContacts {
			if updatedCustomerContact.Type == oldCustomerContact.Type {
				foundCorrespondingUpdatedContact = true
			}
		}
		if !foundCorrespondingUpdatedContact {
			newCustomerContact := oldCustomerContact
			newCustomerContacts = append(newCustomerContacts, newCustomerContact)
		}
	}
	return newCustomerContacts
}

func paymentRequestCheckAllowableFieldCheck(serviceItemData *updateMTOServiceItemData) bool {

	allowableFieldChange, disallowedFieldChange := false, false

	// Fields allowed to change when service item has a payment request
	if serviceItemData.updatedServiceItem.SITDepartureDate != nil ||
		serviceItemData.updatedServiceItem.SITRequestedDelivery != nil ||
		serviceItemData.updatedServiceItem.SITCustomerContacted != nil {
		allowableFieldChange = true
	}

	// Fields not allowed to change when service item has a payment request
	if serviceItemData.updatedServiceItem.ReService.Code.String() != "" &&
		serviceItemData.updatedServiceItem.ReService.Code.String() != serviceItemData.oldServiceItem.ReService.Code.String() ||
		serviceItemData.updatedServiceItem.SITEntryDate != nil ||
		serviceItemData.updatedServiceItem.Reason != nil ||
		serviceItemData.updatedServiceItem.SITPostalCode != nil ||
		serviceItemData.updatedServiceItem.RequestedApprovalsRequestedStatus != nil {
		disallowedFieldChange = true
	}

	if allowableFieldChange && !disallowedFieldChange {
		return true
	}

	return false
}

// checkDuplicateServiceCodes checks if the move or shipment has a duplicate service item with the same code as the one requested
func (o *mtoServiceItemCreator) checkDuplicateServiceCodes(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem) error {
	var duplicateServiceItem models.MTOServiceItem

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("move_id", "=", serviceItem.MoveTaskOrderID),
		query.NewQueryFilter("re_service_id", "=", serviceItem.ReServiceID),
	}
	if serviceItem.MTOShipmentID != nil {
		queryFilters = append(queryFilters, query.NewQueryFilter("mto_shipment_id", "=", serviceItem.MTOShipmentID))
	}

	// We DON'T want to find this service item:
	err := o.builder.FetchOne(appCtx, &duplicateServiceItem, queryFilters)
	if err == nil && duplicateServiceItem.ID != uuid.Nil {
		return apperror.NewConflictError(duplicateServiceItem.ID,
			fmt.Sprintf("for creating a service item. A service item with reServiceCode %s already exists for this move and/or shipment.", serviceItem.ReService.Code))
	} else if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		return err
	}

	return nil
}

// checkSITEntryDateAndFADD checks that the SIT entry date is not before the First Available Delivery Date
func (o *mtoServiceItemCreator) checkSITEntryDateAndFADD(serviceItem *models.MTOServiceItem) error {
	if serviceItem.SITEntryDate == nil {
		return nil
	}

	for _, contact := range serviceItem.CustomerContacts {
		if !contact.FirstAvailableDeliveryDate.IsZero() &&
			serviceItem.SITEntryDate.Before(contact.FirstAvailableDeliveryDate) {
			return apperror.NewUnprocessableEntityError(fmt.Sprintf("the SIT Entry Date (%s) cannot be before the First Available Delivery Date (%s)",
				serviceItem.SITEntryDate.Format("2006-01-02"), contact.FirstAvailableDeliveryDate.Format("2006-01-02")))
		}
	}

	return nil
}
