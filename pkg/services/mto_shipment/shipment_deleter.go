package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentDeleter struct {
	checks []validator
}

// NewShipmentDeleter creates a new struct with the service dependencies
func NewShipmentDeleter() services.ShipmentDeleter {
	return &shipmentDeleter{[]validator{checkDeleteAllowed()}}
}

// NewPrimeShipmentDeleter creates a new struct with the service dependencies
func NewPrimeShipmentDeleter() services.ShipmentDeleter {
	return &shipmentDeleter{[]validator{checkPrimeDeleteAllowed()}}
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
		return nil
	})

	if transactionError != nil {
		return uuid.Nil, transactionError
	}

	return shipment.MoveTaskOrderID, err
}
