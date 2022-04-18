package mtoshipment

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/db/utilities"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type sitExtensionApprover struct {
	moveRouter services.MoveRouter
}

// NewSITExtensionApprover creates a new struct with the service dependencies
func NewSITExtensionApprover(moveRouter services.MoveRouter) services.SITExtensionApprover {
	return &sitExtensionApprover{moveRouter}
}

// ApproveSITExtension approves the SIT Extension and also updates the shipment's SIT days allowance
func (f *sitExtensionApprover) ApproveSITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, approvedDays int, officeRemarks *string, eTag string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	sitExtension, err := f.findSITExtension(appCtx, sitExtensionID)
	if err != nil {
		return nil, err
	}

	if sitExtension.MTOShipmentID != shipment.ID {
		return nil, apperror.NewNotFoundError(shipmentID, "while looking for SITExtension's shipment ID")
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.approveSITExtension(appCtx, *shipment, *sitExtension, approvedDays, officeRemarks)
}

func (f *sitExtensionApprover) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope()).EagerPreload("MoveTaskOrder").Find(&shipment, shipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipmentID, "while looking for shipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}

func (f *sitExtensionApprover) findSITExtension(appCtx appcontext.AppContext, sitExtensionID uuid.UUID) (*models.SITExtension, error) {
	var sitExtension models.SITExtension
	err := appCtx.DB().Q().Find(&sitExtension, sitExtensionID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(sitExtensionID, "while looking for SIT extension")
		default:
			return nil, apperror.NewQueryError("SITExtension", err, "")
		}
	}

	return &sitExtension, nil
}

func (f *sitExtensionApprover) approveSITExtension(appCtx appcontext.AppContext, shipment models.MTOShipment, sitExtension models.SITExtension, approvedDays int, officeRemarks *string) (*models.MTOShipment, error) {
	var returnedShipment models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateSITExtension(txnAppCtx, sitExtension, approvedDays, officeRemarks); err != nil {
			return err
		}

		updatedShipment, err := f.updateSitDaysAllowance(txnAppCtx, shipment, approvedDays)
		if err != nil {
			return err
		}

		if _, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, shipment.MoveTaskOrder); err != nil {
			return err
		}

		returnedShipment = *updatedShipment

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedShipment, nil
}

func (f *sitExtensionApprover) updateSITExtension(appCtx appcontext.AppContext, sitExtension models.SITExtension, approvedDays int, officeRemarks *string) error {
	sitExtension.ApprovedDays = &approvedDays
	sitExtension.OfficeRemarks = officeRemarks
	sitExtension.Status = models.SITExtensionStatusApproved
	now := time.Now()
	sitExtension.DecisionDate = &now

	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *sitExtensionApprover) updateSitDaysAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment, approvedDays int) (*models.MTOShipment, error) {
	if shipment.SITDaysAllowance != nil {
		sda := approvedDays + int(*shipment.SITDaysAllowance)
		shipment.SITDaysAllowance = &sda
	} else {
		shipment.SITDaysAllowance = &approvedDays
	}
	verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
	if e := f.handleError(shipment.ID, verrs, err); e != nil {
		return &shipment, e
	}

	err = appCtx.DB().Q().EagerPreload("SITExtensions").Find(&shipment, shipment.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipment.ID, "looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}

func (f *sitExtensionApprover) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
