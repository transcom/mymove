package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	ppmdocumentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestGetPPMDocumentsHandlerUnit() {
	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			UserUploader: userUploader,
		})

		ppmShipment.WeightTickets = append(
			ppmShipment.WeightTickets,
			testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
				ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
				PPMShipment:   ppmShipment,
				UserUploader:  userUploader,
			}),
		)

		for i := 1; i < 3; i++ {
			ppmShipment.MovingExpenses = append(
				ppmShipment.MovingExpenses,
				testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
					ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					PPMShipment:   ppmShipment,
					UserUploader:  userUploader,
				}),
			)

		}

		for i := 1; i < 4; i++ {
			ppmShipment.ProgearExpenses = append(
				ppmShipment.ProgearExpenses,
				testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
					ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					PPMShipment:   ppmShipment,
					UserUploader:  userUploader,
				}),
			)
		}

	})

	setUpRequestAndParams := func() ppmdocumentops.GetPPMDocumentsParams {
		endpoint := fmt.Sprintf("/shipments/%s/ppm-documents", ppmShipment.Shipment.ID.String())

		req := httptest.NewRequest("GET", endpoint, nil)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := ppmdocumentops.GetPPMDocumentsParams{
			HTTPRequest: req,
			ShipmentID:  handlers.FmtUUIDValue(ppmShipment.Shipment.ID),
		}

		return params
	}

	setUpMockPPMDocumentFetcher := func(returnValues ...interface{}) services.PPMDocumentFetcher {
		mockPPMDocumentFetcher := &mocks.PPMDocumentFetcher{}

		mockPPMDocumentFetcher.On("GetPPMDocuments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(returnValues...)

		return mockPPMDocumentFetcher
	}

	setUpHandler := func(ppmDocumentFetcher services.PPMDocumentFetcher) GetPPMDocumentsHandler {
		return GetPPMDocumentsHandler{
			suite.createS3HandlerConfig(),
			ppmDocumentFetcher,
		}
	}

	suite.Run("Returns an error if the request is not coming from the office app", func() {
		params := setUpRequestAndParams()

		params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		ppmDocumentsFetcher := setUpMockPPMDocumentFetcher(ppmShipment.ShipmentID, nil)

		handler := setUpHandler(ppmDocumentsFetcher)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.GetPPMDocumentsForbidden{}, response) {
			payload := response.(*ppmdocumentops.GetPPMDocumentsForbidden).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.True(strings.HasPrefix(*payload.Message, "Instance: "))
		}
	})

	serverErrorCases := map[string]error{
		"issues retrieving ppm documents": apperror.NewQueryError("PPMDocument", nil, "Unable to find PPMDocuments"),
		"unexpected error":                apperror.NewConflictError(uuid.Nil, "Unexpected error"),
	}

	for errorDetail, fetcherError := range serverErrorCases {
		errorDetail := errorDetail
		fetcherError := fetcherError

		suite.Run(fmt.Sprintf("Returns a server error if there is an %s", errorDetail), func() {
			params := setUpRequestAndParams()

			ppmDocumentFetcher := setUpMockPPMDocumentFetcher(nil, fetcherError)

			handler := setUpHandler(ppmDocumentFetcher)

			response := handler.Handle(params)

			if suite.IsType(&ppmdocumentops.GetPPMDocumentsInternalServerError{}, response) {
				payload := response.(*ppmdocumentops.GetPPMDocumentsInternalServerError).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.True(strings.HasPrefix(*payload.Message, "Instance:"))
			}
		})
	}

	suite.Run("Returns 200 when PPM documents are found", func() {
		params := setUpRequestAndParams()

		ppmDocuments := models.PPMDocuments{
			WeightTickets:   ppmShipment.WeightTickets,
			MovingExpenses:  ppmShipment.MovingExpenses,
			ProgearExpenses: ppmShipment.ProgearExpenses,
		}

		ppmDocumentFetcher := setUpMockPPMDocumentFetcher(&ppmDocuments, nil)

		handler := setUpHandler(ppmDocumentFetcher)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.GetPPMDocumentsOK{}, response) {
			okResponse := response.(*ppmdocumentops.GetPPMDocumentsOK)
			returnedPPMDocuments := okResponse.Payload

			suite.NoError(returnedPPMDocuments.Validate(strfmt.Default))

			suite.Equal(len(ppmShipment.WeightTickets), len(returnedPPMDocuments.WeightTickets))
			suite.Equal(len(ppmShipment.ProgearExpenses), len(returnedPPMDocuments.ProGearWeightTickets))
			suite.Equal(len(ppmShipment.MovingExpenses), len(returnedPPMDocuments.MovingExpenses))

			for i, returnedWeightTicket := range returnedPPMDocuments.WeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID.String(), returnedWeightTicket.ID.String())
			}

			for i, returnedMovingExpense := range returnedPPMDocuments.MovingExpenses {
				suite.Equal(ppmShipment.MovingExpenses[i].ID.String(), returnedMovingExpense.ID.String())
			}

			for i, returnedProGearWeightTicket := range returnedPPMDocuments.ProGearWeightTickets {
				suite.Equal(ppmShipment.ProgearExpenses[i].ID.String(), returnedProGearWeightTicket.ID.String())
			}
		}
	})
}

