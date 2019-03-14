package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"

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
		h.Logger().Info("No client certificate provided")
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
	}
	if !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
	}
	if !verifyOrdersReadAccess(models.Issuer(params.Issuer), clientCert, h.Logger(), true) {
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
	}

	var err error

	orders, err := models.FetchElectronicOrderByIssuerAndOrdersNum(h.DB(), params.Issuer, params.OrdersNum)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumNotFound()
	} else if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by Issuer and Orders Num")
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumInternalServerError()
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ordersoperations.NewGetOrdersByIssuerAndOrdersNumOK().WithPayload(ordersPayload)
}
