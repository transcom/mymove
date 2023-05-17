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

type sitAddressUpdateRequestRejector struct {
	checks []sitAddressUpdateValidator
}

// NewSITAddressUpdateRequestRejector creates a new struct with the service dependencies
func NewSITAddressUpdateRequestRejector() services.SITAddressUpdateRequestRejector {
	return &sitAddressUpdateRequestRejector{
		checks: []sitAddressUpdateValidator{
			checkRequiredFields(),
			checkTOORequiredFields(),
		},
	}
}

// RejectSITAddressUpdateRequest rejects the update request
func (f *sitAddressUpdateRequestRejector) RejectSITAddressUpdateRequest(appCtx appcontext.AppContext, serviceItemID uuid.UUID, sitAddressUpdateRequestID uuid.UUID, officeRemarks *string, eTag string) (*models.SITAddressUpdate, error) {
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

	return f.rejectSITAddressUpdateRequest(appCtx, *sitAddressUpdateRequest, officeRemarks)
}

// Find the service item the prime is requesting to update
func (f *sitAddressUpdateRequestRejector) findServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
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

// Find SIT address update request being rejected
func (f *sitAddressUpdateRequestRejector) findSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequestID uuid.UUID) (*models.SITAddressUpdate, error) {
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

func (f *sitAddressUpdateRequestRejector) rejectSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequest models.SITAddressUpdate, officeRemarks *string) (*models.SITAddressUpdate, error) {
	var updatedSITAddressUpdateRequest models.SITAddressUpdate

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		returnedSITAddressUpdateRequest, err := f.updateSITAddressUpdateRequest(txnAppCtx, sitAddressUpdateRequest, officeRemarks)
		if err != nil {
			return err
		}

		updatedSITAddressUpdateRequest = *returnedSITAddressUpdateRequest

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &updatedSITAddressUpdateRequest, nil
}

func (f *sitAddressUpdateRequestRejector) updateSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequest models.SITAddressUpdate, officeRemarks *string) (*models.SITAddressUpdate, error) {
	sitAddressUpdateRequest.OfficeRemarks = officeRemarks
	sitAddressUpdateRequest.Status = models.SITAddressUpdateStatusRejected

	err := validateSITAddressUpdate(appCtx, &sitAddressUpdateRequest, f.checks...)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(sitAddressUpdateRequest)
	if error := f.handleError(sitAddressUpdateRequest.ID, verrs, err); error != nil {
		return nil, error
	}

	return &sitAddressUpdateRequest, nil
}

func (f *sitAddressUpdateRequestRejector) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
