package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type DeliverAndPriceShipment struct {
	DB     *pop.Connection
	Engine *rateengine.RateEngine
}

// Call delivers a Shipment and prices associated line items
func (c DeliverAndPriceShipment) Call(deliveryDate time.Time, shipment *models.Shipment) (*validate.Errors, error) {
	err := shipment.Deliver(deliveryDate)
	if err != nil {
		return validate.NewErrors(), err
	}

	return PriceShipment{DB: c.DB, Engine: c.Engine}.Call(shipment, ShipmentPriceNEW)
}
