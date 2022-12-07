package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	weightticketops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestGetWeightTicketsHandlerUnit() {
	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			UserUploader: userUploader,
		})

		for i := 1; i < 3; i++ {
			ppmShipment.WeightTickets = append(
				ppmShipment.WeightTickets,
				testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
					ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					PPMShipment:   ppmShipment,
					UserUploader:  userUploader,
				}),
			)
		}
	})

	setUpRequestAndParams := func() weightticketops.GetWeightTicketsParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-tickets", ppmShipment.ID.String())

		req := httptest.NewRequest("GET", endpoint, nil)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := weightticketops.GetWeightTicketsParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
		}

		return params
	}

	setUpMockWeightTicketFetcher := func(returnValues ...interface{}) services.WeightTicketFetcher {
		mockWeightTicketFetcher := &mocks.WeightTicketFetcher{}

		mockWeightTicketFetcher.On("ListWeightTickets",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(returnValues...)

		return mockWeightTicketFetcher
	}

	setUpHandler := func(weightTicketFetcher services.WeightTicketFetcher) GetWeightTicketsHandler {
		return GetWeightTicketsHandler{
			suite.createS3HandlerConfig(),
			weightTicketFetcher,
		}
	}

	suite.Run("Returns an error if the request is not coming from the office app", func() {
		params := setUpRequestAndParams()

		params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		weightTicketFetcher := setUpMockWeightTicketFetcher(ppmShipment.WeightTickets, nil)

		handler := setUpHandler(weightTicketFetcher)

		response := handler.Handle(params)

		if suite.IsType(&weightticketops.GetWeightTicketsForbidden{}, response) {
			payload := response.(*weightticketops.GetWeightTicketsForbidden).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.True(strings.HasPrefix(*payload.Message, "Instance:"))
		}
	})

	suite.Run("Returns a forbidden error if the fetcher returns a forbidden error", func() {
		params := setUpRequestAndParams()

		fetcherError := apperror.NewForbiddenError("forbidden")
		weightTicketFetcher := setUpMockWeightTicketFetcher(nil, fetcherError)

		handler := setUpHandler(weightTicketFetcher)

		response := handler.Handle(params)

		if suite.IsType(&weightticketops.GetWeightTicketsForbidden{}, response) {
			payload := response.(*weightticketops.GetWeightTicketsForbidden).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.True(strings.HasPrefix(*payload.Message, "Instance:"))
		}
	})

	serverErrorCases := map[string]error{
		"issues retrieving weight tickets": apperror.NewQueryError("WeightTicket", nil, "Unable to find WeightTickets"),
		"unexpected error":                 apperror.NewConflictError(uuid.Nil, "Unexpected error"),
	}

	for errorDetail, fetcherError := range serverErrorCases {
		errorDetail := errorDetail
		fetcherError := fetcherError

		suite.Run(fmt.Sprintf("Returns a server error if there is an %s", errorDetail), func() {
			params := setUpRequestAndParams()

			weightTicketFetcher := setUpMockWeightTicketFetcher(nil, fetcherError)

			handler := setUpHandler(weightTicketFetcher)

			response := handler.Handle(params)

			if suite.IsType(&weightticketops.GetWeightTicketsInternalServerError{}, response) {
				payload := response.(*weightticketops.GetWeightTicketsInternalServerError).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.True(strings.HasPrefix(*payload.Message, "Instance:"))
			}
		})
	}

	suite.Run("Returns 200 when weight tickets are found", func() {
		params := setUpRequestAndParams()

		weightTicketFetcher := setUpMockWeightTicketFetcher(ppmShipment.WeightTickets, nil)

		handler := setUpHandler(weightTicketFetcher)

		response := handler.Handle(params)

		if suite.IsType(&weightticketops.GetWeightTicketsOK{}, response) {
			okResponse := response.(*weightticketops.GetWeightTicketsOK)
			returnedWeightTickets := okResponse.Payload

			suite.NoError(returnedWeightTickets.Validate(strfmt.Default))

			suite.Equal(len(ppmShipment.WeightTickets), len(returnedWeightTickets))

			for i, returnedWeightTicket := range returnedWeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID.String(), returnedWeightTicket.ID.String())
			}
		}
	})
}

func (suite *HandlerSuite) TestGetWeightTicketsHandlerIntegration() {
	weightTicketFetcher := weightticket.NewWeightTicketFetcher()

	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			UserUploader: userUploader,
		})

		for i := 1; i < 3; i++ {
			ppmShipment.WeightTickets = append(
				ppmShipment.WeightTickets,
				testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
					ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					PPMShipment:   ppmShipment,
					UserUploader:  userUploader,
				}),
			)
		}
	})

	setUpParamsAndHandler := func() (weightticketops.GetWeightTicketsParams, GetWeightTicketsHandler) {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-tickets", ppmShipment.ID.String())

		req := httptest.NewRequest("GET", endpoint, nil)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := weightticketops.GetWeightTicketsParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
		}

		handler := GetWeightTicketsHandler{
			suite.createS3HandlerConfig(),
			weightTicketFetcher,
		}

		return params, handler
	}

	suite.Run("Returns 200 when weight tickets are found", func() {
		params, handler := setUpParamsAndHandler()

		response := handler.Handle(params)

		if suite.IsType(&weightticketops.GetWeightTicketsOK{}, response) {
			okResponse := response.(*weightticketops.GetWeightTicketsOK)
			returnedWeightTickets := okResponse.Payload

			suite.NoError(returnedWeightTickets.Validate(strfmt.Default))

			suite.Equal(len(ppmShipment.WeightTickets), len(returnedWeightTickets))

			for i, returnedWeightTicket := range returnedWeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID.String(), returnedWeightTicket.ID.String())
			}
		}
	})
}

func (suite *HandlerSuite) TestUpdateWeightTicketHandler() {
	// Reusable objects
	ppmShipmentUpdater := mocks.PPMShipmentUpdater{}
	weightTicketFetcher := weightticket.NewWeightTicketFetcher()
	weightTicketUpdater := weightticket.NewCustomerWeightTicketUpdater(weightTicketFetcher, &ppmShipmentUpdater)

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
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		eTag := etag.GenerateEtag(subtestData.weightTicket.UpdatedAt)

		officeUser := testdatagen.MakeDefaultOfficeUser(db)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
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

		// An upload must exist if trailer is owned and qualifies to be claimed
		testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &subtestData.weightTicket.ProofOfTrailerOwnershipDocumentID,
				Document:   subtestData.weightTicket.ProofOfTrailerOwnershipDocument,
			},
		})

		// Add vehicleDescription
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
}
