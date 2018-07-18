package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	publictspsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/tsps"
)

// TSPIndexHandler returns a list of all the TSPs
type TSPIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPIndexHandler) Handle(params publictspsop.IndexTSPsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexTSPs has not yet been implemented")
}

// TSPShipmentsHandler lists all the shipments that belong to a tsp
type TSPShipmentsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPShipmentsHandler) Handle(params publictspsop.TspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}

// TSPBlackoutsHandler lists all the shipments that belong to a tsp
type TSPBlackoutsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPBlackoutsHandler) Handle(params publictspsop.TspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}
