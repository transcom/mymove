package primeapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
	//TODO: add in checks for params
	//if len(params.Body.MoveTaskOrderID) == 0 {
	//	logger.Error("no MoveTaskID on params")
	//}

	serviceItemIDs := params.Body.ServiceItemIDs
	var popServiceItemIDs []uuid.UUID
	for _, id := range serviceItemIDs {
		// params receives swagger version of UUID and it needs to convert to pop version
		properUUIDTypeID, err := uuid.FromString(id.String())
		if err != nil {
			logger.Error("Error saving payment request", zap.Error(err))
			return paymentrequestop.NewCreatePaymentRequestInternalServerError()
		}

		popServiceItemIDs = append(popServiceItemIDs, properUUIDTypeID)
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

	isFinal := *params.Body.IsFinal

	paymentRequest := models.PaymentRequest{
		ID:              uuid.UUID{},
		IsFinal:         isFinal,
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrderID,
		ServiceItemIDs:  popServiceItemIDs,
		RejectionReason: "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}
	createdPaymentRequest, verrs, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if err != nil || verrs != nil {
		logger.Error("Error saving payment request", zap.Error(verrs))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	serviceItemIDStrings := []string{}
	for _, id := range serviceItemIDs {
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
