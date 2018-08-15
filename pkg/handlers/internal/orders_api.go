package internal

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers/utils"
)

// GetOrdersHandler returns Orders by uuid
type GetOrdersHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h GetOrdersHandler) Handle(params ordersoperations.GetOrdersParams) middleware.Responder {
	return middleware.NotImplemented("operation .getOrders has not yet been implemented")
}

// IndexOrdersHandler returns a list of Orders matching the provided search parameters
type IndexOrdersHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h IndexOrdersHandler) Handle(params ordersoperations.IndexOrdersParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexOrders has not yet been implemented")
}

// PostRevisionHandler adds a Revision to Orders matching the provided search parameters
type PostRevisionHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h PostRevisionHandler) Handle(params ordersoperations.PostRevisionParams) middleware.Responder {
	return middleware.NotImplemented("operation .postRevision has not yet been implemented")
}

// PostRevisionToOrdersHandler adds a Revision to Orders by uuid
type PostRevisionToOrdersHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h PostRevisionToOrdersHandler) Handle(params ordersoperations.PostRevisionToOrdersParams) middleware.Responder {
	return middleware.NotImplemented("operation .postRevisionToOrders has not yet been implemented")
}
