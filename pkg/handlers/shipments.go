package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/operations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForShipmentModel(s models.PossiblyAwardedShipment) *messages.ShipmentPayload {
	shipmentPayload := &messages.ShipmentPayload{
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

// IndexShipmentsHandler returns a list of all shipments
func IndexShipmentsHandler(p shipmentop.IndexShipmentsParams) middleware.Responder {
	var response middleware.Responder

	shipments, err := models.FetchPossiblyAwardedShipments(dbConnection)

	if err != nil {
		zap.L().Error("DB Query", zap.Error(err))
		response = shipmentop.NewIndexShipmentsBadRequest()
	} else {
		isp := make(messages.IndexShipmentsPayload, len(shipments))
		for i, s := range shipments {
			isp[i] = payloadForShipmentModel(s)
		}
		response = shipmentop.NewIndexShipmentsOK().WithPayload(isp)
	}
	return response
}
