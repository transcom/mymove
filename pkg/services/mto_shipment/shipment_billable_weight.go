package mtoshipment

import (
	"github.com/gofrs/uuid"

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
// if a shipment has a reweigh weight and an original weight, it returns the lowest weight
// if there's a billableWeightCap set that takes precedence
// The reweigh is assumed to have been loaded for the passed in shipment for this service to
// guarantee that the correct calculated weight is returned.
func (f *shipmentBillableWeightCalculator) CalculateShipmentBillableWeight(shipment *models.MTOShipment) (services.BillableWeightInputs, error) {
	var calculatedWeight *unit.Pound
	var reweighWeight *unit.Pound
	if shipment.Reweigh == nil {
		return services.BillableWeightInputs{}, services.NewConflictError(shipment.ID, "Invalid shipment, must have Reweigh eager loaded")
	}
	if shipment.Reweigh != nil && shipment.Reweigh.ID != uuid.Nil {
		if shipment.Reweigh.Weight != nil && shipment.PrimeActualWeight != nil {
			reweighWeight = shipment.Reweigh.Weight
			if int(*shipment.PrimeActualWeight) < int(*reweighWeight) {
				calculatedWeight = shipment.PrimeActualWeight
			} else {
				calculatedWeight = reweighWeight
			}
		}
	} else if shipment.BillableWeightCap == nil {
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
	}, nil
}
