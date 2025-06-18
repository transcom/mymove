package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	gunSafeops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	internalmessages "github.com/transcom/mymove/pkg/gen/internalmessages"
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
	gunSafeCreator := gunsafe.NewCustomerGunSafeWeightTicketCreator()

	type gunSafeCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      gunSafeops.CreateGunSafeWeightTicketParams
		handler     CreateGunSafeWeightTicketHandler
	}
	makeCreateSubtestData := func(authenticateRequest bool) (subtestData gunSafeCreateSubtestData) {
		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/gun-safe-weight-tickets", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = gunSafeops.CreateGunSafeWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateGunSafeWeightTicketHandler{
			suite.HandlerConfig(),
			gunSafeCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&gunSafeops.CreateGunSafeWeightTicketCreated{}, response)

		createdgunSafe := response.(*gunSafeops.CreateGunSafeWeightTicketCreated).Payload

		suite.NotEmpty(createdgunSafe.ID.String())
		suite.NotNil(createdgunSafe.DocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.CreateGunSafeWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure - 404 - not found - wrong service member", func() {
		subtestData := makeCreateSubtestData(false)

		unauthorizedUser := factory.BuildServiceMember(suite.DB(), nil, nil)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, unauthorizedUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&gunSafeops.CreateGunSafeWeightTicketNotFound{}, response)
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

		suite.IsType(&gunSafeops.CreateGunSafeWeightTicketInternalServerError{}, response)
	})

	suite.Run("Fails to Create GunSafe Weight Ticket When FF is Toggled OFF - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		gunSafeCreator := gunsafe.NewCustomerGunSafeWeightTicketCreator()

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
			gunSafeCreator,
		}

		response := handler.Handle(subtestData.params)
		suite.IsType(&gunSafeops.CreateGunSafeWeightTicketForbidden{}, response)
	})
}

//
// UPDATE test
//

func (suite *HandlerSuite) TestUpdateGunSafeWeightTicketHandler() {
	// Reusable objects
	gunSafeUpdater := gunsafe.NewCustomerGunSafeWeightTicketUpdater()

	type gunSafeUpdateSubtestData struct {
		ppmShipment models.PPMShipment
		gunSafe     models.GunSafeWeightTicket
		params      gunSafeops.UpdateGunSafeWeightTicketParams
		handler     UpdateGunSafeWeightTicketHandler
	}
	makeUpdateSubtestData := func(authenticateRequest bool) (subtestData gunSafeUpdateSubtestData) {
		// Use fake data:
		subtestData.gunSafe = factory.BuildGunSafeWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.gunSafe.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/gun-safe-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.gunSafe.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		eTag := etag.GenerateEtag(subtestData.gunSafe.UpdatedAt)
		subtestData.params = gunSafeops.UpdateGunSafeWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			GunSafeWeightTicketID: *handlers.FmtUUID(subtestData.gunSafe.ID),
			IfMatch:               eTag,
		}

		subtestData.handler = UpdateGunSafeWeightTicketHandler{
			suite.createS3HandlerConfig(),
			gunSafeUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		subtestData := makeUpdateSubtestData(true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.gunSafe.Document,
				LinkOnly: true,
			},
		}, nil)

		gunSafeDes := "Gun safe desctription"
		hasWeightTickets := true
		params.UpdateGunSafeWeightTicket = &internalmessages.UpdateGunSafeWeightTicket{
			Description:      gunSafeDes,
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(4000),
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketOK{}, response)

		updatedgunSafe := response.(*gunSafeops.UpdateGunSafeWeightTicketOK).Payload
		suite.Equal(subtestData.gunSafe.ID.String(), updatedgunSafe.ID.String())
		suite.Equal(params.UpdateGunSafeWeightTicket.Description, *updatedgunSafe.Description)
	})

	suite.Run("PATCH failure -400 - nil body", func() {
		subtestData := makeUpdateSubtestData(true)
		subtestData.params.UpdateGunSafeWeightTicket = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketBadRequest{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		gunSafeDes := "Gun safe desctription"
		hasWeightTickets := true
		params.UpdateGunSafeWeightTicket = &internalmessages.UpdateGunSafeWeightTicket{
			Description:      gunSafeDes,
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(0),
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateGunSafeWeightTicket = &internalmessages.UpdateGunSafeWeightTicket{}
		// This test should fail because of the wrong ID
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("f20d9c9b-2de5-4860-ad31-fd5c10e739f6"))
		params.GunSafeWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateGunSafeWeightTicket = &internalmessages.UpdateGunSafeWeightTicket{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.GunSafeWeightTicketUpdater{}
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		gunSafeDes := "Gun safe desctription"
		hasWeightTickets := true
		params.UpdateGunSafeWeightTicket = &internalmessages.UpdateGunSafeWeightTicket{
			Description:      gunSafeDes,
			Weight:           handlers.FmtInt64(1),
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

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketInternalServerError{}, response)
		errResponse := response.(*gunSafeops.UpdateGunSafeWeightTicketInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.Run("UPDATE failure - Fails when FF is toggled OFF", func() {
		subtestData := makeUpdateSubtestData(true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.gunSafe.Document,
				LinkOnly: true,
			},
		}, nil)

		gunSafeDes := "Gun safe desctription"
		hasWeightTickets := true
		params.UpdateGunSafeWeightTicket = &internalmessages.UpdateGunSafeWeightTicket{
			Description:      gunSafeDes,
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(4000),
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
			gunSafeUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&gunSafeops.UpdateGunSafeWeightTicketForbidden{}, response)
	})
}

// DELETE test
func (suite *HandlerSuite) TestDeleteGunSafeWeightTicketHandler() {
	// Create Reusable objects
	gunSafeWeightTicketDeleter := gunsafe.NewGunSafeWeightTicketDeleter()

	type gunSafeWeightTicketDeleteSubtestData struct {
		ppmShipment         models.PPMShipment
		gunSafeWeightTicket models.GunSafeWeightTicket
		params              gunSafeops.DeleteGunSafeWeightTicketParams
		handler             DeleteGunSafeWeightTicketHandler
	}
	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData gunSafeWeightTicketDeleteSubtestData) {
		// Fake data:
		subtestData.gunSafeWeightTicket = factory.BuildGunSafeWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.gunSafeWeightTicket.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/gun-safe-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.gunSafeWeightTicket.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = gunSafeops.
			DeleteGunSafeWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			GunSafeWeightTicketID: *handlers.FmtUUID(subtestData.gunSafeWeightTicket.ID),
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = DeleteGunSafeWeightTicketHandler{
			suite.createS3HandlerConfig(),
			gunSafeWeightTicketDeleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete gun-safe Weight Ticket - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong application / user", func() {
		subtestData := makeDeleteSubtestData(false)

		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketForbidden{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong service member user", func() {
		subtestData := makeDeleteSubtestData(false)

		otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, otherServiceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketForbidden{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.GunSafeWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketNotFound{}, response)
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
		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketNotFound{}, response)
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

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketInternalServerError{}, response)
	})

	suite.Run("DELETE failure - Fails to Create GunSafe Weight Ticket When FF is Toggled OFF - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true)

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

		gunSafeWeightTicketDeleter := gunsafe.NewGunSafeWeightTicketDeleter()
		handler := DeleteGunSafeWeightTicketHandler{
			handlerConfig,
			gunSafeWeightTicketDeleter,
		}

		params := subtestData.params
		response := handler.Handle(params)

		suite.IsType(&gunSafeops.DeleteGunSafeWeightTicketForbidden{}, response)
	})
}
