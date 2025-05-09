// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetPPMActualWeightHandlerFunc turns a function with the right signature into a get p p m actual weight handler
type GetPPMActualWeightHandlerFunc func(GetPPMActualWeightParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPPMActualWeightHandlerFunc) Handle(params GetPPMActualWeightParams) middleware.Responder {
	return fn(params)
}

// GetPPMActualWeightHandler interface for that can handle valid get p p m actual weight params
type GetPPMActualWeightHandler interface {
	Handle(GetPPMActualWeightParams) middleware.Responder
}

// NewGetPPMActualWeight creates a new http.Handler for the get p p m actual weight operation
func NewGetPPMActualWeight(ctx *middleware.Context, handler GetPPMActualWeightHandler) *GetPPMActualWeight {
	return &GetPPMActualWeight{Context: ctx, Handler: handler}
}

/*
	GetPPMActualWeight swagger:route GET /ppm-shipments/{ppmShipmentId}/actual-weight ppm getPPMActualWeight

# Get the actual weight for a PPM shipment

Retrieves the actual weight for the specified PPM shipment.
*/
type GetPPMActualWeight struct {
	Context *middleware.Context
	Handler GetPPMActualWeightHandler
}

func (o *GetPPMActualWeight) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetPPMActualWeightParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
