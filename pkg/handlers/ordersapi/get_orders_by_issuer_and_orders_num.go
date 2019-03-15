package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetOrdersByIssuerAndOrdersNumHandler returns Orders for a specific issuer by ordersNum
type GetOrdersByIssuerAndOrdersNumHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersByIssuerAndOrdersNumHandler) responds to GET /issuers/{issuer}/orders/{ordersNum}
func (h GetOrdersByIssuerAndOrdersNumHandler) Handle(params ordersoperations.GetOrdersByIssuerAndOrdersNumParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}
	if !verifyOrdersReadAccess(models.Issuer(params.Issuer), clientCert) {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read orders from this issuer"))
	}

	orders, err := models.FetchElectronicOrderByIssuerAndOrdersNum(h.DB(), params.Issuer, params.OrdersNum)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ordersoperations.NewGetOrdersByIssuerAndOrdersNumOK().WithPayload(ordersPayload)
}
