package mtoshipment

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type sitExtensionApprover struct {
}

// NewSITExtensionApprover creates a new struct with the service dependencies
func NewSITExtensionApprover() services.SITExtension {
	return &sitExtensionApprover{}
}

// ApproveSITExtension approves the SIT Extension and also updates the shipment's SIT days allowance
func (f *sitExtensionApprover) ApproveSITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, approvedDays *int, officeRemarks *string, eTag string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	sitExtension, err := f.findSITExtension(appCtx, sitExtensionID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, services.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.approveSITExtension(appCtx, *shipment, *sitExtension, approvedDays, officeRemarks)
}

func (f *sitExtensionApprover) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (f *sitExtensionApprover) findSITExtension(appCtx appcontext.AppContext, sitExtensionID uuid.UUID) (*models.SITExtension, error) {
	var sitExtension models.SITExtension
	err := appCtx.DB().Q().Find(&sitExtension, sitExtensionID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(sitExtensionID, "while looking for SIT extension")
	} else if err != nil {
		return nil, err
	}

	return &sitExtension, nil
}

func (f *sitExtensionApprover) approveSITExtension(appCtx appcontext.AppContext, shipment models.MTOShipment, sitExtension models.SITExtension, approvedDays *int, officeRemarks *string) (*models.MTOShipment, error) {
	var returnedShipment models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateApprovedDays(txnAppCtx, sitExtension, approvedDays); err != nil {
			return err
		}

		if officeRemarks != nil {
			if err := f.updateOfficeRemarks(txnAppCtx, sitExtension, officeRemarks); err != nil {
				return err
			}
		}

		if err := f.updateSITExtensionStatusToApproved(txnAppCtx, sitExtension); err != nil {
			return err
		}

		if err := f.updateSITExtensionUpdatedAt(txnAppCtx, sitExtension); err != nil {
			return err
		}

		updatedShipment, err := f.updateSitDaysAllowance(txnAppCtx, shipment, approvedDays)
		if err != nil {
			return err
		}

		// TODO: does shipment.UpdatedAt need to be updated here? Separately?
		returnedShipment = *updatedShipment

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedShipment, nil
}

func (f *sitExtensionApprover) updateApprovedDays(appCtx appcontext.AppContext, sitExtension models.SITExtension, approvedDays *int) error {
	sitExtension.ApprovedDays = approvedDays
	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *sitExtensionApprover) updateOfficeRemarks(appCtx appcontext.AppContext, sitExtension models.SITExtension, officeRemarks *string) error {
	sitExtension.OfficeRemarks = officeRemarks
	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *sitExtensionApprover) updateSITExtensionStatusToApproved(appCtx appcontext.AppContext, sitExtension models.SITExtension) error {
	sitExtension.Status = models.SITExtensionStatusApproved
	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *sitExtensionApprover) updateSITExtensionUpdatedAt(appCtx appcontext.AppContext, sitExtension models.SITExtension) error {
	now := time.Now()
	sitExtension.UpdatedAt = now
	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *sitExtensionApprover) updateSitDaysAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment, approvedDays *int) (*models.MTOShipment, error) {
	if shipment.SITDaysAllowance != nil {
		sda := int(*approvedDays) + int(*shipment.SITDaysAllowance)
		shipment.SITDaysAllowance = &sda
	} else {
		shipment.SITDaysAllowance = approvedDays
	}
	verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
	if e := f.handleError(shipment.ID, verrs, err); e != nil {
		return &shipment, e
	}

	return &shipment, nil
}

func (f *sitExtensionApprover) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return services.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
