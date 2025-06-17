package ghcapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
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
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestGetPPMDocumentsHandlerUnit() {
	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)

		ppmShipment.WeightTickets = append(
			ppmShipment.WeightTickets,
			factory.BuildWeightTicket(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		for i := 1; i < 3; i++ {
			ppmShipment.MovingExpenses = append(
				ppmShipment.MovingExpenses,
				factory.BuildMovingExpense(suite.DB(), []factory.Customization{
					{
						Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
						LinkOnly: true,
					},
					{
						Model:    ppmShipment,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{},
						ExtendedParams: &factory.UserUploadExtendedParams{
							UserUploader: userUploader,
							AppContext:   suite.AppContextForTest(),
						},
					},
				}, nil),
			)

		}

		for i := 1; i < 4; i++ {
			ppmShipment.ProgearWeightTickets = append(
				ppmShipment.ProgearWeightTickets,
				factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
					{
						Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
						LinkOnly: true,
					},
					{
						Model:    ppmShipment,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{},
						ExtendedParams: &factory.UserUploadExtendedParams{
							UserUploader: userUploader,
							AppContext:   suite.AppContextForTest(),
						},
					},
				}, nil),
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

		ppmShipment = factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)

		ppmShipment.WeightTickets = append(
			ppmShipment.WeightTickets,
			factory.BuildWeightTicket(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		for i := 1; i < 3; i++ {
			ppmShipment.MovingExpenses = append(
				ppmShipment.MovingExpenses,
				factory.BuildMovingExpense(suite.DB(), []factory.Customization{
					{
						Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
						LinkOnly: true,
					},
					{
						Model:    ppmShipment,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{},
						ExtendedParams: &factory.UserUploadExtendedParams{
							UserUploader: userUploader,
							AppContext:   suite.AppContextForTest(),
						},
					},
				}, nil),
			)

		}

		for i := 1; i < 4; i++ {
			ppmShipment.ProgearWeightTickets = append(
				ppmShipment.ProgearWeightTickets,
				factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
					{
						Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
						LinkOnly: true,
					},
					{
						Model:    ppmShipment,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{},
						ExtendedParams: &factory.UserUploadExtendedParams{
							UserUploader: userUploader,
							AppContext:   suite.AppContextForTest(),
						},
					},
				}, nil),
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

		suite.IsType(&ppmdocumentops.GetPPMDocumentsOK{}, response)
		okResponse := response.(*ppmdocumentops.GetPPMDocumentsOK)
		returnedPPMDocuments := okResponse.Payload

		suite.NoError(returnedPPMDocuments.Validate(strfmt.Default))

		suite.Equal(len(ppmShipment.WeightTickets), len(returnedPPMDocuments.WeightTickets))
		suite.Equal(len(ppmShipment.ProgearWeightTickets), len(returnedPPMDocuments.ProGearWeightTickets))
		suite.Equal(len(ppmShipment.MovingExpenses), len(returnedPPMDocuments.MovingExpenses))

		// extract and sort IDs for WeightTickets to avoid flaky failures
		expectedWeightTicketIDs := make([]string, len(ppmShipment.WeightTickets))
		for i, wt := range ppmShipment.WeightTickets {
			expectedWeightTicketIDs[i] = wt.ID.String()
		}
		actualWeightTicketIDs := make([]string, len(returnedPPMDocuments.WeightTickets))
		for i, wt := range returnedPPMDocuments.WeightTickets {
			actualWeightTicketIDs[i] = wt.ID.String()
		}
		sort.Strings(expectedWeightTicketIDs)
		sort.Strings(actualWeightTicketIDs)
		suite.Equal(expectedWeightTicketIDs, actualWeightTicketIDs, "WeightTicket IDs should match")

		// extract and sort IDs for ProGearWeightTickets to avoid flaky failures
		expectedProGearIDs := make([]string, len(ppmShipment.ProgearWeightTickets))
		for i, pt := range ppmShipment.ProgearWeightTickets {
			expectedProGearIDs[i] = pt.ID.String()
		}
		actualProGearIDs := make([]string, len(returnedPPMDocuments.ProGearWeightTickets))
		for i, pt := range returnedPPMDocuments.ProGearWeightTickets {
			actualProGearIDs[i] = pt.ID.String()
		}
		sort.Strings(expectedProGearIDs)
		sort.Strings(actualProGearIDs)
		suite.Equal(expectedProGearIDs, actualProGearIDs, "ProGearWeightTicket IDs should match")

		// extract and sort IDs for MovingExpenses to avoid flaky failures
		expectedMovingExpenseIDs := make([]string, len(ppmShipment.MovingExpenses))
		for i, me := range ppmShipment.MovingExpenses {
			expectedMovingExpenseIDs[i] = me.ID.String()
		}
		actualMovingExpenseIDs := make([]string, len(returnedPPMDocuments.MovingExpenses))
		for i, me := range returnedPPMDocuments.MovingExpenses {
			actualMovingExpenseIDs[i] = me.ID.String()
		}
		sort.Strings(expectedMovingExpenseIDs)
		sort.Strings(actualMovingExpenseIDs)
		suite.Equal(expectedMovingExpenseIDs, actualMovingExpenseIDs, "MovingExpense IDs should match")
	})

}

