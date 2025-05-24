package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	progearops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CREATE TEST
func (suite *HandlerSuite) TestCreateProGearWeightTicketHandler() {
	// Reusable objects
	progearCreator := progear.NewOfficeProgearWeightTicketCreator()

	type progearCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      progearops.CreateProGearWeightTicketParams
		handler     CreateProGearWeightTicketHandler
	}
	makeCreateSubtestData := func(authenticateRequest bool) (subtestData progearCreateSubtestData) {
		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
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

	suite.Run("Successfully Create ProGear Weight Ticket - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.CreateProGearWeightTicketCreated{}, response)

		createdProgear := response.(*progearops.CreateProGearWeightTicketCreated).Payload

		suite.NotEmpty(createdProgear.ID.String())
		suite.NotNil(createdProgear.DocumentID.String())
	})

	suite.Run("DELETE failure - 404- Create not found", func() {
		subtestData := makeCreateSubtestData(true)
		params := subtestData.params

		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.PpmShipmentID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketNotFound{}, response)
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		subtestData := makeCreateSubtestData(false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.CreateProGearWeightTicketUnauthorized{}, response)
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

// UPDATE Customer test
func (suite *HandlerSuite) TestUpdateProGearWeightTicketHandler() {
	// Reusable objects
	progearUpdater := progear.NewCustomerProgearWeightTicketUpdater()

	type progearUpdateSubtestData struct {
		ppmShipment models.PPMShipment
		progear     models.ProgearWeightTicket
		params      progearops.UpdateProGearWeightTicketParams
		handler     UpdateProgearWeightTicketHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, _ bool) (subtestData progearUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.progear = factory.BuildProgearWeightTicket(db, nil, nil)
		subtestData.ppmShipment = subtestData.progear.PPMShipment

		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.progear.ID.String())
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("PATCH", endpoint, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		eTag := etag.GenerateEtag(subtestData.progear.UpdatedAt)

		subtestData.params = progearops.UpdateProGearWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			ProGearWeightTicketID: *handlers.FmtUUID(subtestData.progear.ID),
			IfMatch:               eTag,
		}

		subtestData.handler = UpdateProgearWeightTicketHandler{
			suite.createS3HandlerConfig(),
			progearUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update ProGear Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.progear.Document,
				LinkOnly: true,
			},
		}, nil)

		hasWeightTickets := true
		belongsToSelf := true
		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(4000),
			BelongsToSelf:    belongsToSelf,
		}

		// Validate incoming payload: no body to validate
		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketOK{}, response)

		updatedProgear := response.(*progearops.UpdateProGearWeightTicketOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedProgear.Validate(strfmt.Default))

		suite.Equal(subtestData.progear.ID.String(), updatedProgear.ID.String())
		suite.Equal(params.UpdateProGearWeightTicket.Weight, updatedProgear.Weight)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{}
		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("3ce0e367-a337-46e3-b4cf-f79aebc4f6c8"))
		params.ProGearWeightTicketID = *wrongUUIDString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{}
		params.IfMatch = "wrong-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ProgearWeightTicketUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		ownsTrailer := true
		hasWeightTickets := true

		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{
			Weight:           handlers.FmtInt64(1000),
			BelongsToSelf:    ownsTrailer,
			HasWeightTickets: hasWeightTickets,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateProgearWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.ProgearWeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateProgearWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketInternalServerError{}, response)
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
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.progearWeightTicket.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
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

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.ProGearWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.DeleteProGearWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and proGear ID don't match", func() {
		subtestData := makeDeleteSubtestData(false)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
				LinkOnly: true,
			},
		}, nil)

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
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
