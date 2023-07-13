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
	moveRouter         services.MoveRouter
}

// NewSITAddressUpdateRequestApprover creates a new struct with the service dependencies
func NewSITAddressUpdateRequestApprover(serviceItemUpdater services.MTOServiceItemUpdater, moveRouter services.MoveRouter) services.SITAddressUpdateRequestApprover {
	return &sitAddressUpdateRequestApprover{
		serviceItemUpdater: serviceItemUpdater,
		checks: []sitAddressUpdateValidator{
			checkAndValidateRequiredFields(),
			checkTOORequiredFields(),
		},
		moveRouter: moveRouter,
	}
}

// ApproveSITAddressUpdateRequest approves the update request and updates the service item's final address
func (f *sitAddressUpdateRequestApprover) ApproveSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequestID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOServiceItem, error) {
	sitAddressUpdateRequest, err := f.findSITAddressUpdateRequest(appCtx, sitAddressUpdateRequestID)
	if err != nil {
		return nil, err
	}

	serviceItem, err := f.findServiceItem(appCtx, sitAddressUpdateRequest.MTOServiceItemID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(sitAddressUpdateRequest.UpdatedAt)
	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(sitAddressUpdateRequest.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.approveSITAddressUpdateRequest(appCtx, *serviceItem, *sitAddressUpdateRequest, officeRemarks)
}

// Find the service item the prime is requesting to update
func (f *sitAddressUpdateRequestApprover) findServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
	serviceItem, err := models.FetchServiceItem(appCtx.DB(), serviceItemID)

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
	SITAddressUpdateRequest, err := models.FetchSITAddressUpdate(appCtx.DB(), sitAddressUpdateRequestID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(sitAddressUpdateRequestID, "while looking for SIT address update request")
		default:
			return nil, apperror.NewQueryError("SITAddressUpdate", err, "unable to approve SIT address update request.")
		}
	}

	return &SITAddressUpdateRequest, nil
}

// Final approval for SIT address update request
func (f *sitAddressUpdateRequestApprover) approveSITAddressUpdateRequest(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem, sitAddressUpdateRequest models.SITAddressUpdate, officeRemarks *string) (*models.MTOServiceItem, error) {
	var returnedServiceItem models.MTOServiceItem

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		//Grabbing the associated move to update its status
		var move models.Move
		err := txnAppCtx.DB().Find(&move, serviceItem.MoveTaskOrderID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(serviceItem.MoveTaskOrderID, "looking for Move")
			default:
				return apperror.NewQueryError("Move", err, "unable to retrieve move")
			}
		}

		// Updating the status of the request as well as office remarks
		err = f.updateSITAddressUpdateRequest(txnAppCtx, sitAddressUpdateRequest, officeRemarks)
		if err != nil {
			return err
		}

		//Update the final address on the service item
		updatedServiceItem, err := f.updateServiceItemFinalAddress(txnAppCtx, serviceItem, sitAddressUpdateRequest)
		if err != nil {
			return err
		}

		// Clear APPROVALS_REQUESTED status on move
		_, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, move)
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

	verrs, err := appCtx.DB().ValidateAndUpdate(&sitAddressUpdateRequest)
	return f.handleError(sitAddressUpdateRequest.ID, verrs, err)
}

func (f *sitAddressUpdateRequestApprover) updateServiceItemFinalAddress(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem, sitAddressUpdateRequest models.SITAddressUpdate) (*models.MTOServiceItem, error) {
	serviceItem.SITDestinationFinalAddressID = &sitAddressUpdateRequest.NewAddressID
	serviceItem.SITDestinationFinalAddress = &sitAddressUpdateRequest.NewAddress

	updatedServiceItem, err := f.serviceItemUpdater.UpdateMTOServiceItemBasic(appCtx, &serviceItem, etag.GenerateEtag(serviceItem.UpdatedAt))
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
