package ghcapi

// import (
// 	"errors"
// 	"fmt"
// 	"net/http/httptest"

// 	"github.com/go-openapi/strfmt"
// 	"github.com/stretchr/testify/mock"

// 	"github.com/transcom/mymove/pkg/appcontext"
// 	"github.com/transcom/mymove/pkg/etag"
// 	"github.com/transcom/mymove/pkg/factory"
// 	gunsafeops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
// 	"github.com/transcom/mymove/pkg/gen/ghcmessages"
// 	"github.com/transcom/mymove/pkg/handlers"
// 	"github.com/transcom/mymove/pkg/models"
// 	gunsafe "github.com/transcom/mymove/pkg/services/gunsafe_weight_ticket"
// 	"github.com/transcom/mymove/pkg/services/mocks"
// 	"github.com/transcom/mymove/pkg/testdatagen"
// )

// // CREATE TEST
// func (suite *HandlerSuite) TestCreateGunsafeWeightTicketHandler() {
// 	// Reusable objects
// 	gunsafeCreator := gunsafe.NewOfficeGunsafeWeightTicketCreator()

// 	type gunsafeCreateSubtestData struct {
// 		ppmShipment models.PPMShipment
// 		params      gunsafeops.CreateGunsafeWeightTicketParams
// 		handler     CreateGunsafeWeightTicketHandler
// 	}
// 	makeCreateSubtestData := func(authenticateRequest bool) (subtestData gunsafeCreateSubtestData) {
// 		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
// 		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets", subtestData.ppmShipment.ID.String())
// 		req := httptest.NewRequest("POST", endpoint, nil)
// 		officeUser := factory.BuildOfficeUser(nil, nil, nil)
// 		if authenticateRequest {
// 			req = suite.AuthenticateOfficeRequest(req, officeUser)
// 		}
// 		subtestData.params = gunsafeops.CreateGunsafeWeightTicketParams{
// 			HTTPRequest:   req,
// 			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
// 		}

// 		subtestData.handler = CreateGunsafeWeightTicketHandler{
// 			suite.HandlerConfig(),
// 			gunsafeCreator,
// 		}

// 		return subtestData
// 	}

// 	suite.Run("Successfully Create Gunsafe Weight Ticket - Integration Test", func() {
// 		subtestData := makeCreateSubtestData(true)

// 		response := subtestData.handler.Handle(subtestData.params)

// 		suite.IsType(&gunsafeops.CreateGunsafeWeightTicketCreated{}, response)

// 		createdGunsafe := response.(*gunsafeops.CreateGunsafeWeightTicketCreated).Payload

// 		suite.NotEmpty(createdGunsafe.ID.String())
// 		suite.NotNil(createdGunsafe.DocumentID.String())
// 	})

// 	suite.Run("DELETE failure - 404- Create not found", func() {
// 		subtestData := makeCreateSubtestData(true)
// 		params := subtestData.params

// 		// Wrong ID provided
// 		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
// 		params.PpmShipmentID = *uuidString

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.CreateGunsafeWeightTicketNotFound{}, response)
// 	})

// 	suite.Run("POST failure - 400- bad request", func() {
// 		subtestData := makeCreateSubtestData(true)
// 		// Missing PPM Shipment ID
// 		params := subtestData.params

// 		params.PpmShipmentID = ""

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.CreateGunsafeWeightTicketBadRequest{}, response)
// 	})

// 	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
// 		subtestData := makeCreateSubtestData(false)

// 		response := subtestData.handler.Handle(subtestData.params)

// 		suite.IsType(&gunsafeops.CreateGunsafeWeightTicketUnauthorized{}, response)
// 	})

// 	suite.Run("Post failure - 500 - Server Error", func() {
// 		mockCreator := mocks.GunsafeWeightTicketCreator{}

// 		subtestData := makeCreateSubtestData(true)
// 		params := subtestData.params
// 		serverErr := errors.New("ServerError")

// 		mockCreator.On("CreateGunsafeWeightTicket",
// 			mock.AnythingOfType("*appcontext.appContext"),
// 			mock.AnythingOfType("uuid.UUID"),
// 		).Return(nil, serverErr)

// 		handler := CreateGunsafeWeightTicketHandler{
// 			suite.HandlerConfig(),
// 			&mockCreator,
// 		}

// 		response := handler.Handle(params)

// 		suite.IsType(&gunsafeops.CreateGunsafeWeightTicketInternalServerError{}, response)
// 	})
// }

