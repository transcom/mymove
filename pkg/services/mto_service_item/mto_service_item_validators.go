package mtoserviceitem

import (
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdateMTOServiceItemBaseValidator is the key for generic validation on the MTO Service Item
const UpdateMTOServiceItemBaseValidator string = "UpdateMTOServiceItemBaseValidator"

// UpdateMTOServiceItemPrimeValidator is the key for validating the MTO Service Item for the Prime contractor
const UpdateMTOServiceItemPrimeValidator string = "UpdateMTOServiceItemPrimeValidator"

// UpdateMTOServiceItemValidators is the map connecting the constant keys to the correct validator
var UpdateMTOServiceItemValidators = map[string]updateMTOServiceItemValidator{
	UpdateMTOServiceItemBaseValidator:  new(baseUpdateMTOServiceItemValidator),
	UpdateMTOServiceItemPrimeValidator: new(primeUpdateMTOServiceItemValidator),
}

type updateMTOServiceItemValidator interface {
	validate(agentData *updateMTOServiceItemData) error
}

// baseUpdateMTOServiceItemValidator is the type for validation that should happen no matter who uses this service object
type baseUpdateMTOServiceItemValidator struct{}

func (v *baseUpdateMTOServiceItemValidator) validate(serviceItemData *updateMTOServiceItemData) error {
	err := serviceItemData.checkLinkedIDs()
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

func (v *primeUpdateMTOServiceItemValidator) validate(serviceItemData *updateMTOServiceItemData) error {
	err := serviceItemData.checkLinkedIDs()
	if err != nil {
		return err
	}

	err = serviceItemData.checkPrimeAvailability()
	if err != nil {
		return err
	}

	err = serviceItemData.checkNonPrimeFields()
	if err != nil {
		return err
	}

	err = serviceItemData.getVerrs()
	if err != nil {
		return err
	}

	return nil
}

// updateMTOServiceItemData represents the data needed to validate an update on an MTOServiceItem
type updateMTOServiceItemData struct {
	updatedServiceItem models.MTOServiceItem
	oldServiceItem     models.MTOServiceItem
	builder            mtoServiceItemQueryBuilder
	verrs              *validate.Errors
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

// checkPrimeAvailability checks that the service item is connected to a Prime-available move
func (v *updateMTOServiceItemData) checkPrimeAvailability() error {
	// NOTE: We cannot use the MoveTaskOrderChecker here because this service uses QueryBuilder and doesn't have access
	// to the DB
	var move models.Move
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", v.updatedServiceItem.MoveTaskOrderID),
	}
	err := v.builder.FetchOne(&move, queryFilters)
	if err != nil {
		return services.NewNotFoundError(v.updatedServiceItem.MoveTaskOrderID, "for a move connected to this service item")
	}

	if move.AvailableToPrimeAt == nil || move.AvailableToPrimeAt.IsZero() {
		return services.NewNotFoundError(v.updatedServiceItem.ID, "while looking for Prime-available MTOServiceItem")
	}

	return nil
}

// checkSITDeparture checks that
func (v *updateMTOServiceItemData) checkSITDeparture() error {
	// Check that the service item is actually a SIT departure service:
	if v.oldServiceItem.ReService.Code != models.ReServiceCodeDDDSIT && v.oldServiceItem.ReService.Code != models.ReServiceCodeDOPSIT {
		// Return an error if the user tried to update the departure date for a non-departure SIT service:
		if v.updatedServiceItem.SITDepartureDate != nil && v.updatedServiceItem.SITDepartureDate != v.oldServiceItem.SITDepartureDate {
			return services.NewConflictError(v.updatedServiceItem.ID,
				"- SIT Departure Date may only be manually updated for DDDSIT and DOPSIT service items.")
		}

		return nil // no need to continue for a non-departure SIT service
	}

	// Check that there is no existing payment request for this service item:
	var paymentServiceItem models.PaymentServiceItem
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("mto_service_item_id", "=", v.updatedServiceItem.ID),
	}

	err := v.builder.FetchOne(&paymentServiceItem, queryFilters)
	if err == nil && paymentServiceItem.ID != uuid.Nil {
		return services.NewConflictError(v.updatedServiceItem.ID,
			"- cannot update the SIT Departure Date for a service item with an existing payment request.")
	} else if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		return err
	}

	// TODO do we need to check anything else?

	return nil
}

// checkNonPrimeFields checks that no fields were modified that are not allowed to be updated by the Prime
func (v *updateMTOServiceItemData) checkNonPrimeFields() error {
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

// getVerrs looks for any validation errors and returns a formatted InvalidInputError if any are found.
// Should only be called after the other check methods have been called.
func (v *updateMTOServiceItemData) getVerrs() error {
	if v.verrs.HasAny() {
		return services.NewInvalidInputError(v.updatedServiceItem.ID, nil, v.verrs,
			"Invalid input found while validating the service item.")
	}

	return nil
}

// setNewMTOServiceItem compares updatedServiceItem and oldServiceItem and updates a new MTOServiceItem instance with
// all data (changed and unchanged) filled in. Does not return an error, data must be checked for validation before
// this step.
func (v *updateMTOServiceItemData) setNewMTOServiceItem() *models.MTOServiceItem {
	newMTOServiceItem := v.oldServiceItem

	if v.updatedServiceItem.Reason != nil {
		newMTOServiceItem.Reason = v.updatedServiceItem.Reason

		if *v.updatedServiceItem.Reason == "" {
			newMTOServiceItem.Reason = nil
		}
	}

	if v.updatedServiceItem.Description != nil {
		newMTOServiceItem.Description = v.updatedServiceItem.Description

		if *v.updatedServiceItem.Description == "" {
			newMTOServiceItem.Description = nil
		}
	}

	if v.updatedServiceItem.Status != "" {
		newMTOServiceItem.Status = v.updatedServiceItem.Status
	}

	if v.updatedServiceItem.RejectionReason != nil {
		newMTOServiceItem.RejectionReason = v.updatedServiceItem.RejectionReason

		if *v.updatedServiceItem.RejectionReason == "" {
			newMTOServiceItem.RejectionReason = nil
		}
	}

	// TODO must be IsZero? How to nullify dates reliably? approved, rejected, sit entry, sit departure
	if v.updatedServiceItem.ApprovedAt != nil {
		newMTOServiceItem.ApprovedAt = v.updatedServiceItem.ApprovedAt
	}
	if v.updatedServiceItem.RejectedAt != nil {
		newMTOServiceItem.RejectedAt = v.updatedServiceItem.RejectedAt
	}

	// TODO are we going to remove this field from the model at some point?
	if v.updatedServiceItem.PickupPostalCode != nil {
		newMTOServiceItem.PickupPostalCode = v.updatedServiceItem.PickupPostalCode

		if *v.updatedServiceItem.PickupPostalCode == "" {
			newMTOServiceItem.PickupPostalCode = nil
		}
	}

	if v.updatedServiceItem.SITPostalCode != nil {
		newMTOServiceItem.SITPostalCode = v.updatedServiceItem.SITPostalCode

		if *v.updatedServiceItem.SITPostalCode == "" {
			newMTOServiceItem.SITPostalCode = nil
		}
	}

	if v.updatedServiceItem.SITEntryDate != nil {
		newMTOServiceItem.SITEntryDate = v.updatedServiceItem.SITEntryDate
	}

	if v.updatedServiceItem.SITDepartureDate != nil {
		newMTOServiceItem.SITDepartureDate = v.updatedServiceItem.SITDepartureDate
	}

	return &newMTOServiceItem
}
