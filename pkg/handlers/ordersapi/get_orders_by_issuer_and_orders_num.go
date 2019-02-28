package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
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
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
	}

	var err error

	orders, err := models.FetchElectronicOrderByIssuerAndOrdersNum(h.DB(), params.Issuer, params.OrdersNum)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumNotFound()
	} else if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by Issuer and Orders Num")
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumInternalServerError()
	}

	if orders.Issuer == ordersmessages.IssuerAirForce {
		if !clientCert.AllowAirForceOrdersRead {
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerArmy {
		if !clientCert.AllowArmyOrdersRead {
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerCoastGuard {
		if !clientCert.AllowCoastGuardOrdersRead {
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerMarineCorps {
		if !clientCert.AllowMarineCorpsOrdersRead {
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerNavy {
		if !clientCert.AllowNavyOrdersRead {
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumUnauthorized()
		}
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ordersoperations.NewGetOrdersByIssuerAndOrdersNumOK().WithPayload(ordersPayload)
}
