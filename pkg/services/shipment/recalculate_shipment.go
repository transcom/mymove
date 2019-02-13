package shipment

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
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

// Call recalculates a Shipment
func (c RecalculateShipment) Call(shipment *models.Shipment) (*validate.Errors, error) {
	c.Logger.Info("Recalculate Shipment: ", zap.Any("shipment.ID", shipment.ID),
		zap.Any("shipment.Status", shipment.Status))

	// Re-price Shipment
	return PriceShipment{DB: c.DB, Engine: c.Engine}.Call(shipment, ShipmentPriceRECALCULATE)
}
