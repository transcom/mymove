package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mobilehomeshipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type mtoShipmentObjects struct {
	builder    *query.Builder
	fetcher    services.Fetcher
	moveRouter services.MoveRouter
}

func (suite *HandlerSuite) setUpMTOShipmentObjects() *mtoShipmentObjects {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())

	return &mtoShipmentObjects{
		builder:    builder,
		fetcher:    fetcher,
		moveRouter: moveRouter,
	}
}

//
// CREATE
//

func (suite *HandlerSuite) TestCreateMTOShipmentHandlerV1() {
	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeCreateSubtestData.
	testMTOShipmentObjects := suite.setUpMTOShipmentObjects()
	addressCreator := address.NewAddressCreator()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV1(testMTOShipmentObjects.builder, testMTOShipmentObjects.fetcher, testMTOShipmentObjects.moveRouter, addressCreator)
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

	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		testMTOShipmentObjects.builder,
		mtoserviceitem.NewMTOServiceItemCreator(planner, testMTOShipmentObjects.builder, testMTOShipmentObjects.moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		testMTOShipmentObjects.moveRouter,
		setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
	)
	shipmentCreator := shipmentorchestrator.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, boatShipmentCreator, mobileHomeShipmentCreator, shipmentRouter, moveTaskOrderUpdater)

	type mtoCreateSubtestData struct {
		serviceMember models.ServiceMember
		pickupAddress models.Address
		mtoShipment   models.MTOShipment
		params        mtoshipmentops.CreateMTOShipmentParams
		handler       CreateMTOShipmentHandler
	}

	makeCreateSubtestData := func(ubFeatureFlag bool, checkForFF bool) (subtestData mtoCreateSubtestData) {
		subtestData.serviceMember = factory.BuildServiceMember(suite.DB(), nil, nil)

		mto := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.serviceMember,
				LinkOnly: true,
			},
		}, nil)

		subtestData.pickupAddress = factory.BuildAddress(suite.DB(), nil, nil)
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})

		subtestData.mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)
		subtestData.mtoShipment.MoveTaskOrderID = mto.ID

		mtoAgent := factory.BuildMTOAgent(suite.DB(), nil, nil)
		agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
			FirstName: mtoAgent.FirstName,
			LastName:  mtoAgent.LastName,
			Email:     mtoAgent.Email,
			Phone:     mtoAgent.Phone,
			AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
		}}

		customerRemarks := "I have some grandfather clocks."

		req := httptest.NewRequest("POST", "/mto_shipments", nil)
		req = suite.AuthenticateRequest(req, subtestData.serviceMember)

		shipmentType := internalmessages.MTOShipmentTypeHHG

		subtestData.params = mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &internalmessages.CreateShipment{
				MoveTaskOrderID: handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID),
				Agents:          agents,
				CustomerRemarks: &customerRemarks,
				PickupAddress: &internalmessages.Address{
					City:           &subtestData.pickupAddress.City,
					PostalCode:     &subtestData.pickupAddress.PostalCode,
					State:          &subtestData.pickupAddress.State,
					StreetAddress1: &subtestData.pickupAddress.StreetAddress1,
					StreetAddress2: subtestData.pickupAddress.StreetAddress2,
					StreetAddress3: subtestData.pickupAddress.StreetAddress3,
				},
				SecondaryPickupAddress: &internalmessages.Address{
					City:           &secondaryPickupAddress.City,
					PostalCode:     &secondaryPickupAddress.PostalCode,
					State:          &secondaryPickupAddress.State,
					StreetAddress1: &secondaryPickupAddress.StreetAddress1,
					StreetAddress2: secondaryPickupAddress.StreetAddress2,
					StreetAddress3: secondaryPickupAddress.StreetAddress3,
				},
				DestinationAddress: &internalmessages.Address{
					City:           &destinationAddress.City,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
				},
				SecondaryDeliveryAddress: &internalmessages.Address{
					City:           &secondaryDeliveryAddress.City,
					PostalCode:     &secondaryDeliveryAddress.PostalCode,
					State:          &secondaryDeliveryAddress.State,
					StreetAddress1: &secondaryDeliveryAddress.StreetAddress1,
					StreetAddress2: secondaryDeliveryAddress.StreetAddress2,
					StreetAddress3: secondaryDeliveryAddress.StreetAddress3,
				},
				RequestedPickupDate:   strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate),
				RequestedDeliveryDate: strfmt.Date(*subtestData.mtoShipment.RequestedDeliveryDate),
				ShipmentType:          &shipmentType,
			},
		}

		if checkForFF {
			ubFF := services.FeatureFlag{
				Key:   "unaccompanied_baggage",
				Match: false,
			}

			handlerConfig := suite.HandlerConfig()
			if !ubFeatureFlag {
				mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
				mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
					mock.Anything,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("string"),
					mock.Anything,
				).Return(ubFF, nil)
				handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
			}

			subtestData.handler = CreateMTOShipmentHandler{
				handlerConfig,
				shipmentCreator,
			}
		} else {
			subtestData.handler = CreateMTOShipmentHandler{
				suite.HandlerConfig(),
				shipmentCreator,
			}
		}

		return subtestData
	}

	suite.Run("Successful POST - Integration Test - HHG", func() {
		subtestData := makeCreateSubtestData(false, false)

		params := subtestData.params

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(internalmessages.MTOShipmentTypeHHG, createdShipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusSubmitted, models.MTOShipmentStatus(createdShipment.Status))
		suite.Equal(*params.Body.CustomerRemarks, *createdShipment.CustomerRemarks)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *createdShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *createdShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *createdShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *createdShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), createdShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.RequestedDeliveryDate.String(), createdShipment.RequestedDeliveryDate.String())

		suite.Equal(params.Body.Agents[0].FirstName, createdShipment.Agents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, createdShipment.Agents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, createdShipment.Agents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, createdShipment.Agents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, createdShipment.Agents[0].AgentType)
		suite.Equal(createdShipment.ID.String(), string(createdShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(createdShipment.Agents[0].ID)
	})

	suite.Run("Successful POST - Integration Test - UB", func() {
		// const ubFlag = "true"
		// suite.T().Setenv("FEATURE_FLAG_UNACCOMPANIED_BAGGAGE", ubFlag)

		subtestData := makeCreateSubtestData(true, true)

		params := subtestData.params
		params.Body.ShipmentType = internalmessages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE.Pointer()

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(internalmessages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE, createdShipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusSubmitted, models.MTOShipmentStatus(createdShipment.Status))
		suite.Equal(*params.Body.CustomerRemarks, *createdShipment.CustomerRemarks)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *createdShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *createdShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *createdShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *createdShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), createdShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.RequestedDeliveryDate.String(), createdShipment.RequestedDeliveryDate.String())

		suite.Equal(params.Body.Agents[0].FirstName, createdShipment.Agents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, createdShipment.Agents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, createdShipment.Agents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, createdShipment.Agents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, createdShipment.Agents[0].AgentType)
		suite.Equal(createdShipment.ID.String(), string(createdShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(createdShipment.Agents[0].ID)
	})

	suite.Run("Unsuccessful POST if UB FF off - Integration Test - UB", func() {
		// const ubFlag = "false"

		// suite.T().Setenv("FEATURE_FLAG_UNACCOMPANIED_BAGGAGE", ubFlag)

		subtestData := makeCreateSubtestData(false, true)
		params := subtestData.params
		params.Body.ShipmentType = internalmessages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE.Pointer()

		response := subtestData.handler.Handle(subtestData.params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "Unaccompanied baggage shipments can't be created unless the unaccompanied_baggage feature flag is enabled.")
	})

	suite.Run("Successful POST - Integration Test - PPM required fields", func() {
		subtestData := makeCreateSubtestData(false, false)

		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM

		// create puckupAddress
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})

		// pointers
		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		sitExpected := false
		// reset Body params to have PPM fields
		params.Body = &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID),
			PpmShipment: &internalmessages.CreatePPMShipment{
				ExpectedDepartureDate: &expectedDepartureDate,
				SitExpected:           &sitExpected,
				PickupAddress: &internalmessages.Address{
					City:           &pickupAddress.City,
					PostalCode:     &pickupAddress.PostalCode,
					State:          &pickupAddress.State,
					StreetAddress1: &pickupAddress.StreetAddress1,
					StreetAddress2: pickupAddress.StreetAddress2,
					StreetAddress3: pickupAddress.StreetAddress3,
				},
				DestinationAddress: &internalmessages.PPMDestinationAddress{
					City:           &destinationAddress.City,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
				},
			},
			ShipmentType: &ppmShipmentType,
		}

		// When a customer first creates a move, there is not enough data to calculate an incentive yet.
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil, nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		suite.Nil(params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload
		suite.NoError(createdShipment.Validate(strfmt.Default))

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(internalmessages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusDraft, models.MTOShipmentStatus(createdShipment.Status))
		suite.Equal(*params.Body.MoveTaskOrderID, createdShipment.MoveTaskOrderID)
		suite.Equal(*params.Body.PpmShipment.ExpectedDepartureDate, *createdShipment.PpmShipment.ExpectedDepartureDate)
		suite.Equal(*params.Body.PpmShipment.SitExpected, *createdShipment.PpmShipment.SitExpected)
		suite.Equal(*params.Body.PpmShipment.PickupAddress.StreetAddress1, *createdShipment.PpmShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.PpmShipment.DestinationAddress.StreetAddress1, *createdShipment.PpmShipment.DestinationAddress.StreetAddress1)
	})

	suite.Run("Successful POST - Integration Test - PPM optional fields", func() {
		subtestData := makeCreateSubtestData(false, false)

		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM
		// pointers
		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		sitExpected := false

		// create  PPM addressed
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		secondaryDestinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		// reset Body params to have PPM fields
		params.Body = &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID),
			PpmShipment: &internalmessages.CreatePPMShipment{
				ExpectedDepartureDate: &expectedDepartureDate,
				SitExpected:           &sitExpected,
				PickupAddress: &internalmessages.Address{
					City:           &pickupAddress.City,
					PostalCode:     &pickupAddress.PostalCode,
					State:          &pickupAddress.State,
					StreetAddress1: &pickupAddress.StreetAddress1,
					StreetAddress2: pickupAddress.StreetAddress2,
					StreetAddress3: pickupAddress.StreetAddress3,
				},
				DestinationAddress: &internalmessages.PPMDestinationAddress{
					City:           &destinationAddress.City,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
				},
				SecondaryPickupAddress: &internalmessages.Address{
					City:           &secondaryPickupAddress.City,
					PostalCode:     &secondaryPickupAddress.PostalCode,
					State:          &secondaryPickupAddress.State,
					StreetAddress1: &secondaryPickupAddress.StreetAddress1,
					StreetAddress2: secondaryPickupAddress.StreetAddress2,
					StreetAddress3: secondaryPickupAddress.StreetAddress3,
				},
				SecondaryDestinationAddress: &internalmessages.Address{
					City:           &secondaryDestinationAddress.City,
					PostalCode:     &secondaryDestinationAddress.PostalCode,
					State:          &secondaryDestinationAddress.State,
					StreetAddress1: &secondaryDestinationAddress.StreetAddress1,
					StreetAddress2: secondaryDestinationAddress.StreetAddress2,
					StreetAddress3: secondaryDestinationAddress.StreetAddress3,
				},
			},
			ShipmentType: &ppmShipmentType,
		}

		// When a customer first creates a move, there is not enough data to calculate an incentive yet.
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil, nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		suite.Nil(params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload
		suite.NoError(createdShipment.Validate(strfmt.Default))

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(internalmessages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusDraft, models.MTOShipmentStatus(createdShipment.Status))
		suite.Equal(*params.Body.MoveTaskOrderID, createdShipment.MoveTaskOrderID)
		suite.Equal(*params.Body.PpmShipment.ExpectedDepartureDate, *createdShipment.PpmShipment.ExpectedDepartureDate)
		suite.Equal(*params.Body.PpmShipment.SitExpected, *createdShipment.PpmShipment.SitExpected)
		suite.Equal(*params.Body.PpmShipment.PickupAddress.StreetAddress1, *createdShipment.PpmShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.PpmShipment.DestinationAddress.StreetAddress1, *createdShipment.PpmShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.PpmShipment.SecondaryPickupAddress.StreetAddress1, *createdShipment.PpmShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.PpmShipment.SecondaryDestinationAddress.StreetAddress1, *createdShipment.PpmShipment.SecondaryDestinationAddress.StreetAddress1)
	})

	suite.Run("Successful POST - Integration Test - NTS-Release", func() {
		subtestData := makeCreateSubtestData(false, false)

		params := subtestData.params

		// Set fields appropriately for NTS-Release
		ntsrShipmentType := internalmessages.MTOShipmentTypeHHGOUTOFNTSDOMESTIC
		params.Body.ShipmentType = &ntsrShipmentType
		params.Body.RequestedPickupDate = strfmt.Date(time.Time{})
		params.Body.PickupAddress = nil
		params.Body.SecondaryPickupAddress = nil

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(ntsrShipmentType, createdShipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusSubmitted, models.MTOShipmentStatus(createdShipment.Status))
		suite.Equal(*params.Body.CustomerRemarks, *createdShipment.CustomerRemarks)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *createdShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *createdShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Nil(createdShipment.RequestedPickupDate)
		suite.Equal(params.Body.RequestedDeliveryDate.String(), createdShipment.RequestedDeliveryDate.String())

		suite.Equal(params.Body.Agents[0].FirstName, createdShipment.Agents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, createdShipment.Agents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, createdShipment.Agents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, createdShipment.Agents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, createdShipment.Agents[0].AgentType)
		suite.Equal(createdShipment.ID.String(), string(createdShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(createdShipment.Agents[0].ID)
	})

	suite.Run("Successful POST - Integration Test - Boat", func() {
		subtestData := makeCreateSubtestData(false, false)

		params := subtestData.params

		boatShipmentType := internalmessages.MTOShipmentTypeBOATHAULAWAY

		boatShipment := &internalmessages.CreateBoatShipment{
			Type:           stringPtr(string(models.BoatShipmentTypeHaulAway)),
			Year:           models.Int64Pointer(1000),
			Make:           models.StringPointer("Boat Make"),
			Model:          models.StringPointer("Boat Model"),
			LengthInInches: models.Int64Pointer(300),
			WidthInInches:  models.Int64Pointer(108),
			HeightInInches: models.Int64Pointer(72),
			HasTrailer:     models.BoolPointer(true),
			IsRoadworthy:   models.BoolPointer(false),
		}
		params.Body.ShipmentType = &boatShipmentType
		params.Body.BoatShipment = boatShipment

		params.Body.RequestedPickupDate = strfmt.Date(time.Time{})

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(boatShipmentType, createdShipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusDraft, models.MTOShipmentStatus(createdShipment.Status))
		suite.Equal(*params.Body.CustomerRemarks, *createdShipment.CustomerRemarks)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *createdShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *createdShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *createdShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *createdShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Nil(createdShipment.RequestedPickupDate)
		suite.Equal(params.Body.RequestedDeliveryDate.String(), createdShipment.RequestedDeliveryDate.String())

		suite.Equal(*params.Body.BoatShipment.Type, *createdShipment.BoatShipment.Type)
		suite.Equal(*params.Body.BoatShipment.Year, *createdShipment.BoatShipment.Year)
		suite.Equal(*params.Body.BoatShipment.Make, *createdShipment.BoatShipment.Make)
		suite.Equal(*params.Body.BoatShipment.Model, *createdShipment.BoatShipment.Model)
		suite.Equal(*params.Body.BoatShipment.LengthInInches, *createdShipment.BoatShipment.LengthInInches)
		suite.Equal(*params.Body.BoatShipment.WidthInInches, *createdShipment.BoatShipment.WidthInInches)
		suite.Equal(*params.Body.BoatShipment.HeightInInches, *createdShipment.BoatShipment.HeightInInches)
		suite.Equal(*params.Body.BoatShipment.HasTrailer, *createdShipment.BoatShipment.HasTrailer)
		suite.Equal(*params.Body.BoatShipment.IsRoadworthy, *createdShipment.BoatShipment.IsRoadworthy)

	})

	suite.Run("POST failure - 400 - invalid input, missing pickup address", func() {
		subtestData := makeCreateSubtestData(false, false)

		badParams := subtestData.params
		badParams.Body.PickupAddress = nil

		response := subtestData.handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 401- permission denied - not authenticated", func() {
		subtestData := makeCreateSubtestData(false, false)

		unauthorizedReq := httptest.NewRequest("POST", "/mto_shipments", nil)
		shipmentType := internalmessages.MTOShipmentTypeHHG
		unauthorizedParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: unauthorizedReq,
			Body: &internalmessages.CreateShipment{
				MoveTaskOrderID: handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID),
				Agents:          internalmessages.MTOAgents{},
				CustomerRemarks: nil,
				PickupAddress: &internalmessages.Address{
					City:           &subtestData.pickupAddress.City,
					PostalCode:     &subtestData.pickupAddress.PostalCode,
					State:          &subtestData.pickupAddress.State,
					StreetAddress1: &subtestData.pickupAddress.StreetAddress1,
					StreetAddress2: subtestData.pickupAddress.StreetAddress2,
					StreetAddress3: subtestData.pickupAddress.StreetAddress3,
				},
				RequestedPickupDate:   strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate),
				RequestedDeliveryDate: strfmt.Date(*subtestData.mtoShipment.RequestedDeliveryDate),
				ShipmentType:          &shipmentType,
			},
		}

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("POST failure - 403 - unauthorized - wrong application", func() {
		subtestData := makeCreateSubtestData(false, false)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("POST failure - 404 - not found - wrong SM does not match move", func() {
		subtestData := makeCreateSubtestData(false, false)

		sm := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateUserRequest(req, sm.User)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.Run("POST failure - 404 -- not found", func() {
		subtestData := makeCreateSubtestData(false, false)

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := subtestData.params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := subtestData.handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		subtestData := makeCreateSubtestData(false, false)

		otherParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: subtestData.params.HTTPRequest,
		}
		response := subtestData.handler.Handle(otherParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})

	suite.Run("POST failure - 400 -- missing required field to Create PPM", func() {
		subtestData := makeCreateSubtestData(false, false)

		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM

		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		sitExpected := false
		badID, _ := uuid.NewV4()

		// reset Body params to have PPM fields
		params.Body = &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(badID),
			PpmShipment: &internalmessages.CreatePPMShipment{
				ExpectedDepartureDate: &expectedDepartureDate,
				SitExpected:           &sitExpected,
			},
			ShipmentType: &ppmShipmentType,
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentNotFound).Payload
		suite.Equal(handlers.NotFoundMessage, *errResponse.Title)

		// Check Error details
		suite.Contains(*errResponse.Detail, "not found for move")
	})

	suite.Run("POST failure - 500", func() {
		subtestData := makeCreateSubtestData(false, false)

		mockShipmentCreator := mocks.ShipmentCreator{}

		err := errors.New("ServerError")

		mockShipmentCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil, err)

		handler := CreateMTOShipmentHandler{
			suite.HandlerConfig(),
			&mockShipmentCreator,
		}

		response := handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})
}

//
// UPDATE
//

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in getDefaultMTOShipmentAndParams or
	// getDefaultPPMShipmentAndParams.
	testMTOShipmentObjects := suite.setUpMTOShipmentObjects()

	planner := &routemocks.Planner{}

	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	moveWeights := moverouter.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(testMTOShipmentObjects.builder)

	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)

	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()
	mtoShipmentUpdater := mtoshipment.NewCustomerMTOShipmentUpdater(testMTOShipmentObjects.builder, testMTOShipmentObjects.fetcher, planner, testMTOShipmentObjects.moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

	ppmEstimator := mocks.PPMEstimator{}

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
	boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
	mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)

	authRequestAndSetUpHandlerAndParams := func(originalShipment models.MTOShipment, mockShipmentUpdater *mocks.ShipmentUpdater) (UpdateMTOShipmentHandler, mtoshipmentops.UpdateMTOShipmentParams) {
		endpoint := fmt.Sprintf("/mto-shipments/%s", originalShipment.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		req = suite.AuthenticateRequest(req, originalShipment.MoveTaskOrder.Orders.ServiceMember)

		eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(originalShipment.ID),
			IfMatch:       eTag,
		}

		shipmentUpdaterSO := shipmentUpdater
		if mockShipmentUpdater != nil {
			shipmentUpdaterSO = mockShipmentUpdater
		}

		handler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdaterSO,
		}

		return handler, params
	}

	type mtoUpdateSubtestData struct {
		mtoShipment *models.MTOShipment
		params      mtoshipmentops.UpdateMTOShipmentParams
		handler     UpdateMTOShipmentHandler
	}

	// getDefaultMTOShipmentAndParams generates a set of default params and an MTOShipment
	getDefaultMTOShipmentAndParams := func(mockShipmentUpdater *mocks.ShipmentUpdater) *mtoUpdateSubtestData {
		originalShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		pickupAddress.StreetAddress1 = "123 Fake Test St NW"

		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		secondaryPickupAddress.StreetAddress1 = "89999 Other Test St NW"

		destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"

		secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)
		secondaryDeliveryAddress.StreetAddress1 = "9999 Test Fake Rd SE"

		mtoAgent := factory.BuildMTOAgent(suite.DB(), nil, nil)
		agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
			FirstName: mtoAgent.FirstName,
			LastName:  mtoAgent.LastName,
			Email:     mtoAgent.Email,
			Phone:     mtoAgent.Phone,
			AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
		}}

		customerRemarks := ""

		handler, params := authRequestAndSetUpHandlerAndParams(originalShipment, mockShipmentUpdater)

		params.Body = &internalmessages.UpdateShipment{
			Agents:          agents,
			CustomerRemarks: &customerRemarks,
			DestinationAddress: &internalmessages.Address{
				City:           &destinationAddress.City,
				PostalCode:     &destinationAddress.PostalCode,
				State:          &destinationAddress.State,
				StreetAddress1: &destinationAddress.StreetAddress1,
				StreetAddress2: destinationAddress.StreetAddress2,
				StreetAddress3: destinationAddress.StreetAddress3,
			},
			SecondaryDeliveryAddress: &internalmessages.Address{
				City:           &secondaryDeliveryAddress.City,
				PostalCode:     &secondaryDeliveryAddress.PostalCode,
				State:          &secondaryDeliveryAddress.State,
				StreetAddress1: &secondaryDeliveryAddress.StreetAddress1,
				StreetAddress2: secondaryDeliveryAddress.StreetAddress2,
				StreetAddress3: secondaryDeliveryAddress.StreetAddress3,
			},
			HasSecondaryDeliveryAddress: handlers.FmtBool(true),
			PickupAddress: &internalmessages.Address{
				City:           &pickupAddress.City,
				PostalCode:     &pickupAddress.PostalCode,
				State:          &pickupAddress.State,
				StreetAddress1: &pickupAddress.StreetAddress1,
				StreetAddress2: pickupAddress.StreetAddress2,
				StreetAddress3: pickupAddress.StreetAddress3,
			},
			SecondaryPickupAddress: &internalmessages.Address{
				City:           &secondaryPickupAddress.City,
				PostalCode:     &secondaryPickupAddress.PostalCode,
				State:          &secondaryPickupAddress.State,
				StreetAddress1: &secondaryPickupAddress.StreetAddress1,
				StreetAddress2: secondaryPickupAddress.StreetAddress2,
				StreetAddress3: secondaryPickupAddress.StreetAddress3,
			},
			HasSecondaryPickupAddress: handlers.FmtBool(true),
			RequestedPickupDate:       handlers.FmtDatePtr(originalShipment.RequestedPickupDate),
			RequestedDeliveryDate:     handlers.FmtDatePtr(originalShipment.RequestedDeliveryDate),
			ShipmentType:              internalmessages.MTOShipmentTypeHHG,
			ActualProGearWeight:       handlers.FmtInt64(1860),
			ActualSpouseProGearWeight: handlers.FmtInt64(202),
		}

		return &mtoUpdateSubtestData{
			mtoShipment: &originalShipment,
			params:      params,
			handler:     handler,
		}
	}

	suite.Run("Successful PATCH - Integration Test", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)
		params := subtestData.params

		response := subtestData.handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		suite.Equal(subtestData.mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(*params.Body.CustomerRemarks, *updatedShipment.CustomerRemarks)
		suite.Equal(*params.Body.ActualProGearWeight, *updatedShipment.ActualProGearWeight)
		suite.Equal(*params.Body.ActualSpouseProGearWeight, *updatedShipment.ActualSpouseProGearWeight)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *updatedShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *updatedShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), updatedShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.RequestedDeliveryDate.String(), updatedShipment.RequestedDeliveryDate.String())

		suite.Equal(params.Body.Agents[0].FirstName, updatedShipment.Agents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, updatedShipment.Agents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, updatedShipment.Agents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, updatedShipment.Agents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, updatedShipment.Agents[0].AgentType)
		suite.Equal(subtestData.mtoShipment.ID.String(), string(updatedShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(updatedShipment.Agents[0].ID)
	})

	suite.Run("Successful PATCH with PPMShipment - Integration Test", func() {

		// checkDatesAndLocationsDidntChange - ensures dates and locations fields didn't change
		checkDatesAndLocationsDidntChange := func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment) {
			suite.EqualDatePtr(&originalShipment.PPMShipment.ExpectedDepartureDate, updatedShipment.PpmShipment.ExpectedDepartureDate)
			suite.Equal(originalShipment.PPMShipment.SITExpected, updatedShipment.PpmShipment.SitExpected)
		}

		// checkEstimatedWeightsDidntChange - ensures estimated weights fields didn't change
		checkEstimatedWeightsDidntChange := func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment) {
			suite.EqualPoundPointers(originalShipment.PPMShipment.EstimatedWeight, updatedShipment.PpmShipment.EstimatedWeight)
			suite.Equal(originalShipment.PPMShipment.HasProGear, updatedShipment.PpmShipment.HasProGear)
			suite.EqualPoundPointers(originalShipment.PPMShipment.ProGearWeight, updatedShipment.PpmShipment.ProGearWeight)
			suite.EqualPoundPointers(originalShipment.PPMShipment.SpouseProGearWeight, updatedShipment.PpmShipment.SpouseProGearWeight)
		}

		// checkAdvanceRequestedFieldsDidntChange - ensures advance requested fields didn't change
		checkAdvanceRequestedFieldsDidntChange := func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment) {
			suite.Equal(originalShipment.PPMShipment.HasRequestedAdvance, updatedShipment.PpmShipment.HasRequestedAdvance)
			suite.EqualCentsPointers(originalShipment.PPMShipment.AdvanceAmountRequested, updatedShipment.PpmShipment.AdvanceAmountRequested)
		}

		type setUpOriginalPPMFunc func() models.PPMShipment
		type runChecksFunc func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment)

		// Address fields
		street1 := "123 main street"
		city := "New York"
		state := "NY"
		zipcode := "90210"

		ppmUpdateTestCases := map[string]struct {
			setUpOriginalPPM   setUpOriginalPPMFunc
			desiredShipment    internalmessages.UpdatePPMShipment
			estimatedIncentive *unit.Cents
			runChecks          runChecksFunc
		}{
			"Edit estimated dates & locations": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								ExpectedDepartureDate: time.Date(testdatagen.GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
								SITExpected:           models.BoolPointer(true),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					ExpectedDepartureDate: handlers.FmtDate(time.Date(testdatagen.GHCTestYear, time.April, 27, 0, 0, 0, 0, time.UTC)),
					SitExpected:           handlers.FmtBool(false),
				},
				estimatedIncentive: nil,
				runChecks: func(updatedShipment *internalmessages.MTOShipment, _ models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check all fields changed as expected
					desiredShipment.ExpectedDepartureDate.Equal(*updatedShipment.PpmShipment.ExpectedDepartureDate)

					suite.Equal(desiredShipment.SitExpected, updatedShipment.PpmShipment.SitExpected)
				},
			},
			"Add estimated weights - no pro gear": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					EstimatedWeight: handlers.FmtInt64(3500),
					HasProGear:      handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check base fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.EstimatedWeight, updatedShipment.PpmShipment.EstimatedWeight)
					suite.Equal(desiredShipment.HasProGear, updatedShipment.PpmShipment.HasProGear)
					suite.Nil(updatedShipment.PpmShipment.ProGearWeight)
					suite.Nil(updatedShipment.PpmShipment.SpouseProGearWeight)
				},
			},
			"Add estimated weights - yes pro gear": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					EstimatedWeight:     handlers.FmtInt64(3500),
					HasProGear:          handlers.FmtBool(true),
					ProGearWeight:       handlers.FmtInt64(1860),
					SpouseProGearWeight: handlers.FmtInt64(160),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check base fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.EstimatedWeight, updatedShipment.PpmShipment.EstimatedWeight)
					suite.Equal(desiredShipment.HasProGear, updatedShipment.PpmShipment.HasProGear)
					suite.Equal(desiredShipment.ProGearWeight, updatedShipment.PpmShipment.ProGearWeight)
					suite.Equal(desiredShipment.SpouseProGearWeight, updatedShipment.PpmShipment.SpouseProGearWeight)
				},
			},
			"Remove pro gear": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:     models.PoundPointer(4000),
								HasProGear:          models.BoolPointer(true),
								ProGearWeight:       models.PoundPointer(1250),
								SpouseProGearWeight: models.PoundPointer(150),
								EstimatedIncentive:  models.CentPointer(unit.Cents(500000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					HasProGear: handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(300000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)

					suite.EqualPoundPointers(originalShipment.PPMShipment.EstimatedWeight, updatedShipment.PpmShipment.EstimatedWeight)

					// check expected fields were updated
					suite.Equal(desiredShipment.HasProGear, updatedShipment.PpmShipment.HasProGear)
					suite.Nil(updatedShipment.PpmShipment.ProGearWeight)
					suite.Nil(updatedShipment.PpmShipment.SpouseProGearWeight)
				},
			},
			"Add advance requested info - no advance": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:    models.PoundPointer(4000),
								HasProGear:         models.BoolPointer(false),
								EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					HasRequestedAdvance: handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.HasRequestedAdvance, updatedShipment.PpmShipment.HasRequestedAdvance)
					suite.Nil(updatedShipment.PpmShipment.AdvanceAmountRequested)
				},
			},
			"Add advance requested info - yes advance": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:    models.PoundPointer(4000),
								HasProGear:         models.BoolPointer(false),
								EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					HasRequestedAdvance:    handlers.FmtBool(true),
					AdvanceAmountRequested: handlers.FmtInt64(200000),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.HasRequestedAdvance, updatedShipment.PpmShipment.HasRequestedAdvance)
					suite.Equal(desiredShipment.AdvanceAmountRequested, updatedShipment.PpmShipment.AdvanceAmountRequested)
				},
			},
			"Remove advance requested": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:        models.PoundPointer(4000),
								HasProGear:             models.BoolPointer(false),
								EstimatedIncentive:     models.CentPointer(unit.Cents(500000)),
								HasRequestedAdvance:    models.BoolPointer(true),
								AdvanceAmountRequested: models.CentPointer(unit.Cents(200000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					HasRequestedAdvance: handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.HasRequestedAdvance, updatedShipment.PpmShipment.HasRequestedAdvance)
					suite.Nil(updatedShipment.PpmShipment.AdvanceAmountRequested)
				},
			},
			"Add actual zips and advance info - no advance": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:        models.PoundPointer(4000),
								HasProGear:             models.BoolPointer(false),
								EstimatedIncentive:     models.CentPointer(unit.Cents(500000)),
								HasRequestedAdvance:    models.BoolPointer(true),
								AdvanceAmountRequested: models.CentPointer(unit.Cents(200000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					ActualPickupPostalCode:      handlers.FmtString("90210"),
					ActualDestinationPostalCode: handlers.FmtString("90210"),
					HasReceivedAdvance:          handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)
					checkAdvanceRequestedFieldsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.ActualPickupPostalCode, updatedShipment.PpmShipment.ActualPickupPostalCode)
					suite.Equal(desiredShipment.ActualDestinationPostalCode, updatedShipment.PpmShipment.ActualDestinationPostalCode)
					suite.Equal(desiredShipment.HasReceivedAdvance, updatedShipment.PpmShipment.HasReceivedAdvance)
					suite.Nil(updatedShipment.PpmShipment.AdvanceAmountReceived)
				},
			},
			"Add actual zips and advance info - yes advance": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:        models.PoundPointer(4000),
								HasProGear:             models.BoolPointer(false),
								EstimatedIncentive:     models.CentPointer(unit.Cents(500000)),
								HasRequestedAdvance:    models.BoolPointer(true),
								AdvanceAmountRequested: models.CentPointer(unit.Cents(200000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					ActualPickupPostalCode:      handlers.FmtString("90210"),
					ActualDestinationPostalCode: handlers.FmtString("90210"),
					HasReceivedAdvance:          handlers.FmtBool(true),
					AdvanceAmountReceived:       handlers.FmtInt64(250000),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)
					checkAdvanceRequestedFieldsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.ActualPickupPostalCode, updatedShipment.PpmShipment.ActualPickupPostalCode)
					suite.Equal(desiredShipment.ActualDestinationPostalCode, updatedShipment.PpmShipment.ActualDestinationPostalCode)
					suite.Equal(desiredShipment.HasReceivedAdvance, updatedShipment.PpmShipment.HasReceivedAdvance)
					suite.Equal(desiredShipment.AdvanceAmountReceived, updatedShipment.PpmShipment.AdvanceAmountReceived)
				},
			},
			"Add W2 Address": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:        models.PoundPointer(4000),
								HasProGear:             models.BoolPointer(false),
								EstimatedIncentive:     models.CentPointer(unit.Cents(500000)),
								HasRequestedAdvance:    models.BoolPointer(true),
								AdvanceAmountRequested: models.CentPointer(unit.Cents(200000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					W2Address: &internalmessages.Address{
						StreetAddress1: &street1,
						City:           &city,
						State:          &state,
						PostalCode:     &zipcode,
					},
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)
					checkAdvanceRequestedFieldsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.W2Address.StreetAddress1, updatedShipment.PpmShipment.W2Address.StreetAddress1)
					suite.Equal(desiredShipment.W2Address.City, updatedShipment.PpmShipment.W2Address.City)
					suite.Equal(desiredShipment.W2Address.PostalCode, updatedShipment.PpmShipment.W2Address.PostalCode)
					suite.Equal(desiredShipment.W2Address.State, updatedShipment.PpmShipment.W2Address.State)
				},
			},
			"Allows updates to W2 Address": {
				setUpOriginalPPM: func() models.PPMShipment {
					buildAddress := factory.BuildAddress(suite.DB(), nil, nil)
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model:    buildAddress,
							LinkOnly: true,
							Type:     &factory.Addresses.W2Address,
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					W2Address: &internalmessages.Address{
						ID:             "92c9d4db-1ae4-41b1-991e-3ed645ee910a",
						StreetAddress1: &street1,
						City:           &city,
						State:          &state,
						PostalCode:     &zipcode,
					},
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)
					checkAdvanceRequestedFieldsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.W2Address.StreetAddress1, updatedShipment.PpmShipment.W2Address.StreetAddress1)
					suite.Equal(desiredShipment.W2Address.City, updatedShipment.PpmShipment.W2Address.City)
					suite.Equal(desiredShipment.W2Address.PostalCode, updatedShipment.PpmShipment.W2Address.PostalCode)
					suite.Equal(desiredShipment.W2Address.State, updatedShipment.PpmShipment.W2Address.State)
					suite.Equal(originalShipment.PPMShipment.W2Address.ID, uuid.FromStringOrNil(updatedShipment.PpmShipment.W2Address.ID.String()))
				},
			},
			"Prevents arbitrary address updates": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					W2Address: &internalmessages.Address{
						ID:             "92c9d4db-1ae4-41b1-991e-3ed645ee910a",
						StreetAddress1: &street1,
						City:           &city,
						State:          &state,
						PostalCode:     &zipcode,
					},
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)
					checkAdvanceRequestedFieldsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.W2Address.StreetAddress1, updatedShipment.PpmShipment.W2Address.StreetAddress1)
					suite.Equal(desiredShipment.W2Address.City, updatedShipment.PpmShipment.W2Address.City)
					suite.Equal(desiredShipment.W2Address.PostalCode, updatedShipment.PpmShipment.W2Address.PostalCode)
					suite.Equal(desiredShipment.W2Address.State, updatedShipment.PpmShipment.W2Address.State)
					suite.NotEqual(desiredShipment.W2Address.ID, updatedShipment.PpmShipment.W2Address.ID)
				},
			},
			"Remove actual advance": {
				setUpOriginalPPM: func() models.PPMShipment {
					return factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
						{
							Model: models.PPMShipment{
								EstimatedWeight:             models.PoundPointer(4000),
								HasProGear:                  models.BoolPointer(false),
								EstimatedIncentive:          models.CentPointer(unit.Cents(500000)),
								HasRequestedAdvance:         models.BoolPointer(true),
								AdvanceAmountRequested:      models.CentPointer(unit.Cents(200000)),
								ActualPickupPostalCode:      models.StringPointer("90210"),
								ActualDestinationPostalCode: models.StringPointer("90210"),
								HasReceivedAdvance:          models.BoolPointer(true),
								AdvanceAmountReceived:       models.CentPointer(unit.Cents(250000)),
							},
						},
					}, nil)
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					HasReceivedAdvance: handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)
					checkAdvanceRequestedFieldsDidntChange(updatedShipment, originalShipment)

					suite.Equal(originalShipment.PPMShipment.ActualPickupPostalCode, updatedShipment.PpmShipment.ActualPickupPostalCode)
					suite.Equal(originalShipment.PPMShipment.ActualDestinationPostalCode, updatedShipment.PpmShipment.ActualDestinationPostalCode)

					// check expected fields were updated
					suite.Equal(desiredShipment.HasReceivedAdvance, updatedShipment.PpmShipment.HasReceivedAdvance)
					suite.Nil(updatedShipment.PpmShipment.AdvanceAmountReceived)
				},
			},
		}

		for name, tc := range ppmUpdateTestCases {
			name := name
			tc := tc

			suite.Run(name, func() {
				ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("models.PPMShipment"),
					mock.AnythingOfType("*models.PPMShipment")).
					Return(tc.estimatedIncentive, nil, nil).Once()

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

				originalPPMShipment := tc.setUpOriginalPPM()

				handler, params := authRequestAndSetUpHandlerAndParams(originalPPMShipment.Shipment, nil)

				params.Body = &internalmessages.UpdateShipment{
					ShipmentType: internalmessages.MTOShipmentTypePPM,
					PpmShipment:  &tc.desiredShipment,
				}

				response := handler.Handle(params)

				suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

				updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

				suite.NoError(updatedShipment.Validate(strfmt.Default))

				// Check that existing fields are not updated
				suite.Equal(originalPPMShipment.ShipmentID.String(), updatedShipment.ID.String())

				suite.EqualCentsPointers(tc.estimatedIncentive, updatedShipment.PpmShipment.EstimatedIncentive)

				tc.runChecks(updatedShipment, originalPPMShipment.Shipment, tc.desiredShipment)
			})
		}
	})

	suite.Run("Successful PATCH - Can update shipment status", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		expectedStatus := internalmessages.MTOShipmentStatusSUBMITTED

		subtestData.params.Body.Status = expectedStatus

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)

		suite.Equal(expectedStatus, updatedResponse.Payload.Status)
	})

	suite.Run("PATCH failure - 400 -- nil body", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		subtestData.params.Body = nil

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.Run("PATCH failure - 400 -- invalid requested status update", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		subtestData.params.Body.Status = internalmessages.MTOShipmentStatusREJECTED

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.Run("PATCH failure - 401- permission denied - not authenticated", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		updateURI := "/mto-shipments/" + subtestData.mtoShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		subtestData.params.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("PATCH failure - 403- permission denied - wrong application / user", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		subtestData := getDefaultMTOShipmentAndParams(nil)

		updateURI := "/mto-shipments/" + subtestData.mtoShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		unauthorizedReq = suite.AuthenticateOfficeRequest(unauthorizedReq, officeUser)
		subtestData.params.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentForbidden{}, response)
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("d874d002-5582-4a91-97d3-786e8f66c763"))
		subtestData.params.MtoShipmentID = *uuidString

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.Run("POST failure - 404 - not found - wrong SM does not match move", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		sm := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateUserRequest(req, sm.User)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		subtestData := getDefaultMTOShipmentAndParams(nil)

		subtestData.params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ShipmentUpdater{}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		subtestData := getDefaultMTOShipmentAndParams(&mockUpdater)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")
	})
}

