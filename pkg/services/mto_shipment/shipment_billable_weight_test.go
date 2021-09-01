package mtoshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestShipmentBillableWeightCalculator() {
	billableWeightCalculator := NewShipmentBillableWeightCalculator()

	suite.T().Run("If the shipment has a lower reweigh weight and a higher original weight and no set billable weight cap, it should return the reweigh weight", func(t *testing.T) {
		reweighWeight := unit.Pound(900)
		originalWeight := unit.Pound(1000)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PrimeActualWeight: &originalWeight,
			},
		})

		testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, reweighWeight)

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(suite.TestAppContext(), shipment.ID)

		suite.NoError(err)
		if shipment.Reweigh != nil {
			suite.Equal(shipment.Reweigh.Weight, billableWeightCalculations.CalculatedBillableWeight)
		}
	})

	suite.T().Run("If the shipment has a lower original weight and a higher reweigh weight and no set billable weight cap, it should return the original weight", func(t *testing.T) {
		reweighWeight := unit.Pound(1000)
		originalWeight := unit.Pound(900)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PrimeActualWeight: &originalWeight,
			},
		})

		testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, reweighWeight)

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(suite.TestAppContext(), shipment.ID)

		suite.NoError(err)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("If the billable weight cap is set it should be returned", func(t *testing.T) {
		reweighWeight := unit.Pound(1000)
		originalWeight := unit.Pound(900)
		billableWeight := unit.Pound(950)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PrimeActualWeight: &originalWeight,
				BillableWeightCap: &billableWeight,
			},
		})

		testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, reweighWeight)

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(suite.TestAppContext(), shipment.ID)

		suite.NoError(err)
		suite.Equal(shipment.BillableWeightCap, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("Passing in a bad shipment id returns a Not Found error", func(t *testing.T) {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := billableWeightCalculator.CalculateShipmentBillableWeight(suite.TestAppContext(), badShipmentID)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})
}
