package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type shipmentDeliverAndPricer struct {
	db             *pop.Connection
	engine         *rateengine.RateEngine
	planner        route.Planner
	shipmentPricer services.ShipmentPricer
}

// Call delivers a Shipment (and its SITs) and prices associated line items
func (c *shipmentDeliverAndPricer) DeliverAndPriceShipment(deliveryDate time.Time, shipment *models.Shipment) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	err := c.db.Transaction(func(db *pop.Connection) error {
		var transactionError error
		verrs, transactionError = shipment.Deliver(db, deliveryDate)
		if transactionError != nil || verrs.HasAny() {
			return transactionError
		}

		verrs, transactionError = c.shipmentPricer.PriceShipment(shipment, ShipmentPriceNEW)
		if transactionError != nil || verrs.HasAny() {
			return transactionError
		}

		return nil
	})

	return verrs, err
}

func NewShipmentDeliverAndPricer(
	db *pop.Connection,
	engine *rateengine.RateEngine,
	planner route.Planner,
	shipmentPricer services.ShipmentPricer,
) services.ShipmentDeliverAndPricer {
	return &shipmentDeliverAndPricer{
		db:             db,
		engine:         engine,
		planner:        planner,
		shipmentPricer: shipmentPricer}
}
