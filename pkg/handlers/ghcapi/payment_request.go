package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForPaymentRequestModel(pr models.PaymentRequest) *ghcmessages.PaymentRequest {

	return &ghcmessages.PaymentRequest{
		ID:              *handlers.FmtUUID(pr.ID),
		IsFinal:         &pr.IsFinal,
		RejectionReason: pr.RejectionReason,
	}
}

type ListPaymentRequestsHandler struct {
	handlers.HandlerContext
	services.PaymentRequestListFetcher
}

func (h ListPaymentRequestsHandler) Handle(params paymentrequestop.ListPaymentRequestsParams) middleware.Responder {
	// TODO: add authorizations
	logger := h.LoggerFromRequest(params.HTTPRequest)

	paymentRequests, err := h.FetchPaymentRequestList()
	if err != nil {
		logger.Error("Error listing payment requests err", zap.Error(err))
		return paymentrequestop.NewListPaymentRequestsInternalServerError()
	}

	paymentRequestsList := make(ghcmessages.PaymentRequests, len(*paymentRequests))
	for i, paymentRequest := range *paymentRequests {
		paymentRequestsList[i] = payloadForPaymentRequestModel(paymentRequest)
	}

	return paymentrequestop.NewListPaymentRequestsOK().WithPayload(paymentRequestsList)
}
