package shipment

import (
	"errors"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	sitservice "github.com/transcom/mymove/pkg/services/storage_in_transit"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type DeliverAndPriceShipment struct {
	DB      *pop.Connection
	Engine  *rateengine.RateEngine
	Planner route.Planner
}

// Call delivers a Shipment and prices associated line items
func (c DeliverAndPriceShipment) Call(deliveryDate time.Time, shipment *models.Shipment, tspID uuid.UUID) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	var err error
	c.DB.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("rollback")

		err = shipment.Deliver(deliveryDate)
		if err != nil {
			return err
		}

		verrs, err = PriceShipment{DB: c.DB, Engine: c.Engine, Planner: c.Planner}.Call(shipment, ShipmentPriceNEW)
		if err != nil || verrs.HasAny() {
			return transactionError
		}

		sitDeliverer := sitservice.NewStorageInTransitsDeliverer(c.DB)
		_, verrs, err = sitDeliverer.DeliverStorageInTransits(shipment.ID, tspID)
		if err != nil || verrs.HasAny() {
			return transactionError
		}

		return nil
	})

	return verrs, err
}
