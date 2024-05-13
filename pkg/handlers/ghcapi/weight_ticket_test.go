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
	weightticketops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
)

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
	makeUpdateSubtestData := func(authenticateRequest bool) (subtestData weightTicketUpdateSubtestData) {

		// Use fake data:
		subtestData.ppmShipment = factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), nil, nil)
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
			suite.HandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteWeightTicketHandler() {
	// Reusable objects
	estimator := mocks.PPMEstimator{}
	Deleter := weightticket.NewWeightTicketDeleter(
		weightticket.NewWeightTicketFetcher(),
		&estimator,
	)

	type DeleteSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.DeleteWeightTicketParams
		handler      DeleteWeightTicketHandler
	}
	makeDeleteSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData DeleteSubtestData) {
		// Use fake data:
		subtestData.ppmShipment = factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), nil, nil)
		subtestData.weightTicket = subtestData.ppmShipment.WeightTickets[0]

		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		req := httptest.NewRequest("DELETE", endpoint, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		subtestData.params = weightticketops.DeleteWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
		}

		subtestData.handler = DeleteWeightTicketHandler{
			suite.createS3HandlerConfig(),
			Deleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete  Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()
		mockDeleter := mocks.WeightTicketDeleter{}

		subtestData := makeDeleteSubtestData(appCtx, true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.weightTicket.FullDocument,
				LinkOnly: true,
			},
		}, nil)

		mockDeleter.On("DeleteWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil)

		mockDeleter.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(subtestData.ppmShipment.FinalIncentive)

		handler := DeleteWeightTicketHandler{
			suite.HandlerConfig(),
			&mockDeleter,
		}

		// Validate incoming payload: no body to validate
		response := handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketNoContent{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, true)
		params := subtestData.params
		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("3ce0e367-a337-46e3-b4cf-f79aebc4f6c8"))
		params.WeightTicketID = *wrongUUIDString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 500", func() {
		mockDeleter := mocks.WeightTicketDeleter{}
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(err)

		handler := DeleteWeightTicketHandler{
			suite.HandlerConfig(),
			&mockDeleter,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketInternalServerError{}, response)
	})
}
