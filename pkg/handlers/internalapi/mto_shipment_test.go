package internalapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"

	"net/http/httptest"

	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

//
// CREATE
//

type mtoCreateSubtestData struct {
	serviceMember models.ServiceMember
	pickupAddress models.Address
	mtoShipment   models.MTOShipment
	builder       *query.Builder
	params        mtoshipmentops.CreateMTOShipmentParams
}

func (suite *HandlerSuite) makeCreateSubtestData() (subtestData *mtoCreateSubtestData) {
	subtestData = &mtoCreateSubtestData{}
	db := suite.DB()
	mto := testdatagen.MakeDefaultMove(db)

	subtestData.serviceMember = testdatagen.MakeDefaultServiceMember(db)

	subtestData.pickupAddress = testdatagen.MakeDefaultAddress(db)
	secondaryPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})

	destinationAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})

	subtestData.mtoShipment = testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move:        mto,
		MTOShipment: models.MTOShipment{},
	})
	subtestData.mtoShipment.MoveTaskOrderID = mto.ID

	mtoAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
	}}

	customerRemarks := "I have some grandfather clocks."

	subtestData.builder = query.NewQueryBuilder()

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

	return subtestData
}

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	moveRouter := moverouter.NewMoveRouter()

	suite.Run("Successful POST - Integration Test - HHG", func() {
		subtestData := suite.makeCreateSubtestData()
		params := subtestData.params
		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}
		response := handler.Handle(subtestData.params)

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
		subtestData := suite.makeCreateSubtestData()
		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM
		// pointers
		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		pickupPostal := "11111"
		destinationPostalCode := "41414"
		sitExpected := false
		// Reset Shipment Type to PPM from default (HHG)
		params.Body.ShipmentType = &ppmShipmentType
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

		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		response := handler.Handle(params)
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
		subtestData := suite.makeCreateSubtestData()
		params := subtestData.params

		// Set fields appropriately for NTS-Release
		ntsrShipmentType := internalmessages.MTOShipmentTypeHHGOUTOFNTSDOMESTIC
		params.Body.ShipmentType = &ntsrShipmentType
		params.Body.RequestedPickupDate = strfmt.Date(time.Time{})

		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}
		response := handler.Handle(subtestData.params)

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
		subtestData := suite.makeCreateSubtestData()
		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		badParams := subtestData.params
		badParams.Body.PickupAddress = nil

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 401- permission denied - not authenticated", func() {
		subtestData := suite.makeCreateSubtestData()
		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)

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

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("POST failure - 403- permission denied - wrong application", func() {
		subtestData := suite.makeCreateSubtestData()
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("POST failure - 404 -- not found", func() {
		subtestData := suite.makeCreateSubtestData()

		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := subtestData.params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		subtestData := suite.makeCreateSubtestData()
		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		otherParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: subtestData.params.HTTPRequest,
		}
		response := handler.Handle(otherParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})

	suite.Run("POST failure - 400 -- missing required field to Create PPM", func() {
		subtestData := suite.makeCreateSubtestData()
		fetcher := fetch.NewFetcher(subtestData.builder)
		creator := mtoshipment.NewMTOShipmentCreator(subtestData.builder, fetcher, moveRouter)
		ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(creator)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			ppmShipmentCreator,
		}

		//otherParams := mtoshipmentops.CreateMTOShipmentParams{
		//	HTTPRequest: subtestData.params.HTTPRequest,
		//}

		params := subtestData.params
		ppmShipmentType := internalmessages.MTOShipmentTypePPM
		// pointers
		expectedDepartureDate := strfmt.Date(*subtestData.mtoShipment.RequestedPickupDate)
		pickupPostal := "11111"
		destinationPostalCode := "41414"
		sitExpected := false
		badID, _ := uuid.NewV4()
		//reason := "invalid memory address or nil pointer dereference"
		params.Body.ShipmentType = &ppmShipmentType
		// reset Body params to have PPM fields
		params.Body = &internalmessages.CreateShipment{
			//MoveTaskOrderID: handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID),
			MoveTaskOrderID: handlers.FmtUUID(badID),
			PpmShipment: &internalmessages.CreatePPMShipment{
				ExpectedDepartureDate: &expectedDepartureDate,
				PickupPostalCode:      &pickupPostal,
				DestinationPostalCode: &destinationPostalCode,
				SitExpected:           &sitExpected,
			},
			ShipmentType: &ppmShipmentType,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentNotFound).Payload
		suite.Equal(handlers.NotFoundMessage, *errResponse.Title)

		// Check Error details
		suite.Contains(*errResponse.Detail, "not found for move")
	})

	suite.Run("POST failure - 500", func() {
		subtestData := suite.makeCreateSubtestData()
		mockCreator := mocks.MTOShipmentCreator{}
		mockPPMShipmentCreator := mocks.PPMShipmentCreator{}
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockCreator,
			&mockPPMShipmentCreator,
		}

		err := errors.New("ServerError")

		mockCreator.On("CreateMTOShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		response := handler.Handle(subtestData.params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})
}

//
// UPDATE
//

// getDefaultMTOShipmentAndParams generates a set of default params and an MTOShipment
func (suite *HandlerSuite) getDefaultMTOShipmentAndParams() (mtoshipmentops.UpdateMTOShipmentParams, *models.MTOShipment) {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	originalShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
	})

	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	pickupAddress.StreetAddress1 = "123 Fake Test St NW"

	secondaryPickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	secondaryPickupAddress.StreetAddress1 = "89999 Other Test St NW"

	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"

	secondaryDeliveryAddress := testdatagen.MakeDefaultAddress(suite.DB())
	secondaryDeliveryAddress.StreetAddress1 = "9999 Test Fake Rd SE"

	mtoAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
	}}

	customerRemarks := ""

	req := httptest.NewRequest("PATCH", "/mto-shipments/"+originalShipment.ID.String(), nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

	payload := internalmessages.UpdateShipment{
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

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(originalShipment.ID),
		Body:          &payload,
		IfMatch:       eTag,
	}

	return params, &originalShipment
}

