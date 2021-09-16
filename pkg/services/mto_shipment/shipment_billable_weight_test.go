package mtoshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)

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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)

		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("If the shipment has an original weight and reweigh with no id and no set billable weight cap, it should return the original weight", func(t *testing.T) {
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			ID: uuid.Nil,
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("If the shipment has an original weight and no reweigh weight and no set billable weight cap, it should return the original weight", func(t *testing.T) {
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			ID: uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d985"),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("If the shipment has an original weight and 0 for the reweigh weight and no set billable weight cap, it should return the original weight", func(t *testing.T) {
		reweighWeight := unit.Pound(0)
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d985"),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
		suite.Equal(shipment.BillableWeightCap, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("Returns an error if the shipment passed in doesnt have a Reweigh eager loaded", func(t *testing.T) {
		originalWeight := unit.Pound(900)
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
		}

		_, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	suite.T().Run("Eagerly loaded reweigh where a reweigh exists", func(t *testing.T) {
		actualWeight := unit.Pound(3100)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PrimeActualWeight: &actualWeight,
			},
		})
		reweigh := testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, unit.Pound(3000))

		var dbShipment models.MTOShipment
		err := suite.DB().Eager("Reweigh").Find(&dbShipment, shipment.ID)
		suite.FatalNoError(err)

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&dbShipment)
		suite.NoError(err)
		suite.Equal(reweigh.Weight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("Eagerly loaded reweigh where a reweight does not exist", func(t *testing.T) {
		actualWeight := unit.Pound(3100)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PrimeActualWeight: &actualWeight,
			},
		})
		// No reweigh

		var dbShipment models.MTOShipment
		err := suite.DB().Eager("Reweigh").Find(&dbShipment, shipment.ID)
		suite.FatalNoError(err)

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&dbShipment)
		suite.NoError(err)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.T().Run("Did not eagerly load reweigh even though one exists", func(t *testing.T) {
		actualWeight := unit.Pound(3100)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PrimeActualWeight: &actualWeight,
			},
		})
		testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, unit.Pound(3000))

		var dbShipment models.MTOShipment
		err := suite.DB().Find(&dbShipment, shipment.ID)
		suite.FatalNoError(err)

		_, err = billableWeightCalculator.CalculateShipmentBillableWeight(&dbShipment)
		suite.Error(err)
		suite.Contains(err.Error(), "Invalid shipment, must have Reweigh eager loaded")
	})
}
