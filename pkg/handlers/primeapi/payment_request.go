package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForPaymentRequestModel(pr models.PaymentRequest) *primemessages.PaymentRequest {
	var serviceItemIDs = []strfmt.UUID{}

	for _, id := range pr.ServiceItemIDs {
		serviceItemIDs = append(serviceItemIDs, *handlers.FmtUUID(id))
	}

	return &primemessages.PaymentRequest{
		ID: *handlers.FmtUUID(pr.ID),
		MoveTaskOrderID: *handlers.FmtUUID(pr.MoveTaskOrderID),
		ServiceItemIDs: serviceItemIDs,
		IsFinal: &pr.IsFinal,
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

	var serviceItemIDs = []uuid.UUID{}
	for _, id := range params.Body.ServiceItemIDs {
		serviceItemIDs = append(serviceItemIDs, handlers.FmtToPopUUID(id))
	}
	paymentRequest := models.PaymentRequest{
		IsFinal:         *params.Body.IsFinal,
		MoveTaskOrderID: handlers.FmtToPopUUID(params.Body.MoveTaskOrderID),
		ServiceItemIDs:  serviceItemIDs,
	}

	moveTaskOrder := models.MoveTaskOrder{}
	err := h.DB().Find(&moveTaskOrder, paymentRequest.MoveTaskOrderID)
	if err != nil {
		logger.Error("Error saving payment request", zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestNotFound()
	}

	createdPaymentRequest, verrs, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if err != nil || verrs != nil {
		logger.Error("Error saving payment request", zap.Error(verrs))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	serviceItemIDStrings := []string{}
	for _, id := range params.Body.ServiceItemIDs {
		idString := id.String()
		serviceItemIDStrings = append(serviceItemIDStrings, idString)
	}
	logger.Info("Payment Request params",
		zap.Bool("is_final", *params.Body.IsFinal),
		zap.String("move_order_id", (params.Body.MoveTaskOrderID).String()),
		zap.Strings("service_item_ids", serviceItemIDStrings),
		// TODO add ProofOfService object to log
	)
	returnPayload := payloadForPaymentRequestModel(*createdPaymentRequest)

	return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload)
}
