package public

import (
	"github.com/go-openapi/runtime/middleware"

	publictspsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/tsps"
	"github.com/transcom/mymove/pkg/handlers/utils"
)

// PublicTspsIndexTSPsHandler returns a list of all the TSPs
type PublicTspsIndexTSPsHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicTspsIndexTSPsHandler) Handle(params publictspsop.IndexTSPsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexTSPs has not yet been implemented")
}

// PublicTspsGetTspShipmentsHandler lists all the shipments that belong to a tsp
type PublicTspsGetTspShipmentsHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicTspsGetTspShipmentsHandler) Handle(params publictspsop.GetTspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}

// PublicTspsGetTspBlackoutsHandler lists all the shipments that belong to a tsp
type PublicTspsGetTspBlackoutsHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicTspsGetTspBlackoutsHandler) Handle(params publictspsop.GetTspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}
