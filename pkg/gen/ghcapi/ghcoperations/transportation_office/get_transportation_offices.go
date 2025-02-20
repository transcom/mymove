// Code generated by go-swagger; DO NOT EDIT.

package transportation_office

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTransportationOfficesHandlerFunc turns a function with the right signature into a get transportation offices handler
type GetTransportationOfficesHandlerFunc func(GetTransportationOfficesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTransportationOfficesHandlerFunc) Handle(params GetTransportationOfficesParams) middleware.Responder {
	return fn(params)
}

// GetTransportationOfficesHandler interface for that can handle valid get transportation offices params
type GetTransportationOfficesHandler interface {
	Handle(GetTransportationOfficesParams) middleware.Responder
}

// NewGetTransportationOffices creates a new http.Handler for the get transportation offices operation
func NewGetTransportationOffices(ctx *middleware.Context, handler GetTransportationOfficesHandler) *GetTransportationOffices {
	return &GetTransportationOffices{Context: ctx, Handler: handler}
}

/*
	GetTransportationOffices swagger:route GET /transportation-offices transportationOffice getTransportationOffices

# Returns the transportation offices matching the search query that is enabled for PPM closeout

Returns the transportation offices matching the search query that is enabled for PPM closeout
*/
type GetTransportationOffices struct {
	Context *middleware.Context
	Handler GetTransportationOfficesHandler
}

func (o *GetTransportationOffices) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTransportationOfficesParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
