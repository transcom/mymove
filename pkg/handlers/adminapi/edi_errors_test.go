package adminapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	edierrorsop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/e_d_i_errors"
	singleedierrorop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/single_e_d_i_error"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
)

func (suite *HandlerSuite) TestFetchEdiErrorsHandler() {
	suite.Run("Successfully fetches list of EDI errors", func() {

		createdAt := time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC)
		ediCode858 := "FailureForEDI858"
		ediDescription858 := "Failed processing for 858"
		paymentRequestID858 := uuid.Must(uuid.NewV4())

		ediCode824 := "FailureForEDI824"
		ediDescription824 := "Failed processing for 824"
		paymentRequestID824 := uuid.Must(uuid.NewV4())

		ediCode997 := "FailureFor997"
		ediDescription997 := "Failed processing for 997"
		paymentRequestID997 := uuid.Must(uuid.NewV4())

		expectedEdiErrors := models.EdiErrors{
			{
				ID:               uuid.Must(uuid.NewV4()),
				EDIType:          models.EDIType858,
				Code:             &ediCode858,
				Description:      &ediDescription858,
				PaymentRequestID: paymentRequestID858,
				CreatedAt:        createdAt,
			},
			{
				ID:               uuid.Must(uuid.NewV4()),
				EDIType:          models.EDIType824,
				Code:             &ediCode824,
				Description:      &ediDescription824,
				PaymentRequestID: paymentRequestID824,
				CreatedAt:        createdAt,
			},
			{
				ID:               uuid.Must(uuid.NewV4()),
				EDIType:          models.EDIType997,
				Code:             &ediCode997,
				Description:      &ediDescription997,
				PaymentRequestID: paymentRequestID997,
				CreatedAt:        createdAt,
			},
		}

		mockFetcher := &mocks.EDIErrorFetcher{}
		mockFetcher.On("FetchEdiErrors", mock.Anything, mock.Anything).Return(expectedEdiErrors, len(expectedEdiErrors), nil)

		handler := FetchEdiErrorsHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			ediErrorFetcher: mockFetcher,
			NewPagination:   pagination.NewPagination,
		}

		params := edierrorsop.FetchEdiErrorsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/edi-errors"),
		}

		response := handler.Handle(params)
		suite.IsType(&edierrorsop.FetchEdiErrorsOK{}, response)

		okResp := response.(*edierrorsop.FetchEdiErrorsOK)
		suite.Len(okResp.Payload, 3)

		suite.Equal(ediCode858, *okResp.Payload[0].Code)
		suite.Equal(ediDescription858, *okResp.Payload[0].Description)
		suite.Equal(paymentRequestID858.String(), okResp.Payload[0].PaymentRequestID.String())
		suite.Equal(models.EDIType858.String(), *okResp.Payload[0].EdiType)
		suite.NotNil(okResp.Payload[0].CreatedAt)

		suite.Equal(ediCode824, *okResp.Payload[1].Code)
		suite.Equal(ediDescription824, *okResp.Payload[1].Description)
		suite.Equal(paymentRequestID824.String(), okResp.Payload[1].PaymentRequestID.String())
		suite.Equal(models.EDIType824.String(), *okResp.Payload[1].EdiType)
		suite.NotNil(okResp.Payload[1].CreatedAt)

		suite.Equal(ediCode997, *okResp.Payload[2].Code)
		suite.Equal(ediDescription997, *okResp.Payload[2].Description)
		suite.Equal(paymentRequestID997.String(), okResp.Payload[2].PaymentRequestID.String())
		suite.Equal(models.EDIType997.String(), *okResp.Payload[2].EdiType)
		suite.NotNil(okResp.Payload[2].CreatedAt)
	})
}

func (suite *HandlerSuite) TestFetchEdiErrorsHandlerFailure() {
	mockFetcher := &mocks.EDIErrorFetcher{}
	expectedErr := apperror.NewQueryError("payment_requests", errors.New("DB failure"), "Could not find payment requests with EDI_ERROR status")

	mockFetcher.On("FetchEdiErrors", mock.Anything, mock.Anything).Return(models.EdiErrors{}, 0, expectedErr)

	handler := FetchEdiErrorsHandler{
		HandlerConfig:   suite.NewHandlerConfig(),
		ediErrorFetcher: mockFetcher,
		NewPagination:   pagination.NewPagination,
	}

	req := suite.setupAuthenticatedRequest("GET", "/edi-errors")
	params := edierrorsop.FetchEdiErrorsParams{
		HTTPRequest: req,
	}

	response := handler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)
	suite.Equal(http.StatusInternalServerError, errResponse.Code)
	suite.Contains(errResponse.Err.Error(), "Could not find payment requests with EDI_ERROR status")
}

func (suite *HandlerSuite) TestGetEdiErrorHandler() {
	suite.Run("Successfully fetches a single EDI error by ID", func() {
		createdAt := time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC)
		ediCode := "FailureFor858"
		ediDescription := "Description for 858"
		paymentRequestID := uuid.Must(uuid.NewV4())
		ediErrorID := uuid.Must(uuid.NewV4())

		paymentRequestNumber := "1234-4567-8"
		expected := models.EdiError{
			ID:               ediErrorID,
			EDIType:          models.EDIType858,
			Code:             &ediCode,
			Description:      &ediDescription,
			PaymentRequestID: paymentRequestID,
			CreatedAt:        createdAt,
			PaymentRequest: models.PaymentRequest{
				PaymentRequestNumber: paymentRequestNumber,
			},
		}

		mockFetcher := &mocks.EDIErrorFetcher{}
		mockFetcher.On("FetchEdiErrorByID", mock.Anything, ediErrorID).Return(expected, nil)

		handler := GetEdiErrorHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			ediErrorFetcher: mockFetcher,
		}

		req := suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/edi-errors/%s", ediErrorID.String()))
		params := singleedierrorop.GetEdiErrorParams{
			HTTPRequest: req,
			EdiErrorID:  strfmt.UUID(ediErrorID.String()),
		}

		response := handler.Handle(params)

		suite.IsType(&singleedierrorop.GetEdiErrorOK{}, response)
		okResp := response.(*singleedierrorop.GetEdiErrorOK)
		suite.Equal(expected.ID.String(), okResp.Payload.ID.String())
		suite.Equal(*expected.Code, *okResp.Payload.Code)
		suite.Equal(*expected.Description, *okResp.Payload.Description)
		suite.Equal(expected.PaymentRequestID.String(), okResp.Payload.PaymentRequestID.String())
		suite.Equal(paymentRequestNumber, okResp.Payload.PaymentRequestNumber)
		suite.Equal(expected.EDIType.String(), *okResp.Payload.EdiType)
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		missingID := uuid.Must(uuid.NewV4())

		mockFetcher := &mocks.EDIErrorFetcher{}
		mockFetcher.On("FetchEdiErrorByID", mock.Anything, mock.MatchedBy(func(id uuid.UUID) bool {
			return id.String() == missingID.String()
		})).Return(models.EdiError{}, apperror.NewNotFoundError(missingID, "EDIError not found"))

		handler := GetEdiErrorHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			ediErrorFetcher: mockFetcher,
		}

		req := suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/edi-errors/%s", missingID.String()))
		params := singleedierrorop.GetEdiErrorParams{
			HTTPRequest: req,
			EdiErrorID:  strfmt.UUID(missingID.String()),
		}

		response := handler.Handle(params)

		suite.IsType(&handlers.ErrResponse{}, response)
		errResp := response.(*handlers.ErrResponse)
		suite.Contains(errResp.Err.Error(), "EDIError not found")
	})

}
