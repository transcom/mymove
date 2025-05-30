// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// SubmitPPMShipmentDocumentationHandlerFunc turns a function with the right signature into a submit p p m shipment documentation handler
type SubmitPPMShipmentDocumentationHandlerFunc func(SubmitPPMShipmentDocumentationParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SubmitPPMShipmentDocumentationHandlerFunc) Handle(params SubmitPPMShipmentDocumentationParams) middleware.Responder {
	return fn(params)
}

// SubmitPPMShipmentDocumentationHandler interface for that can handle valid submit p p m shipment documentation params
type SubmitPPMShipmentDocumentationHandler interface {
	Handle(SubmitPPMShipmentDocumentationParams) middleware.Responder
}

// NewSubmitPPMShipmentDocumentation creates a new http.Handler for the submit p p m shipment documentation operation
func NewSubmitPPMShipmentDocumentation(ctx *middleware.Context, handler SubmitPPMShipmentDocumentationHandler) *SubmitPPMShipmentDocumentation {
	return &SubmitPPMShipmentDocumentation{Context: ctx, Handler: handler}
}

/*
	SubmitPPMShipmentDocumentation swagger:route POST /ppm-shipments/{ppmShipmentId}/submit-ppm-shipment-documentation ppm submitPPMShipmentDocumentation

# Saves signature and routes PPM shipment to service counselor

Routes the PPM shipment to the service
counselor PPM Closeout queue for review.
*/
type SubmitPPMShipmentDocumentation struct {
	Context *middleware.Context
	Handler SubmitPPMShipmentDocumentationHandler
}

func (o *SubmitPPMShipmentDocumentation) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewSubmitPPMShipmentDocumentationParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
