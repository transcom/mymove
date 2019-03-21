package shipment

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
)

// PricingType describe the type of pricing to do for a shipment
type PricingType string

const (
	// ShipmentPriceNEW captures enum value "NEW"
	ShipmentPriceNEW PricingType = "NEW"
	// ShipmentPriceRECALCULATE captures enum value "RECALCULATE"
	ShipmentPriceRECALCULATE PricingType = "RECALCULATE"
)

// PriceShipment is a service object to price a Shipment
type PriceShipment struct {
	DB      *pop.Connection
	Engine  *rateengine.RateEngine
	Planner route.Planner
}

// Call prices a Shipment
func (c PriceShipment) Call(shipment *models.Shipment, price PricingType) (*validate.Errors, error) {
	origin := shipment.PickupAddress
	if origin == nil || origin.ID == uuid.Nil {
		return validate.NewErrors(), errors.New("PickupAddress not provided")
	}

	destination := shipment.Move.Orders.NewDutyStation.Address
	if shipment.DestinationAddressOnAcceptance != nil {
		destination = *shipment.DestinationAddressOnAcceptance
	}

	if destination.ID == uuid.Nil {
		return validate.NewErrors(), errors.New("Destination address not provided")
	}

	distanceCalculation, err := models.NewDistanceCalculation(c.Planner, *origin, destination)
	if err != nil {
		return validate.NewErrors(), errors.Wrap(err, "Error creating DistanceCalculation model")
	}

	// Delivering a shipment is a trigger to populate several shipment line items in the database.  First
	// calculate charges, then submit the updated shipment record and line items in a DB transaction.
	shipmentCost, err := c.Engine.HandleRunOnShipment(*shipment, distanceCalculation)
	if err != nil {
		return validate.NewErrors(), err
	}

	if price == ShipmentPriceRECALCULATE {
		// Delete Base Shipment Line Items for repricing
		err = shipment.DeleteBaseShipmentLineItems(c.DB)
		if err != nil {
			return validate.NewErrors(), err
		}

		// Save and validate Shipment after deleting Base Shipment Line Items
		verrs, err := models.SaveShipment(c.DB, shipment)
		if verrs.HasAny() || err != nil {
			saveError := errors.Wrap(err, "Error saving shipment for ShipmentPriceRECALCULATE")
			return verrs, saveError
		}
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

	verrs, err := shipment.SaveShipmentAndPricingInfo(c.DB, lineItems, preApprovals, distanceCalculation)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	if price == ShipmentPriceRECALCULATE {
		log := models.ShipmentRecalculateLog{ShipmentID: shipment.ID}
		log.SaveShipmentRecalculateLog(c.DB)
	}

	return validate.NewErrors(), nil
}
