package primeapiv3

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapiv3/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mobilehomeshipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	fetcher := fetch.NewFetcher(builder)
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV2(builder, fetcher, moveRouter, addressCreator)
	ppmEstimator := mocks.PPMEstimator{}
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
	boatShipmentCreator := boatshipment.NewBoatShipmentCreator()
	mobileHomeShipmentCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
	shipmentRouter := mtoshipment.NewShipmentRouter()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	vLocationServices := address.NewVLocation()

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}
	mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		builder,
		mtoServiceItemCreator,
		moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
	)
	mockSender := suite.TestNotificationSender()
	waf := entitlements.NewWeightAllotmentFetcher()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(), waf, mockSender)
	shipmentCreator := shipmentorchestrator.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, boatShipmentCreator, mobileHomeShipmentCreator, shipmentRouter, moveTaskOrderUpdater, moveWeights)
	mockCreator := mocks.ShipmentCreator{}

	var pickupAddress primev3messages.Address
	var secondaryPickupAddress primev3messages.Address
	var tertiaryPickupAddress primev3messages.Address
	var destinationAddress primev3messages.Address
	var ppmDestinationAddress primev3messages.PPMDestinationAddress
	var secondaryDestinationAddress primev3messages.Address
	var tertiaryDestinationAddress primev3messages.Address
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
	boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
	mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()
	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)

	setupAddresses := func() {
		// Make stubbed addresses just to collect address data for payload
		newAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		pickupAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		secondaryPickupAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		tertiaryPickupAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		newAddress = factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
		destinationAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		secondaryDestinationAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		tertiaryDestinationAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
	}

	setupTestData := func(boatFeatureFlag bool, ubFeatureFlag bool) (CreateMTOShipmentHandler, models.Move) {
		vLocationServices := address.NewVLocation()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		handlerConfig := suite.HandlerConfig()
		expectedFeatureFlag := services.FeatureFlag{
			Key:   "",
			Match: true,
		}
		if !boatFeatureFlag {
			expectedFeatureFlag.Key = "boat"
			mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}

			mockFeatureFlagFetcher.On("GetBooleanFlag",
				mock.Anything,
				mock.Anything,
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(expectedFeatureFlag, nil)
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(services.FeatureFlag{}, errors.New("Some error"))
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		}
		if ubFeatureFlag {
			expectedFeatureFlag.Key = "unaccompanied_baggage"
			expectedFeatureFlag.Match = true
			mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}

			mockFeatureFlagFetcher.On("GetBooleanFlag",
				mock.Anything,                 // context.Context
				mock.Anything,                 // *zap.Logger
				mock.AnythingOfType("string"), // entityID (userID)
				mock.AnythingOfType("string"), // key
				mock.Anything,                 // flagContext (map[string]string)
			).Return(expectedFeatureFlag, nil)
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(expectedFeatureFlag, nil)
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		} else {
			mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
			mockFeatureFlagFetcher.On("GetBooleanFlag",
				mock.Anything,                 // context.Context
				mock.Anything,                 // *zap.Logger
				mock.AnythingOfType("string"), // entityID (userID)
				mock.AnythingOfType("string"), // key
				mock.Anything,                 // flagContext (map[string]string)
			).Return(services.FeatureFlag{}, errors.New("Some error"))
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(services.FeatureFlag{}, errors.New("Some error"))

			mockFeatureFlagFetcher.On("GetBooleanFlag",
				mock.Anything,
				mock.Anything,
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(func(_ context.Context, _ *zap.Logger, _ string, key string, flagContext map[string]string) (services.FeatureFlag, error) {
				return services.FeatureFlag{}, errors.New("Some error")
			})
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		}
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			mtoChecker,
			vLocationServices,
		}

		setupAddresses()
		return handler, move
	}

	setupTestDataWithoutFF := func() (CreateMTOShipmentHandler, models.Move) {
		vLocationServices := address.NewVLocation()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		handler := CreateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentCreator,
			mtoChecker,
			vLocationServices,
		}

		setupAddresses()
		return handler, move
	}

	suite.Run("Successful POST - Integration Test", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Successful submission, status should be APPROVED
		handler, move := setupTestData(false, true)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createMTOShipmentPayload.Validate(strfmt.Default))

		// check that the mto shipment status is submitted
		suite.Require().Equal(createMTOShipmentPayload.Status, primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.PrimeEstimatedWeight, params.Body.PrimeEstimatedWeight)
	})

	suite.Run("Successful POST - Integration Test - Unaccompanied Baggage", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Successful submission, status should be APPROVED

		suite.T().Setenv("FEATURE_FLAG_UNACCOMPANIED_BAGGAGE", "true") // Set to true in order to test UB shipments can be created with UB flag on

		handler, move := setupTestData(true, true)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createMTOShipmentPayload.Validate(strfmt.Default))

		// check that the mto shipment status is Submitted
		suite.Require().Equal(primev3messages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE, createMTOShipmentPayload.ShipmentType, "MTO Shipment type should be unaccompanied baggage")
		suite.Require().Equal(primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, createMTOShipmentPayload.Status, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.PrimeEstimatedWeight, params.Body.PrimeEstimatedWeight)
	})

	suite.Run("Successful POST/PATCH - Integration Test (PPM)", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:      Create a PPM shipment on an available move
		// Expected:   Successful submission, status should be APPROVED
		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		counselorRemarks := "Some counselor remarks"
		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		sitExpected := true
		sitLocation := primev3messages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1500)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3200)
		hasProGear := true
		proGearWeight := unit.Pound(400)
		spouseProGearWeight := unit.Pound(250)
		estimatedIncentive := 123456
		sitEstimatedCost := 67500

		address1 := models.Address{
			StreetAddress1: "some address",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}
		address2 := models.Address{
			StreetAddress1: "some address",
			City:           "Scott Afb",
			State:          "IL",
			PostalCode:     "62225",
		}

		address3 := models.Address{
			StreetAddress1: "some address",
			City:           "Suffolk",
			State:          "VA",
			PostalCode:     "23435",
		}

		expectedPickupAddress := address1
		pickupAddress = primev3messages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedSecondaryPickupAddress := address2
		secondaryPickupAddress = primev3messages.Address{
			City:           &expectedSecondaryPickupAddress.City,
			PostalCode:     &expectedSecondaryPickupAddress.PostalCode,
			State:          &expectedSecondaryPickupAddress.State,
			StreetAddress1: &expectedSecondaryPickupAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryPickupAddress.StreetAddress2,
			StreetAddress3: expectedSecondaryPickupAddress.StreetAddress3,
		}

		expectedTertiaryPickupAddress := address3
		tertiaryDestinationAddress = primev3messages.Address{
			City:           &expectedTertiaryPickupAddress.City,
			PostalCode:     &expectedTertiaryPickupAddress.PostalCode,
			State:          &expectedTertiaryPickupAddress.State,
			StreetAddress1: &expectedTertiaryPickupAddress.StreetAddress1,
			StreetAddress2: expectedTertiaryPickupAddress.StreetAddress2,
			StreetAddress3: expectedTertiaryPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := address1
		destinationAddress = primev3messages.Address{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}
		ppmDestinationAddress = primev3messages.PPMDestinationAddress{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}

		expectedSecondaryDestinationAddress := address2
		secondaryDestinationAddress = primev3messages.Address{
			City:           &expectedSecondaryDestinationAddress.City,
			PostalCode:     &expectedSecondaryDestinationAddress.PostalCode,
			State:          &expectedSecondaryDestinationAddress.State,
			StreetAddress1: &expectedSecondaryDestinationAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryDestinationAddress.StreetAddress2,
			StreetAddress3: expectedSecondaryDestinationAddress.StreetAddress3,
		}

		expectedTertiaryDestinationAddress := address3
		tertiaryDestinationAddress = primev3messages.Address{
			City:           &expectedTertiaryDestinationAddress.City,
			PostalCode:     &expectedTertiaryDestinationAddress.PostalCode,
			State:          &expectedTertiaryDestinationAddress.State,
			StreetAddress1: &expectedTertiaryDestinationAddress.StreetAddress1,
			StreetAddress2: expectedTertiaryDestinationAddress.StreetAddress2,
			StreetAddress3: expectedTertiaryDestinationAddress.StreetAddress3,
		}

		// Need a logged in user
		lgu := uuid.Must(uuid.NewV4()).String()
		user := models.User{
			OktaID:    lgu,
			OktaEmail: "email@example.com",
		}
		suite.MustSave(&user)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          user.ID,
			IDToken:         "fake token",
			Roles:           roles.Roles{},
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req.WithContext(ctx),
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:  handlers.FmtUUID(move.ID),
				ShipmentType:     primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypePPM),
				CounselorRemarks: &counselorRemarks,
				PpmShipment: &primev3messages.CreatePPMShipment{
					ExpectedDepartureDate:  handlers.FmtDate(expectedDepartureDate),
					PickupAddress:          struct{ primev3messages.Address }{pickupAddress},
					SecondaryPickupAddress: struct{ primev3messages.Address }{secondaryPickupAddress},
					DestinationAddress: struct {
						primev3messages.PPMDestinationAddress
					}{ppmDestinationAddress},
					SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
					SitExpected:                 &sitExpected,
					SitLocation:                 &sitLocation,
					SitEstimatedWeight:          handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:       handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate:   handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:             handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                  &hasProGear,
					ProGearWeight:               handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:         handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createdShipment := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createdShipment.Validate(strfmt.Default))

		createdPPM := createdShipment.PpmShipment

		suite.Equal(move.ID.String(), createdShipment.MoveTaskOrderID.String())
		suite.Equal(primev3messages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, createdShipment.Status)
		suite.Equal(&counselorRemarks, createdShipment.CounselorRemarks)

		suite.Equal(createdShipment.ID.String(), createdPPM.ShipmentID.String())
		suite.Equal(primev3messages.PPMShipmentStatusSUBMITTED, createdPPM.Status)
		suite.Equal(handlers.FmtDatePtr(&expectedDepartureDate), createdPPM.ExpectedDepartureDate)
		suite.Equal(address1.PostalCode, *createdPPM.PickupAddress.PostalCode)
		suite.Equal(address1.PostalCode, *createdPPM.DestinationAddress.PostalCode)
		suite.Equal(address2.PostalCode, *createdPPM.SecondaryPickupAddress.PostalCode)
		suite.Equal(address2.PostalCode, *createdPPM.SecondaryDestinationAddress.PostalCode)
		suite.Equal(&sitExpected, createdPPM.SitExpected)
		suite.Equal(&sitLocation, createdPPM.SitLocation)
		suite.True(*models.BoolPointer(*createdPPM.HasSecondaryPickupAddress))
		suite.True(*models.BoolPointer(*createdPPM.HasSecondaryDestinationAddress))
		suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), createdPPM.SitEstimatedWeight)
		suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), createdPPM.SitEstimatedEntryDate)
		suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), createdPPM.SitEstimatedDepartureDate)
		suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), createdPPM.EstimatedWeight)
		suite.Equal(handlers.FmtBool(hasProGear), createdPPM.HasProGear)
		suite.Equal(handlers.FmtPoundPtr(&proGearWeight), createdPPM.ProGearWeight)
		suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), createdPPM.SpouseProGearWeight)
		suite.Equal(int64(estimatedIncentive), *createdPPM.EstimatedIncentive)
		suite.Equal(int64(sitEstimatedCost), *createdPPM.SitEstimatedCost)

		// ************
		// PATCH TESTS
		// ************
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Times(4)

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", createdPPM.ShipmentID.String()), nil)

		var mtoShipment models.MTOShipment
		err := suite.DB().Find(&mtoShipment, createdPPM.ShipmentID)
		suite.NoError(err)
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: createdPPM.ShipmentID,
			IfMatch:       eTag,
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			ShipmentType: primev3messages.MTOShipmentTypePPM,
		}
		// *************************************************************************************
		// *************************************************************************************
		// Run it without any flags, no deletes should occur for secondary addresses
		// *************************************************************************************
		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			HasProGear: &hasProGear,
		}

		// Validate incoming payload
		suite.NoError(patchParams.Body.Validate(strfmt.Default))

		patchResponse := patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment := okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM := updatedShipment.PpmShipment
		suite.NotNil(updatedPPM.SecondaryPickupAddress)
		suite.NotNil(updatedPPM.SecondaryDestinationAddress)
		suite.True(*models.BoolPointer(*updatedPPM.HasSecondaryPickupAddress))
		suite.True(*models.BoolPointer(*updatedPPM.HasSecondaryDestinationAddress))

		// *************************************************************************************
		// *************************************************************************************
		// Run it second time, but really delete secondary addresses with has flags set to false
		// *************************************************************************************
		err = suite.DB().Find(&mtoShipment, createdPPM.ShipmentID)
		suite.NoError(err)
		eTag = etag.GenerateEtag(mtoShipment.UpdatedAt)
		patchParams.IfMatch = eTag
		patchParams.MtoShipmentID = createdPPM.ShipmentID
		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			HasProGear:                     &hasProGear,
			HasSecondaryPickupAddress:      models.BoolPointer(false),
			HasSecondaryDestinationAddress: models.BoolPointer(false),
		}
		suite.NoError(patchParams.Body.Validate(strfmt.Default))
		patchResponse = patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse = patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment = okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM = updatedShipment.PpmShipment
		// secondary should be all nils
		suite.Nil(updatedPPM.SecondaryPickupAddress)
		suite.Nil(updatedPPM.SecondaryDestinationAddress)
		suite.False(*models.BoolPointer(*updatedPPM.HasSecondaryPickupAddress))
		suite.False(*models.BoolPointer(*updatedPPM.HasSecondaryDestinationAddress))

		// *************************************************************************************
		// *************************************************************************************
		// Run it third time, but really add secondary addresses with has flags set to true
		// *************************************************************************************
		eTag = etag.GenerateEtag(time.Time(updatedShipment.UpdatedAt))
		patchParams.IfMatch = eTag
		patchParams.MtoShipmentID = createdPPM.ShipmentID

		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			HasProGear:                     &hasProGear,
			HasSecondaryPickupAddress:      models.BoolPointer(true),
			HasSecondaryDestinationAddress: models.BoolPointer(true),
			SecondaryPickupAddress:         struct{ primev3messages.Address }{secondaryPickupAddress},
			SecondaryDestinationAddress:    struct{ primev3messages.Address }{secondaryDestinationAddress},
		}

		// Validate incoming payload
		suite.NoError(patchParams.Body.Validate(strfmt.Default))

		patchResponse = patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse = patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment = okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM = updatedShipment.PpmShipment
		suite.NotNil(updatedPPM.SecondaryPickupAddress)
		suite.NotNil(updatedPPM.SecondaryDestinationAddress)
		suite.True(*models.BoolPointer(*updatedPPM.HasSecondaryPickupAddress))
		suite.True(*models.BoolPointer(*updatedPPM.HasSecondaryDestinationAddress))

		// *************************************************************************************
		// *************************************************************************************
		// Run it fourth time, but really add tertiary addresses with has flags set to true
		// *************************************************************************************
		eTag = etag.GenerateEtag(time.Time(updatedShipment.UpdatedAt))
		patchParams.IfMatch = eTag
		patchParams.MtoShipmentID = updatedPPM.ShipmentID
		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			HasProGear:                    &hasProGear,
			HasTertiaryPickupAddress:      models.BoolPointer(true),
			HasTertiaryDestinationAddress: models.BoolPointer(true),
			TertiaryPickupAddress:         struct{ primev3messages.Address }{tertiaryPickupAddress},
			TertiaryDestinationAddress:    struct{ primev3messages.Address }{tertiaryDestinationAddress},
		}

		// Validate incoming payload
		suite.NoError(patchParams.Body.Validate(strfmt.Default))

		patchResponse = patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse = patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment = okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM = updatedShipment.PpmShipment
		suite.NotNil(updatedPPM.TertiaryPickupAddress)
		suite.NotNil(updatedPPM.TertiaryDestinationAddress)
		suite.True(*models.BoolPointer(*updatedPPM.HasTertiaryPickupAddress))
		suite.True(*models.BoolPointer(*updatedPPM.HasTertiaryDestinationAddress))
	})

	suite.Run("Successful POST/PATCH - Integration Test (PPM) - Destination address street 1 OPTIONAL", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:      Create a PPM shipment on an available move
		// Expected:   Successful submission, status should be SUBMITTED
		handler, move := setupTestData(true, true)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		counselorRemarks := "Some counselor remarks"
		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		sitExpected := true
		sitLocation := primev3messages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1500)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3200)
		hasProGear := true
		proGearWeight := unit.Pound(400)
		spouseProGearWeight := unit.Pound(250)
		estimatedIncentive := 123456
		sitEstimatedCost := 67500

		address1 := models.Address{
			StreetAddress1: "some address",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}
		addressWithEmptyStreet1 := models.Address{
			StreetAddress1: "",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}

		expectedPickupAddress := address1
		pickupAddress = primev3messages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := address1
		destinationAddress = primev3messages.Address{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}
		ppmDestinationAddress = primev3messages.PPMDestinationAddress{
			City:           &addressWithEmptyStreet1.City,
			PostalCode:     &addressWithEmptyStreet1.PostalCode,
			State:          &addressWithEmptyStreet1.State,
			StreetAddress1: &addressWithEmptyStreet1.StreetAddress1,
			StreetAddress2: addressWithEmptyStreet1.StreetAddress2,
			StreetAddress3: addressWithEmptyStreet1.StreetAddress3,
		}

		// Need a logged in user
		lgu := uuid.Must(uuid.NewV4()).String()
		user := models.User{
			OktaID:    lgu,
			OktaEmail: "email@example.com",
		}
		suite.MustSave(&user)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          user.ID,
			IDToken:         "fake token",
			Roles:           roles.Roles{},
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req.WithContext(ctx),
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:  handlers.FmtUUID(move.ID),
				ShipmentType:     primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypePPM),
				CounselorRemarks: &counselorRemarks,
				PpmShipment: &primev3messages.CreatePPMShipment{
					ExpectedDepartureDate: handlers.FmtDate(expectedDepartureDate),
					PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
					DestinationAddress: struct {
						primev3messages.PPMDestinationAddress
					}{ppmDestinationAddress},
					SitExpected:               &sitExpected,
					SitLocation:               &sitLocation,
					SitEstimatedWeight:        handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:     handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate: handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:           handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                &hasProGear,
					ProGearWeight:             handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:       handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createdShipment := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createdShipment.Validate(strfmt.Default))

		createdPPM := createdShipment.PpmShipment

		suite.Equal(move.ID.String(), createdShipment.MoveTaskOrderID.String())
		suite.Equal(primev3messages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, createdShipment.Status)

		suite.Equal(createdShipment.ID.String(), createdPPM.ShipmentID.String())
		suite.Equal(addressWithEmptyStreet1.StreetAddress1, *createdPPM.DestinationAddress.StreetAddress1)
		suite.True(len(*createdPPM.DestinationAddress.StreetAddress1) == 0)

		// ************
		// PATCH TESTS
		// ************
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Times(2)

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", createdPPM.ShipmentID.String()), nil)

		var mtoShipment models.MTOShipment
		err := suite.DB().Find(&mtoShipment, createdPPM.ShipmentID)
		suite.NoError(err)
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: createdPPM.ShipmentID,
			IfMatch:       eTag,
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			ShipmentType: primev3messages.MTOShipmentTypePPM,
		}
		// *************************************************************************************
		// *************************************************************************************
		// Run with whitespace in destination street 1. Whitespace will be trimmed and seen as
		// as empty on the server side.
		// *************************************************************************************
		ppmDestinationAddressOptionalStreet1ContainingWhitespaces := primev3messages.PPMDestinationAddress{
			City:           models.StringPointer("Beverly Hills"),
			Country:        models.StringPointer("US"),
			PostalCode:     models.StringPointer("90210"),
			State:          models.StringPointer("CA"),
			StreetAddress1: models.StringPointer("  "), //whitespace
		}
		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			DestinationAddress: struct {
				primev3messages.PPMDestinationAddress
			}{ppmDestinationAddressOptionalStreet1ContainingWhitespaces},
		}

		// Validate incoming payload
		suite.NoError(patchParams.Body.Validate(strfmt.Default))

		patchResponse := patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment := okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM := updatedShipment.PpmShipment
		suite.Equal(ppmDestinationAddressOptionalStreet1ContainingWhitespaces.City, updatedPPM.DestinationAddress.City)
		// test whitespace has been trimmed. it should not be equal after update
		suite.NotEqual(ppmDestinationAddressOptionalStreet1ContainingWhitespaces.StreetAddress1, updatedPPM.DestinationAddress.StreetAddress1)
		// verify street address1 is returned as empty string
		suite.True(len(*updatedPPM.DestinationAddress.StreetAddress1) == 0)
	})

	suite.Run("Successful POST with Shuttle service items without primeEstimatedWeight - Integration Test", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Successful submission, status should be SUBMITTED
		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Reason: models.StringPointer("not applicable"),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
		}, nil)
		serviceItem.ID = uuid.Nil

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:     handlers.FmtUUID(move.ID),
				Agents:              nil,
				CustomerRemarks:     nil,
				PointOfContact:      "John Doe",
				RequestedPickupDate: handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:        primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:       struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:  struct{ primev3messages.Address }{destinationAddress},
			},
		}

		mtoServiceItems := models.MTOServiceItems{serviceItem}
		params.Body.SetMtoServiceItems(*payloads.MTOServiceItems(&mtoServiceItems))

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createMTOShipmentPayload.Validate(strfmt.Default))

		// check that the mto shipment status is Submitted
		suite.Require().Equal(createMTOShipmentPayload.Status, primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, "MTO Shipment should have been submitted")
	})

	suite.Run("POST failure - 500", func() {
		// Under Test: CreateMTOShipmentHandler
		// Mocked:     CreateMTOShipment creator
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 500 response
		// Expected:   500 Response returned
		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		// Create a handler with the mocked creator
		handler.ShipmentCreator = &mockCreator

		err := errors.New("ServerError")

		mockCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil, err)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.Run("POST failure - 422 -- Bad agent IDs set on shipment", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with an agent that doesn't really exist, handler should return unprocessable entity
		// Expected:   422 Unprocessable Entity Response returned

		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		badID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		agent := &primev3messages.MTOAgent{
			ID:            badID,
			MtoShipmentID: badID,
			FirstName:     handlers.FmtString("Mary"),
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				Agents:               primev3messages.MTOAgents{agent},
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.NotEmpty(typedResponse.Payload.InvalidFields)
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:id")
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:mtoShipmentID")
	})

	suite.Run("POST failure - 422 - invalid input, missing pickup address", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with missing pickup address, handler should return unprocessable entity
		// Expected:   422 Unprocessable Entity Response returned

		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}
		params.Body.PickupAddress.Address.StreetAddress1 = nil

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		unprocessableEntity := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		// TODO: Can't validate the response because of the issue noted below. Figure out a way to
		//   either alter the service or relax the swagger requirements.
		// suite.NoError(unprocessableEntity.Payload.Validate(strfmt.Default))
		// CreateShipment is returning apperror.InvalidInputError without any validation errors
		// so InvalidFields won't be added to the payload.

		suite.Contains(*unprocessableEntity.Payload.Detail, "PickupAddress is required")
	})

	suite.Run("POST failure - 404 -- not found", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment on a non-existent move
		// Expected:   404 Not Found returned
		handler, _ := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		// Generate a unique id
		badID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      &badID,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.CreateMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a request with no data in the body
		// Expected:   422 Unprocessable Entity Response returned

		handler, _ := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		paramsNilBody := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}

		// Validate incoming payload: nil body (the point of this test)

		response := handler.Handle(paramsNilBody)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
		responsePayload := response.(*mtoshipmentops.CreateMTOShipmentBadRequest).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 404 -- MTO is not available to Prime", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment on an unavailable move, prime cannot update these
		// Expected:   404 Not found returned

		handler, _ := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		unavailableMove := factory.BuildMove(suite.DB(), nil, nil)
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(unavailableMove.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentNotFound)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*typedResponse.Payload.Detail, unavailableMove.ID.String())
	})

	suite.Run("POST failure - 500 - App Event Internal DTOD Server Error", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with DTOD outage simulated or bad zip
		// Expected:   500 Internal Server Error returned

		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)
		handler.ShipmentCreator = &mockCreator

		err := apperror.EventError{}

		mockCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:     handlers.FmtUUID(move.ID),
				Agents:              nil,
				CustomerRemarks:     nil,
				PointOfContact:      "John Doe",
				RequestedPickupDate: handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:        primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:       struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:  struct{ primev3messages.Address }{destinationAddress},
			},
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)
		suite.Contains(*typedResponse.Payload.Detail, "An internal server error has occurred")
	})

	suite.Run("POST failure - 422 - MTO Shipment object not formatted correctly", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with service items that don't match the modeltype
		// Expected:   422 Unprocessable Entity returned

		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)
		handler.ShipmentCreator = &mockCreator

		err := apperror.NotFoundError{}

		mockCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
				BoatShipment:         &primev3messages.CreateBoatShipment{}, // Empty boat shipment will trigger validation error on MTO Shipment creation
			},
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*typedResponse.Payload.Detail, "The MTO shipment object is invalid.")
	})

	suite.Run("POST failure - 422 - modelType() not supported", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with service items that don't match the modeltype
		// Expected:   422 Unprocessable Entity returned

		handler, move := setupTestData(true, false)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)
		handler.ShipmentCreator = &mockCreator

		err := apperror.NotFoundError{}

		mockCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		// Create a service item that doesn't match the modeltype
		mtoServiceItems := models.MTOServiceItems{
			models.MTOServiceItem{
				MoveTaskOrderID:  move.ID,
				MTOShipmentID:    &uuid.Nil,
				ReService:        models.ReService{Code: models.ReServiceCodeMS},
				Reason:           nil,
				PickupPostalCode: nil,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
			},
		}

		params.Body.SetMtoServiceItems(*payloads.MTOServiceItems(&mtoServiceItems))

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*typedResponse.Payload.Detail, "MTOServiceItem modelType() not allowed")
	})

	suite.Run("POST failure - Error when feature flag fetcher fails and a boat shipment is passed in.", func() {
		// Under Test: CreateMTOShipmentHandler
		// Mocked:     CreateMTOShipment creator
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 500 response
		// Expected:   500 Response returned
		suite.T().Setenv("FEATURE_FLAG_BOAT", "true") // Set to true in order to test that it will default to "false" if flag fetcher errors out.

		handler, move := setupTestData(false, false)

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeBOATHAULAWAY),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "Boat shipment type was used but the feature flag is not enabled.")
	})

	suite.Run("POST failure - Error when UB FF is off and UB shipment is passed in.", func() {
		// Under Test: CreateMTOShipmentHandler
		// Mocked:     CreateMTOShipment creator
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 500 response
		// Expected:   500 Response returned
		suite.T().Setenv("FEATURE_FLAG_UNACCOMPANIED_BAGGAGE", "false") // Set to true in order to test that it will default to "false" if flag fetcher errors out.

		handler, move := setupTestData(false, false)

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "Unaccompanied baggage shipments can't be created unless the unaccompanied_baggage feature flag is enabled.")
	})

	suite.Run("POST failure - Error creating a mto shipment contains tertiary destination address no secondary destination address.", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 422 response
		// Expected:   422 Response returned

		handler, move := setupTestData(false, false)

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		newAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		destinationAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		destinationAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		tertiaryDestinationAddress = primev3messages.Address{
			City:           &newAddress.City,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:            handlers.FmtUUID(move.ID),
				Agents:                     nil,
				CustomerRemarks:            nil,
				PointOfContact:             "John Doe",
				PrimeEstimatedWeight:       handlers.FmtInt64(1200),
				RequestedPickupDate:        handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:               primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:              struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:         struct{ primev3messages.Address }{destinationAddress},
				TertiaryDestinationAddress: struct{ primev3messages.Address }{tertiaryDestinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "Invalid input found while validating the MTO shipment")
		suite.Contains(errResponse.Payload.InvalidFields["error validating mto shipment"][0], "Shipment cannot have a third address without a second address present")
	})

	suite.Run("POST failure - Error creating an mto shipment with ppm shipment contains tertiary pickup address no secondary pickup address.", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   If underlying UpdateMTOShipment returns error, handler should return 422 response
		// Expected:   422 Response returned

		handler, move := setupTestData(false, false)

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		sitExpected := true
		estimatedWeight := unit.Pound(1500)
		hasProGear := true
		ppmShipmentDestinationAddress := primev3messages.PPMDestinationAddress{
			City:           destinationAddress.City,
			PostalCode:     destinationAddress.PostalCode,
			State:          destinationAddress.State,
			StreetAddress1: destinationAddress.StreetAddress1,
			StreetAddress2: destinationAddress.StreetAddress2,
			StreetAddress3: destinationAddress.StreetAddress3,
		}
		ppmShipmentParams := primev3messages.CreatePPMShipment{
			PickupAddress: struct{ primev3messages.Address }{pickupAddress},
			DestinationAddress: struct {
				primev3messages.PPMDestinationAddress
			}{ppmShipmentDestinationAddress},
			TertiaryPickupAddress: struct{ primev3messages.Address }{tertiaryPickupAddress},
			ExpectedDepartureDate: handlers.FmtDate(expectedDepartureDate),
			EstimatedWeight:       handlers.FmtPoundPtr(&estimatedWeight),
			HasProGear:            &hasProGear,
			SitExpected:           &sitExpected,
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypePPM),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
				PpmShipment:          &ppmShipmentParams,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "The MTO shipment object is invalid. Shipment cannot have a third pickup address without a second pickup address present")
	})

	suite.Run("POST failure - Error creating mto shipment containing a ppm shipment contains tertiary destination address no secondary destination address.", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   If underlying UpdateMTOShipment returns error, handler should return 422 response
		// Expected:   422 Response returned

		handler, move := setupTestData(false, false)

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		sitExpected := true
		estimatedWeight := unit.Pound(1500)
		hasProGear := true
		ppmShipmentDestinationAddress := primev3messages.PPMDestinationAddress{
			City:           destinationAddress.City,
			PostalCode:     destinationAddress.PostalCode,
			State:          destinationAddress.State,
			StreetAddress1: destinationAddress.StreetAddress1,
			StreetAddress2: destinationAddress.StreetAddress2,
			StreetAddress3: destinationAddress.StreetAddress3,
		}
		ppmShipmentParams := primev3messages.CreatePPMShipment{
			PickupAddress: struct{ primev3messages.Address }{pickupAddress},
			DestinationAddress: struct {
				primev3messages.PPMDestinationAddress
			}{ppmShipmentDestinationAddress},
			TertiaryDestinationAddress: struct{ primev3messages.Address }{tertiaryDestinationAddress},
			ExpectedDepartureDate:      handlers.FmtDate(expectedDepartureDate),
			EstimatedWeight:            handlers.FmtPoundPtr(&estimatedWeight),
			HasProGear:                 &hasProGear,
			SitExpected:                &sitExpected,
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypePPM),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
				PpmShipment:          &ppmShipmentParams,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "The MTO shipment object is invalid. Shipment cannot have a third destination address without a second destination address present")
	})
	suite.Run("PATCH failure - Error updating an mto shipment contains tertiary pickup address no secondary pickup address.", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 422 response
		// Expected:   422 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMove.MTOShipments[0].UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMove.MTOShipments[0].ID.String()),
			IfMatch:       eTag,
		}
		tertiaryAddress := GetTestAddress()
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{tertiaryAddress},
		}
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.Contains(*errResponse.Payload.Detail, "Invalid input found while validating the MTO shipment")
		suite.Contains(errResponse.Payload.InvalidFields["error validating mto shipment"][0], "Shipment cannot have a third address without a second address present")

	})
	suite.Run("PATCH failure - Error updating an ppm shipment contains tertiary destination address no secondary destination address", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Mocked:     UpdateMTOShipment creator
		// Setup:   If underlying UpdateMTOShipment returns error, handler should return 422 response
		// Expected:   422 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{}, nil)
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID.String()), nil)

		eTag := etag.GenerateEtag(testMove.MTOShipments[0].UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMove.MTOShipments[0].ID.String()),
			IfMatch:       eTag,
		}
		tertiaryAddress := GetTestAddress()
		ppmShipmentParamSetup := primev3messages.UpdatePPMShipment{
			HasTertiaryDestinationAddress: models.BoolPointer(true),
			TertiaryDestinationAddress:    struct{ primev3messages.Address }{tertiaryAddress},
		}
		mtoShipmentParamSetup := primev3messages.UpdateMTOShipment{
			PpmShipment: &ppmShipmentParamSetup,
		}

		patchParams.Body = &mtoShipmentParamSetup
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.Contains(*errResponse.Payload.Detail, "Invalid input found while validating the PPM shipment")
		suite.Contains(errResponse.Payload.InvalidFields["error validating ppm shipment"][0], "Shipment cannot have a third address without a second address present")

	})

	suite.Run("PATCH sucess - updating an mto shipment contains tertiary pickup and secondary pickup address.", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 422 response
		// Expected:   422 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryDeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMtoShipment.ID.String()),
			IfMatch:       eTag,
		}
		tertiaryAddress := GetTestAddress()
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{tertiaryAddress},
		}
		patchResponse := patchHandler.Handle(patchParams)
		response := patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
	})

	suite.Run("PATCH failure - Invalid address.", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   Set an invalid zip
		// Expected:   422 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryDeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryDeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMtoShipment.ID.String()),
			IfMatch:       eTag,
		}
		tertiaryAddress := GetTestAddress()
		tertiaryAddress.PostalCode = handlers.FmtString("99999")
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{tertiaryAddress},
		}
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, errResponse)
	})

	suite.Run("PATCH failure - Internal Server error GetLocationsByZipCityState", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   Mock location to return an error
		// Expected:   500 Response returned
		handler, move := setupTestData(false, true)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		expectedError := models.ErrFetchNotFound
		vLocationFetcher := &mocks.VLocation{}
		vLocationFetcher.On("GetLocationsByZipCityState",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler.VLocation = vLocationFetcher
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
	})

	suite.Run("PATCH success - valid AK address FF is on", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   Set an valid AK address but turn FF on
		// Expected:   200 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryDeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMtoShipment.ID.String()),
			IfMatch:       eTag,
		}
		alaskaAddress := primev3messages.Address{
			City:           handlers.FmtString("Juneau"),
			PostalCode:     handlers.FmtString("99801"),
			State:          handlers.FmtString("AK"),
			StreetAddress1: handlers.FmtString("Some AK street"),
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{alaskaAddress},
		}

		// setting the AK flag to true
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_alaska",
			Match: true,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		patchHandler.HandlerConfig = handlerConfig
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, errResponse)
	})

	suite.Run("PATCH success - valid HI address FF is on", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   Set an valid HI address but turn FF on
		// Expected:   200 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryDeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMtoShipment.ID.String()),
			IfMatch:       eTag,
		}
		hawaiiAddress := primev3messages.Address{
			City:           handlers.FmtString("HONOLULU"),
			PostalCode:     handlers.FmtString("96835"),
			State:          handlers.FmtString("HI"),
			StreetAddress1: handlers.FmtString("Some HI street"),
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{hawaiiAddress},
		}

		// setting the HI flag to true
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_hawaii",
			Match: true,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		patchHandler.HandlerConfig = handlerConfig
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, errResponse)
	})

	suite.Run("PATCH failure - valid AK address FF is off", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   Set an valid AK address but turn FF off
		// Expected:   422 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryDeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryDeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMtoShipment.ID.String()),
			IfMatch:       eTag,
		}
		alaskaAddress := primev3messages.Address{
			City:           handlers.FmtString("Juneau"),
			PostalCode:     handlers.FmtString("99801"),
			State:          handlers.FmtString("AK"),
			StreetAddress1: handlers.FmtString("Some AK street"),
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{alaskaAddress},
		}

		// setting the AK flag to false
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_alaska",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		patchHandler.HandlerConfig = handlerConfig
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, errResponse)
	})

	suite.Run("PATCH failure - valid HI address FF is off", func() {
		// Under Test: UpdateMTOShipmentHandler
		// Setup:   Set an valid HI address but turn FF off
		// Expected:   422 Response returned

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			planner,
			vLocationServices,
		}

		now := time.Now()
		mto_shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third pickup address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryPickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some second delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.SecondaryDeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some third delivery address",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     "90210",
				},
				Type: &factory.Addresses.TertiaryDeliveryAddress,
			},
		}, nil)
		move := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		var testMove models.Move
		err := suite.DB().EagerPreload("MTOShipments.PPMShipment").Find(&testMove, move.ID)
		suite.NoError(err)
		var testMtoShipment models.MTOShipment
		err = suite.DB().Find(&testMtoShipment, mto_shipment.ID)
		suite.NoError(err)
		testMtoShipment.MoveTaskOrderID = testMove.ID
		testMtoShipment.MoveTaskOrder = testMove
		err = suite.DB().Save(&testMtoShipment)
		suite.NoError(err)
		testMove.MTOShipments = append(testMove.MTOShipments, mto_shipment)
		err = suite.DB().Save(&testMove)
		suite.NoError(err)

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", testMove.MTOShipments[0].ID), nil)

		eTag := etag.GenerateEtag(testMtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: strfmt.UUID(testMtoShipment.ID.String()),
			IfMatch:       eTag,
		}
		hawaiiAddress := primev3messages.Address{
			City:           handlers.FmtString("HONOLULU"),
			PostalCode:     handlers.FmtString("HI"),
			State:          handlers.FmtString("96835"),
			StreetAddress1: handlers.FmtString("Some HI street"),
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			TertiaryDeliveryAddress: struct{ primev3messages.Address }{hawaiiAddress},
		}

		// setting the HI flag to false
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_hawaii",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		patchHandler.HandlerConfig = handlerConfig
		patchResponse := patchHandler.Handle(patchParams)
		errResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, errResponse)
	})

	suite.Run("POST failure - 422 - Invalid address", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Failure, invalid address
		handler, move := setupTestDataWithoutFF()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:             handlers.FmtUUID(move.ID),
				Agents:                      nil,
				CustomerRemarks:             nil,
				PointOfContact:              "John Doe",
				PrimeEstimatedWeight:        handlers.FmtInt64(1200),
				RequestedPickupDate:         handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:                primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
				SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
				TertiaryPickupAddress:       struct{ primev3messages.Address }{tertiaryPickupAddress},
				DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
				SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
				TertiaryDestinationAddress:  struct{ primev3messages.Address }{tertiaryDestinationAddress},
			},
		}

		// set bad data for address so the validation fails
		params.Body.PickupAddress.City = handlers.FmtString("Bad City")

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 422 - Doesn't return results for valid AK address if FF returns false", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Failure, valid AK address but AK FF off, no results
		handler, move := setupTestDataWithoutFF()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:             handlers.FmtUUID(move.ID),
				Agents:                      nil,
				CustomerRemarks:             nil,
				PointOfContact:              "John Doe",
				PrimeEstimatedWeight:        handlers.FmtInt64(1200),
				RequestedPickupDate:         handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:                primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
				SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
				TertiaryPickupAddress:       struct{ primev3messages.Address }{tertiaryPickupAddress},
				DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
				SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
				TertiaryDestinationAddress:  struct{ primev3messages.Address }{tertiaryDestinationAddress},
			},
		}

		// setting the AK flag to false and use a valid address
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_alaska",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlag",
			mock.Anything,                 // context.Context
			mock.Anything,                 // *zap.Logger
			mock.AnythingOfType("string"), // entityID (userID)
			mock.AnythingOfType("string"), // key
			mock.Anything,                 // flagContext (map[string]string)
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler.HandlerConfig = handlerConfig
		params.Body.PickupAddress.City = handlers.FmtString("JUNEAU")
		params.Body.PickupAddress.State = handlers.FmtString("AK")
		params.Body.PickupAddress.PostalCode = handlers.FmtString("99801")

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 422 - Doesn't return results for valid HI address if FF returns false", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Failure, valid HI address but HI FF off, no results
		handler, move := setupTestDataWithoutFF()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:             handlers.FmtUUID(move.ID),
				Agents:                      nil,
				CustomerRemarks:             nil,
				PointOfContact:              "John Doe",
				PrimeEstimatedWeight:        handlers.FmtInt64(1200),
				RequestedPickupDate:         handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:                primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
				SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
				TertiaryPickupAddress:       struct{ primev3messages.Address }{tertiaryPickupAddress},
				DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
				SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
				TertiaryDestinationAddress:  struct{ primev3messages.Address }{tertiaryDestinationAddress},
			},
		}

		// setting the HI flag to false and use a valid address
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_hawaii",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlag",
			mock.Anything,                 // context.Context
			mock.Anything,                 // *zap.Logger
			mock.AnythingOfType("string"), // entityID (userID)
			mock.AnythingOfType("string"), // key
			mock.Anything,                 // flagContext (map[string]string)
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler.HandlerConfig = handlerConfig
		params.Body.PickupAddress.City = handlers.FmtString("HONOLULU")
		params.Body.PickupAddress.State = handlers.FmtString("HI")
		params.Body.PickupAddress.PostalCode = handlers.FmtString("96835")

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("POST success - 200 - valid AK address if FF ON", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Success, valid AK address AK FF ON
		handler, move := setupTestData(false, true)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:             handlers.FmtUUID(move.ID),
				Agents:                      nil,
				CustomerRemarks:             nil,
				PointOfContact:              "John Doe",
				PrimeEstimatedWeight:        handlers.FmtInt64(1200),
				RequestedPickupDate:         handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:                primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
				SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
				TertiaryPickupAddress:       struct{ primev3messages.Address }{tertiaryPickupAddress},
				DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
				SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
				TertiaryDestinationAddress:  struct{ primev3messages.Address }{tertiaryDestinationAddress},
			},
		}

		// setting the AK flag to false and use a valid address
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_alaska",
			Match: true,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlag",
			mock.Anything,                 // context.Context
			mock.Anything,                 // *zap.Logger
			mock.AnythingOfType("string"), // entityID (userID)
			mock.AnythingOfType("string"), // key
			mock.Anything,                 // flagContext (map[string]string)
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler.HandlerConfig = handlerConfig
		params.Body.PickupAddress.City = handlers.FmtString("JUNEAU")
		params.Body.PickupAddress.State = handlers.FmtString("AK")
		params.Body.PickupAddress.PostalCode = handlers.FmtString("99801")

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
	})

	suite.Run("POST success - 200 - valid HI address if FF ON", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Success, valid HI address HI FF ON
		handler, move := setupTestData(false, true)
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:             handlers.FmtUUID(move.ID),
				Agents:                      nil,
				CustomerRemarks:             nil,
				PointOfContact:              "John Doe",
				PrimeEstimatedWeight:        handlers.FmtInt64(1200),
				RequestedPickupDate:         handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:                primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
				SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
				TertiaryPickupAddress:       struct{ primev3messages.Address }{tertiaryPickupAddress},
				DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
				SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
				TertiaryDestinationAddress:  struct{ primev3messages.Address }{tertiaryDestinationAddress},
			},
		}

		// setting the HI flag to false and use a valid address
		handlerConfig := suite.HandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_hawaii",
			Match: true,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlag",
			mock.Anything,                 // context.Context
			mock.Anything,                 // *zap.Logger
			mock.AnythingOfType("string"), // entityID (userID)
			mock.AnythingOfType("string"), // key
			mock.Anything,                 // flagContext (map[string]string)
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler.HandlerConfig = handlerConfig
		params.Body.PickupAddress.City = handlers.FmtString("HONOLULU")
		params.Body.PickupAddress.State = handlers.FmtString("HI")
		params.Body.PickupAddress.PostalCode = handlers.FmtString("96835")

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
	})

	suite.Run("Failure POST - 422 - Invalid address (PPM)", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:      Create a PPM shipment on an available move
		// Expected:   Failure, returns an invalid address error
		handler, move := setupTestDataWithoutFF()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		counselorRemarks := "Some counselor remarks"
		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		sitExpected := true
		sitLocation := primev3messages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1500)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3200)
		hasProGear := true
		proGearWeight := unit.Pound(400)
		spouseProGearWeight := unit.Pound(250)
		estimatedIncentive := 123456
		sitEstimatedCost := 67500

		address1 := models.Address{
			StreetAddress1: "some address",
			City:           "Bad City",
			State:          "CA",
			PostalCode:     "90210",
		}

		expectedPickupAddress := address1
		pickupAddress = primev3messages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := address1
		destinationAddress = primev3messages.Address{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}
		ppmDestinationAddress = primev3messages.PPMDestinationAddress{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:  handlers.FmtUUID(move.ID),
				ShipmentType:     primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypePPM),
				CounselorRemarks: &counselorRemarks,
				PpmShipment: &primev3messages.CreatePPMShipment{
					ExpectedDepartureDate:  handlers.FmtDate(expectedDepartureDate),
					PickupAddress:          struct{ primev3messages.Address }{pickupAddress},
					SecondaryPickupAddress: struct{ primev3messages.Address }{secondaryPickupAddress},
					TertiaryPickupAddress:  struct{ primev3messages.Address }{tertiaryPickupAddress},
					DestinationAddress: struct {
						primev3messages.PPMDestinationAddress
					}{ppmDestinationAddress},
					SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
					TertiaryDestinationAddress:  struct{ primev3messages.Address }{tertiaryDestinationAddress},
					SitExpected:                 &sitExpected,
					SitLocation:                 &sitLocation,
					SitEstimatedWeight:          handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:       handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate:   handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:             handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                  &hasProGear,
					ProGearWeight:               handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:         handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})
}
func GetTestAddress() primev3messages.Address {
	newAddress := factory.BuildAddress(nil, []factory.Customization{
		{
			Model: models.Address{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
	}, nil)
	return primev3messages.Address{
		City:           &newAddress.City,
		PostalCode:     &newAddress.PostalCode,
		State:          &newAddress.State,
		StreetAddress1: &newAddress.StreetAddress1,
		StreetAddress2: newAddress.StreetAddress2,
		StreetAddress3: newAddress.StreetAddress3,
	}
}
