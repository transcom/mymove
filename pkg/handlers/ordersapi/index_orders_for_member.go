package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// IndexOrdersForMemberHandler returns a list of Orders matching the provided search parameters
type IndexOrdersForMemberHandler struct {
	handlers.HandlerContext
}

// Handle (IndexOrdersForMemberHandler) responds to GET /edipis/{edipi}/orders
func (h IndexOrdersForMemberHandler) Handle(params ordersoperations.IndexOrdersForMemberParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil {
		h.Logger().Info("No client certificate provided")
		return ordersoperations.NewIndexOrdersForMemberUnauthorized()
	}
	if !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewIndexOrdersForMemberForbidden()
	}

	var err error

	orders, err := models.FetchElectronicOrdersByEdipi(h.DB(), params.Edipi)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewIndexOrdersForMemberOK().WithPayload([]*ordersmessages.Orders{})
	} else if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by EDIPI")
		return ordersoperations.NewIndexOrdersForMemberInternalServerError()
	}

	var ordersPayloads []*ordersmessages.Orders
	for _, o := range orders {
		// only return orders that the client is permitted to see
		if o.Issuer == ordersmessages.IssuerAirForce {
			if !clientCert.AllowAirForceOrdersRead {
				continue
			}
		} else if o.Issuer == ordersmessages.IssuerArmy {
			if !clientCert.AllowArmyOrdersRead {
				continue
			}
		} else if o.Issuer == ordersmessages.IssuerCoastGuard {
			if !clientCert.AllowCoastGuardOrdersRead {
				continue
			}
		} else if o.Issuer == ordersmessages.IssuerMarineCorps {
			if !clientCert.AllowMarineCorpsOrdersRead {
				continue
			}
		} else if o.Issuer == ordersmessages.IssuerNavy {
			if !clientCert.AllowNavyOrdersRead {
				continue
			}
		}

		ordersPayload, err := payloadForElectronicOrderModel(o)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
		ordersPayloads = append(ordersPayloads, ordersPayload)
	}

	return ordersoperations.NewIndexOrdersForMemberOK().WithPayload(ordersPayloads)
}