//
// GET ALL
//

type mtoListSubtestData struct {
	shipments models.MTOShipments
	params    mtoshipmentops.ListMTOShipmentsParams
}

func (suite *HandlerSuite) makeListSubtestData() (subtestData *mtoListSubtestData) {
	subtestData = &mtoListSubtestData{}
	mto := factory.BuildMove(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)

	pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
	secondaryPickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "123 Nowhere",
				StreetAddress2: models.StringPointer("P.O. Box 5555"),
				StreetAddress3: models.StringPointer("c/o Some Other Person"),
				City:           "El Paso",
				State:          "TX",
				PostalCode:     "79916",
			},
		},
	}, nil)

	deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
	secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "5432 Everywhere",
				StreetAddress2: models.StringPointer("P.O. Box 111"),
				StreetAddress3: models.StringPointer("c/o Some Other Person"),
				City:           "Portsmouth",
				State:          "NH",
				PostalCode:     "03801",
			},
		},
	}, nil)

	mtoShipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:              models.MTOShipmentStatusSubmitted,
				RequestedPickupDate: &requestedPickupDate,
			},
		},
		{
			Model:    pickupAddress,
			Type:     &factory.Addresses.PickupAddress,
			LinkOnly: true,
		},
		{
			Model:    secondaryPickupAddress,
			Type:     &factory.Addresses.SecondaryPickupAddress,
			LinkOnly: true,
		},
		{
			Model:    deliveryAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
		{
			Model:    secondaryDeliveryAddress,
			Type:     &factory.Addresses.SecondaryDeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	ppmShipment2 := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, []factory.Trait{factory.GetTraitApprovedPPMShipment})

	advanceAmountRequested := unit.Cents(10000)
	ppmShipment3 := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				AdvanceAmountRequested: &advanceAmountRequested,
			},
		},
	}, nil)

	boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

	subtestData.shipments = models.MTOShipments{mtoShipment, mtoShipment2, ppmShipment.Shipment, ppmShipment2.Shipment, ppmShipment3.Shipment, boatShipment.Shipment}

	req := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/mto_shipments", mto.ID.String()), nil)
	req = suite.AuthenticateRequest(req, mto.Orders.ServiceMember)

	subtestData.params = mtoshipmentops.ListMTOShipmentsParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
	}

	return subtestData

}

