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
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData weightTicketUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(db, testdatagen.Assertions{})
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
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		ppmShipmentUpdater.On(
			"UpdatePPMShipmentWithDefaultCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, nil)

		// Add full and empty weights
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{
			EmptyWeight: handlers.FmtInt64(1),
			FullWeight:  handlers.FmtInt64(4000),
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
	})

	suite.Run("PATCH failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{}
		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("cde78daf-802f-491f-a230-fc1fdcfe6595"))
		params.WeightTicketID = *wrongUUIDString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{}
		params.IfMatch = "wrong-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.WeightTicketUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
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
