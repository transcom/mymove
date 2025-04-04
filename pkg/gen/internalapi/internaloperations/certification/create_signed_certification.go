// Code generated by go-swagger; DO NOT EDIT.

package certification

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CreateSignedCertificationHandlerFunc turns a function with the right signature into a create signed certification handler
type CreateSignedCertificationHandlerFunc func(CreateSignedCertificationParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateSignedCertificationHandlerFunc) Handle(params CreateSignedCertificationParams) middleware.Responder {
	return fn(params)
}

// CreateSignedCertificationHandler interface for that can handle valid create signed certification params
type CreateSignedCertificationHandler interface {
	Handle(CreateSignedCertificationParams) middleware.Responder
}

// NewCreateSignedCertification creates a new http.Handler for the create signed certification operation
func NewCreateSignedCertification(ctx *middleware.Context, handler CreateSignedCertificationHandler) *CreateSignedCertification {
	return &CreateSignedCertification{Context: ctx, Handler: handler}
}

/*
	CreateSignedCertification swagger:route POST /moves/{moveId}/signed_certifications certification createSignedCertification

# Submits signed certification for the given move ID

Create an instance of signed_certification tied to the move ID
*/
type CreateSignedCertification struct {
	Context *middleware.Context
	Handler CreateSignedCertificationHandler
}

func (o *CreateSignedCertification) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCreateSignedCertificationParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
