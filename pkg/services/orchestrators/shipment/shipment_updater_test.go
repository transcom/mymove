package shipment

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ShipmentSuite) TestUpdateShipment() {

	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeSubtestData.

	updateMTOShipmentMethodName := "UpdateMTOShipment"
	updatePPMShipmentMethodName := "UpdatePPMShipmentWithDefaultCheck"
	updateBoatShipmentMethodName := "UpdateBoatShipmentWithDefaultCheck"
	updateMobileHomeShipmentMethodName := "UpdateMobileHomeShipmentWithDefaultCheck"

	type subtestDataObjects struct {
		mockMTOShipmentUpdater        *mocks.MTOShipmentUpdater
		mockPPMShipmentUpdater        *mocks.PPMShipmentUpdater
		mockBoatShipmentUpdater       *mocks.BoatShipmentUpdater
		mockMobileHomeShipmentUpdater *mocks.MobileHomeShipmentUpdater
		shipmentUpdaterOrchestrator   services.ShipmentUpdater

		fakeError error
	}

	makeSubtestData := func(returnErrorForMTOShipment bool, returnErrorForPPMShipment bool, returnErrorForBoatShipment bool, returnErrorForMobileHomeShipment bool) (subtestData subtestDataObjects) {
		mockMTOShipmentUpdater := mocks.MTOShipmentUpdater{}
		subtestData.mockMTOShipmentUpdater = &mockMTOShipmentUpdater

		mockPPMShipmentUpdater := mocks.PPMShipmentUpdater{}
		subtestData.mockPPMShipmentUpdater = &mockPPMShipmentUpdater

		mockBoatShipmentUpdater := mocks.BoatShipmentUpdater{}
		subtestData.mockBoatShipmentUpdater = &mockBoatShipmentUpdater

		mockMobileHomeShipmentUpdater := mocks.MobileHomeShipmentUpdater{}
		subtestData.mockMobileHomeShipmentUpdater = &mockMobileHomeShipmentUpdater

		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater, subtestData.mockBoatShipmentUpdater, subtestData.mockMobileHomeShipmentUpdater)

		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater, subtestData.mockBoatShipmentUpdater, subtestData.mockMobileHomeShipmentUpdater)

		if returnErrorForMTOShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Pickup date missing")

			subtestData.mockMTOShipmentUpdater.
				On(
					updateMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string")).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockMTOShipmentUpdater.
				On(
					updateMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string")).
				Return(
					&models.MTOShipment{
						ID: uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
					}, nil)
		}

		if returnErrorForPPMShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Invalid input found while validating the PPM shipment.")

			subtestData.mockPPMShipmentUpdater.
				On(
					updatePPMShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.PPMShipment"),
					mock.AnythingOfType("uuid.UUID"),
				).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockPPMShipmentUpdater.
				On(
					updatePPMShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.PPMShipment"),
					mock.AnythingOfType("uuid.UUID"),
					mock.AnythingOfType("string"),
				).
				Return(
					func(_ appcontext.AppContext, ship *models.PPMShipment, _ uuid.UUID) *models.PPMShipment {
						// Mimicking how the PPMShipment updater actually returns a new pointer so that we can test
						// a bit more realistically while still using mocks.
						updatedShip := *ship

						return &updatedShip
					},
					func(_ appcontext.AppContext, _ *models.PPMShipment, _ uuid.UUID) error {
						return nil
					},
				)
		}

		if returnErrorForBoatShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Invalid input found while validating the Boat shipment.")

			subtestData.mockBoatShipmentUpdater.
				On(
					updateBoatShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.BoatShipment"),
					mock.AnythingOfType("uuid.UUID"),
				).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockBoatShipmentUpdater.
				On(
					updateBoatShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.BoatShipment"),
					mock.AnythingOfType("uuid.UUID"),
					mock.AnythingOfType("string"),
				).
				Return(
					func(_ appcontext.AppContext, ship *models.BoatShipment, _ uuid.UUID) *models.BoatShipment {
						updatedShip := *ship

						return &updatedShip
					},
					func(_ appcontext.AppContext, _ *models.BoatShipment, _ uuid.UUID) error {
						return nil
					},
				)
		}
		if returnErrorForMobileHomeShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Invalid input found while validating the Mobile Home shipment.")

			subtestData.mockMobileHomeShipmentUpdater.
				On(
					updateMobileHomeShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MobileHome"),
					mock.AnythingOfType("uuid.UUID"),
				).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockMobileHomeShipmentUpdater.
				On(
					updateMobileHomeShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MobileHome"),
					mock.AnythingOfType("uuid.UUID"),
					mock.AnythingOfType("string"),
				).
				Return(
					func(_ appcontext.AppContext, ship *models.MobileHome, _ uuid.UUID) *models.MobileHome {
						updatedShip := *ship

						return &updatedShip
					},
					func(_ appcontext.AppContext, _ *models.MobileHome, _ uuid.UUID) error {
						return nil
					},
				)
		}

		return subtestData
	}

	suite.Run("Returns an InvalidInputError if there is an error with the shipment info that was input", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(false, false, false, false)

		shipment := factory.BuildMTOShipment(appCtx.DB(), nil, nil)

		// Set invalid data, can't pass in blank to the generator above (it'll default to HHG if blank) so we're setting it afterward.
		shipment.ShipmentType = ""

		updatedShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test")

		suite.Nil(updatedShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the shipment.")
	})

	shipmentTypeTestCases := []models.MTOShipmentType{

		models.MTOShipmentTypeHHG,
		models.MTOShipmentTypeHHGIntoNTSDom,
		models.MTOShipmentTypeHHGOutOfNTSDom,
		models.MTOShipmentTypePPM,
	}

	for _, shipmentType := range shipmentTypeTestCases {
		shipmentType := shipmentType

		suite.Run(fmt.Sprintf("Calls necessary service objects for %s shipments", shipmentType), func() {
			appCtx := suite.AppContextForTest()

			subtestData := makeSubtestData(false, false, false, false)

			var shipment models.MTOShipment

			isPPMShipment := shipmentType == models.MTOShipmentTypePPM

			if isPPMShipment {
				ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, nil)

				shipment = ppmShipment.Shipment
			} else {
				shipment = factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
					{
						Model: models.MTOShipment{
							ID:           uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
							ShipmentType: shipmentType,
						},
					},
				}, nil)
			}

			eTag := etag.GenerateEtag(shipment.UpdatedAt)

			// Need to start a transaction so we can assert the call with the correct appCtx
			err := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
				mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(txAppCtx, &shipment, eTag, "test")

				suite.NoError(err)
				suite.NotNil(mtoShipment)

				subtestData.mockMTOShipmentUpdater.AssertCalled(
					suite.T(),
					updateMTOShipmentMethodName,
					txAppCtx,
					&shipment,
					eTag,
					"test",
				)

				if isPPMShipment {
					subtestData.mockPPMShipmentUpdater.AssertCalled(
						suite.T(),
						updatePPMShipmentMethodName,
						txAppCtx,
						shipment.PPMShipment,
						uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
					)
				} else {
					subtestData.mockPPMShipmentUpdater.AssertNotCalled(
						suite.T(),
						updatePPMShipmentMethodName,
						mock.AnythingOfType("*appcontext.appContext"),
						mock.AnythingOfType("*models.PPMShipment"),
						mock.AnythingOfType("uuid.UUID"),
						mock.AnythingOfType("string"),
					)
				}

				return nil
			})

			suite.NoError(err) // just making golangci-lint happy
		})
	}

	suite.Run("Sets MTOShipment info on PPMShipment and updated PPMShipment back on MTOShipment", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(false, false, false, false)

		ppmShipment := factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					HasProGear:            models.BoolPointer(true),
					ProGearWeight:         models.PoundPointer(unit.Pound(1900)),
					SpouseProGearWeight:   models.PoundPointer(unit.Pound(300)),
					AdvanceAmountReceived: nil,
					HasReceivedAdvance:    models.BoolPointer(false),
				},
			},
		}, nil)
		shipment := ppmShipment.Shipment

		// set new field to update
		shipment.PPMShipment.HasProGear = models.BoolPointer(false)
		shipment.PPMShipment.AdvanceAmountReceived = models.CentPointer(unit.Cents(55000))
		shipment.PPMShipment.HasReceivedAdvance = models.BoolPointer(true)

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test")

		suite.NoError(err)
		suite.NotNil(mtoShipment)

		// check we got the latest version back
		suite.NotEqual(&ppmShipment, mtoShipment.PPMShipment)
		suite.False(*mtoShipment.PPMShipment.HasProGear)
		suite.Equal(*mtoShipment.PPMShipment.AdvanceAmountReceived, *shipment.PPMShipment.AdvanceAmountReceived)
		suite.True(*mtoShipment.PPMShipment.HasReceivedAdvance)
	})

	suite.Run("Sets MTOShipment info on BoatShipment and updated BoatShipment back on MTOShipment", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(false, false, false, false)

		boatShipment := factory.BuildBoatShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.BoatShipment{
					Type:           models.BoatShipmentTypeHaulAway,
					Year:           models.IntPointer(2000),
					Make:           models.StringPointer("Boat Make"),
					Model:          models.StringPointer("Boat Model"),
					LengthInInches: models.IntPointer(300),
					WidthInInches:  models.IntPointer(108),
					HeightInInches: models.IntPointer(72),
					HasTrailer:     models.BoolPointer(true),
					IsRoadworthy:   models.BoolPointer(false),
				},
			},
		}, nil)
		shipment := boatShipment.Shipment

		// set new field to update
		shipment.BoatShipment.Year = models.IntPointer(1991)
		shipment.BoatShipment.LengthInInches = models.IntPointer(20)
		shipment.BoatShipment.HasTrailer = models.BoolPointer(false)

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test")

		suite.NoError(err)

		suite.NotNil(mtoShipment)

		// check that the BoatShipment has the expected MTOShipment fields
		suite.Equal(mtoShipment.ID, mtoShipment.BoatShipment.ShipmentID)
		suite.Equal(*mtoShipment, mtoShipment.BoatShipment.Shipment)

		// check we got the latest version back
		suite.NotEqual(&boatShipment, mtoShipment.BoatShipment)
		suite.Equal(*mtoShipment.BoatShipment.Year, *shipment.BoatShipment.Year)
		suite.Equal(*mtoShipment.BoatShipment.LengthInInches, *shipment.BoatShipment.LengthInInches)
		suite.False(*mtoShipment.BoatShipment.HasTrailer)
	})

	serviceObjectErrorTestCases := map[string]struct {
		shipmentType                     models.MTOShipmentType
		returnErrorForMTOShipment        bool
		returnErrorForPPMShipment        bool
		returnErrorForBoatShipment       bool
		returnErrorForMobileHomeShipment bool
	}{
		"error updating MTOShipment": {
			shipmentType:               models.MTOShipmentTypeHHG,
			returnErrorForMTOShipment:  true,
			returnErrorForPPMShipment:  false,
			returnErrorForBoatShipment: false,
		},
		"error updating PPMShipment": {
			shipmentType:               models.MTOShipmentTypePPM,
			returnErrorForMTOShipment:  false,
			returnErrorForPPMShipment:  true,
			returnErrorForBoatShipment: false,
		},
		"error updating BoatShipment": {
			shipmentType:               models.MTOShipmentTypeBoatHaulAway,
			returnErrorForMTOShipment:  false,
			returnErrorForPPMShipment:  false,
			returnErrorForBoatShipment: true,
		},
	}

	for name, tc := range serviceObjectErrorTestCases {
		name := name
		tc := tc

		suite.Run(fmt.Sprintf("Returns transaction error if there is an %s", name), func() {
			appCtx := suite.AppContextForTest()

			subtestData := makeSubtestData(tc.returnErrorForMTOShipment, tc.returnErrorForPPMShipment, tc.returnErrorForBoatShipment, tc.returnErrorForMobileHomeShipment)

			var shipment models.MTOShipment

			isBoatShipmentType := tc.shipmentType == models.MTOShipmentTypeBoatHaulAway || tc.shipmentType == models.MTOShipmentTypeBoatTowAway

			if tc.shipmentType == models.MTOShipmentTypePPM {
				ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, nil)

				shipment = ppmShipment.Shipment
			} else if isBoatShipmentType {
				boatShipment := factory.BuildBoatShipment(appCtx.DB(), nil, nil)

				shipment = boatShipment.Shipment
			} else if tc.shipmentType == models.MTOShipmentTypeMobileHome {
				mobileHomeShipment := factory.BuildMobileHomeShipment(appCtx.DB(), nil, nil)

				shipment = mobileHomeShipment.Shipment
			} else {
				shipment = factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
					{
						Model: models.MTOShipment{
							ShipmentType: tc.shipmentType,
						},
					},
				}, nil)
			}

			mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test")

			suite.Nil(mtoShipment)

			suite.Error(err)
			suite.Equal(subtestData.fakeError, err)
		})
	}

	suite.Run("Returns error early if MTOShipment can't be updated", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(true, false, false, false)

		shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, eTag, "test")

		suite.Nil(mtoShipment)

		suite.Error(err)
		suite.Equal(subtestData.fakeError, err)

		subtestData.mockMTOShipmentUpdater.AssertCalled(
			suite.T(),
			updateMTOShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			&shipment,
			eTag,
			"test",
		)

		subtestData.mockPPMShipmentUpdater.AssertNotCalled(
			suite.T(),
			updatePPMShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
			mock.AnythingOfType("uuid.UUID"),
		)
	})
}
