package services

import (
	"time"

	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// PricingType describe the type of pricing to do for a shipment
type PricingType string

type ShipmentDeliverAndPricer interface {
	DeliverAndPriceShipment(deliveryDate time.Time, shipment *models.Shipment) (*validate.Errors, error)
}

type ShipmentPricer interface {
	PriceShipment(shipment *models.Shipment, price PricingType) (*validate.Errors, error)
}