func (suite *HandlerSuite) TestListMTOShipmentsHandler() {
	suite.Run("Successful list fetch - 200 - Integration Test", func() {
		subtestData := suite.makeListSubtestData()
		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mtoshipment.NewMTOShipmentFetcher(),
		}

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsOK{}, response)

		okResponse := response.(*mtoshipmentops.ListMTOShipmentsOK)
		suite.Len(okResponse.Payload, 5)

		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		for i, returnedShipment := range okResponse.Payload {
			expectedShipment := subtestData.shipments[i]

			// we expect the shipment that was created first to come first in the response
			suite.EqualUUID(expectedShipment.ID, returnedShipment.ID)

			suite.Equal(expectedShipment.Status, models.MTOShipmentStatus(returnedShipment.Status))

			if expectedShipment.ShipmentType == models.MTOShipmentTypePPM {
				suite.EqualUUID(expectedShipment.PPMShipment.ID, returnedShipment.PpmShipment.ID)
				suite.EqualUUID(expectedShipment.PPMShipment.ShipmentID, returnedShipment.PpmShipment.ShipmentID)
				suite.EqualDateTime(expectedShipment.PPMShipment.CreatedAt, returnedShipment.PpmShipment.CreatedAt)
				suite.Equal(string(expectedShipment.PPMShipment.Status), string(returnedShipment.PpmShipment.Status))
				suite.EqualDate(expectedShipment.PPMShipment.ExpectedDepartureDate, *returnedShipment.PpmShipment.ExpectedDepartureDate)
				suite.EqualDatePtr(expectedShipment.PPMShipment.ActualMoveDate, returnedShipment.PpmShipment.ActualMoveDate)
				suite.EqualDateTimePtr(expectedShipment.PPMShipment.SubmittedAt, returnedShipment.PpmShipment.SubmittedAt)
				suite.EqualDateTimePtr(expectedShipment.PPMShipment.ReviewedAt, returnedShipment.PpmShipment.ReviewedAt)
				suite.EqualDateTimePtr(expectedShipment.PPMShipment.ApprovedAt, returnedShipment.PpmShipment.ApprovedAt)
				suite.Equal(*expectedShipment.PPMShipment.SITExpected, *returnedShipment.PpmShipment.SitExpected)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.EstimatedWeight, returnedShipment.PpmShipment.EstimatedWeight)
				suite.Equal(expectedShipment.PPMShipment.HasProGear, returnedShipment.PpmShipment.HasProGear)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.ProGearWeight, returnedShipment.PpmShipment.ProGearWeight)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.SpouseProGearWeight, returnedShipment.PpmShipment.SpouseProGearWeight)
				suite.Equal(expectedShipment.PPMShipment.HasRequestedAdvance, returnedShipment.PpmShipment.HasRequestedAdvance)
				suite.EqualCentsPointers(expectedShipment.PPMShipment.AdvanceAmountRequested, returnedShipment.PpmShipment.AdvanceAmountRequested)

				if expectedShipment.PPMShipment.EstimatedIncentive != nil {
					suite.Equal(expectedShipment.PPMShipment.EstimatedIncentive.Int64(), *returnedShipment.PpmShipment.EstimatedIncentive)
				} else {
					suite.Nil(returnedShipment.PpmShipment.EstimatedIncentive)
				}

				continue // PPM Shipments won't have the rest of the fields below.
			}

			if expectedShipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway || expectedShipment.ShipmentType == models.MTOShipmentTypeBoatTowAway {
				suite.EqualUUID(expectedShipment.BoatShipment.ID, returnedShipment.BoatShipment.ID)
				suite.EqualUUID(expectedShipment.BoatShipment.ShipmentID, returnedShipment.BoatShipment.ShipmentID)
				suite.EqualDateTime(expectedShipment.BoatShipment.CreatedAt, returnedShipment.BoatShipment.CreatedAt)
				suite.Equal(string(expectedShipment.BoatShipment.Type), string(*returnedShipment.BoatShipment.Type))
				suite.Equal(expectedShipment.BoatShipment.Year, *returnedShipment.BoatShipment.Year)
				suite.Equal(expectedShipment.BoatShipment.Make, *returnedShipment.BoatShipment.Make)
				suite.Equal(expectedShipment.BoatShipment.Model, *returnedShipment.BoatShipment.Model)
				suite.Equal(expectedShipment.BoatShipment.LengthInInches, *returnedShipment.BoatShipment.LengthInInches)
				suite.Equal(expectedShipment.BoatShipment.WidthInInches, *returnedShipment.BoatShipment.WidthInInches)
				suite.Equal(expectedShipment.BoatShipment.HeightInInches, *returnedShipment.BoatShipment.HeightInInches)
				suite.Equal(expectedShipment.BoatShipment.HasTrailer, *returnedShipment.BoatShipment.HasTrailer)
				suite.Equal(expectedShipment.BoatShipment.IsRoadworthy, *returnedShipment.BoatShipment.IsRoadworthy)
			}

			suite.EqualDatePtr(expectedShipment.RequestedPickupDate, returnedShipment.RequestedPickupDate)

			suite.Equal(expectedShipment.PickupAddress.StreetAddress1, *returnedShipment.PickupAddress.StreetAddress1)
			suite.Equal(*expectedShipment.PickupAddress.StreetAddress2, *returnedShipment.PickupAddress.StreetAddress2)
			suite.Equal(*expectedShipment.PickupAddress.StreetAddress3, *returnedShipment.PickupAddress.StreetAddress3)
			suite.Equal(expectedShipment.PickupAddress.City, *returnedShipment.PickupAddress.City)
			suite.Equal(expectedShipment.PickupAddress.State, *returnedShipment.PickupAddress.State)
			suite.Equal(expectedShipment.PickupAddress.PostalCode, *returnedShipment.PickupAddress.PostalCode)

			if expectedShipment.SecondaryPickupAddress != nil {
				suite.Equal(expectedShipment.SecondaryPickupAddress.StreetAddress1, *returnedShipment.SecondaryPickupAddress.StreetAddress1)
				suite.Equal(*expectedShipment.SecondaryPickupAddress.StreetAddress2, *returnedShipment.SecondaryPickupAddress.StreetAddress2)
				suite.Equal(*expectedShipment.SecondaryPickupAddress.StreetAddress3, *returnedShipment.SecondaryPickupAddress.StreetAddress3)
				suite.Equal(expectedShipment.SecondaryPickupAddress.City, *returnedShipment.SecondaryPickupAddress.City)
				suite.Equal(expectedShipment.SecondaryPickupAddress.State, *returnedShipment.SecondaryPickupAddress.State)
				suite.Equal(expectedShipment.SecondaryPickupAddress.PostalCode, *returnedShipment.SecondaryPickupAddress.PostalCode)
			}

			suite.Equal(expectedShipment.DestinationAddress.StreetAddress1, *returnedShipment.DestinationAddress.StreetAddress1)
			suite.Equal(*expectedShipment.DestinationAddress.StreetAddress2, *returnedShipment.DestinationAddress.StreetAddress2)
			suite.Equal(*expectedShipment.DestinationAddress.StreetAddress3, *returnedShipment.DestinationAddress.StreetAddress3)
			suite.Equal(expectedShipment.DestinationAddress.City, *returnedShipment.DestinationAddress.City)
			suite.Equal(expectedShipment.DestinationAddress.State, *returnedShipment.DestinationAddress.State)
			suite.Equal(expectedShipment.DestinationAddress.PostalCode, *returnedShipment.DestinationAddress.PostalCode)

			if expectedShipment.SecondaryDeliveryAddress != nil {
				suite.Equal(expectedShipment.SecondaryDeliveryAddress.StreetAddress1, *returnedShipment.SecondaryDeliveryAddress.StreetAddress1)
				suite.Equal(*expectedShipment.SecondaryDeliveryAddress.StreetAddress2, *returnedShipment.SecondaryDeliveryAddress.StreetAddress2)
				suite.Equal(*expectedShipment.SecondaryDeliveryAddress.StreetAddress3, *returnedShipment.SecondaryDeliveryAddress.StreetAddress3)
				suite.Equal(expectedShipment.SecondaryDeliveryAddress.City, *returnedShipment.SecondaryDeliveryAddress.City)
				suite.Equal(expectedShipment.SecondaryDeliveryAddress.State, *returnedShipment.SecondaryDeliveryAddress.State)
				suite.Equal(expectedShipment.SecondaryDeliveryAddress.PostalCode, *returnedShipment.SecondaryDeliveryAddress.PostalCode)
			}
		}
	})

	suite.Run("POST failure - 400 - Bad Request", func() {
		subtestData := suite.makeListSubtestData()
		emtpyMTOID := mtoshipmentops.ListMTOShipmentsParams{
			HTTPRequest:     subtestData.params.HTTPRequest,
			MoveTaskOrderID: "",
		}
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}
		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mockMTOShipmentFetcher,
		}

		response := handler.Handle(emtpyMTOID)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsBadRequest{}, response)
	})

	suite.Run("POST failure - 401 - permission denied - not authenticated", func() {
		subtestData := suite.makeListSubtestData()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		unauthorizedReq := suite.AuthenticateOfficeRequest(subtestData.params.HTTPRequest, officeUser)
		unauthorizedParams := mtoshipmentops.ListMTOShipmentsParams{
			HTTPRequest:     unauthorizedReq,
			MoveTaskOrderID: *handlers.FmtUUID(subtestData.shipments[0].MoveTaskOrderID),
		}
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}
		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mockMTOShipmentFetcher,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsUnauthorized{}, response)
	})

	suite.Run("Failure list fetch - 404 Not Found - service member user not authorized", func() {
		subtestData := suite.makeListSubtestData()
		unauthorizedUser := factory.BuildServiceMember(suite.DB(), nil, nil)
		unauthorizedReq := suite.AuthenticateRequest(subtestData.params.HTTPRequest, unauthorizedUser)
		unauthorizedParams := mtoshipmentops.ListMTOShipmentsParams{
			HTTPRequest:     unauthorizedReq,
			MoveTaskOrderID: *handlers.FmtUUID(subtestData.shipments[0].MoveTaskOrderID),
		}

		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mtoshipment.NewMTOShipmentFetcher(),
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsNotFound{}, response)
	})

	suite.Run("Failure list fetch - 500 Internal Server Error", func() {
		subtestData := suite.makeListSubtestData()
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}
		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mockMTOShipmentFetcher,
		}

		internalServerErr := errors.New("ServerError")

		mockMTOShipmentFetcher.On("ListMTOShipments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return([]models.MTOShipment{}, internalServerErr)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsInternalServerError{}, response)
	})
}

