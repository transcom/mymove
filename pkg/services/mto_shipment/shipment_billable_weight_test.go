package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)

		if shipment.Reweigh != nil {
			suite.Equal(shipment.Reweigh.Weight, billableWeightCalculations.CalculatedBillableWeight)
		}
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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)

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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
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

		billableWeightCalculations, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.NoError(err)
		suite.Equal(shipment.BillableWeightCap, billableWeightCalculations.CalculatedBillableWeight)
	})

	suite.Run("Returns an error if the shipment passed in doesnt have a Reweigh eager loaded", func() {
		originalWeight := unit.Pound(900)
		shipment := models.MTOShipment{
			PrimeActualWeight: &originalWeight,
		}

		_, err := billableWeightCalculator.CalculateShipmentBillableWeight(&shipment)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Eagerly loaded reweigh where a reweigh exists", func() {
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

	suite.Run("Eagerly loaded reweigh where a reweigh does not exist", func() {
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

	suite.Run("Did not eagerly load reweigh even though one exists", func() {
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
