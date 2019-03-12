package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
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
	if clientCert == nil {
		h.Logger().Info("No client certificate provided")
		return ordersoperations.NewGetOrdersUnauthorized()
	}
	if !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewGetOrdersForbidden()
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

	if orders.Issuer == models.IssuerAirForce {
		if !clientCert.AllowAirForceOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Air Force Orders")
			return ordersoperations.NewGetOrdersForbidden()
		}
	} else if orders.Issuer == models.IssuerArmy {
		if !clientCert.AllowArmyOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Army Orders")
			return ordersoperations.NewGetOrdersForbidden()
		}
	} else if orders.Issuer == models.IssuerCoastGuard {
		if !clientCert.AllowCoastGuardOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Coast Guard Orders")
			return ordersoperations.NewGetOrdersForbidden()
		}
	} else if orders.Issuer == models.IssuerMarineCorps {
		if !clientCert.AllowMarineCorpsOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Marine Corps Orders")
			return ordersoperations.NewGetOrdersForbidden()
		}
	} else if orders.Issuer == models.IssuerNavy {
		if !clientCert.AllowNavyOrdersRead {
			h.Logger().Info("Client certificate is not permitted to read Navy Orders")
			return ordersoperations.NewGetOrdersForbidden()
		}
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return ordersoperations.NewGetOrdersOK().WithPayload(ordersPayload)
}
