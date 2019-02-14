package shipment

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
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
	DB     *pop.Connection
	Engine *rateengine.RateEngine
}

// Call prices a Shipment
func (c PriceShipment) Call(shipment *models.Shipment, price PricingType) (*validate.Errors, error) {
	// Delivering a shipment is a trigger to populate several shipment line items in the database.  First
	// calculate charges, then submit the updated shipment record and line items in a DB transaction.
	shipmentCost, err := c.Engine.HandleRunOnShipment(*shipment)
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

	verrs, err := shipment.SaveShipmentAndLineItems(c.DB, lineItems, preApprovals)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	if price == ShipmentPriceRECALCULATE {
		log := models.ShipmentRecalculateLog{ShipmentID: shipment.ID}
		log.SaveShipmentRecalculateLog(c.DB)
	}

	return validate.NewErrors(), nil
}
