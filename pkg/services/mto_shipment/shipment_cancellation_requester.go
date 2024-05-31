package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentCancellationRequester struct {
	router     services.ShipmentRouter
	moveRouter services.MoveRouter
}

// NewShipmentCancellationRequester creates a new struct with the service dependencies
func NewShipmentCancellationRequester(router services.ShipmentRouter, moveRouter services.MoveRouter) services.ShipmentCancellationRequester {
	return &shipmentCancellationRequester{
		router,
		moveRouter,
	}
}

// RequestShipmentCancellation Requests the shipment diversion
func (f *shipmentCancellationRequester) RequestShipmentCancellation(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := FindShipment(appCtx, shipmentID, "MoveTaskOrder")
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	today := time.Now()
	// Cancellation Request can only be made before todays date
	if shipment.ActualPickupDate.After(today) || shipment.ActualPickupDate == &today {
		return &models.MTOShipment{}, apperror.NewUpdateError(shipmentID, "cancellation request date cannot be on or after actual pick update")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		// this changes the shipment status to "CANCELLATION_REQUESTED" but only on an approved shipment
		err = f.router.RequestCancellation(appCtx, shipment)
		if err != nil {
			return err
		}

		// save the shipment to the db
		verrs, saveErr := appCtx.DB().ValidateAndSave(shipment)
		if verrs != nil && verrs.HasAny() {
			invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "Could not validate shipment while requesting the shipment cancellation.")

			return invalidInputError
		}
		if saveErr != nil {
			return err
		}

		// checking if the move still requires action by the TOO
		// if no action is needed, then the move status will stay in APPROVED or APPROVALS_REQUESTED
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
		return nil, transactionError
	}

	return shipment, err
}
