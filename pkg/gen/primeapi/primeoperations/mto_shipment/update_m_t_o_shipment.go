// Code generated by go-swagger; DO NOT EDIT.

package mto_shipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateMTOShipmentHandlerFunc turns a function with the right signature into a update m t o shipment handler
type UpdateMTOShipmentHandlerFunc func(UpdateMTOShipmentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateMTOShipmentHandlerFunc) Handle(params UpdateMTOShipmentParams) middleware.Responder {
	return fn(params)
}

// UpdateMTOShipmentHandler interface for that can handle valid update m t o shipment params
type UpdateMTOShipmentHandler interface {
	Handle(UpdateMTOShipmentParams) middleware.Responder
}

// NewUpdateMTOShipment creates a new http.Handler for the update m t o shipment operation
func NewUpdateMTOShipment(ctx *middleware.Context, handler UpdateMTOShipmentHandler) *UpdateMTOShipment {
	return &UpdateMTOShipment{Context: ctx, Handler: handler}
}

/*
	UpdateMTOShipment swagger:route PATCH /mto-shipments/{mtoShipmentID} mtoShipment updateMTOShipment

updateMTOShipment

_[Deprecated: this endpoint was deprecated on August 5th, 2024]_
Please use the new endpoint at `/prime/v3/updateMTOShipment` instead.
*/
type UpdateMTOShipment struct {
	Context *middleware.Context
	Handler UpdateMTOShipmentHandler
}

func (o *UpdateMTOShipment) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateMTOShipmentParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}