func (suite *HandlerSuite) TestFinishPPMDocumentsReviewHandlerUnit() {
	var ppmShipment models.PPMShipment

	setUpPPMShipment := func() models.PPMShipment {
		ppmShipment = factory.BuildPPMShipmentWithApprovedDocuments(nil)

		move := factory.BuildMove(suite.DB(), nil, nil)

		ppmShipment.ID = uuid.Must(uuid.NewV4())
		ppmShipment.CreatedAt = time.Now()
		ppmShipment.UpdatedAt = ppmShipment.CreatedAt.AddDate(0, 0, 5)
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Must(uuid.NewV4())
		ppmShipment.Shipment.MoveTaskOrderID = move.ID

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
		expectedPPMShipment.Status = models.PPMShipmentStatusCloseoutComplete

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

	setUpSignedCertificationCreatorMock := func(signedCert models.SignedCertification) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(&signedCert, nil)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(signedCert models.SignedCertification) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(&signedCert, nil)

		return mockUpdater
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	ppmEstimator := &mocks.PPMEstimator{}
	SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(mockPPMCloseoutFetcher, ppmEstimator)
	mockPPMCloseoutFetcher.On("GetActualWeight", mock.AnythingOfType("*models.PPMShipment")).Return(unit.Pound(1000))

	setUpParamsAndHandler := func(ppmShipment models.PPMShipment, officeUser models.OfficeUser, signedCert models.SignedCertification) (ppmdocumentops.FinishDocumentReviewParams, FinishDocumentReviewHandler) {
		reviewer := ppmshipment.NewPPMShipmentReviewDocuments(ppmShipmentRouter, setUpSignedCertificationCreatorMock(signedCert), setUpSignedCertificationUpdaterMock(signedCert), SSWPPMComputer)
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

		certType := models.SignedCertificationTypePPMPAYMENT
		signedCert := models.SignedCertification{
			CertificationType: &certType,
			CertificationText: "LEGAL",
			Signature:         "ACCEPT",
			Date:              time.Now(),
		}

		params, handler := setUpParamsAndHandler(shipmentWithUnknownID, officeUser, signedCert)

		response := handler.Handle(params)

		suite.IsType(&ppmdocumentops.FinishDocumentReviewNotFound{}, response)
	})

	suite.Run("Returns an error if the PPM shipment is not awaiting payment review", func() {
		draftPpmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		draftPpmShipment.Status = models.PPMShipmentStatusDraft
		suite.NoError(suite.DB().Save(&draftPpmShipment))
		certType := models.SignedCertificationTypePPMPAYMENT
		signedCert := models.SignedCertification{
			CertificationType: &certType,
			CertificationText: "LEGAL",
			Signature:         "ACCEPT",
			Date:              time.Now(),
		}
		params, handler := setUpParamsAndHandler(draftPpmShipment, officeUser, signedCert)

		response := handler.Handle(params)

		suite.IsType(&ppmdocumentops.FinishDocumentReviewConflict{}, response)
	})

	suite.Run("Can successfully submit a PPM shipment for close out", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		certType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT

		signedCert := models.SignedCertification{
			PpmID:             &ppmShipment.ID,
			MoveID:            ppmShipment.Shipment.MoveTaskOrderID,
			CertificationType: &certType,
			CertificationText: "LEGAL",
			Signature:         "ACCEPT",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Date:              time.Now(),
		}
		shipmentCert := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment.Shipment.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.SignedCertification{
					PpmID:             &ppmShipment.ID,
					CertificationType: &certType,
				},
			},
		}, nil)
		ppmShipment.SignedCertification = &shipmentCert
		suite.NoError(suite.DB().Save(&ppmShipment))
		params, handler := setUpParamsAndHandler(ppmShipment, officeUser, signedCert)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.FinishDocumentReviewOK{}, response) {
			okResponse := response.(*ppmdocumentops.FinishDocumentReviewOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(ppmShipment.ID, returnedPPMShipment.ID)
			suite.Equal(string(models.PPMShipmentStatusCloseoutComplete), string(returnedPPMShipment.Status))
		}
	})

	suite.Run("Sets PPM to CLOSEOUT COMPLETE if there are rejected documents", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsToBeResubmitted(suite.DB(), nil, nil)
		certType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
		signedCert := models.SignedCertification{
			PpmID:             &ppmShipment.ID,
			MoveID:            ppmShipment.Shipment.MoveTaskOrderID,
			CertificationType: &certType,
			CertificationText: "LEGAL",
			Signature:         "ACCEPT",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Date:              time.Now(),
		}
		shipmentCert := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment.Shipment.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.SignedCertification{
					PpmID:             &ppmShipment.ID,
					CertificationType: &certType,
				},
			},
		}, nil)
		ppmShipment.SignedCertification = &shipmentCert
		ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout
		suite.NoError(suite.DB().Save(&ppmShipment))
		params, handler := setUpParamsAndHandler(ppmShipment, officeUser, signedCert)

		response := handler.Handle(params)

		if suite.IsType(&ppmdocumentops.FinishDocumentReviewOK{}, response) {
			okResponse := response.(*ppmdocumentops.FinishDocumentReviewOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(ppmShipment.ID, returnedPPMShipment.ID)
			suite.Equal(string(models.PPMShipmentStatusCloseoutComplete), string(returnedPPMShipment.Status))
		}
	})
}

