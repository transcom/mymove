package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentTermination struct {
}

// NewShipmentTermination creates a new shipmentTermination service
func NewShipmentTermination() services.ShipmentTermination {
	return &shipmentTermination{}
}

// TerminateShipment terminates a shipment
// updates the shipment status to TERMINATED_FOR_CAUSE
// and updates terminated_at with the current timestamp
// updates termination_comments with provided required comments with a TERMINATED FOR CAUSE prefix
func (f *shipmentTermination) TerminateShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, terminationComments string) (*models.MTOShipment, error) {
	shipment, err := FindShipment(appCtx, shipmentID, "MoveTaskOrder", "PPMShipment")
	if shipment == nil {
		notFoundError := apperror.NewNotFoundError(shipmentID, "shipment not found when trying to terminate for cause")
		appCtx.Logger().Error(notFoundError.Error())
		return nil, notFoundError
	}
	if err != nil {
		return nil, err
	}

	if shipment.Status == models.MTOShipmentStatusTerminatedForCause {
		return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipment in TERMINATED FOR CAUSE status cannot be terminated for cause again")
	}

	if shipment.ActualPickupDate != nil {
		return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipment cannot have an actual pickup date set in order to terminate for cause")
	}

	if shipment.Status != models.MTOShipmentStatusApproved {
		return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipment must be in APPROVED status in order to terminate for cause")
	}

	if shipment.PPMShipment != nil && shipment.PPMShipment.ID != uuid.Nil {
		// This shipment is tied to a PPM, it shouldn't be possible to terminate
		return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipments tied to PPMs do not qualify for termination")
	}

	terminatedAt := time.Now()
	shipment.TerminatedAt = &terminatedAt
	comments := "TERMINATED FOR CAUSE - " + terminationComments
	shipment.TerminationComments = &comments
	shipment.Status = models.MTOShipmentStatusTerminatedForCause

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := appCtx.DB().ValidateAndUpdate(shipment)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(shipment.ID, err, verrs, "Invalid input found while terminating the shipment")
		} else if err != nil {
			return apperror.NewQueryError("MTOShipments", err, "")
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return shipment, err
}
