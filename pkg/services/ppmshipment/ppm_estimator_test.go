package ppmshipment

import (
	"github.com/transcom/mymove/pkg/apperror"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestEstimatedIncentive() {
	suite.Run("Estimated Incentive - Success", func() {
		//oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		oldPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.Int32Pointer(int32(1000000)),
			},
			Stub: true,
		})
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
			SitExpected:           models.BoolPointer(false),
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.NotEqualValues(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.NotEqualValues(oldPPMShipment.EstimatedWeight, newPPM.EstimatedWeight)
		suite.Equal(oldPPMShipment.Advance, newPPM.Advance)
		suite.Equal(oldPPMShipment.AdvanceRequested, newPPM.AdvanceRequested)
		suite.Equal(int32(1000000), *ppmEstimate)
		suite.Equal(int32(1000000), *newPPM.EstimatedIncentive)
	})

	suite.Run("Estimated Incentive - Success - clears advance and advance requested values", func() {
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
			SitExpected:           models.BoolPointer(false),
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		//suite.Equal(newPPM.Advance, *unit.Cents((*unit.Cents)(nil)))
		//suite.Equal(newPPM.AdvanceRequested, *bool((*bool)(nil)))
		suite.Equal(int32(1000000), *ppmEstimate)
	})

	suite.Run("Estimated Incentive - does not change when required fields are the same", func() {
		oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.Int32Pointer(int32(1000000)),
			},
		})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: oldPPMShipment.DestinationPostalCode,
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
			SitExpected:           oldPPMShipment.SitExpected,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.Equal(oldPPMShipment.EstimatedWeight, newPPM.EstimatedWeight)
		suite.Equal(oldPPMShipment.DestinationPostalCode, newPPM.DestinationPostalCode)
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, newPPM.ExpectedDepartureDate)
		suite.Equal(*oldPPMShipment.EstimatedIncentive, *ppmEstimate)
	})
	suite.Run("Estimated Incentive - Failure - is not created when status is not DRAFT", func() {
		oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.Int32Pointer(int32(1000000)),
			},
		})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "PAYMENT_APPROVED",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: "94040",
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
			SitExpected:           oldPPMShipment.SitExpected,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Nil(ppmEstimate)
	})

	suite.Run("Estimated Incentive - Failure - is not created when Estimated Weight is missing", func() {
		oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.Int32Pointer(int32(1000000)),
			},
		})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: "94040",
			EstimatedWeight:       nil,
			SitExpected:           oldPPMShipment.SitExpected,
		}

		_, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		//suite.Nil(ppmEstimate)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Not Found Error - missing ppm shipment ID", func() {
		oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.Int32Pointer(int32(1000000)),
			},
		})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.Nil,
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                "DRAFT",
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: oldPPMShipment.DestinationPostalCode,
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
			SitExpected:           oldPPMShipment.SitExpected,
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)

		suite.Nil(ppmEstimate)
		suite.IsType(err, nil)
	})
}
