package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// shipmentBillableWeightCalculator handles the db connection
type shipmentBillableWeightCalculator struct {
}

// NewShipmentBillableWeightCalculator updates the address for an MTO Shipment
func NewShipmentBillableWeightCalculator() services.ShipmentBillableWeightCalculator {
	return &shipmentBillableWeightCalculator{}
}

// CalculateShipmentBillableWeight calculates a shipment's billable weight.
// If there's a billableWeightCap set that takes precedence.
// Warning: The reweigh object is assumed to have been loaded for the passed in shipment for this service to
// guarantee that the correct calculated weight is returned!
// Without reweigh EagerPreload(ed) there is the risk of miscalculation.
// Due to the nature of EagerPreload, we can no longer tell if Reweigh was NOT preloaded
// OR if the shipment does not have a Reweigh (there is no good way to do this as of this PR)
// https://github.com/transcom/mymove/pull/10780
func (f *shipmentBillableWeightCalculator) CalculateShipmentBillableWeight(shipment *models.MTOShipment) services.BillableWeightInputs {
	var calculatedWeight *unit.Pound
	var reweighWeight *unit.Pound
	var primeActualWeight *unit.Pound
	// Warning: This function assumes that the shipment Reweigh was eager loaded!
	if shipment.Reweigh != nil && shipment.Reweigh.ID != uuid.Nil {
		if shipment.Reweigh.Weight != nil && shipment.PrimeActualWeight != nil {
			reweighWeight = shipment.Reweigh.Weight
			primeActualWeight = shipment.PrimeActualWeight
			if int(*primeActualWeight) < int(*reweighWeight) {
				calculatedWeight = primeActualWeight
			} else if int(*reweighWeight) > 0 {
				// Only use the reweigh weight if it's greater than 0
				calculatedWeight = reweighWeight
			} else {
				// If the prime actual weight is not lower than the reweigh weight, but the
				// reweigh weight is 0, use the prime actual weight.
				calculatedWeight = primeActualWeight
			}
		} else if shipment.Reweigh.Weight == nil && shipment.PrimeActualWeight != nil {
			// if there is no reweigh weight, use the prime actual weight if it is not nil.
			calculatedWeight = shipment.PrimeActualWeight
		}

	} else if shipment.BillableWeightCap == nil {
		calculatedWeight = shipment.PrimeActualWeight
	}

	//Take the lowest between 110% prime estimated and the actual weight, unless shipment is NTSR in which case
	//it should take lowest between 110% prime estimated weight and ntsRecordedWeight
	if shipment.PrimeEstimatedWeight != nil && shipment.PrimeActualWeight != nil && calculatedWeight != nil {
		adjustedEstimatedWeight := unit.Pound(shipment.PrimeEstimatedWeight.Float64() * float64(1.1))
		if shipment.ShipmentType != "HHG_OUTOF_NTS_DOMESTIC" {
			if adjustedEstimatedWeight < *calculatedWeight {
				calculatedWeight = &adjustedEstimatedWeight
			}
		} else {
			if shipment.NTSRecordedWeight != nil {
				if adjustedEstimatedWeight < *shipment.NTSRecordedWeight {
					calculatedWeight = &adjustedEstimatedWeight
				} else {
					calculatedWeight = shipment.NTSRecordedWeight
				}
			}
		}
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
