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
		logger.Error("Invalid payment request: params Body is nil")
		return paymentrequestop.NewCreatePaymentRequestBadRequest()
	}

	paymentRequest := models.PaymentRequest{
		IsFinal: *params.Body.IsFinal,
	}

	createdPaymentRequest, verrs, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if verrs.HasAny() {
		logger.Error("Error creating payment request verrs", zap.Error(verrs))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	if err != nil {
		logger.Error("Error creating payment request err", zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	logger.Info("Payment Request params",
		zap.Bool("is_final", *params.Body.IsFinal),
		// TODO: add MoveTaskOrder
		// TODO: add ServiceItems
		// TODO add ProofOfService object to log
	)
	returnPayload := payloadForPaymentRequestModel(*createdPaymentRequest)

	return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload)
}
