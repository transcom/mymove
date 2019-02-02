package shipment

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"go.uber.org/zap"
)

// RecalculateShipment is a service object to re-price a Shipment
type RecalculateShipment struct {
	DB     *pop.Connection
	Logger *zap.Logger
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

	/* TODO: This doesn't seem necessary, line items are not returned to the viewer and are therefore
	   TODO: not actionable until the recalcuation is finished

	// Save Status and temporarily set new status for Shipment
	saveStatus := shipment.Status
	err := c.saveShipmentStatus(models.ShipmentStatusRECALCULATE, shipment)
	if err != nil {
		return validate.NewErrors(), err
	}
	*/

	c.Logger.Info("Recalculate ShipmentID: ",
		zap.Any("shipment.ID", shipment.ID))
	// Re-price Shipment
	//verrs, err := PriceShipment{DB: c.DB, Engine: c.Engine}.Call(shipment, ShipmentPriceRECALCULATE)
	return PriceShipment{DB: c.DB, Engine: c.Engine}.Call(shipment, ShipmentPriceRECALCULATE)

	/* TODO: This doesn't seem necessary, line items are not returned to the viewer and are therefore
	   TODO: not actionable until the recalcuation is finished

	var finalError error
	// Revert ShipmentStatus back to original status (i.e., DELIVERED or COMPLETED)
	err2 := c.saveShipmentStatus(saveStatus, shipment)


	// Wrap error to return both errors (if any) back from PriceShipment{}.Call() and saveShipmentStatus()
	if err2 != nil {
		finalError = errors.Wrapf(err2, "Failed to revert shipment status back to %s", saveStatus)
		if err != nil {
			finalError = errors.Wrap(finalError, err.Error())
		}
	}

	return verrs, err
	*/
}
