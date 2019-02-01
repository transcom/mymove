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

	return PriceShipment{DB: c.DB, Engine: c.Engine}.Call(shipment, ShipmentPriceRECALCULATE)
}
