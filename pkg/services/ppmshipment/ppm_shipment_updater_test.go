package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {

	// One-time test setup
	ppmEstimator := NewEstimatePPM()
	ppmShipmentUpdater := NewPPMShipmentUpdater(ppmEstimator)

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				ExpectedDepartureDate: testdatagen.NextValidMoveDate,
				PickupPostalCode:      "79912",
				DestinationPostalCode: "90909",
				SitExpected:           models.BoolPointer(false),
			},
		})

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek),
			PickupPostalCode:      "79906",
			DestinationPostalCode: "94303",
			SitExpected:           models.BoolPointer(true),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that should now be updated
		newPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(newPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(newPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(newPPM.SitExpected, updatedPPM.SitExpected)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations - add secondary zips", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		newPPM := models.PPMShipment{
			SecondaryPickupPostalCode:      models.StringPointer("79906"),
			SecondaryDestinationPostalCode: models.StringPointer("94303"),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.SecondaryPickupPostalCode, *updatedPPM.SecondaryPickupPostalCode)
		suite.Equal(*newPPM.SecondaryDestinationPostalCode, *updatedPPM.SecondaryDestinationPostalCode)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations - remove secondary zips", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SecondaryPickupPostalCode:      models.StringPointer("79906"),
				SecondaryDestinationPostalCode: models.StringPointer("94303"),
			},
		})

		newPPM := models.PPMShipment{
			SecondaryPickupPostalCode:      models.StringPointer(""),
			SecondaryDestinationPostalCode: models.StringPointer(""),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)

		// Fields that should now be updated
		suite.Nil(updatedPPM.SecondaryPickupPostalCode)
		suite.Nil(updatedPPM.SecondaryDestinationPostalCode)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add estimated weights - no pro gear", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		newPPM := models.PPMShipment{
			EstimatedWeight: models.PoundPointer(4000),
			HasProGear:      models.BoolPointer(false),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.NotNil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add estimated weights - has pro gear", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		newPPM := models.PPMShipment{
			EstimatedWeight:     models.PoundPointer(4000),
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       models.PoundPointer(1000),
			SpouseProGearWeight: models.PoundPointer(0),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*newPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*newPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.NotNil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated weights - pro gear no to yes", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight: models.PoundPointer(4000),
				HasProGear:      models.BoolPointer(false),
			},
		})

		newPPM := models.PPMShipment{
			EstimatedWeight:     models.PoundPointer(4500),
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       models.PoundPointer(1000),
			SpouseProGearWeight: models.PoundPointer(0),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*newPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*newPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.NotNil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated weights - pro gear yes to no", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:     models.PoundPointer(4000),
				HasProGear:          models.BoolPointer(true),
				ProGearWeight:       models.PoundPointer(1000),
				SpouseProGearWeight: models.PoundPointer(0),
			},
		})

		newPPM := models.PPMShipment{
			EstimatedWeight: models.PoundPointer(4500),
			HasProGear:      models.BoolPointer(false),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.NotNil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add advance info - no advance", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:    models.PoundPointer(4000),
				HasProGear:         models.BoolPointer(false),
				EstimatedIncentive: models.CentPointer(unit.Cents(1000000)),
			},
		})

		newPPM := models.PPMShipment{
			AdvanceRequested: models.BoolPointer(false),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.AdvanceRequested, *updatedPPM.AdvanceRequested)
		suite.Nil(updatedPPM.Advance)
	})

	suite.Run("Can successfully update a PPMShipment - add advance info - yes advance", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:    models.PoundPointer(4000),
				HasProGear:         models.BoolPointer(false),
				EstimatedIncentive: models.CentPointer(unit.Cents(1000000)),
			},
		})

		newPPM := models.PPMShipment{
			AdvanceRequested: models.BoolPointer(true),
			Advance:          models.CentPointer(unit.Cents(300000)),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.AdvanceRequested, *updatedPPM.AdvanceRequested)
		suite.Equal(*newPPM.Advance, *updatedPPM.Advance)
	})

	suite.Run("Can successfully update a PPMShipment - edit advance - advance requested no to yes", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				AdvanceRequested: models.BoolPointer(false),
				Advance:          nil,
			},
		})

		newPPM := models.PPMShipment{
			AdvanceRequested: models.BoolPointer(true),
			Advance:          models.CentPointer(unit.Cents(400000)),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.AdvanceRequested, *updatedPPM.AdvanceRequested)
		suite.Equal(*newPPM.Advance, *updatedPPM.Advance)
	})

	suite.Run("Can successfully update a PPMShipment - edit advance - advance requested yes to no", func() {
		appCtx := suite.AppContextForTest()

		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				AdvanceRequested: models.BoolPointer(true),
				Advance:          models.CentPointer(unit.Cents(300000)),
			},
		})

		newPPM := models.PPMShipment{
			AdvanceRequested: models.BoolPointer(false),
		}

		updatedPPM, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SitExpected, updatedPPM.SitExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.AdvanceRequested, *updatedPPM.AdvanceRequested)
		suite.Nil(updatedPPM.Advance)
	})

	suite.Run("Can't update if Shipment can't be found", func() {
		badMTOShipmentID := uuid.Must(uuid.NewV4())

		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &models.PPMShipment{}, badMTOShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", badMTOShipmentID.String()), err.Error())
	})

	suite.Run("Can't update if there is invalid input", func() {
		appCtx := suite.AppContextForTest()

		originalPPMShipment := testdatagen.MakeDefaultPPMShipment(appCtx.DB())

		// Easiest invalid input to trigger is to set an invalid Advance value. The rest are harder to trigger based
		// on how the service object is set up.
		newPPMShipment := models.PPMShipment{
			EstimatedIncentive: models.CentPointer(unit.Cents(1000000)),
			Advance:            models.CentPointer(unit.Cents(3000000)),
		}

		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPMShipment, originalPPMShipment.ShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the PPM shipment.", err.Error())
	})
}
