package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type shipmentDeliverAndPricer struct {
	db      *pop.Connection
	engine  *rateengine.RateEngine
	planner route.Planner
}

// Call delivers a Shipment (and its SITs) and prices associated line items
func (c *shipmentDeliverAndPricer) DeliverAndPriceShipment(deliveryDate time.Time, shipment *models.Shipment) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	err := c.db.Transaction(func(db *pop.Connection) error {
		var transactionError error
		transactionError = shipment.Deliver(deliveryDate)
		if transactionError != nil {
			return transactionError
		}
		verrs, transactionError = db.ValidateAndSave(shipment)
		if transactionError != nil {
			return transactionError
		}
		// force validation errors to fail the transaction...
		if verrs.HasAny() {
			return errors.New("error saving shipment")
		}

		verrs, transactionError = db.ValidateAndSave(shipment.StorageInTransits)
		if transactionError != nil {
			return transactionError
		}
		if verrs.HasAny() {
			return errors.New("error saving storage in transits")
		}
		shipmentPricer := NewShipmentPricer(db, c.engine, c.planner)
		verrs, transactionError = shipmentPricer.PriceShipment(shipment, ShipmentPriceNEW)
		if transactionError != nil {
			return transactionError
		}
		if verrs.HasAny() {
			return errors.New("error saving shipment line items")
		}

		return nil
	})

	return verrs, err
}

func NewShipmentDeliverAndPricer(
	db *pop.Connection,
	engine *rateengine.RateEngine,
	planner route.Planner,
) services.ShipmentDeliverAndPricer {
	return &shipmentDeliverAndPricer{
		db:      db,
		engine:  engine,
		planner: planner,
	}
}
