package ppmshipment

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {

	// One-time test setup

	fakeEstimatedIncentive := models.CentPointer(unit.Cents(1000000))

	type updateSubtestData struct {
		ppmShipmentUpdater services.PPMShipmentUpdater
	}

	// setUpForTests - Sets up objects/mocks that need to be set up on a per-test basis.
	setUpForTests := func(estimatedIncentiveAmount *unit.Cents, sitEstimatedCost *unit.Cents, estimatedIncentiveError error) (subtestData updateSubtestData) {
		ppmEstimator := mocks.PPMEstimator{}

		ppmEstimator.
			On(
				"EstimateIncentiveWithDefaultChecks",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(estimatedIncentiveAmount, sitEstimatedCost, estimatedIncentiveError)

		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		subtestData.ppmShipmentUpdater = NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		return subtestData
	}

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil)

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				ExpectedDepartureDate: testdatagen.NextValidMoveDate,
				PickupPostalCode:      "79912",
				DestinationPostalCode: "90909",
				SITExpected:           models.BoolPointer(false),
			},
		})

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek),
			PickupPostalCode:      "79906",
			DestinationPostalCode: "94303",
			SITExpected:           models.BoolPointer(true),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that should now be updated
		newPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(newPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(newPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(newPPM.SITExpected, updatedPPM.SITExpected)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations - weights already set", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				ExpectedDepartureDate: testdatagen.NextValidMoveDate,
				PickupPostalCode:      "79912",
				DestinationPostalCode: "90909",
				SITExpected:           models.BoolPointer(false),
				EstimatedWeight:       models.PoundPointer(4000),
				HasProGear:            models.BoolPointer(false),
				EstimatedIncentive:    fakeEstimatedIncentive,
			},
		})

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek),
			PickupPostalCode:      "79906",
			DestinationPostalCode: "94303",
			SITExpected:           models.BoolPointer(true),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)

		// Fields that should now be updated
		newPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(newPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(newPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(newPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*newFakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations - add secondary zips", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil)

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		newPPM := models.PPMShipment{
			SecondaryPickupPostalCode:      models.StringPointer("79906"),
			SecondaryDestinationPostalCode: models.StringPointer("94303"),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.SecondaryPickupPostalCode, *updatedPPM.SecondaryPickupPostalCode)
		suite.Equal(*newPPM.SecondaryDestinationPostalCode, *updatedPPM.SecondaryDestinationPostalCode)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations - remove secondary zips", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil)

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

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Nil(updatedPPM.SecondaryPickupPostalCode)
		suite.Nil(updatedPPM.SecondaryDestinationPostalCode)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add estimated weights - no pro gear", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		newPPM := models.PPMShipment{
			EstimatedWeight: models.PoundPointer(4000),
			HasProGear:      models.BoolPointer(false),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.Equal(*fakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add estimated weights - has pro gear", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		newPPM := models.PPMShipment{
			EstimatedWeight:     models.PoundPointer(4000),
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       models.PoundPointer(1000),
			SpouseProGearWeight: models.PoundPointer(0),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*newPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*newPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.Equal(*fakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated weights - pro gear no to yes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:    models.PoundPointer(4000),
				HasProGear:         models.BoolPointer(false),
				EstimatedIncentive: fakeEstimatedIncentive,
			},
		})

		newPPM := models.PPMShipment{
			EstimatedWeight:     models.PoundPointer(4500),
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       models.PoundPointer(1000),
			SpouseProGearWeight: models.PoundPointer(0),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*newPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*newPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*newFakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated weights - pro gear yes to no", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:     models.PoundPointer(4000),
				HasProGear:          models.BoolPointer(true),
				ProGearWeight:       models.PoundPointer(1000),
				SpouseProGearWeight: models.PoundPointer(0),
				EstimatedIncentive:  fakeEstimatedIncentive,
			},
		})

		newPPM := models.PPMShipment{
			EstimatedWeight: models.PoundPointer(4500),
			HasProGear:      models.BoolPointer(false),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)
		suite.Equal(*newFakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add advance info - no advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:    models.PoundPointer(4000),
				HasProGear:         models.BoolPointer(false),
				EstimatedIncentive: fakeEstimatedIncentive,
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance: models.BoolPointer(false),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Nil(updatedPPM.AdvanceAmountRequested)
	})

	suite.Run("Can successfully update a PPMShipment - add advance info - yes advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedWeight:    models.PoundPointer(4000),
				HasProGear:         models.BoolPointer(false),
				EstimatedIncentive: fakeEstimatedIncentive,
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(300000)),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
	})

	suite.Run("Can successfully update a PPMShipment - office user rejects requested advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})
		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive:     fakeEstimatedIncentive,
				HasRequestedAdvance:    models.BoolPointer(true),
				AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(false),
			AdvanceAmountRequested: nil,
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		rejected := models.PPMAdvanceStatusRejected
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Nil(updatedPPM.AdvanceAmountRequested)
		suite.Equal(&rejected, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - office user edits requested advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})
		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive:     fakeEstimatedIncentive,
				HasRequestedAdvance:    models.BoolPointer(true),
				AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(200000)),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		edited := models.PPMAdvanceStatusEdited
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(&edited, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - office user approves requested advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})
		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive:     fakeEstimatedIncentive,
				HasRequestedAdvance:    models.BoolPointer(true),
				AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		approved := models.PPMAdvanceStatusApproved
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(&approved, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - edit advance - advance requested no to yes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive:     fakeEstimatedIncentive,
				HasRequestedAdvance:    models.BoolPointer(false),
				AdvanceAmountRequested: nil,
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
	})

	suite.Run("Can successfully update a PPMShipment - edit advance - advance requested yes to no", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				EstimatedIncentive:     fakeEstimatedIncentive,
				HasRequestedAdvance:    models.BoolPointer(true),
				AdvanceAmountRequested: models.CentPointer(unit.Cents(300000)),
			},
		})

		newPPM := models.PPMShipment{
			HasRequestedAdvance: models.BoolPointer(false),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Nil(updatedPPM.AdvanceAmountRequested)
	})

	suite.Run("Can successfully update a PPMShipment - edit SIT - yes to no", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil)
		sitLocation := models.SITLocationTypeOrigin

		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SITExpected:               models.BoolPointer(true),
				SITLocation:               &sitLocation,
				SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
				SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
				SITEstimatedWeight:        models.PoundPointer(1000),
				SITEstimatedCost:          models.CentPointer(unit.Cents(69900)),
			},
		})

		newPPM := models.PPMShipment{
			SITExpected: models.BoolPointer(false),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SecondaryPickupPostalCode, updatedPPM.SecondaryPickupPostalCode)
		suite.Equal(originalPPM.SecondaryDestinationPostalCode, updatedPPM.SecondaryDestinationPostalCode)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(*originalPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)

		// Fields that should now be updated
		suite.Equal(*newPPM.SITExpected, *updatedPPM.SITExpected)
		suite.Nil(updatedPPM.SITLocation)
		suite.Nil(updatedPPM.SITEstimatedEntryDate)
		suite.Nil(updatedPPM.SITEstimatedDepartureDate)
		suite.Nil(updatedPPM.SITEstimatedWeight)
		suite.Nil(updatedPPM.SITEstimatedCost)
	})

	suite.Run("Can successfully update a PPMShipment - edit SIT - no to yes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForTests(newFakeEstimatedIncentive, newFakeSITEstimatedCost, nil)
		sitLocation := models.SITLocationTypeOrigin

		originalPPM := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				SITExpected: models.BoolPointer(false),
			},
		})

		newPPM := models.PPMShipment{
			SITExpected:               models.BoolPointer(true),
			SITLocation:               &sitLocation,
			SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
			SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
			SITEstimatedWeight:        models.PoundPointer(1000),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		originalPPM.ExpectedDepartureDate.Equal(updatedPPM.ExpectedDepartureDate)
		suite.Equal(originalPPM.PickupPostalCode, updatedPPM.PickupPostalCode)
		suite.Equal(originalPPM.DestinationPostalCode, updatedPPM.DestinationPostalCode)
		suite.Equal(originalPPM.SecondaryPickupPostalCode, updatedPPM.SecondaryPickupPostalCode)
		suite.Equal(originalPPM.SecondaryDestinationPostalCode, updatedPPM.SecondaryDestinationPostalCode)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(*originalPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)

		// Fields that should now be updated
		suite.Equal(*newPPM.SITExpected, *updatedPPM.SITExpected)
		suite.Equal(*newPPM.SITLocation, *updatedPPM.SITLocation)
		suite.Equal(*newPPM.SITEstimatedEntryDate, *updatedPPM.SITEstimatedEntryDate)
		suite.Equal(*newPPM.SITEstimatedDepartureDate, *updatedPPM.SITEstimatedDepartureDate)
		suite.Equal(*newPPM.SITEstimatedWeight, *updatedPPM.SITEstimatedWeight)
		suite.Equal(*newFakeSITEstimatedCost, *updatedPPM.SITEstimatedCost)
	})

	suite.Run("Can't update if Shipment can't be found", func() {
		badMTOShipmentID := uuid.Must(uuid.NewV4())

		subtestData := setUpForTests(nil, nil, nil)

		updatedPPMShipment, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextWithSessionForTest(&auth.Session{}), &models.PPMShipment{}, badMTOShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", badMTOShipmentID.String()), err.Error())
	})

	suite.Run("Can't update if there is invalid input", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil)

		originalPPMShipment := testdatagen.MakeDefaultPPMShipment(appCtx.DB())

		// Easiest invalid input to trigger is to set an invalid AdvanceAmountRequested value. The rest are harder to
		// trigger based on how the service object is set up.
		newPPMShipment := models.PPMShipment{
			AdvanceAmountRequested: models.CentPointer(unit.Cents(3000000)),
		}

		updatedPPMShipment, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPMShipment, originalPPMShipment.ShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the PPM shipment.", err.Error())
	})

	suite.Run("Can't update if there is an error calculating incentive", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		fakeEstimatedIncentiveError := errors.New("failed to calculate incentive")
		subtestData := setUpForTests(nil, nil, fakeEstimatedIncentiveError)

		originalPPMShipment := testdatagen.MakeDefaultPPMShipment(appCtx.DB())

		newPPMShipment := models.PPMShipment{
			HasRequestedAdvance: models.BoolPointer(false),
		}

		updatedPPMShipment, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPMShipment, originalPPMShipment.ShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.Equal(fakeEstimatedIncentiveError, err)
	})

	suite.Run("Can successfully update a PPMShipment - add W-2 address", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "08004"

		newPPM := models.PPMShipment{
			W2Address: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.NotNil(updatedPPM.W2AddressID)
		suite.Equal(streetAddress1, updatedPPM.W2Address.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.W2Address.StreetAddress2)
		suite.Equal(city, updatedPPM.W2Address.City)
		suite.Equal(state, updatedPPM.W2Address.State)
		suite.Equal(postalCode, updatedPPM.W2Address.PostalCode)
	})

	suite.Run("Can successfully update a PPMShipment - add W-2 address with empty strings for optional fields", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil)

		originalPPM := testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())

		streetAddress1 := "1819 S Cedar Street"
		city := "Fayetteville"
		state := "NC"
		postalCode := "28314"

		newPPM := models.PPMShipment{
			W2Address: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: models.StringPointer(""),
				StreetAddress3: models.StringPointer(""),
				City:           city,
				State:          state,
				PostalCode:     postalCode,
				Country:        models.StringPointer(""),
			},
		}
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.NotNil(updatedPPM.W2AddressID)
		suite.Equal(streetAddress1, updatedPPM.W2Address.StreetAddress1)
		suite.Equal(city, updatedPPM.W2Address.City)
		suite.Equal(state, updatedPPM.W2Address.State)
		suite.Equal(postalCode, updatedPPM.W2Address.PostalCode)
		suite.Nil(updatedPPM.W2Address.StreetAddress2)
		suite.Nil(updatedPPM.W2Address.StreetAddress3)
		suite.Nil(updatedPPM.W2Address.Country)
	})

	suite.Run("Can successfully update a PPMShipment - modify W-2 address", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil)

		address := testdatagen.MakeAddress(appCtx.DB(), testdatagen.Assertions{})
		originalPPM := testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				W2Address:   &address,
				W2AddressID: &address.ID,
			},
		})

		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Cookstown"
		state := "NJ"
		postalCode := "08511"

		newPPM := models.PPMShipment{
			W2Address: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.Equal(address.ID, *updatedPPM.W2AddressID)
		suite.Equal(streetAddress1, updatedPPM.W2Address.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.W2Address.StreetAddress2)
		suite.Equal(city, updatedPPM.W2Address.City)
		suite.Equal(state, updatedPPM.W2Address.State)
		suite.Equal(postalCode, updatedPPM.W2Address.PostalCode)
		suite.Equal(*address.StreetAddress3, *updatedPPM.W2Address.StreetAddress3)
		suite.Equal(*address.Country, *updatedPPM.W2Address.Country)
	})
}
