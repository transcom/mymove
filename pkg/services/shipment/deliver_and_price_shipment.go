package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type DeliverAndPriceShipment struct {
	DB            *pop.Connection
	Engine        *rateengine.RateEngine
	Planner       route.Planner
	PriceShipment PriceShipment
}

// Call delivers a Shipment (and its SITs) and prices associated line items
func (c DeliverAndPriceShipment) Call(deliveryDate time.Time, shipment *models.Shipment) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	err := c.DB.Transaction(func(db *pop.Connection) error {
		var transactionError error
		verrs, transactionError = shipment.Deliver(db, deliveryDate)
		if transactionError != nil || verrs.HasAny() {
			return transactionError
		}

		verrs, transactionError = c.PriceShipment.Call(shipment, ShipmentPriceNEW)
		if transactionError != nil || verrs.HasAny() {
			return transactionError
		}

		return nil
	})

	return verrs, err
}
