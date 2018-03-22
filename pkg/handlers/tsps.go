package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
)

// TSPIndexHandler returns a list of all the TSPs
type TSPIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPIndexHandler) Handle(params apioperations.IndexTSPsParams) middleware.Responder {
	return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
}

// TSPShipmentsHandler lists all the shipments that belong to a tsp
type TSPShipmentsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPShipmentsHandler) Handle(params apioperations.TspShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
}

// TSPAcceptShipmentsHandler accepts a shipment offered to a TSP
type TSPAcceptShipmentsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPAcceptShipmentsHandler) Handle(params apioperations.AcceptShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
}

// TSPRejectShipmentsHandler accepts a shipment offered to a TSP
type TSPRejectShipmentsHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h TSPRejectShipmentsHandler) Handle(params apioperations.RejectShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
}
