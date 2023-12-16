package mtoserviceitem

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
)

type mtoServiceItemQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoServiceItemUpdater struct {
	builder          mtoServiceItemQueryBuilder
	createNewBuilder func() mtoServiceItemQueryBuilder
	moveRouter       services.MoveRouter
	shipmentFetcher  services.MTOShipmentFetcher
	addressCreator   services.AddressCreator
}

// NewMTOServiceItemUpdater returns a new mto service item updater
func NewMTOServiceItemUpdater(builder mtoServiceItemQueryBuilder, moveRouter services.MoveRouter, shipmentFetcher services.MTOShipmentFetcher, addressCreator services.AddressCreator) services.MTOServiceItemUpdater {
	// used inside a transaction and mocking		return &mtoServiceItemUpdater{builder: builder}
	createNewBuilder := func() mtoServiceItemQueryBuilder {
		return query.NewQueryBuilder()
	}

	return &mtoServiceItemUpdater{builder, createNewBuilder, moveRouter, shipmentFetcher, addressCreator}
}

func (p *mtoServiceItemUpdater) ApproveOrRejectServiceItem(
	appCtx appcontext.AppContext,
	mtoServiceItemID uuid.UUID,
	status models.MTOServiceItemStatus,
	rejectionReason *string,
	eTag string,
) (*models.MTOServiceItem, error) {
	mtoServiceItem, err := p.findServiceItem(appCtx, mtoServiceItemID)
	if err != nil {
		return &models.MTOServiceItem{}, err
	}

	return p.approveOrRejectServiceItem(appCtx, *mtoServiceItem, status, rejectionReason, eTag, checkMoveStatus(), checkETag())
}

func (p *mtoServiceItemUpdater) ConvertItemToMembersExpense(
	appCtx appcontext.AppContext,
	shipment *models.MTOShipment,
) (*models.MTOServiceItem, error) {
	var DOFSITCodeID uuid.UUID
	reServiceErr := appCtx.DB().RawQuery(`SELECT id FROM re_services WHERE code = 'DOFSIT'`).First(&DOFSITCodeID) // First get uuid for DOFSIT service code
	if reServiceErr != nil {
		return nil, reServiceErr
	}

	// Now get the DOFSIT service item associated with the current mto_shipment
	var SITItem models.MTOServiceItem
	getSITItemErr := appCtx.DB().RawQuery(`SELECT * FROM mto_service_items WHERE re_service_id = ? AND mto_shipment_id = ?`, DOFSITCodeID, shipment.ID).First(&SITItem)
	if getSITItemErr != nil {
		switch getSITItemErr {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipment.ID, "for MTO Service Item")
		default:
			return nil, getSITItemErr
		}
	}

	eTag := etag.GenerateEtag(SITItem.UpdatedAt)

	// Finally, update the mto_service_item with the members_expense flag set to TRUE
	SITItem.MembersExpense = models.BoolPointer(true)
	mtoServiceItem, err := p.findServiceItem(appCtx, SITItem.ID)
	if err != nil {
		return &models.MTOServiceItem{}, err
	}

	return p.convertItemToMembersExpense(appCtx, *mtoServiceItem, eTag, checkETag())
}

func (p *mtoServiceItemUpdater) findServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
	var serviceItem models.MTOServiceItem
	err := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder",
		"SITDestinationFinalAddress",
		"ReService",
	).Find(&serviceItem, serviceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(serviceItemID, "while looking for service item")
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	return &serviceItem, nil
}

func (p *mtoServiceItemUpdater) approveOrRejectServiceItem(
	appCtx appcontext.AppContext,
	serviceItem models.MTOServiceItem,
	status models.MTOServiceItemStatus,
	rejectionReason *string,
	eTag string,
	checks ...validator,
) (*models.MTOServiceItem, error) {
	if verr := validateServiceItem(appCtx, &serviceItem, eTag, checks...); verr != nil {
		return nil, verr
	}

	var returnedServiceItem models.MTOServiceItem

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		updatedServiceItem, err := p.updateServiceItem(txnAppCtx, serviceItem, status, rejectionReason)
		if err != nil {
			return err
		}
		move := serviceItem.MoveTaskOrder

		if _, err = p.moveRouter.ApproveOrRequestApproval(txnAppCtx, move); err != nil {
			return err
		}

		returnedServiceItem = *updatedServiceItem

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedServiceItem, nil
}

