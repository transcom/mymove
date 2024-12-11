package shipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ShipmentSuite) TestUpdateShipment() {

	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeSubtestData.

	updateMTOShipmentMethodName := "UpdateMTOShipment"
	updatePPMShipmentMethodName := "UpdatePPMShipmentWithDefaultCheck"
	updateBoatShipmentMethodName := "UpdateBoatShipmentWithDefaultCheck"
	updateMobileHomeShipmentMethodName := "UpdateMobileHomeShipmentWithDefaultCheck"
	createMtoServiceItemMethodName := "FindEstimatedPrice"
	updateMtoServiceItemMethodName := "UpdateMTOServiceItemBasic"

	type subtestDataObjects struct {
		mockMTOShipmentUpdater        *mocks.MTOShipmentUpdater
		mockPPMShipmentUpdater        *mocks.PPMShipmentUpdater
		mockBoatShipmentUpdater       *mocks.BoatShipmentUpdater
		mockMobileHomeShipmentUpdater *mocks.MobileHomeShipmentUpdater
		mockMTOServiceItemCreator     *mocks.MTOServiceItemCreator
		mockMTOServiceItemUpdater     *mocks.MTOServiceItemUpdater
		shipmentUpdaterOrchestrator   services.ShipmentUpdater

		fakeError error
	}

	makeSubtestData := func(returnErrorForMTOShipment bool, returnErrorForPPMShipment bool, returnErrorForBoatShipment bool, returnErrorForMobileHomeShipment bool, returnErrorForPricingServiceItem bool, returnErrorForMTOServiceItemUpdate bool, mockedReturnPrice unit.Cents) (subtestData subtestDataObjects) {
		mockMTOShipmentUpdater := mocks.MTOShipmentUpdater{}
		subtestData.mockMTOShipmentUpdater = &mockMTOShipmentUpdater

		mockPPMShipmentUpdater := mocks.PPMShipmentUpdater{}
		subtestData.mockPPMShipmentUpdater = &mockPPMShipmentUpdater

		mockBoatShipmentUpdater := mocks.BoatShipmentUpdater{}
		subtestData.mockBoatShipmentUpdater = &mockBoatShipmentUpdater

		mockMobileHomeShipmentUpdater := mocks.MobileHomeShipmentUpdater{}
		subtestData.mockMobileHomeShipmentUpdater = &mockMobileHomeShipmentUpdater

		mockMTOServiceItemCreator := mocks.MTOServiceItemCreator{}
		mockMTOServiceItemUpdater := mocks.MTOServiceItemUpdater{}
		subtestData.mockMTOServiceItemCreator = &mockMTOServiceItemCreator
		subtestData.mockMTOServiceItemUpdater = &mockMTOServiceItemUpdater

		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater, subtestData.mockBoatShipmentUpdater, subtestData.mockMobileHomeShipmentUpdater)

		var returnCents = mockedReturnPrice

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
		if returnErrorForMTOServiceItemUpdate {
			subtestData.mockMTOServiceItemUpdater.On(
				updateMtoServiceItemMethodName,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.MTOServiceItem"),
				mock.AnythingOfType("string"),
			).Return(
				&models.MTOShipment{
					ID: uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
				}, nil)
		} else {
			subtestData.mockMTOServiceItemUpdater.On(
				updateMtoServiceItemMethodName,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.MTOServiceItem"),
				mock.AnythingOfType("string"),
			).Return(
				&models.MTOShipment{
					ID: uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
				}, nil)
		}
		if returnErrorForPricingServiceItem {
			subtestData.mockMTOServiceItemCreator.On(
				createMtoServiceItemMethodName,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.MTOServiceItem"),
				mock.AnythingOfType("models.MTOShipment"),
			).Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockMTOServiceItemCreator.On(
				createMtoServiceItemMethodName,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.MTOServiceItem"),
				mock.AnythingOfType("models.MTOShipment"),
			).Return(returnCents, nil)
		}

		return subtestData
	}

	suite.Run("Returns an InvalidInputError if there is an error with the shipment info that was input", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(false, false, false, false, false, false, 0)

		shipment := factory.BuildMTOShipment(appCtx.DB(), nil, nil)

		// Set invalid data, can't pass in blank to the generator above (it'll default to HHG if blank) so we're setting it afterward.
		shipment.ShipmentType = ""

		updatedShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test", nil)

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

			subtestData := makeSubtestData(false, false, false, false, false, false, 0)

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
				mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(txAppCtx, &shipment, eTag, "test", nil)

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

		subtestData := makeSubtestData(false, false, false, false, false, false, 0)

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

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test", nil)

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

		subtestData := makeSubtestData(false, false, false, false, false, false, 0)

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

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test", nil)

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

			subtestData := makeSubtestData(tc.returnErrorForMTOShipment, tc.returnErrorForPPMShipment, tc.returnErrorForBoatShipment, tc.returnErrorForMobileHomeShipment, false, false, 0)

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

			mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt), "test", nil)

			suite.Nil(mtoShipment)

			suite.Error(err)
			suite.Equal(subtestData.fakeError, err)
		})
	}

	suite.Run("Returns error early if MTOShipment can't be updated", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(true, false, false, false, false, false, 0)

		shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, eTag, "test", nil)

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

	suite.Run("I hate tests", func() {
		appCtx := suite.AppContextForTest()
		subtestData := makeSubtestData(false, false, false, false, false, false, 0)
		setupTestData := func() models.MTOShipment {
			// Set up data to use for all Origin SIT Service Item tests

			move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
			estimatedPrimeWeight := unit.Pound(6000)
			pickupDate := time.Date(2024, time.July, 31, 12, 0, 0, 0, time.UTC)
			pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
			deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

			mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    pickupAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.PickupAddress,
				},
				{
					Model:    deliveryAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedPrimeWeight,
						RequestedPickupDate:  &pickupDate,
					},
				},
			}, nil)

			return mtoShipment
		}
		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reServiceCodeDPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)
		reServiceCodeDDP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)
		reServiceCodeDUPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)
		reServiceCodeDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)

		startDate := time.Now().AddDate(-1, 0, 0)
		endDate := startDate.AddDate(1, 1, 1)
		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		reason := "lorem ipsum"

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		contractYear := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Name:                 "Test Contract Year",
					EscalationCompounded: 1.125,
					StartDate:            startDate,
					EndDate:              endDate,
				},
			})

		serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "945",
					ServicesSchedule: 1,
				},
			})

		serviceAreaDest := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "503",
					ServicesSchedule: 1,
				},
			})

		serviceAreaPriceDOP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDOP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(1234),
		}

		serviceAreaPriceDPK := models.ReDomesticOtherPrice{
			ContractID:   contractYear.Contract.ID,
			ServiceID:    reServiceCodeDPK.ID,
			IsPeakPeriod: true,
			Schedule:     1,
			PriceCents:   unit.Cents(121),
		}

		serviceAreaPriceDDP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDDP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceAreaDest.ID,
			PriceCents:            unit.Cents(482),
		}

		serviceAreaPriceDUPK := models.ReDomesticOtherPrice{
			ContractID:   contractYear.Contract.ID,
			ServiceID:    reServiceCodeDUPK.ID,
			IsPeakPeriod: true,
			Schedule:     1,
			PriceCents:   unit.Cents(945),
		}

		serviceAreaPriceDLH := models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			WeightLower:           500,
			WeightUpper:           10000,
			MilesLower:            1,
			MilesUpper:            10000,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceMillicents:       unit.Millicents(482),
		}

		serviceAreaPriceDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDSH.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(999),
		}

		testdatagen.MakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2025, time.March, 17, 0, 0, 0, 0, time.UTC),
			},
		})

		suite.MustSave(&serviceAreaPriceDOP)
		suite.MustSave(&serviceAreaPriceDPK)
		suite.MustSave(&serviceAreaPriceDDP)
		suite.MustSave(&serviceAreaPriceDUPK)
		suite.MustSave(&serviceAreaPriceDLH)
		suite.MustSave(&serviceAreaPriceDSH)

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceArea,
				Zip3:                "945",
			},
		})

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceAreaDest,
				Zip3:                "503",
			},
		})

		shipment := setupTestData()
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		serviceItemDOP := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDOP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDPK := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDDP := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDDP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDUPK := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDUPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDLH := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDLH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDSH := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDSH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemFSC := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeFSC,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, serviceItemDDP, serviceItemDLH, serviceItemDOP, serviceItemDPK, serviceItemDSH, serviceItemDUPK, serviceItemFSC)
		suite.MustSave(&shipment)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		subtestData.mockMTOShipmentUpdater.On(
			"FetchServiceItemPrice",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.MTOServiceItem"),
			mock.AnythingOfType("*models.MTOShipment")).
			Return(shipment, nil)

		suite.NotEmpty(shipment.MTOServiceItems)

		err := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(txAppCtx, &shipment, eTag, "test", nil)

			suite.NoError(err)
			suite.NotNil(mtoShipment)

			fmt.Printf("%+v", mtoShipment)

			suite.NotEmpty(mtoShipment.MTOServiceItems)

			subtestData.mockMTOShipmentUpdater.AssertCalled(
				suite.T(),
				updateMTOShipmentMethodName,
				txAppCtx,
				&shipment,
				eTag,
				"test",
			)

			return nil
		})

		for _, serviceItem := range shipment.MTOServiceItems {
			suite.Equal(unit.Cents(1), serviceItem.PricingEstimate)
		}
		suite.NoError(err)
	})
}
