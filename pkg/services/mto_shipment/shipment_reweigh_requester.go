package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentReweighRequester struct {
	Sender notifications.NotificationSender
}

// NewShipmentReweighRequester creates a new struct with the service dependencies
func NewShipmentReweighRequester(sender notifications.NotificationSender) services.ShipmentReweighRequester {
	return &shipmentReweighRequester{Sender: sender}
}

// RequestShipmentReweigh Requests the shipment reweigh
func (f *shipmentReweighRequester) RequestShipmentReweigh(appCtx appcontext.AppContext, shipmentID uuid.UUID, requester models.ReweighRequester) (*models.Reweigh, error) {
	shipment, err := FindShipment(appCtx, shipmentID, "Reweigh")
	if err != nil {
		return nil, err
	}

	if shipment.ShipmentType == models.MTOShipmentTypePPM {
		return nil, apperror.NewBadDataError("PPM Shipments are not eligible for reweigh requests")
	}

	reweigh, err := f.createReweigh(appCtx, shipment, requester, checkReweighAllowed())

	/* Don't send emails to BLUEBARK moves */
	if err == nil && shipment.MoveTaskOrder.Orders.CanSendEmailWithOrdersType() && shipment.CanSendReweighEmailForShipmentType() {
		err := f.Sender.SendNotification(appCtx,
			notifications.NewReweighRequested(shipment.MoveTaskOrderID, *shipment),
		)
		if err != nil {
			return reweigh, err
		}
	}

	return reweigh, err
}

func (f *shipmentReweighRequester) createReweigh(appCtx appcontext.AppContext, shipment *models.MTOShipment, requester models.ReweighRequester, checks ...validator) (*models.Reweigh, error) {
	if verr := validateShipment(appCtx, shipment, shipment, checks...); verr != nil {
		return nil, verr
	}

	now := time.Now()
	reweigh := models.Reweigh{
		RequestedBy: requester,
		RequestedAt: now,
		ShipmentID:  shipment.ID,
	}

	verrs, dbErr := appCtx.DB().ValidateAndSave(&reweigh)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "Could not save the reweigh while requesting the reweigh as a TOO.")

		return nil, invalidInputError
	}
	if dbErr != nil {
		return nil, dbErr
	}

	return &reweigh, nil
}
