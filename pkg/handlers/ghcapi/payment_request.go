package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

type ListPaymentRequestsHandler struct {
	handlers.HandlerContext
	services.PaymentRequestListFetcher
}

func (h ListPaymentRequestsHandler) Handle(params paymentrequestop.ListPaymentRequestsParams) (middleware.Responder) {
	// TODO: add authorizations
	logger := h.LoggerFromRequest(params.HTTPRequest)

	paymentRequests, err := h.FetchPaymentRequestList()
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	response := paymentrequestop.NewListPaymentRequestsOK().WithPayload(paymentRequests)

	return response
}
