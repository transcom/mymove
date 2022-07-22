package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"

	weightticketops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
)

//
// CREATE TEST
//

func (suite *HandlerSuite) TestCreateWeightTicketHandler() {
	// Reusable objects
	weightTicketCreator := weightticket.NewCustomerWeightTicketCreator()

	type weightTicketCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      weightticketops.CreateWeightTicketParams
		handler     CreateWeightTicketHandler
	}
	makeCreateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData weightTicketCreateSubtestData) {
		db := appCtx.DB()

		subtestData.ppmShipment = testdatagen.MakePPMShipment(db, testdatagen.Assertions{})
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight_ticket", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = weightticketops.CreateWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateWeightTicketHandler{
			handlers.NewHandlerConfig(appCtx.DB(), appCtx.Logger()),
			weightTicketCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.CreateWeightTicketOK{}, response)

		createdWeightTicket := response.(*weightticketops.CreateWeightTicketOK).Payload

		suite.NotEmpty(createdWeightTicket.ID.String())
		suite.NotNil(createdWeightTicket.EmptyDocumentID.String())
		suite.NotNil(createdWeightTicket.FullDocumentID.String())
		suite.NotNil(createdWeightTicket.ProofOfTrailerOwnershipDocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.CreateWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.CreateWeightTicketUnauthorized{}, response)
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

		suite.IsType(&weightticketops.CreateWeightTicketForbidden{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.WeightTicketCreator{}
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateWeightTicketHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockCreator,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.CreateWeightTicketInternalServerError{}, response)
	})
}

//
// UPDATE test
//

func (suite *HandlerSuite) TestUpdateWeightTicketHandler() {
	// Reusable objects
	weightTicketUpdater := weightticket.NewCustomerWeightTicketUpdater()

	type weightTicketUpdateSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.UpdateWeightTicketParams
		handler      UpdateWeightTicketHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData weightTicketUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.weightTicket = testdatagen.MakeWeightTicket(db, testdatagen.Assertions{})
		subtestData.ppmShipment = subtestData.weightTicket.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		eTag := etag.GenerateEtag(subtestData.weightTicket.UpdatedAt)
		subtestData.params = weightticketops.UpdateWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
			IfMatch:        eTag,
		}

		subtestData.handler = UpdateWeightTicketHandler{
			suite.createHandlerConfig(),
			weightTicketUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		// Add vehicleDescription
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{
			VehicleDescription:       "Subaru",
			EmptyWeight:              handlers.FmtInt64(1),
			MissingEmptyWeightTicket: false,
			FullWeight:               handlers.FmtInt64(4000),
			MissingFullWeightTicket:  false,
			OwnsTrailer:              true,
			TrailerMeetsCriteria:     true,
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketOK{}, response)

		updatedWeightTicket := response.(*weightticketops.UpdateWeightTicketOK).Payload
		suite.Equal(subtestData.weightTicket.ID.String(), updatedWeightTicket.ID.String())
		suite.Equal(params.UpdateWeightTicketPayload.VehicleDescription, *updatedWeightTicket.VehicleDescription)
	})
	// TODO: for failures pick any field, except the bools, and pass in an empty string for the udpate to trigger failure
	suite.Run("PATCH failure -400 - nil body", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		subtestData.params.UpdateWeightTicketPayload = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.UpdateWeightTicketBadRequest{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{
			VehicleDescription:       "Subaru",
			EmptyWeight:              handlers.FmtInt64(0),
			MissingEmptyWeightTicket: false,
			FullWeight:               handlers.FmtInt64(4000),
			MissingFullWeightTicket:  false,
			OwnsTrailer:              true,
			TrailerMeetsCriteria:     true,
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{}
		// This test should fail because of the wrong ID
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("f20d9c9b-2de5-4860-ad31-fd5c10e739f6"))
		params.WeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.WeightTicketUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{
			VehicleDescription:       "Subaru",
			EmptyWeight:              handlers.FmtInt64(1),
			MissingEmptyWeightTicket: false,
			FullWeight:               handlers.FmtInt64(4000),
			MissingFullWeightTicket:  false,
			OwnsTrailer:              true,
			TrailerMeetsCriteria:     true,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.WeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateWeightTicketHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketInternalServerError{}, response)
		errResponse := response.(*weightticketops.UpdateWeightTicketInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})
}
