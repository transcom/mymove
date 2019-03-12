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

	var err error

	orders, err := models.FetchElectronicOrderByIssuerAndOrdersNum(h.DB(), params.Issuer, params.OrdersNum)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumNotFound()
	} else if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by Issuer and Orders Num")
		return ordersoperations.NewGetOrdersByIssuerAndOrdersNumInternalServerError()
	}

	if orders.Issuer == models.IssuerAirForce {
		if !clientCert.AllowAirForceOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Air Force Orders")
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
		}
	} else if orders.Issuer == models.IssuerArmy {
		if !clientCert.AllowArmyOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Army Orders")
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
		}
	} else if orders.Issuer == models.IssuerCoastGuard {
		if !clientCert.AllowCoastGuardOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Coast Guard Orders")
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
		}
	} else if orders.Issuer == models.IssuerMarineCorps {
		if !clientCert.AllowMarineCorpsOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Marine Corps Orders")
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
		}
	} else if orders.Issuer == models.IssuerNavy {
		if !clientCert.AllowNavyOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Navy Orders")
			return ordersoperations.NewGetOrdersByIssuerAndOrdersNumForbidden()
		}
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ordersoperations.NewGetOrdersByIssuerAndOrdersNumOK().WithPayload(ordersPayload)
}
