package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
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
	moveRouter := moverouter.NewMoveRouter()

	return &mtoShipmentObjects{
		builder:    builder,
		fetcher:    fetcher,
		moveRouter: moveRouter,
	}
}

//
// CREATE
//

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeCreateSubtestData.
	testMTOShipmentObjects := suite.setUpMTOShipmentObjects()

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(testMTOShipmentObjects.builder, testMTOShipmentObjects.fetcher, testMTOShipmentObjects.moveRouter)
	ppmEstimator := mocks.PPMEstimator{}
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)

	shipmentCreator := shipmentorchestrator.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator)

	type mtoCreateSubtestData struct {
		serviceMember models.ServiceMember
		pickupAddress models.Address
		mtoShipment   models.MTOShipment
		params        mtoshipmentops.CreateMTOShipmentParams
		handler       CreateMTOShipmentHandler
	}

	makeCreateSubtestData := func(appCtx appcontext.AppContext) (subtestData mtoCreateSubtestData) {
		db := appCtx.DB()

		subtestData.serviceMember = testdatagen.MakeDefaultServiceMember(db)

		mto := testdatagen.MakeMove(db, testdatagen.Assertions{
			Order: models.Order{
				ServiceMember:   subtestData.serviceMember, // have to set here too because MakeOrder isn't checking the right assertion.
				ServiceMemberID: subtestData.serviceMember.ID,
			},
			ServiceMember: subtestData.serviceMember,
		})

		subtestData.pickupAddress = testdatagen.MakeDefaultAddress(db)
		secondaryPickupAddress := testdatagen.MakeAddress2(db, testdatagen.Assertions{})

		destinationAddress := testdatagen.MakeAddress3(db, testdatagen.Assertions{})
		secondaryDeliveryAddress := testdatagen.MakeAddress4(db, testdatagen.Assertions{})

		subtestData.mtoShipment = testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			Move:        mto,
			MTOShipment: models.MTOShipment{},
		})
		subtestData.mtoShipment.MoveTaskOrderID = mto.ID

		mtoAgent := testdatagen.MakeDefaultMTOAgent(db)
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
					Country:        subtestData.pickupAddress.Country,
					PostalCode:     &subtestData.pickupAddress.PostalCode,
					State:          &subtestData.pickupAddress.State,
					StreetAddress1: &subtestData.pickupAddress.StreetAddress1,
					StreetAddress2: subtestData.pickupAddress.StreetAddress2,
					StreetAddress3: subtestData.pickupAddress.StreetAddress3,
				},
				SecondaryPickupAddress: &internalmessages.Address{
					City:           &secondaryPickupAddress.City,
					Country:        secondaryPickupAddress.Country,
					PostalCode:     &secondaryPickupAddress.PostalCode,
					State:          &secondaryPickupAddress.State,
					StreetAddress1: &secondaryPickupAddress.StreetAddress1,
					StreetAddress2: secondaryPickupAddress.StreetAddress2,
					StreetAddress3: secondaryPickupAddress.StreetAddress3,
				},
				DestinationAddress: &internalmessages.Address{
					City:           &destinationAddress.City,
					Country:        destinationAddress.Country,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
				},
				SecondaryDeliveryAddress: &internalmessages.Address{
					City:           &secondaryDeliveryAddress.City,
					Country:        secondaryDeliveryAddress.Country,
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

		subtestData.handler = CreateMTOShipmentHandler{
			handlers.NewHandlerConfig(appCtx.DB(), appCtx.Logger()),
			shipmentCreator,
		}

		return subtestData
	}

	suite.Run("Successful POST - Integration Test - HHG", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		params := subtestData.params

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(internalmessages.MTOShipmentTypeHHG, createdShipment.ShipmentType)
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

	suite.Run("Successful POST - Integration Test - PPM", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM
		// pointers
		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		pickupPostal := "11111"
		destinationPostalCode := "41414"
		sitExpected := false
		// reset Body params to have PPM fields
		params.Body = &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID),
			PpmShipment: &internalmessages.CreatePPMShipment{
				ExpectedDepartureDate: &expectedDepartureDate,
				PickupPostalCode:      &pickupPostal,
				DestinationPostalCode: &destinationPostalCode,
				SitExpected:           &sitExpected,
			},
			ShipmentType: &ppmShipmentType,
		}

		// When a customer first creates a move, there is not enough data to calculate an incentive yet.
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil).Once()

		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(internalmessages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(*params.Body.MoveTaskOrderID, createdShipment.MoveTaskOrderID)
		suite.Equal(*params.Body.PpmShipment.ExpectedDepartureDate, *createdShipment.PpmShipment.ExpectedDepartureDate)
		suite.Equal(*params.Body.PpmShipment.PickupPostalCode, *createdShipment.PpmShipment.PickupPostalCode)
		suite.Equal(*params.Body.PpmShipment.DestinationPostalCode, *createdShipment.PpmShipment.DestinationPostalCode)
		suite.Equal(*params.Body.PpmShipment.SitExpected, *createdShipment.PpmShipment.SitExpected)
	})

	suite.Run("Successful POST - Integration Test - NTS-Release", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		params := subtestData.params

		// Set fields appropriately for NTS-Release
		ntsrShipmentType := internalmessages.MTOShipmentTypeHHGOUTOFNTSDOMESTIC
		params.Body.ShipmentType = &ntsrShipmentType
		params.Body.RequestedPickupDate = strfmt.Date(time.Time{})

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(ntsrShipmentType, createdShipment.ShipmentType)
		suite.Equal(*params.Body.CustomerRemarks, *createdShipment.CustomerRemarks)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *createdShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *createdShipment.SecondaryPickupAddress.StreetAddress1)
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

	suite.Run("POST failure - 400 - invalid input, missing pickup address", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		badParams := subtestData.params
		badParams.Body.PickupAddress = nil

		response := subtestData.handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 401- permission denied - not authenticated", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

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
					Country:        subtestData.pickupAddress.Country,
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

	suite.Run("POST failure - 403- permission denied - wrong application", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("POST failure - 404 -- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := subtestData.params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := subtestData.handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		otherParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: subtestData.params.HTTPRequest,
		}
		response := subtestData.handler.Handle(otherParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})

	suite.Run("POST failure - 400 -- missing required field to Create PPM", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM

		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		pickupPostal := "11111"
		destinationPostalCode := "41414"
		sitExpected := false
		badID, _ := uuid.NewV4()

		// reset Body params to have PPM fields
		params.Body = &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(badID),
			PpmShipment: &internalmessages.CreatePPMShipment{
				ExpectedDepartureDate: &expectedDepartureDate,
				PickupPostalCode:      &pickupPostal,
				DestinationPostalCode: &destinationPostalCode,
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
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx)

		mockShipmentCreator := mocks.ShipmentCreator{}

		err := errors.New("ServerError")

		mockShipmentCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil, err)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerConfig(appCtx.DB(), appCtx.Logger()),
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

	mtoShipmentUpdater := mtoshipment.NewMTOShipmentUpdater(testMTOShipmentObjects.builder, testMTOShipmentObjects.fetcher, planner, testMTOShipmentObjects.moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator)

	ppmEstimator := mocks.PPMEstimator{}

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator)

	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

	authRequestAndSetUpHandlerAndParams := func(appCtx appcontext.AppContext, originalShipment models.MTOShipment, mockShipmentUpdater *mocks.ShipmentUpdater) (UpdateMTOShipmentHandler, mtoshipmentops.UpdateMTOShipmentParams) {
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
			handlers.NewHandlerConfig(appCtx.DB(), appCtx.Logger()),
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
	getDefaultMTOShipmentAndParams := func(appCtx appcontext.AppContext, mockShipmentUpdater *mocks.ShipmentUpdater) *mtoUpdateSubtestData {
		originalShipment := testdatagen.MakeDefaultMTOShipment(appCtx.DB())

		pickupAddress := testdatagen.MakeDefaultAddress(appCtx.DB())
		pickupAddress.StreetAddress1 = "123 Fake Test St NW"

		secondaryPickupAddress := testdatagen.MakeDefaultAddress(appCtx.DB())
		secondaryPickupAddress.StreetAddress1 = "89999 Other Test St NW"

		destinationAddress := testdatagen.MakeDefaultAddress(appCtx.DB())
		destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"

		secondaryDeliveryAddress := testdatagen.MakeDefaultAddress(appCtx.DB())
		secondaryDeliveryAddress.StreetAddress1 = "9999 Test Fake Rd SE"

		mtoAgent := testdatagen.MakeDefaultMTOAgent(appCtx.DB())
		agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
			FirstName: mtoAgent.FirstName,
			LastName:  mtoAgent.LastName,
			Email:     mtoAgent.Email,
			Phone:     mtoAgent.Phone,
			AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
		}}

		customerRemarks := ""

		handler, params := authRequestAndSetUpHandlerAndParams(appCtx, originalShipment, mockShipmentUpdater)

		params.Body = &internalmessages.UpdateShipment{
			Agents:          agents,
			CustomerRemarks: &customerRemarks,
			DestinationAddress: &internalmessages.Address{
				City:           &destinationAddress.City,
				Country:        destinationAddress.Country,
				PostalCode:     &destinationAddress.PostalCode,
				State:          &destinationAddress.State,
				StreetAddress1: &destinationAddress.StreetAddress1,
				StreetAddress2: destinationAddress.StreetAddress2,
				StreetAddress3: destinationAddress.StreetAddress3,
			},
			SecondaryDeliveryAddress: &internalmessages.Address{
				City:           &secondaryDeliveryAddress.City,
				Country:        secondaryDeliveryAddress.Country,
				PostalCode:     &secondaryDeliveryAddress.PostalCode,
				State:          &secondaryDeliveryAddress.State,
				StreetAddress1: &secondaryDeliveryAddress.StreetAddress1,
				StreetAddress2: secondaryDeliveryAddress.StreetAddress2,
				StreetAddress3: secondaryDeliveryAddress.StreetAddress3,
			},
			PickupAddress: &internalmessages.Address{
				City:           &pickupAddress.City,
				Country:        pickupAddress.Country,
				PostalCode:     &pickupAddress.PostalCode,
				State:          &pickupAddress.State,
				StreetAddress1: &pickupAddress.StreetAddress1,
				StreetAddress2: pickupAddress.StreetAddress2,
				StreetAddress3: pickupAddress.StreetAddress3,
			},
			SecondaryPickupAddress: &internalmessages.Address{
				City:           &secondaryPickupAddress.City,
				Country:        secondaryPickupAddress.Country,
				PostalCode:     &secondaryPickupAddress.PostalCode,
				State:          &secondaryPickupAddress.State,
				StreetAddress1: &secondaryPickupAddress.StreetAddress1,
				StreetAddress2: secondaryPickupAddress.StreetAddress2,
				StreetAddress3: secondaryPickupAddress.StreetAddress3,
			},
			RequestedPickupDate:   strfmt.Date(*originalShipment.RequestedPickupDate),
			RequestedDeliveryDate: strfmt.Date(*originalShipment.RequestedDeliveryDate),
			ShipmentType:          internalmessages.MTOShipmentTypeHHG,
		}

		return &mtoUpdateSubtestData{
			mtoShipment: &originalShipment,
			params:      params,
			handler:     handler,
		}
	}

	suite.Run("Successful PATCH - Integration Test", func() {
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)
		params := subtestData.params

		response := subtestData.handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		suite.Equal(subtestData.mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(*params.Body.CustomerRemarks, *updatedShipment.CustomerRemarks)
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
			suite.Equal(originalShipment.PPMShipment.PickupPostalCode, *updatedShipment.PpmShipment.PickupPostalCode)
			suite.Equal(originalShipment.PPMShipment.DestinationPostalCode, *updatedShipment.PpmShipment.DestinationPostalCode)
			suite.Equal(originalShipment.PPMShipment.SitExpected, updatedShipment.PpmShipment.SitExpected)
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
			suite.Equal(originalShipment.PPMShipment.AdvanceRequested, updatedShipment.PpmShipment.AdvanceRequested)
			suite.EqualCentsPointers(originalShipment.PPMShipment.Advance, updatedShipment.PpmShipment.Advance)
		}

		type setUpOriginalPPMFunc func(appCtx appcontext.AppContext) models.PPMShipment
		type runChecksFunc func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment)

		ppmUpdateTestCases := map[string]struct {
			setUpOriginalPPM   setUpOriginalPPMFunc
			desiredShipment    internalmessages.UpdatePPMShipment
			estimatedIncentive *unit.Cents
			runChecks          runChecksFunc
		}{
			"Edit estimated dates & locations": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							ExpectedDepartureDate: time.Date(testdatagen.GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
							PickupPostalCode:      "90808",
							DestinationPostalCode: "79912",
							SitExpected:           models.BoolPointer(true),
						},
					})
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					ExpectedDepartureDate: handlers.FmtDate(time.Date(testdatagen.GHCTestYear, time.April, 27, 0, 0, 0, 0, time.UTC)),
					PickupPostalCode:      handlers.FmtString("90900"),
					DestinationPostalCode: handlers.FmtString("79916"),
					SitExpected:           handlers.FmtBool(false),
				},
				estimatedIncentive: nil,
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check all fields changed as expected
					desiredShipment.ExpectedDepartureDate.Equal(*updatedShipment.PpmShipment.ExpectedDepartureDate)

					suite.Equal(desiredShipment.PickupPostalCode, updatedShipment.PpmShipment.PickupPostalCode)
					suite.Equal(desiredShipment.DestinationPostalCode, updatedShipment.PpmShipment.DestinationPostalCode)
					suite.Equal(desiredShipment.SitExpected, updatedShipment.PpmShipment.SitExpected)
				},
			},
			"Edit estimated dates & locations - add secondary zips": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					SecondaryPickupPostalCode:      handlers.FmtString("90900"),
					SecondaryDestinationPostalCode: handlers.FmtString("79916"),
				},
				estimatedIncentive: nil,
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)

					// check new fields were set
					suite.Equal(desiredShipment.SecondaryPickupPostalCode, updatedShipment.PpmShipment.SecondaryPickupPostalCode)
					suite.Equal(desiredShipment.SecondaryDestinationPostalCode, updatedShipment.PpmShipment.SecondaryDestinationPostalCode)
				},
			},
			"Edit estimated dates & locations - remove secondary zips": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							SecondaryPickupPostalCode:      models.StringPointer("90900"),
							SecondaryDestinationPostalCode: models.StringPointer("79916"),
						},
					})
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					SecondaryPickupPostalCode:      handlers.FmtString(""), // TODO: Update to pass null since that's what the FE will actually send.
					SecondaryDestinationPostalCode: handlers.FmtString(""), // TODO: Update to pass null since that's what the FE will actually send.
				},
				estimatedIncentive: nil,
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Nil(updatedShipment.PpmShipment.SecondaryPickupPostalCode)
					suite.Nil(updatedShipment.PpmShipment.SecondaryDestinationPostalCode)
				},
			},
			"Add estimated weights - no pro gear": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())
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
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalDefaultPPMShipment(appCtx.DB())
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
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:     models.PoundPointer(4000),
							HasProGear:          models.BoolPointer(true),
							ProGearWeight:       models.PoundPointer(1250),
							SpouseProGearWeight: models.PoundPointer(150),
							EstimatedIncentive:  models.CentPointer(unit.Cents(500000)),
						},
					})
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
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:    models.PoundPointer(4000),
							HasProGear:         models.BoolPointer(false),
							EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
						},
					})
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					AdvanceRequested: handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.AdvanceRequested, updatedShipment.PpmShipment.AdvanceRequested)
					suite.Nil(updatedShipment.PpmShipment.Advance)
				},
			},
			"Add advance requested info - yes advance": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:    models.PoundPointer(4000),
							HasProGear:         models.BoolPointer(false),
							EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
						},
					})
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					AdvanceRequested: handlers.FmtBool(true),
					Advance:          handlers.FmtInt64(200000),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.AdvanceRequested, updatedShipment.PpmShipment.AdvanceRequested)
					suite.Equal(desiredShipment.Advance, updatedShipment.PpmShipment.Advance)
				},
			},
			"Remove advance requested": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:    models.PoundPointer(4000),
							HasProGear:         models.BoolPointer(false),
							EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
							AdvanceRequested:   models.BoolPointer(true),
							Advance:            models.CentPointer(unit.Cents(200000)),
						},
					})
				},
				desiredShipment: internalmessages.UpdatePPMShipment{
					AdvanceRequested: handlers.FmtBool(false),
				},
				estimatedIncentive: models.CentPointer(unit.Cents(500000)),
				runChecks: func(updatedShipment *internalmessages.MTOShipment, originalShipment models.MTOShipment, desiredShipment internalmessages.UpdatePPMShipment) {
					// check existing fields didn't change
					checkDatesAndLocationsDidntChange(updatedShipment, originalShipment)
					checkEstimatedWeightsDidntChange(updatedShipment, originalShipment)

					// check expected fields were updated
					suite.Equal(desiredShipment.AdvanceRequested, updatedShipment.PpmShipment.AdvanceRequested)
					suite.Nil(updatedShipment.PpmShipment.Advance)
				},
			},
			"Add actual zips and advance info - no advance": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:    models.PoundPointer(4000),
							HasProGear:         models.BoolPointer(false),
							EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
							AdvanceRequested:   models.BoolPointer(true),
							Advance:            models.CentPointer(unit.Cents(200000)),
						},
					})
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
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:    models.PoundPointer(4000),
							HasProGear:         models.BoolPointer(false),
							EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
							AdvanceRequested:   models.BoolPointer(true),
							Advance:            models.CentPointer(unit.Cents(200000)),
						},
					})
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
			"Remove actual advance": {
				setUpOriginalPPM: func(appCtx appcontext.AppContext) models.PPMShipment {
					return testdatagen.MakeMinimalPPMShipment(appCtx.DB(), testdatagen.Assertions{
						PPMShipment: models.PPMShipment{
							EstimatedWeight:             models.PoundPointer(4000),
							HasProGear:                  models.BoolPointer(false),
							EstimatedIncentive:          models.CentPointer(unit.Cents(500000)),
							AdvanceRequested:            models.BoolPointer(true),
							Advance:                     models.CentPointer(unit.Cents(200000)),
							ActualPickupPostalCode:      models.StringPointer("90210"),
							ActualDestinationPostalCode: models.StringPointer("90210"),
							HasReceivedAdvance:          models.BoolPointer(true),
							AdvanceAmountReceived:       models.CentPointer(unit.Cents(250000)),
						},
					})
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
				appCtx := suite.AppContextForTest()

				ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("models.PPMShipment"),
					mock.AnythingOfType("*models.PPMShipment")).
					Return(tc.estimatedIncentive, nil).Once()

				originalPPMShipment := tc.setUpOriginalPPM(appCtx)

				handler, params := authRequestAndSetUpHandlerAndParams(appCtx, originalPPMShipment.Shipment, nil)

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
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

		expectedStatus := internalmessages.MTOShipmentStatusSUBMITTED

		subtestData.params.Body.Status = expectedStatus

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)

		suite.Equal(expectedStatus, updatedResponse.Payload.Status)
	})

	suite.Run("PATCH failure - 400 -- nil body", func() {
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

		subtestData.params.Body = nil

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.Run("PATCH failure - 400 -- invalid requested status update", func() {
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

		subtestData.params.Body.Status = internalmessages.MTOShipmentStatusREJECTED

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.Run("PATCH failure - 401- permission denied - not authenticated", func() {
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

		updateURI := "/mto-shipments/" + subtestData.mtoShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		subtestData.params.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("PATCH failure - 403- permission denied - wrong application / user", func() {
		appCtx := suite.AppContextForTest()

		officeUser := testdatagen.MakeDefaultOfficeUser(appCtx.DB())

		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

		updateURI := "/mto-shipments/" + subtestData.mtoShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		unauthorizedReq = suite.AuthenticateOfficeRequest(unauthorizedReq, officeUser)
		subtestData.params.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentForbidden{}, response)
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("d874d002-5582-4a91-97d3-786e8f66c763"))
		subtestData.params.MtoShipmentID = *uuidString

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), nil)

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
		).Return(nil, err)

		subtestData := getDefaultMTOShipmentAndParams(suite.AppContextForTest(), &mockUpdater)

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
	mto := testdatagen.MakeDefaultMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})

	requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)

	pickupAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "123 Nowhere",
			StreetAddress2: swag.String("P.O. Box 5555"),
			StreetAddress3: swag.String("c/o Some Other Person"),
			City:           "El Paso",
			State:          "TX",
			PostalCode:     "79916",
			Country:        swag.String("US"),
		},
	})

	deliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "5432 Everywhere",
			StreetAddress2: swag.String("P.O. Box 111"),
			StreetAddress3: swag.String("c/o Some Other Person"),
			City:           "Portsmouth",
			State:          "NH",
			PostalCode:     "03801",
			Country:        swag.String("US"),
		},
	})

	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status:              models.MTOShipmentStatusSubmitted,
			RequestedPickupDate: &requestedPickupDate,
		},
		PickupAddress:            pickupAddress,
		SecondaryPickupAddress:   secondaryPickupAddress,
		DestinationAddress:       deliveryAddress,
		SecondaryDeliveryAddress: secondaryDeliveryAddress,
	})

	ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})

	ppmShipment2 := testdatagen.MakeApprovedPPMShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})

	advance := unit.Cents(10000)
	ppmShipment3 := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		Move:        mto,
		PPMShipment: models.PPMShipment{Advance: &advance},
	})

	subtestData.shipments = models.MTOShipments{mtoShipment, mtoShipment2, ppmShipment.Shipment, ppmShipment2.Shipment, ppmShipment3.Shipment}
	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	req := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/mto_shipments", mto.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

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
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
				suite.EqualDateTime(expectedShipment.PPMShipment.UpdatedAt, returnedShipment.PpmShipment.UpdatedAt)
				suite.Equal(string(expectedShipment.PPMShipment.Status), string(returnedShipment.PpmShipment.Status))
				suite.EqualDate(expectedShipment.PPMShipment.ExpectedDepartureDate, *returnedShipment.PpmShipment.ExpectedDepartureDate)
				suite.EqualDatePtr(expectedShipment.PPMShipment.ActualMoveDate, returnedShipment.PpmShipment.ActualMoveDate)
				suite.EqualDateTimePtr(expectedShipment.PPMShipment.SubmittedAt, returnedShipment.PpmShipment.SubmittedAt)
				suite.EqualDateTimePtr(expectedShipment.PPMShipment.ReviewedAt, returnedShipment.PpmShipment.ReviewedAt)
				suite.EqualDateTimePtr(expectedShipment.PPMShipment.ApprovedAt, returnedShipment.PpmShipment.ApprovedAt)
				suite.Equal(expectedShipment.PPMShipment.PickupPostalCode, *returnedShipment.PpmShipment.PickupPostalCode)
				suite.Equal(expectedShipment.PPMShipment.SecondaryPickupPostalCode, returnedShipment.PpmShipment.SecondaryPickupPostalCode)
				suite.Equal(expectedShipment.PPMShipment.DestinationPostalCode, *returnedShipment.PpmShipment.DestinationPostalCode)
				suite.Equal(expectedShipment.PPMShipment.SecondaryDestinationPostalCode, returnedShipment.PpmShipment.SecondaryDestinationPostalCode)
				suite.Equal(*expectedShipment.PPMShipment.SitExpected, *returnedShipment.PpmShipment.SitExpected)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.EstimatedWeight, returnedShipment.PpmShipment.EstimatedWeight)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.NetWeight, returnedShipment.PpmShipment.NetWeight)
				suite.Equal(expectedShipment.PPMShipment.HasProGear, returnedShipment.PpmShipment.HasProGear)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.ProGearWeight, returnedShipment.PpmShipment.ProGearWeight)
				suite.EqualPoundPointers(expectedShipment.PPMShipment.SpouseProGearWeight, returnedShipment.PpmShipment.SpouseProGearWeight)
				if expectedShipment.PPMShipment.EstimatedIncentive != nil {
					suite.Equal(expectedShipment.PPMShipment.EstimatedIncentive.Int64(), *returnedShipment.PpmShipment.EstimatedIncentive)
				} else {
					suite.Nil(returnedShipment.PpmShipment.EstimatedIncentive)
				}
				if expectedShipment.PPMShipment.Advance != nil {
					suite.Equal(expectedShipment.PPMShipment.Advance.Int64(), *returnedShipment.PpmShipment.Advance)
				} else {
					suite.Nil(returnedShipment.PpmShipment.Advance)
				}
				continue // PPM Shipments won't have the rest of the fields below.
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
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mockMTOShipmentFetcher,
		}

		response := handler.Handle(emtpyMTOID)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsBadRequest{}, response)
	})

	suite.Run("POST failure - 401 - permission denied - not authenticated", func() {
		subtestData := suite.makeListSubtestData()
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		unauthorizedReq := suite.AuthenticateOfficeRequest(subtestData.params.HTTPRequest, officeUser)
		unauthorizedParams := mtoshipmentops.ListMTOShipmentsParams{
			HTTPRequest:     unauthorizedReq,
			MoveTaskOrderID: *handlers.FmtUUID(subtestData.shipments[0].MoveTaskOrderID),
		}
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mockMTOShipmentFetcher,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsUnauthorized{}, response)
	})

	suite.Run("Failure list fetch - 500 Internal Server Error", func() {
		subtestData := suite.makeListSubtestData()
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		order := testdatagen.MakeDefaultOrder(suite.DB())
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
		})

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(shipment.MoveTaskOrderID, nil)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
		order := testdatagen.MakeDefaultOrder(suite.DB())
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
		})

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.NotFoundError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
		order := testdatagen.MakeDefaultOrder(suite.DB())
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
		})

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ConflictError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
		sm1 := testdatagen.MakeStubbedServiceMember(suite.DB())
		sm2 := testdatagen.MakeStubbedServiceMember(suite.DB())
		order := testdatagen.MakeDefaultOrder(suite.DB())
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
			ServiceMember: sm2,
		})

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ForbiddenError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, sm1)

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
		order := testdatagen.MakeDefaultOrder(suite.DB())
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
		})

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.UnprocessableEntityError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
		order := testdatagen.MakeDefaultOrder(suite.DB())
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
			},
		})

		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.InternalServerError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
