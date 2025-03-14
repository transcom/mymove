// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// IndexUsersHandlerFunc turns a function with the right signature into a index users handler
type IndexUsersHandlerFunc func(IndexUsersParams) middleware.Responder

// Handle executing the request and returning a response
func (fn IndexUsersHandlerFunc) Handle(params IndexUsersParams) middleware.Responder {
	return fn(params)
}

// IndexUsersHandler interface for that can handle valid index users params
type IndexUsersHandler interface {
	Handle(IndexUsersParams) middleware.Responder
}

// NewIndexUsers creates a new http.Handler for the index users operation
func NewIndexUsers(ctx *middleware.Context, handler IndexUsersHandler) *IndexUsers {
	return &IndexUsers{Context: ctx, Handler: handler}
}

/*
	IndexUsers swagger:route GET /users Users indexUsers

# List Users

This endpoint returns a list of Users. Do not use this endpoint directly as it
is meant to be used with the Admin UI exclusively.
*/
type IndexUsers struct {
	Context *middleware.Context
	Handler IndexUsersHandler
}

func (o *IndexUsers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewIndexUsersParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
