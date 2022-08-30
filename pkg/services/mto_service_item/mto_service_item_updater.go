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
}

// NewMTOServiceItemUpdater returns a new mto service item updater
func NewMTOServiceItemUpdater(builder mtoServiceItemQueryBuilder, moveRouter services.MoveRouter) services.MTOServiceItemUpdater {
	// used inside a transaction and mocking		return &mtoServiceItemUpdater{builder: builder}
	createNewBuilder := func() mtoServiceItemQueryBuilder {
		return query.NewQueryBuilder()
	}

	return &mtoServiceItemUpdater{builder, createNewBuilder, moveRouter}
}

func (p *mtoServiceItemUpdater) ApproveOrRejectServiceItem(appCtx appcontext.AppContext, mtoServiceItemID uuid.UUID, status models.MTOServiceItemStatus, rejectionReason *string, eTag string) (*models.MTOServiceItem, error) {
	mtoServiceItem, err := p.findServiceItem(appCtx, mtoServiceItemID)
	if err != nil {
		return &models.MTOServiceItem{}, err
	}

	return p.approveOrRejectServiceItem(appCtx, *mtoServiceItem, status, rejectionReason, eTag, checkMoveStatus(), checkETag())
}

func (p *mtoServiceItemUpdater) findServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
	var serviceItem models.MTOServiceItem
	err := appCtx.DB().Q().EagerPreload("MoveTaskOrder").Find(&serviceItem, serviceItemID)
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

func (p *mtoServiceItemUpdater) approveOrRejectServiceItem(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem, status models.MTOServiceItemStatus, rejectionReason *string, eTag string, checks ...validator) (*models.MTOServiceItem, error) {
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
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&serviceItem)
	if e := handleError(serviceItem.ID, verrs, err); e != nil {
		return nil, e
	}

	return &serviceItem, nil
}

// UpdateMTOServiceItemBasic updates the MTO Service Item using base validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemBasic(appCtx appcontext.AppContext, mtoServiceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(appCtx, mtoServiceItem, eTag, UpdateMTOServiceItemBasicValidator)
}

// UpdateMTOServiceItemPrime updates the MTO Service Item using Prime API validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemPrime(appCtx appcontext.AppContext, mtoServiceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(appCtx, mtoServiceItem, eTag, UpdateMTOServiceItemPrimeValidator)
}

// UpdateMTOServiceItem updates the given service item
func (p *mtoServiceItemUpdater) UpdateMTOServiceItem(appCtx appcontext.AppContext, mtoServiceItem *models.MTOServiceItem, eTag string, validatorKey string) (*models.MTOServiceItem, error) {
	oldServiceItem := models.MTOServiceItem{}

	// Find the service item, return error if not found
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoServiceItem.ID),
	}
	err := p.builder.FetchOne(appCtx, &oldServiceItem, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
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
				validServiceItem.SITDestinationFinalAddressID = &validServiceItem.SITDestinationFinalAddress.ID
			} else {
				// If this service item already had a SITDestinationFinalAddress, update that record instead
				// of creating a new one.
				verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(validServiceItem.SITDestinationFinalAddress)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(validServiceItem.ID, updateErr, verrs, "Invalid input found while updating final Destination SIT address for the service item.")
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

	// Get the updated address and return
	updatedServiceItem := models.MTOServiceItem{}
	err = p.builder.FetchOne(appCtx, &updatedServiceItem, queryFilters) // using the same queryFilters set at the beginning
	if err != nil {
		return nil, apperror.NewQueryError("MTOServiceItem", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedServiceItem, nil
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