func (p *mtoServiceItemUpdater) updateServiceItem(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem, status models.MTOServiceItemStatus, rejectionReason *string) (*models.MTOServiceItem, error) {
	serviceItem.Status = status
	now := time.Now()

	if status == models.MTOServiceItemStatusRejected {
		if rejectionReason == nil {
			verrs := validate.NewErrors()
			verrs.Add("rejectionReason", "field must be provided when status is set to REJECTED")
			err := apperror.NewInvalidInputError(serviceItem.ID, nil, verrs, "Invalid input found in the request.")
			return nil, err
		}

		serviceItem.RejectionReason = rejectionReason
		serviceItem.RejectedAt = &now
		// clear field if previously accepted
		serviceItem.ApprovedAt = nil
	} else if status == models.MTOServiceItemStatusApproved {
		// clear fields if previously rejected
		serviceItem.RejectionReason = nil
		serviceItem.RejectedAt = nil
		serviceItem.ApprovedAt = &now

		// Check to see if there is already a SIT Destination Original Address
		// by checking for the ID before trying to set one on the service item.
		// If there isn't one, then we set it. We also make sure that the
		// expression looks for the DDDSIT or DDSFSC service codes and only
		// updates the address fields if the service item is of DDDSIT or
		// DDSFSC.
		if (serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
			serviceItem.ReService.Code == models.ReServiceCodeDDSFSC) &&
			serviceItem.SITDestinationOriginalAddressID == nil {

			// Get the shipment destination address
			mtoShipment, err := p.shipmentFetcher.GetShipment(appCtx, *serviceItem.MTOShipmentID, "DestinationAddress")
			if err != nil {
				return nil, err
			}

			// Set the original address on a service item to the shipment's
			// destination address when approving a SIT service item.
			// Creating a new address record to ensure SITDestinationOriginalAddress
			// doesn't change if shipment destination address is updated
			shipmentDestinationAddress := &models.Address{
				StreetAddress1: mtoShipment.DestinationAddress.StreetAddress1,
				StreetAddress2: mtoShipment.DestinationAddress.StreetAddress2,
				StreetAddress3: mtoShipment.DestinationAddress.StreetAddress3,
				City:           mtoShipment.DestinationAddress.City,
				State:          mtoShipment.DestinationAddress.State,
				PostalCode:     mtoShipment.DestinationAddress.PostalCode,
				Country:        mtoShipment.DestinationAddress.Country,
			}
			shipmentDestinationAddress, err = p.addressCreator.CreateAddress(appCtx, shipmentDestinationAddress)
			if err != nil {
				return nil, err
			}
			serviceItem.SITDestinationOriginalAddressID = &shipmentDestinationAddress.ID
			serviceItem.SITDestinationOriginalAddress = shipmentDestinationAddress

			if serviceItem.SITDestinationFinalAddressID == nil {
				serviceItem.SITDestinationFinalAddressID = &shipmentDestinationAddress.ID
				serviceItem.SITDestinationFinalAddress = shipmentDestinationAddress
			}
		}
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&serviceItem)
	if e := handleError(serviceItem.ID, verrs, err); e != nil {
		return nil, e
	}

	return &serviceItem, nil
}

func (p *mtoServiceItemUpdater) convertItemToMembersExpense(
	appCtx appcontext.AppContext,
	serviceItem models.MTOServiceItem,
	eTag string,
	checks ...validator,
) (*models.MTOServiceItem, error) {
	if verr := validateServiceItem(appCtx, &serviceItem, eTag, checks...); verr != nil {
		return nil, verr
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		serviceItem.MembersExpense = models.BoolPointer(true)
		verrs, err := appCtx.DB().ValidateAndUpdate(&serviceItem)
		e := handleError(serviceItem.ID, verrs, err)
		return e
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &serviceItem, nil
}

// UpdateMTOServiceItemBasic updates the MTO Service Item using base validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemBasic(
	appCtx appcontext.AppContext,
	mtoServiceItem *models.MTOServiceItem,
	eTag string,
) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(appCtx, mtoServiceItem, eTag, UpdateMTOServiceItemBasicValidator)
}

// UpdateMTOServiceItemPrime updates the MTO Service Item using Prime API validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemPrime(
	appCtx appcontext.AppContext,
	mtoServiceItem *models.MTOServiceItem,
	eTag string,
) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(appCtx, mtoServiceItem, eTag, UpdateMTOServiceItemPrimeValidator)
}

