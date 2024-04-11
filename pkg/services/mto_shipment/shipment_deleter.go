package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentDeleter struct {
	checks               []validator
	moveTaskOrderUpdater services.MoveTaskOrderUpdater
	moveRouter           services.MoveRouter
}

// NewShipmentDeleter creates a new struct with the service dependencies
func NewShipmentDeleter(moveTaskOrderUpdater services.MoveTaskOrderUpdater, moveRouter services.MoveRouter) services.ShipmentDeleter {
	return &shipmentDeleter{
		checks:               []validator{checkDeleteAllowed()},
		moveTaskOrderUpdater: moveTaskOrderUpdater,
		moveRouter:           moveRouter,
	}
}

// NewPrimeShipmentDeleter creates a new struct with the service dependencies
func NewPrimeShipmentDeleter(moveTaskOrderUpdater services.MoveTaskOrderUpdater) services.ShipmentDeleter {
	return &shipmentDeleter{
		checks:               []validator{checkPrimeDeleteAllowed()},
		moveTaskOrderUpdater: moveTaskOrderUpdater,
	}
}

// DeleteShipment soft deletes the shipment
func (f *shipmentDeleter) DeleteShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (uuid.UUID, error) {
	shipment, err := FindShipment(appCtx, shipmentID, "MoveTaskOrder", "PPMShipment")
	if err != nil {
		return uuid.Nil, err
	}

	// run the (read-only) validations
	if verr := validateShipment(appCtx, shipment, shipment, f.checks...); verr != nil {
		return uuid.Nil, verr
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = utilities.SoftDestroy(appCtx.DB(), shipment)
		if err != nil {
			switch err.Error() {
			case "error updating model":
				return apperror.NewUnprocessableEntityError("while updating model")
			case "this model does not have deleted_at field":
				return apperror.NewPreconditionFailedError(shipmentID, errors.New("model or sub table missing deleted_at field"))
			default:
				return apperror.NewInternalServerError("failed attempt to soft delete model")
			}
		}
		// Update PPMType once shipment gets created.
		_, err = f.moveTaskOrderUpdater.UpdatePPMType(txnAppCtx, shipment.MoveTaskOrderID)
		if err != nil {
			return err
		}

		// if the shipment had any actions for the TOO we can remove these by checking if the move status should change
		move := shipment.MoveTaskOrder
		if move.Status == models.MoveStatusAPPROVALSREQUESTED || move.Status == models.MoveStatusAPPROVED {
			_, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, shipment.MoveTaskOrder)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if transactionError != nil {
		return uuid.Nil, transactionError
	}

	return shipment.MoveTaskOrderID, err
}