// // UPDATE Customer test
// func (suite *HandlerSuite) TestUpdateGunsafeWeightTicketHandler() {
// 	// Reusable objects
// 	gunsafeUpdater := gunsafe.NewCustomerGunsafeWeightTicketUpdater()

// 	type gunsafeUpdateSubtestData struct {
// 		ppmShipment models.PPMShipment
// 		gunsafe     models.GunsafeWeightTicket
// 		params      gunsafeops.UpdateGunsafeWeightTicketParams
// 		handler     UpdateGunsafeWeightTicketHandler
// 	}
// 	makeUpdateSubtestData := func(appCtx appcontext.AppContext, _ bool) (subtestData gunsafeUpdateSubtestData) {
// 		db := appCtx.DB()

// 		// Use fake data:
// 		subtestData.gunsafe = factory.BuildGunsafeWeightTicket(db, nil, nil)
// 		subtestData.ppmShipment = subtestData.gunsafe.PPMShipment

// 		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.gunsafe.ID.String())
// 		officeUser := factory.BuildOfficeUser(nil, nil, nil)
// 		req := httptest.NewRequest("PATCH", endpoint, nil)
// 		req = suite.AuthenticateOfficeRequest(req, officeUser)
// 		eTag := etag.GenerateEtag(subtestData.gunsafe.UpdatedAt)

// 		subtestData.params = gunsafeops.UpdateGunsafeWeightTicketParams{
// 			HTTPRequest:           req,
// 			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
// 			GunsafeWeightTicketID: *handlers.FmtUUID(subtestData.gunsafe.ID),
// 			IfMatch:               eTag,
// 		}

// 		subtestData.handler = UpdateGunsafeWeightTicketHandler{
// 			suite.createS3HandlerConfig(),
// 			gunsafeUpdater,
// 		}

// 		return subtestData
// 	}

// 	suite.Run("Successfully Update Gunsafe Weight Ticket - Integration Test", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)

// 		params := subtestData.params

// 		factory.BuildUserUpload(suite.DB(), []factory.Customization{
// 			{
// 				Model:    subtestData.gunsafe.Document,
// 				LinkOnly: true,
// 			},
// 		}, nil)

// 		hasWeightTickets := true
// 		belongsToSelf := true
// 		params.UpdateGunsafeWeightTicket = &ghcmessages.UpdateGunsafeWeightTicket{
// 			HasWeightTickets: hasWeightTickets,
// 			Weight:           handlers.FmtInt64(4000),
// 			BelongsToSelf:    belongsToSelf,
// 		}

// 		// Validate incoming payload: no body to validate
// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.UpdateGunsafeWeightTicketOK{}, response)

// 		updatedGunsafe := response.(*gunsafeops.UpdateGunsafeWeightTicketOK).Payload

// 		// Validate outgoing payload
// 		suite.NoError(updatedGunsafe.Validate(strfmt.Default))

// 		suite.Equal(subtestData.gunsafe.ID.String(), updatedGunsafe.ID.String())
// 		suite.Equal(params.UpdateGunsafeWeightTicket.Weight, updatedGunsafe.Weight)
// 	})

// 	suite.Run("PATCH failure - 404- not found", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		params.UpdateGunsafeWeightTicket = &ghcmessages.UpdateGunsafeWeightTicket{}
// 		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("3ce0e367-a337-46e3-b4cf-f79aebc4f6c8"))
// 		params.GunsafeWeightTicketID = *wrongUUIDString

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.UpdateGunsafeWeightTicketNotFound{}, response)
// 	})

// 	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		params.UpdateGunsafeWeightTicket = &ghcmessages.UpdateGunsafeWeightTicket{}
// 		params.IfMatch = "wrong-if-match-header-value"

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.UpdateGunsafeWeightTicketPreconditionFailed{}, response)
// 	})

// 	suite.Run("PATCH failure - 500", func() {
// 		mockUpdater := mocks.GunsafeWeightTicketUpdater{}
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		ownsTrailer := true
// 		hasWeightTickets := true

// 		params.UpdateGunsafeWeightTicket = &ghcmessages.UpdateGunsafeWeightTicket{
// 			Weight:           handlers.FmtInt64(1000),
// 			BelongsToSelf:    ownsTrailer,
// 			HasWeightTickets: hasWeightTickets,
// 		}

// 		err := errors.New("ServerError")

// 		mockUpdater.On("UpdateGunsafeWeightTicket",
// 			mock.AnythingOfType("*appcontext.appContext"),
// 			mock.AnythingOfType("models.GunsafeWeightTicket"),
// 			mock.AnythingOfType("string"),
// 		).Return(nil, err)

