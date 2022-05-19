package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	paymentServiceItemOp "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	paymentServiceItemService "github.com/transcom/mymove/pkg/services/payment_service_item"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trace"
)

type updatePaymentSubtestData struct {
	paymentServiceItem models.PaymentServiceItem
	params             paymentServiceItemOp.UpdatePaymentServiceItemStatusParams
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

	return subtestData
}

func (suite *HandlerSuite) TestUpdatePaymentServiceItemHandler() {
	suite.Run("Successful patch - Approval - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerConfig:                   handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentServiceItemStatusUpdater: paymentServiceItemService.NewPaymentServiceItemStatusUpdater(),
		}
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(subtestData.paymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
	})

	suite.Run("404 - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerConfig:                   handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentServiceItemStatusUpdater: paymentServiceItemService.NewPaymentServiceItemStatusUpdater(),
		}
		subtestData.params.PaymentServiceItemID = uuid.Nil.String()

		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusNotFound{}, response)

	})

	suite.Run("422 - Fails to reject without rejectionReason - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerConfig:                   handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentServiceItemStatusUpdater: paymentServiceItemService.NewPaymentServiceItemStatusUpdater(),
		}

		subtestData.params.Body.Status = ghcmessages.PaymentServiceItemStatusDENIED
		subtestData.params.Body.RejectionReason = nil
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := handler.Handle(subtestData.params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusUnprocessableEntity{}, response)

	})

	suite.Run("Successful patch - Rejection - Integration Test", func() {
		subtestData := suite.makeUpdatePaymentSubtestData()
		paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItem(suite.DB())

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerConfig:                   handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentServiceItemStatusUpdater: paymentServiceItemService.NewPaymentServiceItemStatusUpdater(),
		}
		subtestData.params.IfMatch = etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		subtestData.params.PaymentServiceItemID = paymentServiceItem.ID.String()
		subtestData.params.Body.Status = ghcmessages.PaymentServiceItemStatusDENIED
		subtestData.params.Body.RejectionReason = swag.String("Because reasons")

		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
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

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerConfig:                   handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentServiceItemStatusUpdater: paymentServiceItemService.NewPaymentServiceItemStatusUpdater(),
		}
		subtestData.params.IfMatch = etag.GenerateEtag(deniedPaymentServiceItem.UpdatedAt)
		subtestData.params.PaymentServiceItemID = deniedPaymentServiceItem.ID.String()
		subtestData.params.Body.Status = ghcmessages.PaymentServiceItemStatusAPPROVED

		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
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

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

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

		handler := UpdatePaymentServiceItemStatusHandler{
			HandlerConfig:                   handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			PaymentServiceItemStatusUpdater: paymentServiceItemService.NewPaymentServiceItemStatusUpdater(),
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&paymentServiceItemOp.UpdatePaymentServiceItemStatusOK{}, response)
		okResponse := response.(*paymentServiceItemOp.UpdatePaymentServiceItemStatusOK)
		suite.Equal(availablePaymentServiceItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(ghcmessages.PaymentServiceItemStatusAPPROVED, okResponse.Payload.Status)
		suite.Nil(okResponse.Payload.RejectionReason)
		suite.HasWebhookNotification(availablePaymentServiceItem.PaymentRequestID, traceID)
	})
}
