// Code generated by go-swagger; DO NOT EDIT.

package move

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UploadAdditionalDocumentsHandlerFunc turns a function with the right signature into a upload additional documents handler
type UploadAdditionalDocumentsHandlerFunc func(UploadAdditionalDocumentsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UploadAdditionalDocumentsHandlerFunc) Handle(params UploadAdditionalDocumentsParams) middleware.Responder {
	return fn(params)
}

// UploadAdditionalDocumentsHandler interface for that can handle valid upload additional documents params
type UploadAdditionalDocumentsHandler interface {
	Handle(UploadAdditionalDocumentsParams) middleware.Responder
}

// NewUploadAdditionalDocuments creates a new http.Handler for the upload additional documents operation
func NewUploadAdditionalDocuments(ctx *middleware.Context, handler UploadAdditionalDocumentsHandler) *UploadAdditionalDocuments {
	return &UploadAdditionalDocuments{Context: ctx, Handler: handler}
}

/*
	UploadAdditionalDocuments swagger:route PATCH /moves/{moveID}/uploadAdditionalDocuments move uploadAdditionalDocuments

# Patch the additional documents for a given move

Customers will on occaision need the ability to upload additional supporting documents, for a variety of reasons. This does not include amended order.
*/
type UploadAdditionalDocuments struct {
	Context *middleware.Context
	Handler UploadAdditionalDocumentsHandler
}

func (o *UploadAdditionalDocuments) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUploadAdditionalDocumentsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}