// 		handler := UpdateGunsafeWeightTicketHandler{
// 			suite.HandlerConfig(),
// 			&mockUpdater,
// 		}

// 		response := handler.Handle(params)

// 		suite.IsType(&gunsafeops.UpdateGunsafeWeightTicketInternalServerError{}, response)
// 	})
// }

// // DELETE test
// func (suite *HandlerSuite) TestDeleteGunsafeWeightTicketHandler() {
// 	// Create Reusable objects
// 	gunsafeWeightTicketDeleter := gunsafe.NewGunsafeWeightTicketDeleter()

// 	type gunsafeWeightTicketDeleteSubtestData struct {
// 		ppmShipment         models.PPMShipment
// 		gunsafeWeightTicket models.GunsafeWeightTicket
// 		params              gunsafeops.DeleteGunsafeWeightTicketParams
// 		handler             DeleteGunsafeWeightTicketHandler
// 	}
// 	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData gunsafeWeightTicketDeleteSubtestData) {
// 		// Fake data:
// 		subtestData.gunsafeWeightTicket = factory.BuildGunsafeWeightTicket(suite.DB(), nil, nil)
// 		subtestData.ppmShipment = subtestData.gunsafeWeightTicket.PPMShipment
// 		officeUser := factory.BuildOfficeUser(nil, nil, nil)

// 		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.gunsafeWeightTicket.ID.String())
// 		req := httptest.NewRequest("DELETE", endpoint, nil)
// 		if authenticateRequest {
// 			req = suite.AuthenticateOfficeRequest(req, officeUser)
// 		}
// 		subtestData.params = gunsafeops.
// 			DeleteGunsafeWeightTicketParams{
// 			HTTPRequest:           req,
// 			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
// 			GunsafeWeightTicketID: *handlers.FmtUUID(subtestData.gunsafeWeightTicket.ID),
// 		}

// 		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
// 		subtestData.handler = DeleteGunsafeWeightTicketHandler{
// 			suite.createS3HandlerConfig(),
// 			gunsafeWeightTicketDeleter,
// 		}

// 		return subtestData
// 	}

// 	suite.Run("Successfully Delete Pro-gear Weight Ticket - Integration Test", func() {
// 		subtestData := makeDeleteSubtestData(true)

// 		params := subtestData.params
// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.DeleteGunsafeWeightTicketNoContent{}, response)
// 	})

// 	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
// 		subtestData := makeDeleteSubtestData(false)
// 		response := subtestData.handler.Handle(subtestData.params)

// 		suite.IsType(&gunsafeops.DeleteGunsafeWeightTicketUnauthorized{}, response)
// 	})

// 	suite.Run("DELETE failure - 404- not found", func() {
// 		subtestData := makeDeleteSubtestData(true)
// 		params := subtestData.params
// 		// Wrong ID provided
// 		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
// 		params.GunsafeWeightTicketID = *uuidString

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&gunsafeops.DeleteGunsafeWeightTicketNotFound{}, response)
// 	})

// 	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and gunsafe ID don't match", func() {
// 		subtestData := makeDeleteSubtestData(false)
// 		officeUser := factory.BuildOfficeUser(nil, nil, nil)

// 		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
// 			{
// 				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
// 				LinkOnly: true,
// 			},
// 		}, nil)

// 		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
// 		req := subtestData.params.HTTPRequest
// 		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
// 		unauthorizedParams := subtestData.params
// 		unauthorizedParams.HTTPRequest = unauthorizedReq

// 		response := subtestData.handler.Handle(unauthorizedParams)
// 		suite.IsType(&gunsafeops.DeleteGunsafeWeightTicketNotFound{}, response)
// 	})
// 	suite.Run("DELETE failure - 500 - server error", func() {
// 		mockDeleter := mocks.GunsafeWeightTicketDeleter{}

// 		subtestData := makeDeleteSubtestData(true)
// 		params := subtestData.params

// 		err := errors.New("ServerError")

// 		mockDeleter.On("DeleteGunsafeWeightTicket",
// 			mock.AnythingOfType("*appcontext.appContext"),
// 			mock.AnythingOfType("uuid.UUID"),
// 			mock.AnythingOfType("uuid.UUID"),
// 		).Return(err)

// 		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
// 		handler := DeleteGunsafeWeightTicketHandler{
// 			suite.createS3HandlerConfig(),
// 			&mockDeleter,
// 		}

// 		response := handler.Handle(params)

// 		suite.IsType(&gunsafeops.DeleteGunsafeWeightTicketInternalServerError{}, response)
// 	})
// }