func (suite *HandlerSuite) TestShowAOAPacketHandler() {
	suite.Run("Successful ShowAOAPacketHandler - 200", func() {
		mockSSWPPMComputer := mocks.SSWPPMComputer{}
		mockSSWPPMGenerator := mocks.SSWPPMGenerator{}
		mockAOAPacketCreator := mocks.AOAPacketCreator{}
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), 25*uploader.MB)
		suite.NoError(err)

		ppmshipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(suite.DB(), userUploader)

		handlerConfig := suite.NewHandlerConfig()
		handler := showAOAPacketHandler{
			HandlerConfig:    handlerConfig,
			SSWPPMComputer:   &mockSSWPPMComputer,
			SSWPPMGenerator:  &mockSSWPPMGenerator,
			AOAPacketCreator: &mockAOAPacketCreator,
		}

		mockAOAPacketCreator.On("CreateAOAPacket", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("bool")).Return(nil, "", nil)
		mockAOAPacketCreator.On("CleanupAOAPacketDir", mock.AnythingOfType("string")).Return(nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		ppmshipmentid := ppmshipment.ID
		// ppmshipmentid := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/aoa-packet/", ppmshipmentid), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := ppmdocumentops.ShowAOAPacketParams{
			HTTPRequest:   request,
			PpmShipmentID: ppmshipmentid.String(),
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&ppmdocumentops.ShowAOAPacketOK{}, response)
		contentDisposition := response.(*ppmdocumentops.ShowAOAPacketOK).ContentDisposition

		// Validate filename content disposition formatting
		found := regexp.MustCompile(`inline; filename=\"AOA-\d{14}.pdf\"`).FindString(contentDisposition)
		suite.NotEmpty(found, "filename format invalid: %s", contentDisposition)
	})

	suite.Run("Successful ShowAOAPacketHandler - error generating PDF - 500", func() {
		mockSSWPPMComputer := mocks.SSWPPMComputer{}
		mockSSWPPMGenerator := mocks.SSWPPMGenerator{}
		mockAOAPacketCreator := mocks.AOAPacketCreator{}
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), 25*uploader.MB)
		suite.NoError(err)

		ppmshipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(suite.DB(), userUploader)

		handlerConfig := suite.NewHandlerConfig()
		handler := showAOAPacketHandler{
			HandlerConfig:    handlerConfig,
			SSWPPMComputer:   &mockSSWPPMComputer,
			SSWPPMGenerator:  &mockSSWPPMGenerator,
			AOAPacketCreator: &mockAOAPacketCreator,
		}

		mockAOAPacketCreator.On("CreateAOAPacket", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("bool")).Return(nil, "", errors.New("Mock error"))
		mockAOAPacketCreator.On("CleanupAOAPacketDir", mock.AnythingOfType("string")).Return(nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		ppmshipmentid := ppmshipment.ID
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/order/download", ppmshipmentid), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := ppmdocumentops.ShowAOAPacketParams{
			HTTPRequest:   request,
			PpmShipmentID: ppmshipmentid.String(),
		}
		response := handler.Handle(params)
		showAOAPacketResponse := response.(*ppmdocumentops.ShowAOAPacketInternalServerError)

		suite.Assertions.IsType(&ppmdocumentops.ShowAOAPacketInternalServerError{}, showAOAPacketResponse)
	})

	suite.Run("Successful ShowAOAPacketHandler - Missing/empty/incorrect PPMShipmentId - 400", func() {
		mockSSWPPMComputer := mocks.SSWPPMComputer{}
		mockSSWPPMGenerator := mocks.SSWPPMGenerator{}
		mockAOAPacketCreator := mocks.AOAPacketCreator{}

		handlerConfig := suite.NewHandlerConfig()
		handler := showAOAPacketHandler{
			HandlerConfig:    handlerConfig,
			SSWPPMComputer:   &mockSSWPPMComputer,
			SSWPPMGenerator:  &mockSSWPPMGenerator,
			AOAPacketCreator: &mockAOAPacketCreator,
		}

		mockAOAPacketCreator.On("CreateAOAPacket", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("bool")).Return(nil, "", errors.New("Mock error"))
		mockAOAPacketCreator.On("CleanupAOAPacketDir", mock.AnythingOfType("string")).Return(nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		ppmshipmentid := ""
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/order/download", ppmshipmentid), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := ppmdocumentops.ShowAOAPacketParams{
			HTTPRequest:   request,
			PpmShipmentID: ppmshipmentid,
		}
		response := handler.Handle(params)
		showAOAPacketResponse := response.(*ppmdocumentops.ShowAOAPacketBadRequest)

		suite.Assertions.IsType(&ppmdocumentops.ShowAOAPacketBadRequest{}, showAOAPacketResponse)
	})

}

