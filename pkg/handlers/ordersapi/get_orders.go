package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetOrdersHandler returns Orders by uuid
type GetOrdersHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersHandler) responds to GET /orders/{uuid}
func (h GetOrdersHandler) Handle(params ordersoperations.GetOrdersParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewGetOrdersUnauthorized()
	}

	var err error

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		h.Logger().Info("Not a valid UUID")
		return ordersoperations.NewGetOrdersBadRequest()
	}

	orders, err := models.FetchElectronicOrderByID(h.DB(), id)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewGetOrdersNotFound()
	} else if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by ID")
		return ordersoperations.NewGetOrdersInternalServerError()
	}

	if orders.Issuer == ordersmessages.IssuerAirForce {
		if !clientCert.AllowAirForceOrdersRead {
			return ordersoperations.NewGetOrdersUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerArmy {
		if !clientCert.AllowArmyOrdersRead {
			return ordersoperations.NewGetOrdersUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerCoastGuard {
		if !clientCert.AllowCoastGuardOrdersRead {
			return ordersoperations.NewGetOrdersUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerMarineCorps {
		if !clientCert.AllowMarineCorpsOrdersRead {
			return ordersoperations.NewGetOrdersUnauthorized()
		}
	} else if orders.Issuer == ordersmessages.IssuerNavy {
		if !clientCert.AllowNavyOrdersRead {
			return ordersoperations.NewGetOrdersUnauthorized()
		}
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ordersoperations.NewGetOrdersOK().WithPayload(ordersPayload)
}
