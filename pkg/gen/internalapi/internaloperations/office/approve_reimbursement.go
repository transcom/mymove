// Code generated by go-swagger; DO NOT EDIT.

package office

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ApproveReimbursementHandlerFunc turns a function with the right signature into a approve reimbursement handler
type ApproveReimbursementHandlerFunc func(ApproveReimbursementParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ApproveReimbursementHandlerFunc) Handle(params ApproveReimbursementParams) middleware.Responder {
	return fn(params)
}

// ApproveReimbursementHandler interface for that can handle valid approve reimbursement params
type ApproveReimbursementHandler interface {
	Handle(ApproveReimbursementParams) middleware.Responder
}

// NewApproveReimbursement creates a new http.Handler for the approve reimbursement operation
func NewApproveReimbursement(ctx *middleware.Context, handler ApproveReimbursementHandler) *ApproveReimbursement {
	return &ApproveReimbursement{Context: ctx, Handler: handler}
}

/*
	ApproveReimbursement swagger:route POST /reimbursement/{reimbursementId}/approve office approveReimbursement

# Approves the reimbursement

Sets the status of the reimbursement to APPROVED.
*/
type ApproveReimbursement struct {
	Context *middleware.Context
	Handler ApproveReimbursementHandler
}

func (o *ApproveReimbursement) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewApproveReimbursementParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}