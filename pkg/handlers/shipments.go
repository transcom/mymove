package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForShipmentModel(s models.ShipmentWithOffer) *internalmessages.ShipmentPayload {
	shipmentPayload := &internalmessages.ShipmentPayload{
		ID:                              fmtUUID(s.ID),
		PickupDate:                      fmtDate(time.Now()),
		DeliveryDate:                    fmtDate(time.Now()),
		TrafficDistributionListID:       fmtUUID(s.TrafficDistributionListID),
		TransportationServiceProviderID: fmtUUIDPtr(s.TransportationServiceProviderID),
		AdministrativeShipment:          (s.AdministrativeShipment),
		CreatedAt:                       fmtDateTime(s.CreatedAt),
		UpdatedAt:                       fmtDateTime(s.UpdatedAt),
	}
	return shipmentPayload
}

// IndexShipmentsHandler returns a list of shipments
type IndexShipmentsHandler HandlerContext

// Handle retrieves a list of all shipments
func (h IndexShipmentsHandler) Handle(p shipmentop.IndexShipmentsParams) middleware.Responder {
	var response middleware.Responder

	shipments, err := models.FetchShipments(h.db, false)

	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = shipmentop.NewIndexShipmentsBadRequest()
	} else {
		isp := make(internalmessages.IndexShipmentsPayload, len(shipments))
		for i, s := range shipments {
			isp[i] = payloadForShipmentModel(s)
		}
		response = shipmentop.NewIndexShipmentsOK().WithPayload(isp)
	}
	return response
}
