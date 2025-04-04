// Code generated by go-swagger; DO NOT EDIT.

package mto_shipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateMTOAgentHandlerFunc turns a function with the right signature into a update m t o agent handler
type UpdateMTOAgentHandlerFunc func(UpdateMTOAgentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateMTOAgentHandlerFunc) Handle(params UpdateMTOAgentParams) middleware.Responder {
	return fn(params)
}

// UpdateMTOAgentHandler interface for that can handle valid update m t o agent params
type UpdateMTOAgentHandler interface {
	Handle(UpdateMTOAgentParams) middleware.Responder
}

// NewUpdateMTOAgent creates a new http.Handler for the update m t o agent operation
func NewUpdateMTOAgent(ctx *middleware.Context, handler UpdateMTOAgentHandler) *UpdateMTOAgent {
	return &UpdateMTOAgent{Context: ctx, Handler: handler}
}

/*
	UpdateMTOAgent swagger:route PUT /mto-shipments/{mtoShipmentID}/agents/{agentID} mtoShipment updateMTOAgent

updateMTOAgent

### Functionality
This endpoint is used to **update** the agents for an MTO Shipment. Only the fields being modified need to be sent in the request body.

### Errors:
The agent must always have a name and at least one method of contact (either `email` or `phone`).

The agent must be associated with the MTO shipment passed in the url.

The shipment should be associated with an MTO that is available to the Prime.
If the caller requests an update to an agent, and the shipment is not on an available MTO, the caller will receive a **NotFound** response.
*/
type UpdateMTOAgent struct {
	Context *middleware.Context
	Handler UpdateMTOAgentHandler
}

func (o *UpdateMTOAgent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateMTOAgentParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
