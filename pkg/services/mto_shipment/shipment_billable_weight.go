package mtoshipment

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
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
func (f *shipmentBillableWeightCalculator) CalculateShipmentBillableWeight(appCtx appcontext.AppContext, shipmentID uuid.UUID) (services.BillableWeightInputs, error) {
	var shipment models.MTOShipment
	var calculatedWeight *unit.Pound
	// var reweighWeight *unit.Pound

	err := appCtx.DB().Q().
		Eager("Reweigh").
		Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return services.BillableWeightInputs{}, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return services.BillableWeightInputs{}, err
	}

	if shipment.Reweigh != nil {
		if shipment.Reweigh.Weight != nil && shipment.PrimeActualWeight != nil {
			if int(*shipment.PrimeActualWeight) < int(*shipment.Reweigh.Weight) {
				calculatedWeight = shipment.PrimeActualWeight
			} else {
				calculatedWeight = shipment.Reweigh.Weight
			}
			fmt.Printf("shipment reweigh weight: %v", int(*shipment.Reweigh.Weight))
		}
	}

	if shipment.BillableWeightCap != nil {
		calculatedWeight = shipment.BillableWeightCap
	}

	// hasOverride := shipment.BillableWeightCap != nil
	hasOverride := true
	return services.BillableWeightInputs{
		CalculatedBillableWeight: calculatedWeight,
		OriginalWeight:           shipment.PrimeActualWeight,
		ReweighWeight:            shipment.Reweigh.Weight,
		HadManualOverride:        &hasOverride,
	}, nil
}
