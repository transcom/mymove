// Code generated by go-swagger; DO NOT EDIT.

package uploads

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetUploadStatusHandlerFunc turns a function with the right signature into a get upload status handler
type GetUploadStatusHandlerFunc func(GetUploadStatusParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUploadStatusHandlerFunc) Handle(params GetUploadStatusParams) middleware.Responder {
	return fn(params)
}

// GetUploadStatusHandler interface for that can handle valid get upload status params
type GetUploadStatusHandler interface {
	Handle(GetUploadStatusParams) middleware.Responder
}

// NewGetUploadStatus creates a new http.Handler for the get upload status operation
func NewGetUploadStatus(ctx *middleware.Context, handler GetUploadStatusHandler) *GetUploadStatus {
	return &GetUploadStatus{Context: ctx, Handler: handler}
}

/*
	GetUploadStatus swagger:route GET /uploads/{uploadID}/status uploads getUploadStatus

# Returns status of an upload

Returns status of an upload based on antivirus run
*/
type GetUploadStatus struct {
	Context *middleware.Context
	Handler GetUploadStatusHandler
}

func (o *GetUploadStatus) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetUploadStatusParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
