package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentDiversionApprover struct {
	router     services.ShipmentRouter
	moveRouter services.MoveRouter
}

// NewShipmentDiversionApprover creates a new struct with the service dependencies
func NewShipmentDiversionApprover(router services.ShipmentRouter, moveRouter services.MoveRouter) services.ShipmentDiversionApprover {
	return &shipmentDiversionApprover{
		router, moveRouter,
	}
}

// ApproveShipmentDiversion Approves the shipment diversion
func (f *shipmentDiversionApprover) ApproveShipmentDiversion(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := FindShipment(appCtx, shipmentID, "MoveTaskOrder")
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = f.router.ApproveDiversion(appCtx, shipment)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	verrs, err := appCtx.DB().ValidateAndSave(shipment)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "Could not validate shipment while approving the diversion.")

		return nil, invalidInputError
	}
	if err != nil {
		return nil, err
	}

	transactionError = appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		move := &models.Move{}
		err := txnAppCtx.DB().Find(move, shipment.MoveTaskOrderID)
		if err != nil {
			return apperror.NewQueryError("Move", err, "")
		}

		if move.Status == models.MoveStatusAPPROVALSREQUESTED || move.Status == models.MoveStatusAPPROVED {
			if _, err := f.moveRouter.ApproveOrRequestApproval(txnAppCtx, *move); err != nil {
				return err
			}
		}

		return nil
	})
	if transactionError != nil {
		return nil, transactionError
	}

	return shipment, err
}
