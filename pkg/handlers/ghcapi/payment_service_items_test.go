package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	paymentServiceItemOp "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_service_item"
)

type updatePaymentSubtestData struct {
	paymentServiceItem models.PaymentServiceItem
	params             paymentServiceItemOp.UpdatePaymentServiceItemStatusParams
	input              paymentServiceItemOp.UpdatePaymentServiceItemStatusParams
}

func (suite *HandlerSuite) makeUpdatePaymentSubtestData() (subtestData *updatePaymentSubtestData) {
	subtestData = &updatePaymentSubtestData{}

	mto := testdatagen.MakeDefaultMove(suite.DB())
	paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItem(suite.DB())
	subtestData.paymentServiceItem = paymentServiceItem
	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move-task-orders/%s/payment-service-items/%s/status", mto.ID.String(), paymentServiceItem.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	subtestData.params = paymentServiceItemOp.UpdatePaymentServiceItemStatusParams{
		HTTPRequest:          req,
		IfMatch:              etag.GenerateEtag(paymentServiceItem.UpdatedAt),
		PaymentServiceItemID: paymentServiceItem.ID.String(),
		Body: &ghcmessages.PaymentServiceItem{
			ETag:                     etag.GenerateEtag(paymentServiceItem.UpdatedAt),
			ID:                       *handlers.FmtUUID(paymentServiceItem.ID),
			MtoServiceItemID:         *handlers.FmtUUID(paymentServiceItem.MTOServiceItemID),
			PaymentRequestID:         *handlers.FmtUUID(paymentServiceItem.PaymentRequestID),
			PaymentServiceItemParams: nil,
			PriceCents:               nil,
			RejectionReason:          nil,
			Status:                   ghcmessages.PaymentServiceItemStatusAPPROVED,
		},
	}

	// for 422 and 500 errors
	psi := testdatagen.MakeDefaultPaymentServiceItem(suite.DB())
	subtestData.input = paymentServiceItemOp.UpdatePaymentServiceItemStatusParams{
		HTTPRequest:          req,
		IfMatch:              etag.GenerateEtag(psi.UpdatedAt),
		PaymentServiceItemID: psi.ID.String(),
		Body: &ghcmessages.PaymentServiceItem{
			ETag:                     etag.GenerateEtag(psi.UpdatedAt),
			ID:                       *handlers.FmtUUID(psi.ID),
			MtoServiceItemID:         *handlers.FmtUUID(psi.MTOServiceItemID),
			PaymentRequestID:         *handlers.FmtUUID(psi.PaymentRequestID),
			PaymentServiceItemParams: nil,
			PriceCents:               nil,
			RejectionReason:          nil,
			Status:                   ghcmessages.PaymentServiceItemStatusAPPROVED,
		},
	}

	return subtestData
}

func (suite *HandlerSuite) TestUpdatePaymentServiceItemHandler() {
	suite.Run("Successful patch - Approval - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}

		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(subtestData.paymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
	})

	suite.Run("404 - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}
		subtestData.params.PaymentServiceItemID = uuid.Nil.String()

		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusNotFound{}, response)

	})

	suite.Run("422 - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		newParam := subtestData.input
		newParam.Body.Status = ""

		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}

		response := handler.Handle(newParam)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusUnprocessableEntity{}, response)
	})

	suite.Run("500 - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		reason := "More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters More than 255 characters "
		newParam := subtestData.input
		newParam.Body.Status = ghcmessages.PaymentServiceItemStatusDENIED
		newParam.Body.RejectionReason = &reason

		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}

		response := handler.Handle(newParam)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusInternalServerError{}, response)
	})

	suite.Run("Successful patch - Rejection - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItem(suite.DB())
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}
		subtestData.params.IfMatch = etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		subtestData.params.PaymentServiceItemID = paymentServiceItem.ID.String()
		subtestData.params.Body.Status = ghcmessages.PaymentServiceItemStatusDENIED
		subtestData.params.Body.RejectionReason = swag.String("Because reasons")

		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(paymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusDENIED, okResponse.Payload.Status)
		suite.Equal("Because reasons", *okResponse.Payload.RejectionReason)
	})

	suite.Run("Successful patch - Approval of previously rejected - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		deniedPaymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusDenied,
			},
		})
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}
		subtestData.params.IfMatch = etag.GenerateEtag(deniedPaymentServiceItem.UpdatedAt)
		subtestData.params.PaymentServiceItemID = deniedPaymentServiceItem.ID.String()
		subtestData.params.Body.Status = ghcmessages.PaymentServiceItemStatusAPPROVED

		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(deniedPaymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
		suite.Nil(okResponse.Payload.RejectionReason)
	})

	suite.Run("Successful patch - Approval of Prime available paymentServiceItem", func() {
		availableMTO := testdatagen.MakeAvailableMove(suite.DB())
		availablePaymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			Move: availableMTO,
		})
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move-task-orders/%s/payment-service-items/%s/status", availableMTO.ID.String(), availablePaymentServiceItem.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := paymentServiceItemOp.UpdatePaymentServiceItemStatusParams{
			HTTPRequest:          req,
			IfMatch:              etag.GenerateEtag(availablePaymentServiceItem.UpdatedAt),
			PaymentServiceItemID: availablePaymentServiceItem.ID.String(),
			Body: &ghcmessages.PaymentServiceItem{
				ID:              *handlers.FmtUUID(availablePaymentServiceItem.ID),
				RejectionReason: nil,
				Status:          ghcmessages.PaymentServiceItemStatusAPPROVED,
			},
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		handler.SetTraceID(traceID)

		response := handler.Handle(params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(availablePaymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
		suite.Nil(okResponse.Payload.RejectionReason)
		suite.HasWebhookNotification(availablePaymentServiceItem.PaymentRequestID, traceID)
	})
}
