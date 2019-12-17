package ghcapi

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/query"

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

type GetPaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
}

func (h GetPaymentRequestHandler) Handle(params paymentrequestop.GetPaymentRequestParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestID.String())}

	paymentRequest, err := h.FetchPaymentRequest(filters)

	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching Payment Request with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	if reflect.DeepEqual(paymentRequest, models.PaymentRequest{}) {
		logger.Info(fmt.Sprintf("Could not find a Payment Request with ID: %s", params.PaymentRequestID.String()))
		return paymentrequestop.NewGetPaymentRequestNotFound()
	}

	returnPayload := payloadForPaymentRequestModel(paymentRequest)
	response := paymentrequestop.NewGetPaymentRequestOK().WithPayload(returnPayload)

	return response
}
