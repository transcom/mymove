package handlers

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/gen/messages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/operations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForShipmentModel(shipment models.Shipment) messages.ShipmentPayload {
	shipmentPayload := messages.ShipmentPayload{
		ID:           fmtUUID(shipment.ID),
		PickupDate:   fmtDate(shipment.PickupDate),
		DeliveryDate: fmtDate(shipment.DeliveryDate),
		Name:         &shipment.Name,
		TrafficDistributionListID: fmtUUID(shipment.TrafficDistributionListID),
		CreatedAt:                 fmtDateTime(shipment.CreatedAt),
		UpdatedAt:                 fmtDateTime(shipment.UpdatedAt),
	}
	return shipmentPayload
}

// IndexShipmentsHandler returns a list of all shipments
func IndexShipmentsHandler(params shipmentop.IndexShipmentsParams) middleware.Responder {
	var response middleware.Responder

	shipmentPayloads := make(messages.IndexShipmentsPayload, 3)
	for i := range shipmentPayloads {
		shipment := models.Shipment{
			ID:                        uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440000"),
			CreatedAt:                 time.Now(),
			UpdatedAt:                 time.Now(),
			Name:                      fmt.Sprintf("Shipment number %d", i+1),
			PickupDate:                time.Now(),
			DeliveryDate:              time.Now(),
			TrafficDistributionListID: uuid.FromStringOrNil("123e4567-e89b-12d3-a456-426655440000"),
		}
		shipmentPayload := payloadForShipmentModel(shipment)
		shipmentPayloads[i] = &shipmentPayload
	}
	response = shipmentop.NewIndexShipmentsOK().WithPayload(shipmentPayloads)
	return response
}
