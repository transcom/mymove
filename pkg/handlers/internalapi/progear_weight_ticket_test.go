package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	progearops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	internalmessages "github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CREATE TEST
func (suite *HandlerSuite) TestCreateProGearWeightTicketHandler() {
	// Reusable objects
	progearCreator := progear.NewCustomerProgearWeightTicketCreator()

	type progearCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      progearops.CreateProGearWeightTicketParams
		handler     CreateProGearWeightTicketHandler
	}
	makeCreateSubtestData := func(authenticateRequest bool) (subtestData progearCreateSubtestData) {
		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = progearops.CreateProGearWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateProGearWeightTicketHandler{
			suite.NewHandlerConfig(),
			progearCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.CreateProGearWeightTicketCreated{}, response)

		createdProgear := response.(*progearops.CreateProGearWeightTicketCreated).Payload

		suite.NotEmpty(createdProgear.ID.String())
		suite.NotNil(createdProgear.DocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure - 404 - not found - wrong service member", func() {
		subtestData := makeCreateSubtestData(false)

		unauthorizedUser := factory.BuildServiceMember(suite.DB(), nil, nil)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, unauthorizedUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&progearops.CreateProGearWeightTicketNotFound{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.ProgearWeightTicketCreator{}

		subtestData := makeCreateSubtestData(true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateProgearWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateProGearWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockCreator,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketInternalServerError{}, response)
	})
}

//
// UPDATE Customer test
//

func (suite *HandlerSuite) TestUpdateProGearWeightTicketHandler() {
	// Reusable objects
	progearUpdater := progear.NewCustomerProgearWeightTicketUpdater()

	type progearUpdateSubtestData struct {
		ppmShipment models.PPMShipment
		progear     models.ProgearWeightTicket
		params      progearops.UpdateProGearWeightTicketParams
		handler     UpdateProGearWeightTicketHandler
	}
	makeUpdateSubtestData := func(authenticateRequest bool) (subtestData progearUpdateSubtestData) {
		// Use fake data:
		subtestData.progear = factory.BuildProgearWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.progear.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.progear.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		eTag := etag.GenerateEtag(subtestData.progear.UpdatedAt)
		subtestData.params = progearops.UpdateProGearWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			ProGearWeightTicketID: *handlers.FmtUUID(subtestData.progear.ID),
			IfMatch:               eTag,
		}

		subtestData.handler = UpdateProGearWeightTicketHandler{
			suite.createS3HandlerConfig(),
			progearUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		subtestData := makeUpdateSubtestData(true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.progear.Document,
				LinkOnly: true,
			},
		}, nil)

		progearDes := "Pro gear desctription"
		hasWeightTickets := true
		belongsToSelf := true
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{
			Description:      progearDes,
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(4000),
			BelongsToSelf:    belongsToSelf,
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketOK{}, response)

		updatedProgear := response.(*progearops.UpdateProGearWeightTicketOK).Payload
		suite.Equal(subtestData.progear.ID.String(), updatedProgear.ID.String())
		suite.Equal(params.UpdateProGearWeightTicket.Description, *updatedProgear.Description)
	})

	suite.Run("PATCH failure -400 - nil body", func() {
		subtestData := makeUpdateSubtestData(true)
		subtestData.params.UpdateProGearWeightTicket = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.UpdateProGearWeightTicketBadRequest{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		progearDes := "Pro gear desctription"
		hasWeightTickets := true
		belongsToSelf := true
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{
			Description:      progearDes,
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(0),
			BelongsToSelf:    belongsToSelf,
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{}
		// This test should fail because of the wrong ID
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("f20d9c9b-2de5-4860-ad31-fd5c10e739f6"))
		params.ProGearWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ProgearWeightTicketUpdater{}
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		progearDes := "Pro gear desctription"
		hasWeightTickets := true
		belongsToSelf := true
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{
			Description:      progearDes,
			Weight:           handlers.FmtInt64(1),
			HasWeightTickets: hasWeightTickets,
			BelongsToSelf:    belongsToSelf,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateProgearWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.ProgearWeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateProGearWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketInternalServerError{}, response)
		errResponse := response.(*progearops.UpdateProGearWeightTicketInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})
}

// DELETE test
func (suite *HandlerSuite) TestDeleteProgearWeightTicketHandler() {
	// Create Reusable objects
	progearWeightTicketDeleter := progear.NewProgearWeightTicketDeleter()

	type progearWeightTicketDeleteSubtestData struct {
		ppmShipment         models.PPMShipment
		progearWeightTicket models.ProgearWeightTicket
		params              progearops.DeleteProGearWeightTicketParams
		handler             DeleteProGearWeightTicketHandler
	}
	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData progearWeightTicketDeleteSubtestData) {
		// Fake data:
		subtestData.progearWeightTicket = factory.BuildProgearWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.progearWeightTicket.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.progearWeightTicket.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = progearops.
			DeleteProGearWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			ProGearWeightTicketID: *handlers.FmtUUID(subtestData.progearWeightTicket.ID),
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = DeleteProGearWeightTicketHandler{
			suite.createS3HandlerConfig(),
			progearWeightTicketDeleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete Pro-gear Weight Ticket - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.DeleteProGearWeightTicketNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.DeleteProGearWeightTicketUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong application / user", func() {
		subtestData := makeDeleteSubtestData(false)

		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&progearops.DeleteProGearWeightTicketForbidden{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong service member user", func() {
		subtestData := makeDeleteSubtestData(false)

		otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, otherServiceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&progearops.DeleteProGearWeightTicketForbidden{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.ProGearWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.DeleteProGearWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and moving expense ID don't match", func() {
		subtestData := makeDeleteSubtestData(false)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
				LinkOnly: true,
			},
		}, nil)

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, serviceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)
		suite.IsType(&progearops.DeleteProGearWeightTicketNotFound{}, response)
	})
	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.ProgearWeightTicketDeleter{}

		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteProgearWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(err)

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		handler := DeleteProGearWeightTicketHandler{
			suite.createS3HandlerConfig(),
			&mockDeleter,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.DeleteProGearWeightTicketInternalServerError{}, response)
	})
}
