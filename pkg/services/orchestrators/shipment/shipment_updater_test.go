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
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
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

	type subtestDataObjects struct {
		mockMTOShipmentUpdater        *mocks.MTOShipmentUpdater
		mockMtoServiceItemCreator     *mocks.MTOServiceItemCreator
		mockPPMShipmentUpdater        *mocks.PPMShipmentUpdater
		mockBoatShipmentUpdater       *mocks.BoatShipmentUpdater
		mockMobileHomeShipmentUpdater *mocks.MobileHomeShipmentUpdater
		shipmentUpdaterOrchestrator   services.ShipmentUpdater

		fakeError error
	}

	planner := &routemocks.Planner{}
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	builder := query.NewQueryBuilder()
	mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	makeSubtestData := func(returnErrorForMTOShipment bool, returnErrorForPPMShipment bool, returnErrorForBoatShipment bool, returnErrorForMobileHomeShipment bool) (subtestData subtestDataObjects) {
		mockMTOShipmentUpdater := mocks.MTOShipmentUpdater{}
		subtestData.mockMTOShipmentUpdater = &mockMTOShipmentUpdater

		mockMtoServiceItemCreator := mocks.MTOServiceItemCreator{}
		subtestData.mockMtoServiceItemCreator = &mockMtoServiceItemCreator

		mockPPMShipmentUpdater := mocks.PPMShipmentUpdater{}
		subtestData.mockPPMShipmentUpdater = &mockPPMShipmentUpdater

		mockBoatShipmentUpdater := mocks.BoatShipmentUpdater{}
		subtestData.mockBoatShipmentUpdater = &mockBoatShipmentUpdater

		mockMobileHomeShipmentUpdater := mocks.MobileHomeShipmentUpdater{}
		subtestData.mockMobileHomeShipmentUpdater = &mockMobileHomeShipmentUpdater

		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater, subtestData.mockBoatShipmentUpdater, subtestData.mockMobileHomeShipmentUpdater, mtoServiceItemCreator)
		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater, subtestData.mockBoatShipmentUpdater, subtestData.mockMobileHomeShipmentUpdater, mtoServiceItemCreator)

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
	makeServiceItemSubtestData := func() (subtestData subtestDataObjects) {
		mockMTOShipmentUpdater := mocks.MTOShipmentUpdater{}
		subtestData.mockMTOShipmentUpdater = &mockMTOShipmentUpdater

		mockMtoServiceItemCreator := mocks.MTOServiceItemCreator{}
		subtestData.mockMtoServiceItemCreator = &mockMtoServiceItemCreator

		mockPPMShipmentUpdater := mocks.PPMShipmentUpdater{}
		subtestData.mockPPMShipmentUpdater = &mockPPMShipmentUpdater

		mockBoatShipmentUpdater := mocks.BoatShipmentUpdater{}
		subtestData.mockBoatShipmentUpdater = &mockBoatShipmentUpdater

		mockMobileHomeShipmentUpdater := mocks.MobileHomeShipmentUpdater{}
		subtestData.mockMobileHomeShipmentUpdater = &mockMobileHomeShipmentUpdater

		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater, subtestData.mockBoatShipmentUpdater, subtestData.mockMobileHomeShipmentUpdater, subtestData.mockMtoServiceItemCreator)

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
		models.MTOShipmentTypeHHGIntoNTS,
		models.MTOShipmentTypeHHGOutOfNTS,
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

	suite.Run("Updating weight will update the estimated price of service items", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeServiceItemSubtestData()

		estimatedWeight := unit.Pound(2000)
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ID:                   uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
					Status:               models.MTOShipmentStatusApproved,
					ShipmentType:         models.MTOShipmentTypeHHG,
					PrimeEstimatedWeight: &estimatedWeight,
				},
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
		}, nil)

		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)

		startDate := time.Now().AddDate(-1, 0, 0)
		endDate := startDate.AddDate(1, 1, 1)
		reason := "lorem ipsum"

		testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Name:                 "Test Contract Year",
					EscalationCompounded: 1.125,
					StartDate:            startDate,
					EndDate:              endDate,
				},
			})

		testdatagen.FetchOrMakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2025, time.March, 17, 0, 0, 0, 0, time.UTC),
			},
		})

		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		requestedPickupDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)

		serviceItemFSC := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeFSC,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		shipment.MTOServiceItems = append(shipment.MTOServiceItems, serviceItemFSC)
		suite.MustSave(&shipment)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		subtestData.mockMtoServiceItemCreator.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(800, nil)

		returnCents := unit.Cents(123)

		subtestData.mockMtoServiceItemCreator.On("FindEstimatedPrice",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(returnCents, nil)

		subtestData.mockMTOShipmentUpdater.
			On(
				updateMTOShipmentMethodName,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.MTOShipment"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
			Return(
				&models.MTOShipment{
					ID:                   uuid.Must(uuid.FromString("a5e95c1d-97c3-4f79-8097-c12dd2557ac7")),
					Status:               models.MTOShipmentStatusApproved,
					ShipmentType:         models.MTOShipmentTypeHHG,
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &requestedPickupDate,
					MTOServiceItems:      models.MTOServiceItems{serviceItemFSC},
					PickupAddress:        &pickupAddress,
					DestinationAddress:   &deliveryAddress,
				}, nil)

		// Need to start a transaction so we can assert the call with the correct appCtx
		err := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(txAppCtx, &shipment, eTag, "test")

			suite.NoError(err)
			suite.NotNil(mtoShipment)

			expectedPrice := unit.Cents(123)
			expectedWeight := unit.Pound(2000)
			suite.Equal(expectedWeight, *mtoShipment.MTOServiceItems[0].EstimatedWeight)
			suite.Equal(expectedPrice, *mtoShipment.MTOServiceItems[0].PricingEstimate)

			return nil
		})

		suite.NoError(err) // just making golangci-lint happy
	})
}
