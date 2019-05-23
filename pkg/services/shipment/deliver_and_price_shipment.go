package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/auth"
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
func (c DeliverAndPriceShipment) Call(deliveryDate time.Time, shipment *models.Shipment, session *auth.Session) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	err := shipment.Deliver(deliveryDate)
	if err != nil {
		return validate.NewErrors(), err
	}

	sitDeliverer := sitservice.NewStorageInTransitsDeliverer(c.DB)
	_, verrs, err = sitDeliverer.DeliverStorageInTransits(shipment.ID, &deliveryDate, session)
	if err != nil {
		return verrs, err
	}
	return PriceShipment{DB: c.DB, Engine: c.Engine, Planner: c.Planner}.Call(shipment, ShipmentPriceNEW)
}
