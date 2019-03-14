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
	allowedIssuers := clientCert.GetAllowedOrdersIssuersRead()
	if len(allowedIssuers) == 0 {
		h.Logger().Info("Client certificate is not permitted to read any Orders")
		return ordersoperations.NewIndexOrdersForMemberForbidden()
	}

	orders, err := models.FetchElectronicOrdersByEdipiAndIssuers(h.DB(), params.Edipi, allowedIssuers)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewIndexOrdersForMemberOK().WithPayload([]*ordersmessages.Orders{})
	} else if err != nil {
		h.Logger().Info("Error while fetching electronic Orders by EDIPI")
		return ordersoperations.NewIndexOrdersForMemberInternalServerError()
	}
	ordersPayloads := make([]*ordersmessages.Orders, len(orders))
	for i, o := range orders {
		payload, err := payloadForElectronicOrderModel(o)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
		ordersPayloads[i] = payload
	}

	return ordersoperations.NewIndexOrdersForMemberOK().WithPayload(ordersPayloads)
}
