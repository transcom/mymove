// Code generated by go-swagger; DO NOT EDIT.

package queues

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetMovesQueueHandlerFunc turns a function with the right signature into a get moves queue handler
type GetMovesQueueHandlerFunc func(GetMovesQueueParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMovesQueueHandlerFunc) Handle(params GetMovesQueueParams) middleware.Responder {
	return fn(params)
}

// GetMovesQueueHandler interface for that can handle valid get moves queue params
type GetMovesQueueHandler interface {
	Handle(GetMovesQueueParams) middleware.Responder
}

// NewGetMovesQueue creates a new http.Handler for the get moves queue operation
func NewGetMovesQueue(ctx *middleware.Context, handler GetMovesQueueHandler) *GetMovesQueue {
	return &GetMovesQueue{Context: ctx, Handler: handler}
}

/*
	GetMovesQueue swagger:route GET /queues/moves queues getMovesQueue

# Gets queued list of all customer moves by GBLOC origin

An office TOO user will be assigned a transportation office that will determine which moves are displayed in their queue based on the origin duty location.  GHC moves will show up here onced they have reached the submitted status sent by the customer and have move task orders, shipments, and service items to approve.
*/
type GetMovesQueue struct {
	Context *middleware.Context
	Handler GetMovesQueueHandler
}

func (o *GetMovesQueue) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetMovesQueueParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
