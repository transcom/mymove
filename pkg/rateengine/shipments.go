package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
	"time"
)

// CostByShipment struct containing shipment and cost
type CostByShipment struct {
	Shipment models.Shipment
	Cost     CostComputation
}

// HandleRunRateEngineOnShipment runs the rate engine on a shipment and returns the shipment and cost
func HandleRunRateEngineOnShipment(shipment models.Shipment, engine *RateEngine) (CostByShipment, error) {
	daysInSIT := 0
	var sitDiscount unit.DiscountRate
	sitDiscount = 0.0
	// Apply rate engine to shipment
	var shipmentCost CostByShipment
	cost, err := engine.ComputeShipment(unit.Pound(*shipment.WeightEstimate),
		shipment.PickupAddress.PostalCode,
		shipment.DeliveryAddress.PostalCode,
		time.Time(*shipment.ActualPickupDate),
		daysInSIT, // We don't want any SIT charges
		.4,        // TODO: placeholder: need to get actual linehaul discount
		sitDiscount,
	)
	if err != nil {
		return CostByShipment{}, err
	}

	shipmentCost = CostByShipment{
		Shipment: shipment,
		Cost:     cost,
	}
	return shipmentCost, err
}
