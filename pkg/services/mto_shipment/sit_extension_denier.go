package mtoshipment

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type sitExtensionDenier struct {
	moveRouter services.MoveRouter
}

// NewSITExtensionDenier creates a new struct with the service dependencies
func NewSITExtensionDenier(moveRouter services.MoveRouter) services.SITExtensionDenier {
	return &sitExtensionDenier{moveRouter}
}

// DenySITExtension denies the SIT Extension
func (f *sitExtensionDenier) DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOShipment, error) {
	shipment, err := FindShipment(appCtx, shipmentID, "MoveTaskOrder")
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

	// var updatedShipment models.MTOShipment
	// err = appCtx.DB().Q().Find(&updatedShipment, shipmentID)
	// return &updatedShipment, err

	return f.denySITExtension(appCtx, *shipment, *sitExtension, officeRemarks)
}

func (f *sitExtensionDenier) findSITExtension(appCtx appcontext.AppContext, sitExtensionID uuid.UUID) (*models.SITExtension, error) {
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

func (f *sitExtensionDenier) denySITExtension(appCtx appcontext.AppContext, shipment models.MTOShipment, sitExtension models.SITExtension, officeRemarks *string) (*models.MTOShipment, error) {
	var returnedShipment models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateSITExtension(txnAppCtx, sitExtension, officeRemarks); err != nil {
			return err
		}

		if _, err := f.moveRouter.ApproveOrRequestApproval(txnAppCtx, shipment.MoveTaskOrder); err != nil {
			return err
		}

		if e := txnAppCtx.DB().Q().EagerPreload("SITExtensions").Find(&returnedShipment, shipment.ID); e != nil {
			switch e {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(shipment.ID, "looking for MTOShipment")
			default:
				return apperror.NewQueryError("MTOShipment", e, "")
			}
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedShipment, nil
}

func (f *sitExtensionDenier) updateSITExtension(appCtx appcontext.AppContext, sitExtension models.SITExtension, officeRemarks *string) error {
	if officeRemarks != nil {
		sitExtension.OfficeRemarks = officeRemarks
	}
	sitExtension.Status = models.SITExtensionStatusDenied
	now := time.Now()
	sitExtension.DecisionDate = &now

	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *sitExtensionDenier) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
