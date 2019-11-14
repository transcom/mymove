package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"

)

type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.ServiceItemCreator
	services.NewQueryFilter
}

func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) (m middleware.Responder) {

	return m
}