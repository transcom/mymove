// Code generated by go-swagger; DO NOT EDIT.

package client_certificates

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CreateClientCertificateHandlerFunc turns a function with the right signature into a create client certificate handler
type CreateClientCertificateHandlerFunc func(CreateClientCertificateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateClientCertificateHandlerFunc) Handle(params CreateClientCertificateParams) middleware.Responder {
	return fn(params)
}

// CreateClientCertificateHandler interface for that can handle valid create client certificate params
type CreateClientCertificateHandler interface {
	Handle(CreateClientCertificateParams) middleware.Responder
}

// NewCreateClientCertificate creates a new http.Handler for the create client certificate operation
func NewCreateClientCertificate(ctx *middleware.Context, handler CreateClientCertificateHandler) *CreateClientCertificate {
	return &CreateClientCertificate{Context: ctx, Handler: handler}
}

/*
	CreateClientCertificate swagger:route POST /client-certificates Client certificates createClientCertificate

create a client cert

This endpoint creates a Client Certificate record and returns the
created record in the `201` response. Do not use this endpoint
directly as it is meant to be used with the Admin UI exclusively.
*/
type CreateClientCertificate struct {
	Context *middleware.Context
	Handler CreateClientCertificateHandler
}

func (o *CreateClientCertificate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCreateClientCertificateParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}