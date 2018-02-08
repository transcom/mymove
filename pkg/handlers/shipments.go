package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/operations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

type possiblyAwardedShipment struct {
	ID                              uuid.UUID  `db:"id"`
	CreatedAt                       time.Time  `db:"created_at"`
	UpdatedAt                       time.Time  `db:"updated_at"`
	TrafficDistributionListID       uuid.UUID  `db:"traffic_distribution_list_id"`
	TransportationServiceProviderID nulls.UUID `db:"transportation_service_provider_id"`
	AdministrativeShipment          nulls.Bool `db:"administrative_shipment"`
}

func payloadForShipmentModel(s possiblyAwardedShipment) *messages.ShipmentPayload {
	shipmentPayload := &messages.ShipmentPayload{
		ID:           fmtUUID(s.ID),
		PickupDate:   fmtDate(time.Now()),
		DeliveryDate: fmtDate(time.Now()),
		Name:         stringPointer("Shipment name"),
		TrafficDistributionListID:       fmtUUID(s.TrafficDistributionListID),
		TransportationServiceProviderID: fmtNullUUID(s.TransportationServiceProviderID),
		AdministrativeShipment:          fmtNullBool(s.AdministrativeShipment),
		CreatedAt:                       fmtDateTime(s.CreatedAt),
		UpdatedAt:                       fmtDateTime(s.UpdatedAt),
	}
	return shipmentPayload
}

// IndexShipmentsHandler returns a list of all shipments
func IndexShipmentsHandler(p shipmentop.IndexShipmentsParams) middleware.Responder {
	var response middleware.Responder

	shipments := []possiblyAwardedShipment{}

	// TODO Can Q() be .All(&shipments)
	query := dbConnection.Q().LeftOuterJoin("awarded_shipments", "awarded_shipments.shipment_id=shipments.id")

	sql, args := query.ToSQL(&pop.Model{Value: models.Shipment{}},
		"shipments.id",
		"shipments.created_at",
		"shipments.updated_at",
		"shipments.traffic_distribution_list_id",
		"awarded_shipments.transportation_service_provider_id",
		"awarded_shipments.administrative_shipment",
	)

	if err := dbConnection.RawQuery(sql, args...).All(&shipments); err != nil {
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
