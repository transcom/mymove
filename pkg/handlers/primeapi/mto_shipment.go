package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// ConvertToMTOShipment passes the fields available for update to the MTOShipment payload object so that it can be
// manipulated later on
func ConvertToMTOShipment(updates *primemessages.PutMTOShipment) *primemessages.MTOShipment {
	shipment := &primemessages.MTOShipment{
		ScheduledPickupDate:        updates.ScheduledPickupDate,
		FirstAvailableDeliveryDate: updates.FirstAvailableDeliveryDate,
		PrimeActualWeight:          updates.PrimeActualWeight,
		PrimeEstimatedWeight:       updates.PrimeEstimatedWeight,
		ActualPickupDate:           updates.ActualPickupDate,
		RequiredDeliveryDate:       updates.RequiredDeliveryDate,
		Agents:                     updates.Agents,
		ShipmentType:               updates.ShipmentType,
		PickupAddress:              updates.PickupAddress,
		DestinationAddress:         updates.DestinationAddress,
		SecondaryPickupAddress:     updates.SecondaryPickupAddress,
		SecondaryDeliveryAddress:   updates.SecondaryDeliveryAddress,
	}

	return shipment
}

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	shipmentPayload := ConvertToMTOShipment(params.Body)
	mtoShipment := payloads.MTOShipmentModel(shipmentPayload)
	eTag := params.IfMatch
	logger.Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	mtoShipment, err := h.mtoShipmentUpdater.UpdateMTOShipment(mtoShipment, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError()
		}
	}
	mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload)
}
