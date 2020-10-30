package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

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

func (suite *HandlerSuite) TestUpdatePaymentServiceItemHandler() {
	mto := testdatagen.MakeDefaultMove(suite.DB())
	paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItem(suite.DB())
	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move-task-orders/%s/payment-service-items/%s/status", mto.ID.String(), paymentServiceItem.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := paymentServiceItemOp.UpdatePaymentServiceItemStatusParams{
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

	suite.T().Run("Successful patch - Approval - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}

		response := handler.Handle(params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(paymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
	})

	suite.T().Run("404 - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}
		params.PaymentServiceItemID = uuid.Nil.String()

		response := handler.Handle(params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusNotFound{}, response)

	})

	suite.T().Run("Successful patch - Rejection - Integration Test", func(t *testing.T) {
		paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItem(suite.DB())
		queryBuilder := query.NewQueryBuilder(suite.DB())

		fetcher := fetch.NewFetcher(queryBuilder)

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        fetcher,
			Builder:        *queryBuilder,
		}
		params.IfMatch = etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		params.PaymentServiceItemID = paymentServiceItem.ID.String()
		params.Body.Status = ghcmessages.PaymentServiceItemStatusDENIED
		params.Body.RejectionReason = swag.String("Because reasons")

		response := handler.Handle(params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(paymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusDENIED, okResponse.Payload.Status)
		suite.Equal("Because reasons", *okResponse.Payload.RejectionReason)
	})

	suite.T().Run("Successful patch - Approval of previously rejected - Integration Test", func(t *testing.T) {
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
		params.IfMatch = etag.GenerateEtag(deniedPaymentServiceItem.UpdatedAt)
		params.PaymentServiceItemID = deniedPaymentServiceItem.ID.String()
		params.Body.Status = ghcmessages.PaymentServiceItemStatusAPPROVED

		response := handler.Handle(params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(deniedPaymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
		suite.Nil(okResponse.Payload.RejectionReason)
	})

	suite.T().Run("Successful patch - Approval of Prime available paymentServiceItem", func(t *testing.T) {
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
