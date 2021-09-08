package mtoshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestShipmentBillableWeightCalculator() {
	billableWeightCalculator := NewShipmentBillableWeightCalculator()

	suite.T().Run("If the shipment has a lower reweigh weight and a higher original weight and no set billable weight cap, it should return the reweigh weight", func(t *testing.T) {
		reweighWeight := unit.Pound(900)
		originalWeight := unit.Pound(1000)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)

		if shipment.Reweigh != nil {
			suite.Equal(shipment.Reweigh.Weight, billableWeightCalculations.CalculatedBillableWeight)
		}
	})

	suite.T().Run("If the shipment has a lower original weight and a higher reweigh weight and no set billable weight cap, it should return the original weight", func(t *testing.T) {
		reweighWeight := unit.Pound(1000)
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d985"),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)

		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("If the shipment has an original weight and no reweigh weight and no set billable weight cap, it should return the original weight", func(t *testing.T) {
		originalWeight := unit.Pound(900)

		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)

		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("If the billable weight cap is set it should be returned", func(t *testing.T) {
		reweighWeight := unit.Pound(1000)
		originalWeight := unit.Pound(900)
		billableWeight := unit.Pound(950)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d986"),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			BillableWeightCap: &billableWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)

		suite.Equal(shipment.BillableWeightCap, billableWeightCalculations.CalculatedBillableWeight)
	})
}
