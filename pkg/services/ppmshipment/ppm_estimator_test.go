package ppmshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestEstimatedIncentive() {
	suite.Run("Estimated Incentive - Success", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()

		estimatedWeight := unit.Pound(5000)
		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: time.Time(time.Date(2022, 12, 11, 0, 0, 0, 0, time.UTC)),
			PickupPostalCode:      "20636",
			DestinationPostalCode: "94040",
			EstimatedWeight:       &estimatedWeight,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.NotEqualValues(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.NotEqualValues(oldPPMShipment.EstimatedWeight, newPPM.EstimatedWeight)
		suite.Equal(int32(1000000), *ppmEstimate)
	})
	suite.Run("Estimated Incentive - does not change when required fields are the same", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: oldPPMShipment.DestinationPostalCode,
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.Equal(oldPPMShipment.EstimatedWeight, newPPM.EstimatedWeight)
		suite.Equal(oldPPMShipment.DestinationPostalCode, newPPM.DestinationPostalCode)
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, newPPM.ExpectedDepartureDate)
		suite.Equal(oldPPMShipment.EstimatedIncentive, ppmEstimate)
	})
	suite.Run("Estimated Incentive - does not change when required fields are the same", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: oldPPMShipment.DestinationPostalCode,
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.Equal(oldPPMShipment.EstimatedWeight, newPPM.EstimatedWeight)
		suite.Equal(oldPPMShipment.DestinationPostalCode, newPPM.DestinationPostalCode)
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, newPPM.ExpectedDepartureDate)
		suite.Equal(oldPPMShipment.EstimatedIncentive, ppmEstimate)
	})

	suite.Run("Not Found Error - missing ppm shipment ID", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.Nil,
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: oldPPMShipment.DestinationPostalCode,
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)

		suite.Nil(ppmEstimate)
		suite.IsType(err, nil)
	})
}
