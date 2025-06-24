package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	weightticketops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateWeightTicketHandler() {
	// Reusable objects
	weightTicketCreator := weightticket.NewCustomerWeightTicketCreator()

	type weightTicketCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      weightticketops.CreateWeightTicketParams
		handler     CreateWeightTicketHandler
	}
	makeCreateSubtestData := func(authenticateRequest bool) (subtestData weightTicketCreateSubtestData) {
		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight_ticket", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}
		subtestData.params = weightticketops.CreateWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateWeightTicketHandler{
			suite.NewHandlerConfig(),
			weightTicketCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.CreateWeightTicketOK{}, response)

		createdWeightTicket := response.(*weightticketops.CreateWeightTicketOK).Payload

		suite.NotEmpty(createdWeightTicket.ID.String())
		suite.NotNil(createdWeightTicket.EmptyDocumentID.String())
		suite.NotNil(createdWeightTicket.FullDocumentID.String())
		suite.NotNil(createdWeightTicket.ProofOfTrailerOwnershipDocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.CreateWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		subtestData := makeCreateSubtestData(false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.CreateWeightTicketUnauthorized{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.WeightTicketCreator{}

		subtestData := makeCreateSubtestData(true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockCreator,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.CreateWeightTicketInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateWeightTicketHandler() {
	// Reusable objects
	ppmShipmentUpdater := mocks.PPMShipmentUpdater{}
	weightTicketFetcher := weightticket.NewWeightTicketFetcher()
	weightTicketUpdater := weightticket.NewOfficeWeightTicketUpdater(weightTicketFetcher, &ppmShipmentUpdater)

	type weightTicketUpdateSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.UpdateWeightTicketParams
		handler      UpdateWeightTicketHandler
	}
	makeUpdateSubtestData := func(_ bool) (subtestData weightTicketUpdateSubtestData) {

		// Use fake data:
		subtestData.ppmShipment = factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		subtestData.weightTicket = subtestData.ppmShipment.WeightTickets[0]
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req := httptest.NewRequest("PATCH", endpoint, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		eTag := etag.GenerateEtag(subtestData.weightTicket.UpdatedAt)

		subtestData.params = weightticketops.UpdateWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
			IfMatch:        eTag,
		}

		subtestData.handler = UpdateWeightTicketHandler{
			suite.createS3HandlerConfig(),
			weightTicketUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		subtestData := makeUpdateSubtestData(true)

		params := subtestData.params

		ppmShipmentUpdater.On(
			"UpdatePPMShipmentWithDefaultCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, nil)

		// Add full, empty, and adjusted net weight
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{
			EmptyWeight:       handlers.FmtInt64(1),
			FullWeight:        handlers.FmtInt64(4000),
			AdjustedNetWeight: handlers.FmtInt64(3999),
			NetWeightRemarks:  "adjusted weight",
		}

		// Validate incoming payload: no body to validate

		response := subtestData.handler.Handle(params)
		suite.IsType(&weightticketops.UpdateWeightTicketOK{}, response)
		updatedWeightTicket := response.(*weightticketops.UpdateWeightTicketOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedWeightTicket.Validate(strfmt.Default))

		suite.Equal(subtestData.weightTicket.ID.String(), updatedWeightTicket.ID.String())
		suite.Equal(params.UpdateWeightTicketPayload.FullWeight, updatedWeightTicket.FullWeight)
		suite.Equal(params.UpdateWeightTicketPayload.EmptyWeight, updatedWeightTicket.EmptyWeight)
		suite.Equal(params.UpdateWeightTicketPayload.AdjustedNetWeight, updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(params.UpdateWeightTicketPayload.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{}
		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("cde78daf-802f-491f-a230-fc1fdcfe6595"))
		params.WeightTicketID = *wrongUUIDString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{}
		params.IfMatch = "wrong-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.WeightTicketUpdater{}
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		ownsTrailer := true
		trailerMeetsCriteria := true

		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{
			EmptyWeight:          handlers.FmtInt64(1),
			FullWeight:           handlers.FmtInt64(1000),
			OwnsTrailer:          ownsTrailer,
			TrailerMeetsCriteria: trailerMeetsCriteria,
		}

		err := errors.New("ServerError")

		// Might remove the mocks:
		mockUpdater.On("UpdateWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.WeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteWeightTicketHandler() {
	// Create Reusable objects
	fetcher := weightticket.NewWeightTicketFetcher()
	estimator := mocks.PPMEstimator{}
	weightTicketDeleter := weightticket.NewWeightTicketDeleter(fetcher, &estimator)

	type weightTicketDeleteSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.DeleteWeightTicketParams
		handler      DeleteWeightTicketHandler
	}
	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData weightTicketDeleteSubtestData) {
		// Fake data:
		subtestData.weightTicket = factory.BuildWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.weightTicket.PPMShipment
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}
		subtestData.params = weightticketops.
			DeleteWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = DeleteWeightTicketHandler{
			suite.createS3HandlerConfig(),
			weightTicketDeleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete Weight Ticket - Integration Test", func() {
		mockIncentive := unit.Cents(100000)
		estimator.On("FinalIncentiveWithDefaultChecks", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("models.PPMShipment"), mock.AnythingOfType("*models.PPMShipment")).Return(&mockIncentive, nil)

		subtestData := makeDeleteSubtestData(true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.DeleteWeightTicketUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and weight ticket ID don't match", func() {
		subtestData := makeDeleteSubtestData(true)

		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
				LinkOnly: true,
			},
		}, nil)

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		unauthorizedParams := subtestData.params

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&weightticketops.DeleteWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.WeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.WeightTicketDeleter{}

		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(err)

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		handler := DeleteWeightTicketHandler{
			suite.createS3HandlerConfig(),
			&mockDeleter,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketInternalServerError{}, response)
	})
}
