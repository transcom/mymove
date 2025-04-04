// Code generated by go-swagger; DO NOT EDIT.

package client_certificates

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateClientCertificateHandlerFunc turns a function with the right signature into a update client certificate handler
type UpdateClientCertificateHandlerFunc func(UpdateClientCertificateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateClientCertificateHandlerFunc) Handle(params UpdateClientCertificateParams) middleware.Responder {
	return fn(params)
}

// UpdateClientCertificateHandler interface for that can handle valid update client certificate params
type UpdateClientCertificateHandler interface {
	Handle(UpdateClientCertificateParams) middleware.Responder
}

// NewUpdateClientCertificate creates a new http.Handler for the update client certificate operation
func NewUpdateClientCertificate(ctx *middleware.Context, handler UpdateClientCertificateHandler) *UpdateClientCertificate {
	return &UpdateClientCertificate{Context: ctx, Handler: handler}
}

/*
	UpdateClientCertificate swagger:route PATCH /client-certificates/{clientCertificateId} Client certificates updateClientCertificate

# Updates a client certificate

This endpoint updates a single Client Certificate by ID. Do not use
this endpoint directly as it is meant to be used with the Admin UI
exclusively.
*/
type UpdateClientCertificate struct {
	Context *middleware.Context
	Handler UpdateClientCertificateHandler
}

func (o *UpdateClientCertificate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateClientCertificateParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
