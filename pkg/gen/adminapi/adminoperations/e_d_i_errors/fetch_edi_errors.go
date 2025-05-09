// Code generated by go-swagger; DO NOT EDIT.

package e_d_i_errors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// FetchEdiErrorsHandlerFunc turns a function with the right signature into a fetch edi errors handler
type FetchEdiErrorsHandlerFunc func(FetchEdiErrorsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FetchEdiErrorsHandlerFunc) Handle(params FetchEdiErrorsParams) middleware.Responder {
	return fn(params)
}

// FetchEdiErrorsHandler interface for that can handle valid fetch edi errors params
type FetchEdiErrorsHandler interface {
	Handle(FetchEdiErrorsParams) middleware.Responder
}

// NewFetchEdiErrors creates a new http.Handler for the fetch edi errors operation
func NewFetchEdiErrors(ctx *middleware.Context, handler FetchEdiErrorsHandler) *FetchEdiErrors {
	return &FetchEdiErrors{Context: ctx, Handler: handler}
}

/*
	FetchEdiErrors swagger:route GET /edi-errors EDI Errors fetchEdiErrors

# List of EDI Errors

Returns a list of EDI errors tied to payment requests that are in EDI_ERROR status. This endpoint is for Admin UI use only.
*/
type FetchEdiErrors struct {
	Context *middleware.Context
	Handler FetchEdiErrorsHandler
}

func (o *FetchEdiErrors) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewFetchEdiErrorsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
