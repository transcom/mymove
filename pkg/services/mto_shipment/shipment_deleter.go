package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentDeleter struct {
}

// NewShipmentDeleter creates a new struct with the service dependencies
func NewShipmentDeleter() services.ShipmentDeleter {
	return &shipmentDeleter{}
}

// DeleteShipment soft deletes the shipment
func (f *shipmentDeleter) DeleteShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (uuid.UUID, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return uuid.Nil, err
	}

	err = f.verifyShipmentCanBeDeleted(appCtx, shipment)
	if err != nil {
		return uuid.Nil, err
	}

	now := time.Now()
	shipment.DeletedAt = &now
	err = appCtx.DB().Save(shipment)

	return shipment.MoveTaskOrderID, err
}

func (f *shipmentDeleter) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Eager("MoveTaskOrder").Where("mto_shipments.deleted_at IS NULL").Find(&shipment, shipmentID)

	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
		}
	}

	return &shipment, nil
}

func (f *shipmentDeleter) verifyShipmentCanBeDeleted(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	move := shipment.MoveTaskOrder
	if move.Status != models.MoveStatusDRAFT && move.Status != models.MoveStatusNeedsServiceCounseling {
		return services.NewForbiddenError("A shipment can only be deleted if the move is in Draft or NeedsServiceCounseling")
	}

	return nil
}
