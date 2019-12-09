package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
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

type ListPaymentRequestsHandler struct {
	handlers.HandlerContext
	services.PaymentRequestLister
}

func (h ListPaymentRequestsHandler) Handle(params payment_requests.ListPaymentRequestsParams) (m middleware.Responder) {
	paymentRequests, _, _ := h.PaymentRequestLister.ListPaymentRequests()

	var returnPayload []primemessages.PaymentRequest
	for _, paymentRequest := range *paymentRequests {
		returnPayload = append(returnPayload, *payloadForPaymentRequestModel(paymentRequest))
	}

	return payment_requests.NewListPaymentRequestsOK().WithPayload(returnPayload)
}