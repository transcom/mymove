// Code generated by go-swagger; DO NOT EDIT.

package move_task_order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// AcknowledgeMovesAndShipmentsHandlerFunc turns a function with the right signature into a acknowledge moves and shipments handler
type AcknowledgeMovesAndShipmentsHandlerFunc func(AcknowledgeMovesAndShipmentsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn AcknowledgeMovesAndShipmentsHandlerFunc) Handle(params AcknowledgeMovesAndShipmentsParams) middleware.Responder {
	return fn(params)
}

// AcknowledgeMovesAndShipmentsHandler interface for that can handle valid acknowledge moves and shipments params
type AcknowledgeMovesAndShipmentsHandler interface {
	Handle(AcknowledgeMovesAndShipmentsParams) middleware.Responder
}

// NewAcknowledgeMovesAndShipments creates a new http.Handler for the acknowledge moves and shipments operation
func NewAcknowledgeMovesAndShipments(ctx *middleware.Context, handler AcknowledgeMovesAndShipmentsHandler) *AcknowledgeMovesAndShipments {
	return &AcknowledgeMovesAndShipments{Context: ctx, Handler: handler}
}

/*
	AcknowledgeMovesAndShipments swagger:route PATCH /move-task-orders/acknowledge moveTaskOrder acknowledgeMovesAndShipments

acknowledgeMovesAndShipments

### Functionality
This endpoint **updates** the Moves and Shipments to indicate that the Prime has acknowledged the Moves and Shipments have been received.
The Move and Shipment data is expected to be sent in the request body.
*/
type AcknowledgeMovesAndShipments struct {
	Context *middleware.Context
	Handler AcknowledgeMovesAndShipmentsHandler
}

func (o *AcknowledgeMovesAndShipments) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewAcknowledgeMovesAndShipmentsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
