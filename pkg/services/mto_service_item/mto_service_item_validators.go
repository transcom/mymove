package mtoserviceitem

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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

type updateMTOServiceItemValidator interface {
	validate(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error
}

// basicUpdateMTOServiceItemValidator is the type for validation that should happen no matter who uses this service object
type basicUpdateMTOServiceItemValidator struct{}

func (v *basicUpdateMTOServiceItemValidator) validate(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData) error {
	err := serviceItemData.checkLinkedIDs(appCtx)
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
	err := serviceItemData.checkLinkedIDs(appCtx)
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
	err = serviceItemData.checkPaymentRequests(appCtx)
	if err != nil {
		return err
	}

	// Checks that only SITDepartureDate is only updated for DDDSIT and DOPSIT objects
	err = serviceItemData.checkSITDeparture(appCtx)
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
	// Should use the createSITAddressUpdate endpoint
	err = serviceItemData.checkSITDestinationFinalAddress(appCtx)
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

// checkLinkedIDs checks that the user didn't attempt to change the service item's move, shipment, or reService IDs
func (v *updateMTOServiceItemData) checkLinkedIDs(appCtx appcontext.AppContext) error {
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

// checkPrimeAvailability checks that the service item is connected to a Prime-available move
func (v *updateMTOServiceItemData) checkPrimeAvailability(appCtx appcontext.AppContext) error {
	isAvailable, err := v.availabilityChecker.MTOAvailableToPrime(appCtx, v.oldServiceItem.MoveTaskOrderID)

	if !isAvailable || err != nil {
		return apperror.NewNotFoundError(v.oldServiceItem.ID, "while looking for Prime-available MTOServiceItem")
	}

	return nil
}

// checkNonPrimeFields checks that no fields were modified that are not allowed to be updated by the Prime
func (v *updateMTOServiceItemData) checkNonPrimeFields(appCtx appcontext.AppContext) error {
	if v.updatedServiceItem.Status != "" && v.updatedServiceItem.Status != v.oldServiceItem.Status {
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
func (v *updateMTOServiceItemData) checkSITDeparture(appCtx appcontext.AppContext) error {
	if v.updatedServiceItem.SITDepartureDate == nil || v.updatedServiceItem.SITDepartureDate == v.oldServiceItem.SITDepartureDate {
		return nil // the SITDepartureDate isn't being updated, so we're fine here
	}

	if v.oldServiceItem.ReService.Code == models.ReServiceCodeDDDSIT || v.oldServiceItem.ReService.Code == models.ReServiceCodeDOPSIT {
		return nil // the service item is a SIT departure service, so we're fine
	}

	return apperror.NewConflictError(v.updatedServiceItem.ID,
		fmt.Sprintf("- SIT Departure Date may only be manually updated for %s and %s service items.", models.ReServiceCodeDDDSIT, models.ReServiceCodeDOPSIT))
}

// checkSITDestinationOriginalAddress checks that SITDestinationOriginalAddress isn't being changed
func (v *updateMTOServiceItemData) checkSITDestinationOriginalAddress(appCtx appcontext.AppContext) error {
	if v.updatedServiceItem.SITDestinationOriginalAddress == nil {
		return nil // SITDestinationOriginalAddress isn't being updated, so we're fine here
	}

	if v.oldServiceItem.SITDestinationOriginalAddressID == nil {
		v.verrs.Add("SITDestinationOriginalAddress", "cannot be manually set")
		return nil // returning here to avoid nil pointer dereference error
	}

	if *v.oldServiceItem.SITDestinationOriginalAddressID != uuid.Nil &&
		v.updatedServiceItem.SITDestinationOriginalAddress != nil &&
		v.updatedServiceItem.SITDestinationOriginalAddress.ID != *v.oldServiceItem.SITDestinationOriginalAddressID {
		v.verrs.Add("SITDestinationOriginalAddress", "cannot be updated")
	}

	return nil
}

// checkSITDestinationFinalAddress checks that SITDestinationFinalAddress isn't being changed
func (v *updateMTOServiceItemData) checkSITDestinationFinalAddress(appCtx appcontext.AppContext) error {
	if v.updatedServiceItem.SITDestinationFinalAddress == nil {
		return nil // SITDestinationFinalAddress isn't being updated, so we're fine here
	}

	if v.oldServiceItem.SITDestinationFinalAddressID == nil {
		return nil // the SITDestinationFinalAddress is being created, so we're fine here
	}

	if v.oldServiceItem.ReService.Code != models.ReServiceCodeDDDSIT {
		return apperror.NewConflictError(v.updatedServiceItem.ID,
			fmt.Sprintf("- SIT Destination Final Address may only be manually created for %s service items.", models.ReServiceCodeDDDSIT))
	}

	if *v.oldServiceItem.SITDestinationFinalAddressID != uuid.Nil &&
		v.updatedServiceItem.SITDestinationFinalAddress != nil &&
		v.updatedServiceItem.SITDestinationFinalAddress.ID != *v.oldServiceItem.SITDestinationFinalAddressID {
		v.verrs.Add("SITDestinationFinalAddress", "cannot be updated")
	}

	return nil
}

// checkPaymentRequests looks for any existing payment requests connected to this service item and returns a
// Conflict Error if any are found
func (v *updateMTOServiceItemData) checkPaymentRequests(appCtx appcontext.AppContext) error {
	var paymentServiceItem models.PaymentServiceItem
	err := appCtx.DB().Where("mto_service_item_id = $1", v.updatedServiceItem.ID).First(&paymentServiceItem)

	if err == nil && paymentServiceItem.ID != uuid.Nil {
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

	// Set string fields:
	newMTOServiceItem.Reason = services.SetOptionalStringField(v.updatedServiceItem.Reason, newMTOServiceItem.Reason)

	newMTOServiceItem.Description = services.SetOptionalStringField(
		v.updatedServiceItem.Description, newMTOServiceItem.Description)

	newMTOServiceItem.RejectionReason = services.SetOptionalStringField(
		v.updatedServiceItem.RejectionReason, newMTOServiceItem.RejectionReason)

	newMTOServiceItem.SITPostalCode = services.SetOptionalStringField(
		v.updatedServiceItem.SITPostalCode, newMTOServiceItem.SITPostalCode)

	// TODO are we going to remove this field from the model at some point?
	newMTOServiceItem.PickupPostalCode = services.SetOptionalStringField(
		v.updatedServiceItem.PickupPostalCode, newMTOServiceItem.PickupPostalCode)

	// Set date fields:
	newMTOServiceItem.ApprovedAt = services.SetOptionalDateTimeField(v.updatedServiceItem.ApprovedAt, newMTOServiceItem.ApprovedAt)

	newMTOServiceItem.RejectedAt = services.SetOptionalDateTimeField(v.updatedServiceItem.RejectedAt, newMTOServiceItem.RejectedAt)

	newMTOServiceItem.SITEntryDate = services.SetOptionalDateTimeField(
		v.updatedServiceItem.SITEntryDate, newMTOServiceItem.SITEntryDate)

	newMTOServiceItem.SITDepartureDate = services.SetOptionalDateTimeField(
		v.updatedServiceItem.SITDepartureDate, newMTOServiceItem.SITDepartureDate)

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
