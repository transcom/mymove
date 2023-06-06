package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestShipmentBillableWeightCalculator() {
	billableWeightCalculator := NewShipmentBillableWeightCalculator()

	suite.Run("If the shipment has a lower reweigh weight and a higher original weight and no set billable weight cap, it should return the reweigh weight", func() {
		reweighWeight := unit.Pound(900)
		originalWeight := unit.Pound(1000)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.Must(uuid.NewV4()),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)

		suite.Equal(shipment.Reweigh.Weight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("If the shipment has a lower original weight and a higher reweigh weight and no set billable weight cap, it should return the original weight", func() {
		reweighWeight := unit.Pound(1000)
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.Must(uuid.NewV4()),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)

		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("If the shipment has an original weight and reweigh with no id and no set billable weight cap, it should return the original weight", func() {
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			ID: uuid.Nil,
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("If the shipment has an original weight and no reweigh weight and no set billable weight cap, it should return the original weight", func() {
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			ID: uuid.Must(uuid.NewV4()),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("If the shipment has an original weight and 0 for the reweigh weight and no set billable weight cap, it should return the original weight", func() {
		reweighWeight := unit.Pound(0)
		originalWeight := unit.Pound(900)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.Must(uuid.NewV4()),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("If the billable weight cap is set it should be returned", func() {
		reweighWeight := unit.Pound(1000)
		originalWeight := unit.Pound(900)
		billableWeight := unit.Pound(950)
		reweigh := models.Reweigh{
			Weight: &reweighWeight,
			ID:     uuid.Must(uuid.NewV4()),
		}
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
			BillableWeightCap: &billableWeight,
			Reweigh:           &reweigh,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Equal(shipment.BillableWeightCap, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("Returns the PrimeActualWeight if the shipment passed in doesnt have a Reweigh eager loaded", func() {
		originalWeight := unit.Pound(900)
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
		}

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Equal(&originalWeight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("Eagerly loaded reweigh where a reweigh exists", func() {
		actualWeight := unit.Pound(3100)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeActualWeight: &actualWeight,
				},
			},
		}, nil)
		reweigh := testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, unit.Pound(3000))

		var dbShipment models.MTOShipment
		err := suite.DB().EagerPreload("Reweigh").Find(&dbShipment, shipment.ID)
		suite.FatalNoError(err)

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&dbShipment)
		suite.Equal(reweigh.Weight, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("Eagerly loaded reweigh where a reweigh does not exist", func() {
		actualWeight := unit.Pound(3100)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeActualWeight: &actualWeight,
				},
			},
		}, nil)
		// No reweigh

		var dbShipment models.MTOShipment
		err := suite.DB().EagerPreload("Reweigh").Find(&dbShipment, shipment.ID)
		suite.FatalNoError(err)

		billableWeightCalculations := billableWeightCalculator.CalculateShipmentBillableWeight(&dbShipment)
		suite.Equal(shipment.PrimeActualWeight, billableWeightCalculations.CalculatedBillableWeight)
	})
}
