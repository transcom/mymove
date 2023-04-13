package ghcapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	ppmdocumentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
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
			ppmShipment.ProgearWeightTickets = append(
				ppmShipment.ProgearWeightTickets,
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

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})

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
			WeightTickets:        ppmShipment.WeightTickets,
			MovingExpenses:       ppmShipment.MovingExpenses,
			ProgearWeightTickets: ppmShipment.ProgearWeightTickets,
		}

		ppmDocumentFetcher := setUpMockPPMDocumentFetcher(&ppmDocuments, nil)

		handler := setUpHandler(ppmDocumentFetcher)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.GetPPMDocumentsOK{}, response) {
			okResponse := response.(*ppmdocumentops.GetPPMDocumentsOK)
			returnedPPMDocuments := okResponse.Payload

			suite.NoError(returnedPPMDocuments.Validate(strfmt.Default))

			suite.Equal(len(ppmShipment.WeightTickets), len(returnedPPMDocuments.WeightTickets))
			suite.Equal(len(ppmShipment.ProgearWeightTickets), len(returnedPPMDocuments.ProGearWeightTickets))
			suite.Equal(len(ppmShipment.MovingExpenses), len(returnedPPMDocuments.MovingExpenses))

			for i, returnedWeightTicket := range returnedPPMDocuments.WeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID.String(), returnedWeightTicket.ID.String())
			}

			for i, returnedMovingExpense := range returnedPPMDocuments.MovingExpenses {
				suite.Equal(ppmShipment.MovingExpenses[i].ID.String(), returnedMovingExpense.ID.String())
			}

			for i, returnedProGearWeightTicket := range returnedPPMDocuments.ProGearWeightTickets {
				suite.Equal(ppmShipment.ProgearWeightTickets[i].ID.String(), returnedProGearWeightTicket.ID.String())
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
			ppmShipment.ProgearWeightTickets = append(
				ppmShipment.ProgearWeightTickets,
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

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})

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
			suite.Equal(len(ppmShipment.ProgearWeightTickets), len(returnedPPMDocuments.ProGearWeightTickets))
			suite.Equal(len(ppmShipment.MovingExpenses), len(returnedPPMDocuments.MovingExpenses))

			for i, returnedWeightTicket := range returnedPPMDocuments.WeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID.String(), returnedWeightTicket.ID.String())
			}

			for i, returnedMovingExpense := range returnedPPMDocuments.MovingExpenses {
				suite.Equal(ppmShipment.MovingExpenses[i].ID.String(), returnedMovingExpense.ID.String())
			}

			for i, returnedProGearWeightTicket := range returnedPPMDocuments.ProGearWeightTickets {
				suite.Equal(ppmShipment.ProgearWeightTickets[i].ID.String(), returnedProGearWeightTicket.ID.String())
			}
		}
	})
}

func (suite *HandlerSuite) TestFinishPPMDocumentsReviewHandlerUnit() {
	var ppmShipment models.PPMShipment

	setUpPPMShipment := func() models.PPMShipment {
		ppmShipment = testdatagen.MakePPMShipmentWithApprovedDocuments(
			suite.DB(),
			testdatagen.Assertions{
				Stub: true,
			},
		)

		ppmShipment.ID = uuid.Must(uuid.NewV4())
		ppmShipment.CreatedAt = time.Now()
		ppmShipment.UpdatedAt = ppmShipment.CreatedAt.AddDate(0, 0, 5)
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Must(uuid.NewV4())

		return ppmShipment
	}

	setUpRequestAndParams := func(ppmShipment models.PPMShipment, authUser bool) (*http.Request, ppmdocumentops.FinishDocumentReviewParams) {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/finish-document-review", ppmShipment.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})

		if authUser {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}

		params := ppmdocumentops.FinishDocumentReviewParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
		}

		return req, params
	}

	setUpMockPPMDocumentReviewer := func(returnValues ...interface{}) services.PPMShipmentReviewDocuments {
		mockPPMDocumentReviewer := &mocks.PPMShipmentReviewDocuments{}

		mockPPMDocumentReviewer.On("SubmitReviewedDocuments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(returnValues...)

		return mockPPMDocumentReviewer
	}

	setUpHandler := func(ppmDocumentReviewer services.PPMShipmentReviewDocuments) FinishDocumentReviewHandler {
		return FinishDocumentReviewHandler{
			suite.createS3HandlerConfig(),
			ppmDocumentReviewer,
		}
	}

	suite.Run("Returns an error if the request is not coming from the office app", func() {
		ppmShipment := setUpPPMShipment()

		request, params := setUpRequestAndParams(ppmShipment, false)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params.HTTPRequest = request

		ppmDocumentReviewer := setUpMockPPMDocumentReviewer(ppmShipment.ShipmentID, nil)

		handler := setUpHandler(ppmDocumentReviewer)

		response := handler.Handle(params)

		suite.IsType(&ppmdocumentops.FinishDocumentReviewForbidden{}, response)
	})

	suite.Run("Returns an error if the PPMShipment ID in the url is invalid", func() {
		ppmShipment := setUpPPMShipment()
		ppmShipment.ID = uuid.Nil

		_, params := setUpRequestAndParams(ppmShipment, true)

		handler := setUpHandler(setUpMockPPMDocumentReviewer(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmdocumentops.FinishDocumentReviewBadRequest{}, response)
	})

	suite.Run("Returns 200 when a PPM is reviewed", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true)

		expectedPPMShipment := ppmShipment
		expectedPPMShipment.Status = models.PPMShipmentStatusPaymentApproved

		suite.FatalNotNil(expectedPPMShipment.SubmittedAt)

		handler := setUpHandler(setUpMockPPMDocumentReviewer(&expectedPPMShipment, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.FinishDocumentReviewOK{}, response) {
			okResponse := response.(*ppmdocumentops.FinishDocumentReviewOK)

			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))
			suite.EqualUUID(expectedPPMShipment.ID, returnedPPMShipment.ID)

		}
	})

}