// getDefaultPPMShipmentAndParams generates a set of default params and a PPMShipment
func (suite *HandlerSuite) getDefaultPPMShipmentAndParams() (mtoshipmentops.UpdateMTOShipmentParams, *models.PPMShipment) {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	originalPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
	})
	originalShipment := originalPPMShipment.Shipment

	req := httptest.NewRequest("PATCH", "/mto-shipments/"+originalShipment.ID.String(), nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

	customerRemarks := "testing"
	payload := internalmessages.UpdateShipment{
		CustomerRemarks: &customerRemarks,
	}

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(originalShipment.ID),
		Body:          &payload,
		IfMatch:       eTag,
	}

	return params, &originalPPMShipment
}

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	moveRouter := moverouter.NewMoveRouter()
	moveWeights := moverouter.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator)
	ppmUpdater := ppmshipment.NewPPMShipmentUpdater()

	suite.Run("Successful PATCH - Integration Test", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, oldShipment := suite.getDefaultMTOShipmentAndParams()

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		suite.Equal(oldShipment.ID.String(), updatedShipment.ID.String())
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
		suite.Equal(oldShipment.ID.String(), string(updatedShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(updatedShipment.Agents[0].ID)
	})

	suite.Run("Successful PATCH with PPMShipment - Integration Test", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, existingPPMShipment := suite.getDefaultPPMShipmentAndParams()

		estimatedWeight := int64(6000)
		proGearWeight := int64(1000)
		spouseProGearWeight := int64(250)
		updatedPPM := &internalmessages.UpdatePPMShipment{
			EstimatedWeight:     &estimatedWeight,
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       &proGearWeight,
			SpouseProGearWeight: &spouseProGearWeight,
		}
		params.Body.PpmShipment = updatedPPM

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Check if existing fields are not updated
		suite.Equal(existingPPMShipment.Shipment.ID.String(), updatedShipment.ID.String())
		suite.EqualDate(existingPPMShipment.ExpectedDepartureDate, *updatedShipment.PpmShipment.ExpectedDepartureDate)
		suite.Equal(existingPPMShipment.PickupPostalCode, *updatedShipment.PpmShipment.PickupPostalCode)
		suite.Equal(existingPPMShipment.DestinationPostalCode, *updatedShipment.PpmShipment.DestinationPostalCode)
		suite.Equal(*existingPPMShipment.SitExpected, *updatedShipment.PpmShipment.SitExpected)

		// Check if mto_shipment fields are updated
		suite.Equal(*params.Body.CustomerRemarks, *updatedShipment.CustomerRemarks)

		// Check if ppm_shipment fields are updated
		suite.Equal(*params.Body.PpmShipment.EstimatedWeight, *updatedShipment.PpmShipment.EstimatedWeight)
		suite.Equal(*params.Body.PpmShipment.HasProGear, *updatedShipment.PpmShipment.HasProGear)
		suite.Equal(*params.Body.PpmShipment.ProGearWeight, *updatedShipment.PpmShipment.ProGearWeight)
		suite.Equal(*params.Body.PpmShipment.SpouseProGearWeight, *updatedShipment.PpmShipment.SpouseProGearWeight)
		suite.Equal(int64(1000000), *updatedShipment.PpmShipment.EstimatedIncentive)

		suite.NoError(updatedShipment.Validate(strfmt.Default))
	})

	suite.Run("Successful PATCH - Can update shipment status", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		expectedStatus := internalmessages.MTOShipmentStatusSUBMITTED

		params, _ := suite.getDefaultMTOShipmentAndParams()
		params.Body.Status = expectedStatus

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)

		suite.Equal(expectedStatus, updatedResponse.Payload.Status)
	})

	suite.Run("PATCH failure - 400 -- nil body", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, _ := suite.getDefaultMTOShipmentAndParams()
		params.Body = nil

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.Run("PATCH failure - 400 -- invalid requested status update", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, _ := suite.getDefaultMTOShipmentAndParams()
		params.Body.Status = internalmessages.MTOShipmentStatusREJECTED

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.Run("PATCH failure - 401- permission denied - not authenticated", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, oldShipment := suite.getDefaultMTOShipmentAndParams()
		updateURI := "/mto-shipments/" + oldShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		params.HTTPRequest = unauthorizedReq

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnauthorized{}, response)
	})

	suite.Run("PATCH failure - 403- permission denied - wrong application / user", func() {
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, oldShipment := suite.getDefaultMTOShipmentAndParams()
		updateURI := "/mto-shipments/" + oldShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		unauthorizedReq = suite.AuthenticateOfficeRequest(unauthorizedReq, officeUser)
		params.HTTPRequest = unauthorizedReq

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentForbidden{}, response)
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		uuidString := handlers.FmtUUID(uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763"))
		params, _ := suite.getDefaultMTOShipmentAndParams()
		params.MtoShipmentID = *uuidString

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			updater,
			ppmUpdater,
		}

		params, _ := suite.getDefaultMTOShipmentAndParams()
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	//TODO: when the address bug (MB-3691) gets fixed, this test should pass
	// Update: This test is not passing due swagger validation failing and no server-side validation
	// happening. These changes weren't covered in MB-3691, so we'll need to do addt'l
	// work to fix. Since we have refactoring slated for addresses, we can do then.
	// suite.Run("PATCH failure - 422 -- invalid input", func() {
	// 	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	// 	fetcher := fetch.NewFetcher(builder)
	// 	updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner)
	// 	handler := UpdateMTOShipmentHandler{
	// 		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	// 		updater,
	// 	}

	// 	oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	// 	oldShipment.RequestedPickupDate = nil

	// 	req := httptest.NewRequest("PATCH", "/mto-shipments/"+oldShipment.ID.String(), nil)
	// 	req = suite.AuthenticateRequest(req, serviceMember)
	// 	// invalid zip
	// 	payloadDestinationAddress := &internalmessages.Address{
	// 		City:           swag.String("Stumptown"),
	// 		Country:        swag.String("USA"),
	// 		ID:             "6e07a670-a072-4014-be9f-4926c1389f9a",
	// 		State:          swag.String("CA"),
	// 		StreetAddress1: swag.String("321 Main St."),
	// 		PostalCode:     swag.String("123"),
	// 	}

	// 	payload := internalmessages.UpdateShipment{
	// 		DestinationAddress: payloadDestinationAddress,
	// 	}
	// 	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	// 	params := mtoshipmentops.UpdateMTOShipmentParams{
	// 		HTTPRequest:   req,
	// 		MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
	// 		Body:          &payload,
	// 		IfMatch:       eTag,
	// 	}

	// 	response := handler.Handle(params)
	// 	suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)

	// })

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.MTOShipmentUpdater{}

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockUpdater,
			ppmUpdater,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipmentCustomer",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		params, _ := suite.getDefaultMTOShipmentAndParams()

		response := handler.Handle(params)

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
		queryBuilder := query.NewQueryBuilder()
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			listFetcher,
			fetcher,
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
				suite.EqualInt32Int64Pointers(expectedShipment.PPMShipment.EstimatedIncentive, returnedShipment.PpmShipment.EstimatedIncentive)
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
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockListFetcher,
			&mockFetcher,
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
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockListFetcher,
			&mockFetcher,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsUnauthorized{}, response)
	})

	suite.Run("Failure list fetch - 404 Not Found - Move Task Order ID", func() {
		subtestData := suite.makeListSubtestData()
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockListFetcher,
			&mockFetcher,
		}

		notfound := errors.New("Not found error")

		mockFetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(notfound)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsNotFound{}, response)
	})

	suite.Run("Failure list fetch - 500 Internal Server Error", func() {
		subtestData := suite.makeListSubtestData()
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockListFetcher,
			&mockFetcher,
		}

		internalServerErr := errors.New("ServerError")

		mockFetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockListFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(internalServerErr)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsInternalServerError{}, response)
	})
}

//
// DELETE
//

func (suite *HandlerSuite) TestDeleteShipmentHandler() {
	suite.Run("Returns 204 when all validations pass", func() {
		sm := testdatagen.MakeStubbedServiceMember(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMember: sm,
			},
		})
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
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerContext,
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
		sm := testdatagen.MakeStubbedServiceMember(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMember: sm,
			},
		})
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
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerContext,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteShipmentNotFound{}, response)
	})

	suite.Run("Returns 403 when deleter returns ForbiddenError", func() {
		sm := testdatagen.MakeStubbedServiceMember(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMember: sm,
			},
		})
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

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ForbiddenError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerContext,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteShipmentForbidden{}, response)
	})

	suite.Run("Returns 500 when deleter returns InternalServerError", func() {
		sm := testdatagen.MakeStubbedServiceMember(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMember: sm,
			},
		})
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
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerContext,
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
