package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	shipmentops "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	tspops "github.com/transcom/mymove/pkg/gen/restapi/apioperations/tsps"
)

// TSPIndexHandler returns a list of all the TSPs
type TSPIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPIndexHandler) Handle(params tspops.IndexTSPsParams) middleware.Responder {
	return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
}

// TSPShipmentsHandler lists all the shipments that belong to a tsp
type TSPShipmentsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPShipmentsHandler) Handle(params shipmentops.TspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
}
