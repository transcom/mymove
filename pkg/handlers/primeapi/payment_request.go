package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForPaymentRequestModel(pr models.PaymentRequest) *primemessages.PaymentRequest {

	return &primemessages.PaymentRequest{
		ID:              *handlers.FmtUUID(pr.ID),
		MoveTaskOrderID: *handlers.FmtUUID(pr.MoveTaskOrderID),
		IsFinal:         &pr.IsFinal,
		RejectionReason: &pr.RejectionReason,
	}
}

type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestCreator
}

func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) (m middleware.Responder) {
	// TODO: authorization to create payment request

	logger := h.LoggerFromRequest(params.HTTPRequest)

	if params.Body == nil {
		logger.Error("Error saving payment request")
		return paymentrequestop.NewCreatePaymentRequestBadRequest()
	}
	paymentRequest := models.PaymentRequest{
		IsFinal:         *params.Body.IsFinal,
		MoveTaskOrderID: handlers.FmtToPopUUID(params.Body.MoveTaskOrderID),
	}

	moveTaskOrder := models.MoveTaskOrder{}
	err := h.DB().Find(&moveTaskOrder, paymentRequest.MoveTaskOrderID)
	if err != nil {
		logger.Error("Error saving payment request", zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestNotFound()
	}

	paymentRequest.MoveTaskOrder = moveTaskOrder

	createdPaymentRequest, verrs, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if err != nil || verrs != nil {
		logger.Error("Error saving payment request", zap.Error(verrs))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	logger.Info("Payment Request params",
		zap.Bool("is_final", *params.Body.IsFinal),
		zap.String("move_task_order_id", (params.Body.MoveTaskOrderID).String()),
		// TODO: add ServiceItems
		// TODO add ProofOfService object to log
	)
	returnPayload := payloadForPaymentRequestModel(*createdPaymentRequest)

	return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload)
}
