package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"go.uber.org/zap"
)

type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
}

func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) (m middleware.Responder) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var serviceItemIDStrings []string

	if params.HTTPRequest.ContentLength == 0 {
	return nil // TODO figure out how to return 501
	}

	serviceItemIDs := params.Body.ServiceItemIDs
	for _, id := range serviceItemIDs {
		serviceItemIDStrings = append(serviceItemIDStrings, id.String())
	}

	logger.Info("Payment Request params",
		zap.Bool("is_final", *params.Body.IsFinal),
		zap.String("move_order_id", (params.Body.MoveOrderID).String()),
		zap.Strings("service_item_ids", serviceItemIDStrings),
	)

	//responseWriter := http.ResponseWriter(http.ResponseWriter.WriteHeader(http.StatusNotImplemented), runtime.Producer())
	//	return middleware.Responder.WriteResponse(
	//		responseWriter,
	//		runtime.Producer())
	//return http.StatusNotImplemented
	return nil // TODO figure out how to return 501
}