func (suite *HandlerSuite) TestShowPaymentPacketHandler() {
	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)
	ppmShipment.ID = uuid.Must(uuid.NewV4())
	suite.Run("Successful ShowAOAPacketHandler - 200", func() {

		mockPaymentPacketCreator := mocks.PaymentPacketCreator{}
		handler := ShowPaymentPacketHandler{
			HandlerConfig:        suite.createS3HandlerConfig(),
			PaymentPacketCreator: &mockPaymentPacketCreator,
		}

		mockPaymentPacketCreator.On("GenerateDefault",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID")).Return(nil, "", nil)

		mockPaymentPacketCreator.On("CleanupPaymentPacketDir", mock.AnythingOfType("string")).Return(nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		ppmshipmentid := ppmShipment.ID
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/payment-packet/", ppmshipmentid), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := ppmdocumentops.ShowPaymentPacketParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(ppmshipmentid.String()),
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&ppmdocumentops.ShowPaymentPacketOK{}, response)
		contentDisposition := response.(*ppmdocumentops.ShowPaymentPacketOK).ContentDisposition

		// Validate filename content disposition formatting
		found := regexp.MustCompile(`inline; filename=\"ppm_payment_packet-\d{14}.pdf\"`).FindString(contentDisposition)
		suite.NotEmpty(found, "filename format invalid: %s", contentDisposition)
	})

	suite.Run("Unsuccessful ShowPaymentPacketHandler - InternalServerError", func() {
		mockPaymentPacketCreator := mocks.PaymentPacketCreator{}
		handler := ShowPaymentPacketHandler{
			HandlerConfig:        suite.createS3HandlerConfig(),
			PaymentPacketCreator: &mockPaymentPacketCreator,
		}

		mockPaymentPacketCreator.On("GenerateDefault",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID")).Return(nil, "", errors.New("Mock error"))

		mockPaymentPacketCreator.On("CleanupPaymentPacketDir", mock.AnythingOfType("string")).Return(nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		ppmshipmentid := ppmShipment.ID
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/payment-packet/", ppmshipmentid), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := ppmdocumentops.ShowPaymentPacketParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(ppmshipmentid.String()),
		}
		response := handler.Handle(params)
		showPaymentPacketResponse := response.(*ppmdocumentops.ShowPaymentPacketInternalServerError)

		suite.Assertions.IsType(&ppmdocumentops.ShowPaymentPacketInternalServerError{}, showPaymentPacketResponse)
	})

	suite.Run("Unsuccessful ShowPaymentPacketHandler - NotFoundError", func() {
		mockPaymentPacketCreator := mocks.PaymentPacketCreator{}
		handler := ShowPaymentPacketHandler{
			HandlerConfig:        suite.createS3HandlerConfig(),
			PaymentPacketCreator: &mockPaymentPacketCreator,
		}

		mockPaymentPacketCreator.On("GenerateDefault",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID")).Return(nil, "", apperror.NotFoundError{})

		mockPaymentPacketCreator.On("CleanupPaymentPacketDir", mock.AnythingOfType("string")).Return(nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		ppmshipmentid := ppmShipment.ID
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/payment-packet/", ppmshipmentid), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := ppmdocumentops.ShowPaymentPacketParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(ppmshipmentid.String()),
		}
		response := handler.Handle(params)
		showPaymentPacketResponse := response.(*ppmdocumentops.ShowPaymentPacketNotFound)

		suite.Assertions.IsType(&ppmdocumentops.ShowPaymentPacketNotFound{}, showPaymentPacketResponse)
	})
}