func (suite *HandlerSuite) TestResubmitPPMShipmentDocumentationHandlerIntegration() {
	mtoShipmentRouter := mtoshipment.NewShipmentRouter()
	ppmShipmentRouter := ppmshipment.NewPPMShipmentRouter(mtoShipmentRouter)

	officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

	reviewer := ppmshipment.NewPPMShipmentReviewDocuments(ppmShipmentRouter)

	setUpParamsAndHandler := func(ppmShipment models.PPMShipment, officeUser models.OfficeUser) (ppmdocumentops.FinishDocumentReviewParams, FinishDocumentReviewHandler) {
		endpoint := fmt.Sprintf(
			"/ppm-shipments/%s/finish-document-review",
			ppmShipment.ID.String(),
		)

		request := httptest.NewRequest("PATCH", endpoint, nil)

		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := ppmdocumentops.FinishDocumentReviewParams{
			HTTPRequest:   request,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
		}

		handler := FinishDocumentReviewHandler{
			suite.createS3HandlerConfig(),
			reviewer,
		}

		return params, handler
	}

	suite.Run("Returns an error if the PPM shipment is not found", func() {
		shipmentWithUnknownID := models.PPMShipment{
			ID: uuid.Must(uuid.NewV4()),
		}

		params, handler := setUpParamsAndHandler(shipmentWithUnknownID, officeUser)

		response := handler.Handle(params)

		suite.IsType(&ppmdocumentops.FinishDocumentReviewNotFound{}, response)
	})

	suite.Run("Returns an error if the PPM shipment is not awaiting payment review", func() {
		draftPpmShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusDraft,
			},
		})

		params, handler := setUpParamsAndHandler(draftPpmShipment, officeUser)

		response := handler.Handle(params)

		fmt.Printf("Response type: %T\n", response)

		suite.IsType(&ppmdocumentops.FinishDocumentReviewConflict{}, response)
	})

	suite.Run("Can successfully submit a PPM shipment for close out", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{})

		params, handler := setUpParamsAndHandler(ppmShipment, officeUser)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.FinishDocumentReviewOK{}, response) {
			okResponse := response.(*ppmdocumentops.FinishDocumentReviewOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(ppmShipment.ID, returnedPPMShipment.ID)
			suite.Equal(string(models.PPMShipmentStatusPaymentApproved), string(returnedPPMShipment.Status))
		}
	})

	suite.Run("Sets PPM to await customer if there are rejected documents", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsToBeResubmitted(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
		})

		params, handler := setUpParamsAndHandler(ppmShipment, officeUser)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.FinishDocumentReviewOK{}, response) {
			okResponse := response.(*ppmdocumentops.FinishDocumentReviewOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(ppmShipment.ID, returnedPPMShipment.ID)
			suite.Equal(string(models.PPMShipmentStatusWaitingOnCustomer), string(returnedPPMShipment.Status))
		}
	})
}
