package mtoshipment

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

//shipmentBillableWeightCalculator handles the db connection
type shipmentBillableWeightCalculator struct {
}

// NewShipmentBillableWeightCalculator updates the address for an MTO Shipment
func NewShipmentBillableWeightCalculator() services.ShipmentBillableWeightCalculator {
	return &shipmentBillableWeightCalculator{}
}

// CalculateShipmentBillableWeight calculates a shipment's billable weight
func (f *shipmentBillableWeightCalculator) CalculateShipmentBillableWeight(shipment *models.MTOShipment) services.BillableWeightInputs {
	var calculatedWeight *unit.Pound
	var reweighWeight *unit.Pound
	if shipment.Reweigh != nil {
		if shipment.Reweigh.Weight != nil && shipment.PrimeActualWeight != nil {
			reweighWeight = shipment.Reweigh.Weight
			if int(*shipment.PrimeActualWeight) < int(*reweighWeight) {
				calculatedWeight = shipment.PrimeActualWeight
			} else {
				calculatedWeight = reweighWeight
			}
		}
	} else if shipment.Reweigh == nil && shipment.BillableWeightCap == nil {
		calculatedWeight = shipment.PrimeActualWeight
	}

	if shipment.BillableWeightCap != nil {
		calculatedWeight = shipment.BillableWeightCap
	}

	hasOverride := shipment.BillableWeightCap != nil
	return services.BillableWeightInputs{
		CalculatedBillableWeight: calculatedWeight,
		OriginalWeight:           shipment.PrimeActualWeight,
		ReweighWeight:            reweighWeight,
		HadManualOverride:        &hasOverride,
	}
}
