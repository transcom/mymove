// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ShowAOAPacketHandlerFunc turns a function with the right signature into a show a o a packet handler
type ShowAOAPacketHandlerFunc func(ShowAOAPacketParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ShowAOAPacketHandlerFunc) Handle(params ShowAOAPacketParams) middleware.Responder {
	return fn(params)
}

// ShowAOAPacketHandler interface for that can handle valid show a o a packet params
type ShowAOAPacketHandler interface {
	Handle(ShowAOAPacketParams) middleware.Responder
}

// NewShowAOAPacket creates a new http.Handler for the show a o a packet operation
func NewShowAOAPacket(ctx *middleware.Context, handler ShowAOAPacketHandler) *ShowAOAPacket {
	return &ShowAOAPacket{Context: ctx, Handler: handler}
}

/*
	ShowAOAPacket swagger:route GET /ppm-shipments/{ppmShipmentId}/aoa-packet ppm showAOAPacket

# Downloads AOA Packet form PPMShipment as a PDF

### Functionality
This endpoint downloads all uploaded move order documentation combined with the Shipment Summary Worksheet into a single PDF.
### Errors
* The PPMShipment must have requested an AOA.
* The PPMShipment AOA Request must have been approved.
*/
type ShowAOAPacket struct {
	Context *middleware.Context
	Handler ShowAOAPacketHandler
}

func (o *ShowAOAPacket) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewShowAOAPacketParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
