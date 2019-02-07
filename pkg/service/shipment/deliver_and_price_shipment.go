package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type DeliverAndPriceShipment struct {
	DB      *pop.Connection
	Engine  *rateengine.RateEngine
	Planner route.Planner
}

// Call delivers a Shipment and prices associated line items
func (c DeliverAndPriceShipment) Call(deliveryDate time.Time, shipment *models.Shipment) (*validate.Errors, error) {
	err := shipment.Deliver(deliveryDate)
	if err != nil {
		return validate.NewErrors(), err
	}

	origin := shipment.PickupAddress
	if origin == nil || origin.ID == uuid.Nil {
		return validate.NewErrors(), errors.New("PickupAddress not provided")
	}

	destination := shipment.Move.Orders.NewDutyStation.Address
	if destination.ID == uuid.Nil {
		return validate.NewErrors(), errors.New("New duty station address not provided")
	}

	distanceCalculaton, err := models.NewDistanceCalculation(c.Planner, *origin, destination)
	if err != nil {
		return validate.NewErrors(), errors.Wrap(err, "Error creating DistanceCalculation model")
	}

	// Delivering a shipment is a trigger to populate several shipment line items in the database.  First
	// calculate charges, then submit the updated shipment record and line items in a DB transaction.
	shipmentCost, err := c.Engine.HandleRunOnShipment(*shipment, distanceCalculaton)
	if err != nil {
		return validate.NewErrors(), err
	}

	lineItems, err := rateengine.CreateBaseShipmentLineItems(c.DB, shipmentCost)
	if err != nil {
		return validate.NewErrors(), err
	}

	// When the shipment is delivered we should also price existing approved pre-approval requests
	preApprovals, err := c.Engine.PricePreapprovalRequestsForShipment(*shipment)
	if err != nil {
		return validate.NewErrors(), err
	}

	verrs, err := shipment.SaveShipmentAndPricingInfo(c.DB, lineItems, preApprovals, distanceCalculaton)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return validate.NewErrors(), nil
}
