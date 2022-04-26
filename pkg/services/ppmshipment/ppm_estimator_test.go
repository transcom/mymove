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
		//oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		oldPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
			},
			Stub: true,
		})
		ppmEstimator := NewEstimatePPM()

		estimatedWeight := unit.Pound(5000)
		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                models.PPMShipmentStatusDraft,
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
		suite.Equal(unit.Cents(1000000), *ppmEstimate)
	})

	suite.Run("Estimated Incentive - Success - clears advance and advance requested values", func() {
		oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())
		ppmEstimator := NewEstimatePPM()

		estimatedWeight := unit.Pound(5000)
		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                models.PPMShipmentStatusDraft,
			ExpectedDepartureDate: time.Time(time.Date(2022, 12, 11, 0, 0, 0, 0, time.UTC)),
			PickupPostalCode:      "20636",
			DestinationPostalCode: "94040",
			EstimatedWeight:       &estimatedWeight,
			SitExpected:           models.BoolPointer(false),
			AdvanceRequested:      models.BoolPointer(true),
			Advance:               models.CentPointer(unit.Cents(498700)),
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Nil(newPPM.Advance)
		suite.Nil(newPPM.AdvanceRequested)
		suite.Equal(unit.Cents(1000000), *ppmEstimate)
	})

	suite.Run("Estimated Incentive - does not change when required fields are the same", func() {
		oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
			},
		})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    oldPPMShipment.ID,
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                models.PPMShipmentStatusDraft,
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: oldPPMShipment.DestinationPostalCode,
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
			SitExpected:           oldPPMShipment.SitExpected,
		}

		estimatedIncentive, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.Equal(*oldPPMShipment.EstimatedWeight, *newPPM.EstimatedWeight)
		suite.Equal(oldPPMShipment.DestinationPostalCode, newPPM.DestinationPostalCode)
		suite.True(oldPPMShipment.ExpectedDepartureDate.Equal(newPPM.ExpectedDepartureDate))
		suite.Equal(*oldPPMShipment.EstimatedIncentive, *estimatedIncentive)
	})
	suite.Run("Estimated Incentive - Failure - is not created when status is not DRAFT", func() {
		oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
			},
		})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                models.PPMShipmentStatusPaymentApproved,
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: "94040",
			EstimatedWeight:       oldPPMShipment.EstimatedWeight,
			SitExpected:           oldPPMShipment.SitExpected,
			EstimatedIncentive:    models.CentPointer(unit.Cents(500000)),
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Nil(ppmEstimate)
		suite.Equal(models.CentPointer(unit.Cents(500000)), newPPM.EstimatedIncentive)
	})

	suite.Run("Estimated Incentive - Failure - is not created when Estimated Weight is missing", func() {
		oldPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})
		ppmEstimator := NewEstimatePPM()

		newPPM := models.PPMShipment{
			ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
			ShipmentID:            oldPPMShipment.ShipmentID,
			Status:                models.PPMShipmentStatusDraft,
			ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
			PickupPostalCode:      oldPPMShipment.PickupPostalCode,
			DestinationPostalCode: "94040",
			EstimatedWeight:       nil,
			SitExpected:           oldPPMShipment.SitExpected,
		}

		_, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NoError(err)
		suite.Nil(newPPM.EstimatedIncentive)
	})
}
