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
	gunsafeops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	gunsafe "github.com/transcom/mymove/pkg/services/gunsafe_weight_ticket"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CREATE TEST
func (suite *HandlerSuite) TestCreateGunSafeWeightTicketHandler() {
	// Reusable objects
	gunsafeCreator := gunsafe.NewOfficeGunSafeWeightTicketCreator()

	type gunsafeCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      gunsafeops.CreateGunSafeWeightTicketParams
		handler     CreateGunSafeWeightTicketHandler
	}
	makeCreateSubtestData := func(authenticateRequest bool) (subtestData gunsafeCreateSubtestData) {
		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/gun-safe-weight-tickets", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}
		subtestData.params = gunsafeops.CreateGunSafeWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateGunSafeWeightTicketHandler{
			suite.HandlerConfig(),
			gunsafeCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create GunSafe Weight Ticket - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&gunsafeops.CreateGunSafeWeightTicketCreated{}, response)

		createdGunSafe := response.(*gunsafeops.CreateGunSafeWeightTicketCreated).Payload

		suite.NotEmpty(createdGunSafe.ID.String())
		suite.NotNil(createdGunSafe.DocumentID.String())
	})

	suite.Run("Fails to Create GunSafe Weight Ticket When FF is Toggled OFF - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		// Overwrite handler config in order to return false for FF
		handlerConfig := suite.HandlerConfig()
		gunSafeFF := services.FeatureFlag{
			Key:   "gun_safe",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(gunSafeFF, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

		handler := CreateGunSafeWeightTicketHandler{
			handlerConfig,
			gunsafeCreator,
		}

		response := handler.Handle(subtestData.params)
		suite.IsType(&gunsafeops.CreateGunSafeWeightTicketForbidden{}, response)
	})

	suite.Run("POST failure - 404- Create not found", func() {
		subtestData := makeCreateSubtestData(true)
		params := subtestData.params

		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.PpmShipmentID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.CreateGunSafeWeightTicketNotFound{}, response)
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.CreateGunSafeWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		subtestData := makeCreateSubtestData(false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&gunsafeops.CreateGunSafeWeightTicketUnauthorized{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.GunSafeWeightTicketCreator{}

		subtestData := makeCreateSubtestData(true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateGunSafeWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateGunSafeWeightTicketHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}

		response := handler.Handle(params)

		suite.IsType(&gunsafeops.CreateGunSafeWeightTicketInternalServerError{}, response)
	})
}

// UPDATE Office test
func (suite *HandlerSuite) TestUpdateGunSafeWeightTicketHandler() {
	// Reusable objects
	gunsafeUpdater := gunsafe.NewOfficeGunSafeWeightTicketUpdater()

	type gunsafeUpdateSubtestData struct {
		ppmShipment models.PPMShipment
		gunsafe     models.GunSafeWeightTicket
		params      gunsafeops.UpdateGunSafeWeightTicketParams
		handler     UpdateGunSafeWeightTicketHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, _ bool) (subtestData gunsafeUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.gunsafe = factory.BuildGunSafeWeightTicket(db, nil, nil)
		subtestData.ppmShipment = subtestData.gunsafe.PPMShipment

		endpoint := fmt.Sprintf("/ppm-shipments/%s/gun-safe-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.gunsafe.ID.String())
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("PATCH", endpoint, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		eTag := etag.GenerateEtag(subtestData.gunsafe.UpdatedAt)

		subtestData.params = gunsafeops.UpdateGunSafeWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			GunSafeWeightTicketID: *handlers.FmtUUID(subtestData.gunsafe.ID),
			IfMatch:               eTag,
		}

		subtestData.handler = UpdateGunSafeWeightTicketHandler{
			suite.createS3HandlerConfig(),
			gunsafeUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update GunSafe Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.gunsafe.Document,
				LinkOnly: true,
			},
		}, nil)

		hasWeightTickets := true
		params.UpdateGunSafeWeightTicket = &ghcmessages.UpdateGunSafeWeightTicket{
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(500),
		}

		// Validate incoming payload: no body to validate
		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.UpdateGunSafeWeightTicketOK{}, response)

		updatedGunSafe := response.(*gunsafeops.UpdateGunSafeWeightTicketOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedGunSafe.Validate(strfmt.Default))

		suite.Equal(subtestData.gunsafe.ID.String(), updatedGunSafe.ID.String())
		suite.Equal(params.UpdateGunSafeWeightTicket.Weight, updatedGunSafe.Weight)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateGunSafeWeightTicket = &ghcmessages.UpdateGunSafeWeightTicket{}
		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("3ce0e367-a337-46e3-b4cf-f79aebc4f6c8"))
		params.GunSafeWeightTicketID = *wrongUUIDString

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.UpdateGunSafeWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateGunSafeWeightTicket = &ghcmessages.UpdateGunSafeWeightTicket{}
		params.IfMatch = "wrong-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.UpdateGunSafeWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.GunSafeWeightTicketUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		hasWeightTickets := true

		params.UpdateGunSafeWeightTicket = &ghcmessages.UpdateGunSafeWeightTicket{
			Weight:           handlers.FmtInt64(1000),
			HasWeightTickets: hasWeightTickets,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateGunSafeWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.GunSafeWeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateGunSafeWeightTicketHandler{
			suite.HandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&gunsafeops.UpdateGunSafeWeightTicketInternalServerError{}, response)
	})

	suite.Run("UPDATE failure - Fails when FF is toggled OFF", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.gunsafe.Document,
				LinkOnly: true,
			},
		}, nil)

		hasWeightTickets := true
		params.UpdateGunSafeWeightTicket = &ghcmessages.UpdateGunSafeWeightTicket{
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(500),
		}

		// Overwrite handler config in order to return false for FF
		handlerConfig := suite.HandlerConfig()
		gunSafeFF := services.FeatureFlag{
			Key:   "gun_safe",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(gunSafeFF, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

		handler := UpdateGunSafeWeightTicketHandler{
			handlerConfig,
			gunsafeUpdater,
		}

		// Validate incoming payload: no body to validate
		response := handler.Handle(params)

		suite.IsType(&gunsafeops.UpdateGunSafeWeightTicketForbidden{}, response)
	})
}

