package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// IndexOrdersHandler returns a list of Orders matching the provided search parameters
type IndexOrdersHandler struct {
	handlers.HandlerContext
}

// Handle (IndexOrdersHandler) responds to GET /orders
func (h IndexOrdersHandler) Handle(params ordersoperations.IndexOrdersParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewIndexOrdersUnauthorized()
	}

	return middleware.NotImplemented("operation .indexOrders has not yet been implemented")
}