//
// DELETE
//

func (suite *HandlerSuite) TestDeleteShipmentHandler() {
	suite.Run("Returns 204 when all validations pass", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(shipment.MoveTaskOrderID, nil)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)

		suite.IsType(&mtoshipmentops.DeleteShipmentNoContent{}, response)
	})

	suite.Run("Returns 404 when deleter returns NotFoundError", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.NotFoundError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteShipmentNotFound{}, response)
	})

	suite.Run("Returns 409 - Conflict error", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ConflictError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)

		suite.IsType(&mtoshipmentops.DeleteShipmentConflict{}, response)
	})

	suite.Run("Returns 403 when servicemember ID doesn't match shipment", func() {
		sm1 := factory.BuildServiceMember(nil, nil, []factory.Trait{factory.GetTraitServiceMemberSetIDs})
		sm2 := factory.BuildServiceMember(nil, nil, []factory.Trait{factory.GetTraitServiceMemberSetIDs})
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
			{
				Model:    sm2,
				LinkOnly: true,
			},
		}, nil)

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ForbiddenError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, sm1)

		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteShipmentForbidden{}, response)
	})

	suite.Run("Returns 422 - Unprocessable Enitity error", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.UnprocessableEntityError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteShipmentUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when deleter returns InternalServerError", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.InternalServerError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteShipmentInternalServerError{}, response)
	})
}