// DELETE test
func (suite *HandlerSuite) TestDeleteGunSafeWeightTicketHandler() {
	// Create Reusable objects
	gunsafeWeightTicketDeleter := gunsafe.NewGunSafeWeightTicketDeleter()

	type gunsafeWeightTicketDeleteSubtestData struct {
		ppmShipment         models.PPMShipment
		gunsafeWeightTicket models.GunSafeWeightTicket
		params              gunsafeops.DeleteGunSafeWeightTicketParams
		handler             DeleteGunSafeWeightTicketHandler
	}
	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData gunsafeWeightTicketDeleteSubtestData) {
		// Fake data:
		subtestData.gunsafeWeightTicket = factory.BuildGunSafeWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.gunsafeWeightTicket.PPMShipment
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		endpoint := fmt.Sprintf("/ppm-shipments/%s/gun-safe-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.gunsafeWeightTicket.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}
		subtestData.params = gunsafeops.
			DeleteGunSafeWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			GunSafeWeightTicketID: *handlers.FmtUUID(subtestData.gunsafeWeightTicket.ID),
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = DeleteGunSafeWeightTicketHandler{
			suite.createS3HandlerConfig(),
			gunsafeWeightTicketDeleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete Gun Safe Weight Ticket - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.DeleteGunSafeWeightTicketNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&gunsafeops.DeleteGunSafeWeightTicketUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.GunSafeWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunsafeops.DeleteGunSafeWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and gunsafe ID don't match", func() {
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
		suite.IsType(&gunsafeops.DeleteGunSafeWeightTicketNotFound{}, response)
	})
	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.GunSafeWeightTicketDeleter{}

		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteGunSafeWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(err)

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		handler := DeleteGunSafeWeightTicketHandler{
			suite.createS3HandlerConfig(),
			&mockDeleter,
		}

		response := handler.Handle(params)

		suite.IsType(&gunsafeops.DeleteGunSafeWeightTicketInternalServerError{}, response)
	})

	suite.Run("DELETE failure - Fails to Create GunSafe Weight Ticket When FF is Toggled OFF - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params

		// Overwrite handler config in order to return false for FF
		handlerConfig := suite.HandlerConfig()
		gunSafeFF := services.FeatureFlag{
			Key:   "gun_safe",
			Match: false,
		}

		mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(gunSafeFF, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

		subtestData.handler = DeleteGunSafeWeightTicketHandler{
			handlerConfig,
			gunsafeWeightTicketDeleter,
		}

		response := subtestData.handler.Handle(params)
		suite.IsType(&gunsafeops.DeleteGunSafeWeightTicketForbidden{}, response)
	})
}
