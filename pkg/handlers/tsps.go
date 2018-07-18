package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	publictspsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/tsps"
)

/*
 * ------------------------------------------
 * The code below is for the INTERNAL REST API.
 * ------------------------------------------
 */

// NO CODE YET!

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

// PublicTspsIndexTSPsHandler returns a list of all the TSPs
type PublicTspsIndexTSPsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicTspsIndexTSPsHandler) Handle(params publictspsop.IndexTSPsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexTSPs has not yet been implemented")
}

// PublicTspsGetTspShipmentsHandler lists all the shipments that belong to a tsp
type PublicTspsGetTspShipmentsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicTspsGetTspShipmentsHandler) Handle(params publictspsop.GetTspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}

// PublicTspsGetTspBlackoutsHandler lists all the shipments that belong to a tsp
type PublicTspsGetTspBlackoutsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicTspsGetTspBlackoutsHandler) Handle(params publictspsop.GetTspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .tspShipments has not yet been implemented")
}