func (suite *HandlerSuite) TestGetPPMDocumentsHandlerIntegration() {
	ppmDocumentsFetcher := ppmshipment.NewPPMDocumentFetcher()

	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			UserUploader: userUploader,
		})

		ppmShipment.WeightTickets = append(
			ppmShipment.WeightTickets,
			testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
				ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
				PPMShipment:   ppmShipment,
				UserUploader:  userUploader,
			}),
		)

		for i := 1; i < 3; i++ {
			ppmShipment.MovingExpenses = append(
				ppmShipment.MovingExpenses,
				testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
					ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					PPMShipment:   ppmShipment,
					UserUploader:  userUploader,
				}),
			)

		}

		for i := 1; i < 4; i++ {
			ppmShipment.ProgearExpenses = append(
				ppmShipment.ProgearExpenses,
				testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
					ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					PPMShipment:   ppmShipment,
					UserUploader:  userUploader,
				}),
			)
		}

	})

	setUpParamsAndHandler := func() (ppmdocumentops.GetPPMDocumentsParams, GetPPMDocumentsHandler) {
		endpoint := fmt.Sprintf("/shipments/%s/ppm-documents", ppmShipment.Shipment.ID.String())

		req := httptest.NewRequest("GET", endpoint, nil)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := ppmdocumentops.GetPPMDocumentsParams{
			HTTPRequest: req,
			ShipmentID:  handlers.FmtUUIDValue(ppmShipment.Shipment.ID),
		}

		handler := GetPPMDocumentsHandler{
			suite.createS3HandlerConfig(),
			ppmDocumentsFetcher,
		}

		return params, handler
	}

	suite.Run("Returns 200 when weight tickets are found", func() {
		params, handler := setUpParamsAndHandler()

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.GetPPMDocumentsOK{}, response) {
			okResponse := response.(*ppmdocumentops.GetPPMDocumentsOK)
			returnedPPMDocuments := okResponse.Payload

			suite.NoError(returnedPPMDocuments.Validate(strfmt.Default))

			suite.Equal(len(ppmShipment.WeightTickets), len(returnedPPMDocuments.WeightTickets))
			suite.Equal(len(ppmShipment.ProgearExpenses), len(returnedPPMDocuments.ProGearWeightTickets))
			suite.Equal(len(ppmShipment.MovingExpenses), len(returnedPPMDocuments.MovingExpenses))

			for i, returnedWeightTicket := range returnedPPMDocuments.WeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID.String(), returnedWeightTicket.ID.String())
			}

			for i, returnedMovingExpense := range returnedPPMDocuments.MovingExpenses {
				suite.Equal(ppmShipment.MovingExpenses[i].ID.String(), returnedMovingExpense.ID.String())
			}

			for i, returnedProGearWeightTicket := range returnedPPMDocuments.ProGearWeightTickets {
				suite.Equal(ppmShipment.ProgearExpenses[i].ID.String(), returnedProGearWeightTicket.ID.String())
			}
		}
	})
}
