package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

// ShipmentLineItemFetcher is the service object for fetching shipment line items
type ShipmentLineItemFetcher interface {
	GetShipmentLineItemsByShipmentID(shipmentID uuid.UUID, session *auth.Session) ([]models.ShipmentLineItem, error)
}

// ShipmentLineItemRecalculator is the service object for recalculating shipment line items
type ShipmentLineItemRecalculator interface {
	RecalculateShipmentLineItems(shipmentID uuid.UUID, session *auth.Session, route route.Planner) ([]models.ShipmentLineItem, error)
}
