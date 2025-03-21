// Code generated by go-swagger; DO NOT EDIT.

package payment_request

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdatePaymentRequestStatusHandlerFunc turns a function with the right signature into a update payment request status handler
type UpdatePaymentRequestStatusHandlerFunc func(UpdatePaymentRequestStatusParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdatePaymentRequestStatusHandlerFunc) Handle(params UpdatePaymentRequestStatusParams) middleware.Responder {
	return fn(params)
}

// UpdatePaymentRequestStatusHandler interface for that can handle valid update payment request status params
type UpdatePaymentRequestStatusHandler interface {
	Handle(UpdatePaymentRequestStatusParams) middleware.Responder
}

// NewUpdatePaymentRequestStatus creates a new http.Handler for the update payment request status operation
func NewUpdatePaymentRequestStatus(ctx *middleware.Context, handler UpdatePaymentRequestStatusHandler) *UpdatePaymentRequestStatus {
	return &UpdatePaymentRequestStatus{Context: ctx, Handler: handler}
}

/*
	UpdatePaymentRequestStatus swagger:route PATCH /payment-requests/{paymentRequestID}/status paymentRequest updatePaymentRequestStatus

updatePaymentRequestStatus

Updates status of a payment request to REVIEWED, SENT_TO_GEX, TPPS_RECEIVED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, PAID, EDI_ERROR, or DEPRECATED.

A status of REVIEWED can optionally have a `rejectionReason`.

This is a support endpoint and is not available in production.
*/
type UpdatePaymentRequestStatus struct {
	Context *middleware.Context
	Handler UpdatePaymentRequestStatusHandler
}

func (o *UpdatePaymentRequestStatus) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdatePaymentRequestStatusParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
