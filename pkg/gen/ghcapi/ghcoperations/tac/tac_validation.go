// Code generated by go-swagger; DO NOT EDIT.

package tac

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// TacValidationHandlerFunc turns a function with the right signature into a tac validation handler
type TacValidationHandlerFunc func(TacValidationParams) middleware.Responder

// Handle executing the request and returning a response
func (fn TacValidationHandlerFunc) Handle(params TacValidationParams) middleware.Responder {
	return fn(params)
}

// TacValidationHandler interface for that can handle valid tac validation params
type TacValidationHandler interface {
	Handle(TacValidationParams) middleware.Responder
}

// NewTacValidation creates a new http.Handler for the tac validation operation
func NewTacValidation(ctx *middleware.Context, handler TacValidationHandler) *TacValidation {
	return &TacValidation{Context: ctx, Handler: handler}
}

/*
	TacValidation swagger:route GET /tac/valid tac order tacValidation

# Validation of a TAC value

Returns a boolean based on whether a tac value is valid or not
*/
type TacValidation struct {
	Context *middleware.Context
	Handler TacValidationHandler
}

func (o *TacValidation) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewTacValidationParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}