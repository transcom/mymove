package primeapi

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"go.uber.org/zap"
	"time"
)

func payloadForPaymentRequestModel(s models.PaymentRequest) *primemessages.PaymentRequest {
	return &primemessages.PaymentRequest{
		ID: *handlers.FmtUUID(s.ID),
	}
}

type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestCreator
}

func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) (m middleware.Responder) {
	// TODO: authorization to create payment request

	logger := h.LoggerFromRequest(params.HTTPRequest)

	var serviceItemIDStrings []string
	serviceItemIDs := params.Body.ServiceItemIDs
	for _, id := range serviceItemIDs {
		serviceItemIDStrings = append(serviceItemIDStrings, id.String())
	}

	moveTaskOrderID, err := uuid.FromString(params.Body.MoveTaskOrderID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.Body.MoveTaskOrderID), zap.Error(err))
	}

	moveTaskOrder := models.MoveTaskOrder{}
	err = h.DB().Find(&moveTaskOrder, moveTaskOrderID)
	if err != nil {
		logger.Error("Error saving payment request", zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestNotFound()
	}

	paymentRequest := models.PaymentRequest{
		ID:              uuid.UUID{},
		IsFinal:         false,
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrderID,
		ServiceItemIDs:  nil,
		RejectionReason: "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}
	createdPaymentRequest, verrs, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if err != nil || verrs != nil {
		logger.Error("Error saving payment request", zap.Error(verrs))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
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
