package move

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveCanceler struct{}

func NewMoveCanceler() services.MoveCanceler {
	return &moveCanceler{}
}

func (f *moveCanceler) CancelMove(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Move, error) {
	move := &models.Move{}
	err := appCtx.DB().Find(move, moveID)
	if err != nil {
		return nil, apperror.NewNotFoundError(moveID, "while looking for a move")
	}

	moveDelta := move
	moveDelta.Status = models.MoveStatusCANCELED

	// get all shipments in move for cancellation
	var shipments []models.MTOShipment
	err = appCtx.DB().EagerPreload("Status", "PPMShipment", "PPMShipment.Status").Where("mto_shipments.move_id = $1", move.ID).All(&shipments)
	if err != nil {
		return nil, apperror.NewNotFoundError(moveID, "while looking for shipments")
	}

	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		for _, shipment := range shipments {
			shipmentDelta := shipment
			shipmentDelta.Status = models.MTOShipmentStatusCanceled

			if shipment.PPMShipment != nil {
				if shipment.PPMShipment.Status == models.PPMShipmentStatusCloseoutComplete {
					return apperror.NewConflictError(move.ID, " cannot cancel move with approved shipment.")
				}
				var ppmshipment models.PPMShipment
				qerr := appCtx.DB().Where("id = ?", shipment.PPMShipment.ID).First(&ppmshipment)
				if qerr != nil {
					return apperror.NewNotFoundError(ppmshipment.ID, "while looking for ppm shipment")
				}

				ppmshipment.Status = models.PPMShipmentStatusCanceled

				verrs, err := txnAppCtx.DB().ValidateAndUpdate(&ppmshipment)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(shipment.ID, err, verrs, "Validation errors found while setting shipment status")
				} else if err != nil {
					return apperror.NewQueryError("PPM Shipment", err, "Failed to update status for ppm shipment")
				}
			}

			if shipment.Status != models.MTOShipmentStatusApproved {
				verrs, err := txnAppCtx.DB().ValidateAndUpdate(&shipmentDelta)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(shipment.ID, err, verrs, "Validation errors found while setting shipment status")
				} else if err != nil {
					return apperror.NewQueryError("Shipment", err, "Failed to update status for shipment")
				}
			} else {
				return apperror.NewConflictError(move.ID, " cannot cancel move with approved shipment.")
			}
		}

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(moveDelta)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(
				move.ID, err, verrs, "Validation errors found while setting move status")
		} else if err != nil {
			return apperror.NewQueryError("Move", err, "Failed to update status for move")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return move, nil
}
