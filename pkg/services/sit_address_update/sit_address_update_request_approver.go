package sitaddressupdate

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type sitAddressUpdateRequestApprover struct {
	serviceItemUpdater services.MTOServiceItemUpdater
	checks             []sitAddressUpdateValidator
}

// NewSITAddressUpdateRequestApprover creates a new struct with the service dependencies
func NewSITAddressUpdateRequestApprover(serviceItemUpdater services.MTOServiceItemUpdater) services.SITAddressUpdateRequestApprover {
	return &sitAddressUpdateRequestApprover{
		serviceItemUpdater: serviceItemUpdater,
		checks: []sitAddressUpdateValidator{
			checkAndValidateRequiredFields(),
			checkTOORequiredFields(),
		},
	}
}

// ApproveSITAddressUpdateRequest approves the update request and updates the service item's final address
func (f *sitAddressUpdateRequestApprover) ApproveSITAddressUpdateRequest(appCtx appcontext.AppContext, serviceItemID uuid.UUID, sitAddressUpdateRequestID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOServiceItem, error) {
	serviceItem, err := f.findServiceItem(appCtx, serviceItemID)
	if err != nil {
		return nil, err
	}

	sitAddressUpdateRequest, err := f.findSITAddressUpdateRequest(appCtx, sitAddressUpdateRequestID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(serviceItem.UpdatedAt)
	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(serviceItemID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.approveSITAddressUpdateRequest(appCtx, serviceItem, *sitAddressUpdateRequest, officeRemarks)
}

// Find the service item the prime is requesting to update
func (f *sitAddressUpdateRequestApprover) findServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
	var serviceItem models.MTOServiceItem
	err := appCtx.DB().Eager("SITDestinationFinalAddress").Where("id = ?", serviceItemID).Where("status = ?", models.MTOServiceItemStatusApproved).First(&serviceItem)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(serviceItemID, "while looking for SIT service item")
		default:
			return nil, apperror.NewQueryError("SITServiceItem", err, "unable to retrieve SIT service item.")
		}
	}

	return &serviceItem, nil
}

// Find SIT address update request being approved
func (f *sitAddressUpdateRequestApprover) findSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequestID uuid.UUID) (*models.SITAddressUpdate, error) {
	var SITAddressUpdateRequest models.SITAddressUpdate
	err := appCtx.DB().Q().Find(&SITAddressUpdateRequest, sitAddressUpdateRequestID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(sitAddressUpdateRequestID, "while looking for SIT address update request")
		default:
			return nil, apperror.NewQueryError("SITAddressUpdate", err, "Unable to create SIT address update request.")
		}
	}

	return &SITAddressUpdateRequest, nil
}

// Final approval for SIT address update request
func (f *sitAddressUpdateRequestApprover) approveSITAddressUpdateRequest(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, sitAddressUpdateRequest models.SITAddressUpdate, officeRemarks *string) (*models.MTOServiceItem, error) {
	var returnedServiceItem models.MTOServiceItem

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateSITAddressUpdateRequest(txnAppCtx, sitAddressUpdateRequest, officeRemarks); err != nil {
			return err
		}

		updatedServiceItem, err := f.updateServiceItemFinalAddress(txnAppCtx, *serviceItem, &sitAddressUpdateRequest.NewAddressID, sitAddressUpdateRequest.NewAddress)
		if err != nil {
			return err
		}

		returnedServiceItem = *updatedServiceItem

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &returnedServiceItem, nil
}

func (f *sitAddressUpdateRequestApprover) updateSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequest models.SITAddressUpdate, officeRemarks *string) error {
	sitAddressUpdateRequest.OfficeRemarks = officeRemarks
	sitAddressUpdateRequest.Status = models.SITAddressUpdateStatusApproved

	err := validateSITAddressUpdate(appCtx, &sitAddressUpdateRequest, f.checks...)
	if err != nil {
		return err
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(sitAddressUpdateRequest)
	if error := f.handleError(sitAddressUpdateRequest.ID, verrs, err); error != nil {
		return error
	}

	return nil
}

func (f *sitAddressUpdateRequestApprover) updateServiceItemFinalAddress(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem, newAddressID *uuid.UUID, newAddress models.Address) (*models.MTOServiceItem, error) {
	serviceItem.SITDestinationFinalAddressID = newAddressID
	serviceItem.SITDestinationFinalAddress = &newAddress

	updatedServiceItem, err := f.serviceItemUpdater.UpdateMTOServiceItemPrime(appCtx, &serviceItem, etag.GenerateEtag(serviceItem.UpdatedAt))
	if err != nil {
		return nil, err
	}

	return updatedServiceItem, nil
}

func (f *sitAddressUpdateRequestApprover) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
