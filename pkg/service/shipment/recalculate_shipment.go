package shipment

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
)

// RecalculateShipment is a service object to re-price a Shipment
type RecalculateShipment struct {
	DB     *pop.Connection
	Engine *rateengine.RateEngine
}

func (c RecalculateShipment) saveShipmentStatus(status models.ShipmentStatus, shipment *models.Shipment) error {
	err := shipment.Recalculate(status)
	if err != nil {
		return err
	}

	// Save and validate Shipment after changing Shipment's Status
	verrs, err := models.SaveShipment(c.DB, shipment)
	if verrs.HasAny() || err != nil {
		verrsString := "verrs: " + verrs.String()
		saveError := errors.Wrap(err, "Error saving shipment status for RecalculateShipment"+verrsString)
		return saveError
	}
	return nil
}

// Call recalculates a Shipment
func (c RecalculateShipment) Call(shipment *models.Shipment) (*validate.Errors, error) {

	// Save Status and temporarily set new status for Shipment
	saveStatus := shipment.Status
	err := c.saveShipmentStatus(models.ShipmentStatusRECALCULATE, shipment)
	if err != nil {
		return validate.NewErrors(), err
	}
	defer c.saveShipmentStatus(saveStatus, shipment)

	// Delivering a shipment is a trigger to populate several shipment line items in the database.  First
	// calculate charges, then submit the updated shipment record and line items in a DB transaction.
	shipmentCost, err := c.Engine.HandleRunOnShipment(*shipment)
	if err != nil {
		return validate.NewErrors(), err
	}

	// Delete Base Shipment Line Items for repricing
	err = shipment.DeleteBaseShipmentLineItems(c.DB)
	if err != nil {
		return validate.NewErrors(), err
	}

	// Save and validate Shipment after deleting Base Shipment Line Items
	verrs, err := models.SaveShipment(c.DB, shipment)
	if verrs.HasAny() || err != nil {
		saveError := errors.Wrap(err, "Error saving shipment for RecalculateShipment")
		return verrs, saveError
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

	verrs, err = shipment.SaveShipmentAndLineItems(c.DB, lineItems, preApprovals)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return validate.NewErrors(), nil
}
