// Code generated by go-swagger; DO NOT EDIT.

package uploads

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CreateUploadHandlerFunc turns a function with the right signature into a create upload handler
type CreateUploadHandlerFunc func(CreateUploadParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateUploadHandlerFunc) Handle(params CreateUploadParams) middleware.Responder {
	return fn(params)
}

// CreateUploadHandler interface for that can handle valid create upload params
type CreateUploadHandler interface {
	Handle(CreateUploadParams) middleware.Responder
}

// NewCreateUpload creates a new http.Handler for the create upload operation
func NewCreateUpload(ctx *middleware.Context, handler CreateUploadHandler) *CreateUpload {
	return &CreateUpload{Context: ctx, Handler: handler}
}

/*
	CreateUpload swagger:route POST /uploads uploads createUpload

# Create a new upload

Uploads represent a single digital file, such as a JPEG or PDF.
*/
type CreateUpload struct {
	Context *middleware.Context
	Handler CreateUploadHandler
}

func (o *CreateUpload) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCreateUploadParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
