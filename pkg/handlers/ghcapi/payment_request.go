package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
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

type ShowPaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
}

func (h ShowPaymentRequestHandler) Handle(params paymentrequestop.GetPaymentRequestParams) middleware.Responder {
	paymentRequestID, _ := uuid.FromString(params.PaymentRequestID.String())
	paymentRequest, _, _ := h.FetchPaymentRequest(paymentRequestID)

	returnPayload := payloadForPaymentRequestModel(*paymentRequest)
	response := paymentrequestop.NewGetPaymentRequestOK().WithPayload(returnPayload)

	return response
}
