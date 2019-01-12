package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// GetOrdersHandler returns Orders by uuid
type GetOrdersHandler struct {
	handlers.HandlerContext
}

// Handle simply returns a NotImplementedError
func (h GetOrdersHandler) Handle(params ordersoperations.GetOrdersParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewGetOrdersUnauthorized()
	}

	return middleware.NotImplemented("operation .getOrders has not yet been implemented")
}

// IndexOrdersHandler returns a list of Orders matching the provided search parameters
type IndexOrdersHandler struct {
	handlers.HandlerContext
}

// Handle simply returns a NotImplementedError
func (h IndexOrdersHandler) Handle(params ordersoperations.IndexOrdersParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewIndexOrdersUnauthorized()
	}

	return middleware.NotImplemented("operation .indexOrders has not yet been implemented")
}

// PostRevisionHandler adds a Revision to Orders matching the provided search parameters
type PostRevisionHandler struct {
	handlers.HandlerContext
}

// Handle simply returns a NotImplementedError
func (h PostRevisionHandler) Handle(params ordersoperations.PostRevisionParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewPostRevisionUnauthorized()
	}

	return middleware.NotImplemented("operation .postRevision has not yet been implemented")
}

// PostRevisionToOrdersHandler adds a Revision to Orders by uuid
type PostRevisionToOrdersHandler struct {
	handlers.HandlerContext
}

// Handle simply returns a NotImplementedError
func (h PostRevisionToOrdersHandler) Handle(params ordersoperations.PostRevisionToOrdersParams) middleware.Responder {
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewPostRevisionToOrdersUnauthorized()
	}

	return middleware.NotImplemented("operation .postRevisionToOrders has not yet been implemented")
}
