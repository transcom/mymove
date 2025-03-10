// Code generated by go-swagger; DO NOT EDIT.

package webhook

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReceiveWebhookNotificationHandlerFunc turns a function with the right signature into a receive webhook notification handler
type ReceiveWebhookNotificationHandlerFunc func(ReceiveWebhookNotificationParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ReceiveWebhookNotificationHandlerFunc) Handle(params ReceiveWebhookNotificationParams) middleware.Responder {
	return fn(params)
}

// ReceiveWebhookNotificationHandler interface for that can handle valid receive webhook notification params
type ReceiveWebhookNotificationHandler interface {
	Handle(ReceiveWebhookNotificationParams) middleware.Responder
}

// NewReceiveWebhookNotification creates a new http.Handler for the receive webhook notification operation
func NewReceiveWebhookNotification(ctx *middleware.Context, handler ReceiveWebhookNotificationHandler) *ReceiveWebhookNotification {
	return &ReceiveWebhookNotification{Context: ctx, Handler: handler}
}

/*
	ReceiveWebhookNotification swagger:route POST /webhook-notify webhook receiveWebhookNotification

# Test endpoint for receiving messages from our own webhook-client

This endpoint receives a notification that matches the webhook notification model. This is a test endpoint that represents a receiving server. In production, the Prime will set up a receiving endpoint. In testing, this server accepts notifications at this endpoint and simply responds with success and logs them. The `webhook-client` is responsible for retrieving messages from the webhook_notifications table and sending them to the Prime (this endpoint in our testing case) via an mTLS connection.
*/
type ReceiveWebhookNotification struct {
	Context *middleware.Context
	Handler ReceiveWebhookNotificationHandler
}

func (o *ReceiveWebhookNotification) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewReceiveWebhookNotificationParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
