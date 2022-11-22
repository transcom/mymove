package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
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
	makeCreateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData progearCreateSubtestData) {
		db := appCtx.DB()

		subtestData.ppmShipment = testdatagen.MakePPMShipment(db, testdatagen.Assertions{})
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
			suite.HandlerConfig(),
			progearCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.CreateProGearWeightTicketCreated{}, response)

		createdProgear := response.(*progearops.CreateProGearWeightTicketCreated).Payload

		suite.NotEmpty(createdProgear.ID.String())
		suite.NotNil(createdProgear.DocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure - 403- permission denied - wrong application", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, false)

		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&progearops.CreateProGearWeightTicketForbidden{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.ProgearWeightTicketCreator{}
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateProgearWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateProGearWeightTicketHandler{
			suite.HandlerConfig(),
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
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData progearUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.progear = testdatagen.MakeProgearWeightTicket(db, testdatagen.Assertions{})
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
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &subtestData.progear.DocumentID,
				Document:   subtestData.progear.Document,
			},
		})

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

		suite.IsType(&progearops.UpdateProGearWeightTicketCreated{}, response)

		updatedProgear := response.(*progearops.UpdateProGearWeightTicketCreated).Payload
		suite.Equal(subtestData.progear.ID.String(), updatedProgear.ID.String())
		suite.Equal(params.UpdateProGearWeightTicket.Description, *updatedProgear.Description)
	})

	suite.Run("PATCH failure -400 - nil body", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		subtestData.params.UpdateProGearWeightTicket = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.UpdateProGearWeightTicketBadRequest{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
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
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{}
		// This test should fail because of the wrong ID
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("f20d9c9b-2de5-4860-ad31-fd5c10e739f6"))
		params.ProGearWeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &internalmessages.UpdateProGearWeightTicket{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ProgearWeightTicketUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
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
			suite.HandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketInternalServerError{}, response)
		errResponse := response.(*progearops.UpdateProGearWeightTicketInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})
}
