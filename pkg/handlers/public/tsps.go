package public

import (
	"github.com/go-openapi/runtime/middleware"

	tspsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/tsps"
	"github.com/transcom/mymove/pkg/handlers/utils"
)

// TspsIndexTSPsHandler returns a list of all the TSPs
type TspsIndexTSPsHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h TspsIndexTSPsHandler) Handle(params tspsop.IndexTSPsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexTSPs has not yet been implemented")
}

// TspsGetTspShipmentsHandler lists all the shipments that belong to a tsp
type TspsGetTspShipmentsHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h TspsGetTspShipmentsHandler) Handle(params tspsop.GetTspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}

// TspsGetTspBlackoutsHandler lists all the shipments that belong to a tsp
type TspsGetTspBlackoutsHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h TspsGetTspBlackoutsHandler) Handle(params tspsop.GetTspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}