// UpdateMTOServiceItem updates the given service item
func (p *mtoServiceItemUpdater) UpdateMTOServiceItem(
	appCtx appcontext.AppContext,
	mtoServiceItem *models.MTOServiceItem,
	eTag string,
	validatorKey string,
) (*models.MTOServiceItem, error) {
	// Find the service item, return error if not found
	oldServiceItem, err := models.FetchServiceItem(appCtx.DB(), mtoServiceItem.ID)
	if err != nil {
		switch err {
		case models.ErrFetchNotFound:
			return nil, apperror.NewNotFoundError(mtoServiceItem.ID, "while looking for MTOServiceItem")
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	checker := movetaskorder.NewMoveTaskOrderChecker()
	serviceItemData := updateMTOServiceItemData{
		updatedServiceItem:  *mtoServiceItem,
		oldServiceItem:      oldServiceItem,
		availabilityChecker: checker,
		verrs:               validate.NewErrors(),
	}

	validServiceItem, err := ValidateUpdateMTOServiceItem(appCtx, &serviceItemData, validatorKey)
	if err != nil {
		return nil, err
	}

	// If we have any Customer Contacts we need to make sure that they are associated with
	// all related destination SIT service items. This is especially important if we are creating new Customer Contacts.
	if len(validServiceItem.CustomerContacts) > 0 {
		relatedServiceItems, fetchErr := models.FetchRelatedDestinationSITServiceItems(appCtx.DB(), validServiceItem.ID)
		if fetchErr != nil {
			return nil, fetchErr
		}
		for i := range validServiceItem.CustomerContacts {
			validServiceItem.CustomerContacts[i].MTOServiceItems = relatedServiceItems
		}
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldServiceItem.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(validServiceItem.ID, nil)
	}

	// Create address record (if needed) and update service item in a single transaction
	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if validServiceItem.SITDestinationFinalAddress != nil {
			if validServiceItem.SITDestinationFinalAddressID == nil || *validServiceItem.SITDestinationFinalAddressID == uuid.Nil {
				verrs, createErr := p.builder.CreateOne(txnAppCtx, validServiceItem.SITDestinationFinalAddress)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(
						validServiceItem.ID, createErr, verrs, "Invalid input found while creating a final Destination SIT address for service item.")
				} else if createErr != nil {
					// If the error is something else (this is unexpected), we create a QueryError
					return apperror.NewQueryError("MTOServiceItem", createErr, "")
				}
			}
			validServiceItem.SITDestinationFinalAddressID = &validServiceItem.SITDestinationFinalAddress.ID
		}
		for index := range validServiceItem.CustomerContacts {
			validCustomerContact := &validServiceItem.CustomerContacts[index]
			if validCustomerContact.ID == uuid.Nil {
				verrs, createErr := p.builder.CreateOne(txnAppCtx, validCustomerContact)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(
						validServiceItem.ID, createErr, verrs, "Invalid input found while creating a Customer Contact for service item.")
				} else if createErr != nil {
					// If the error is something else (this is unexpected), we create a QueryError
					return apperror.NewQueryError("MTOServiceItem", createErr, "")
				}
			} else {
				verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(validCustomerContact)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(validServiceItem.ID, updateErr, verrs, "Invalid input found while updating customer contact for the service item.")
				} else if updateErr != nil {
					// If the error is something else (this is unexpected), we create a QueryError
					return apperror.NewQueryError("MTOServiceItem", updateErr, "")
				}
			}
		}

		// Make the update and create a InvalidInputError if there were validation issues
		verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(validServiceItem)

		// If there were validation errors create an InvalidInputError type
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(validServiceItem.ID, updateErr, verrs, "Invalid input found while updating the service item.")
		} else if updateErr != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("MTOServiceItem", updateErr, "")
		}
		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return validServiceItem, nil
}

// ValidateUpdateMTOServiceItem checks the provided serviceItemData struct against the validator indicated by validatorKey.
// Defaults to base validation if the empty string is entered as the key.
// Returns an MTOServiceItem that has been set up for update.
func ValidateUpdateMTOServiceItem(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData, validatorKey string) (*models.MTOServiceItem, error) {
	if validatorKey == "" {
		validatorKey = UpdateMTOServiceItemBasicValidator
	}
	validator, ok := UpdateMTOServiceItemValidators[validatorKey]
	if !ok {
		err := fmt.Errorf("validator key %s was not found in update MTO Service Item validators", validatorKey)
		return nil, err
	}
	err := validator.validate(appCtx, serviceItemData)
	if err != nil {
		return nil, err
	}

	newServiceItem := serviceItemData.setNewMTOServiceItem()

	return newServiceItem, nil
}
