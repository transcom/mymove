// Code generated by go-swagger; DO NOT EDIT.

package move_task_order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateMTOPostCounselingInformationHandlerFunc turns a function with the right signature into a update m t o post counseling information handler
type UpdateMTOPostCounselingInformationHandlerFunc func(UpdateMTOPostCounselingInformationParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateMTOPostCounselingInformationHandlerFunc) Handle(params UpdateMTOPostCounselingInformationParams) middleware.Responder {
	return fn(params)
}

// UpdateMTOPostCounselingInformationHandler interface for that can handle valid update m t o post counseling information params
type UpdateMTOPostCounselingInformationHandler interface {
	Handle(UpdateMTOPostCounselingInformationParams) middleware.Responder
}

// NewUpdateMTOPostCounselingInformation creates a new http.Handler for the update m t o post counseling information operation
func NewUpdateMTOPostCounselingInformation(ctx *middleware.Context, handler UpdateMTOPostCounselingInformationHandler) *UpdateMTOPostCounselingInformation {
	return &UpdateMTOPostCounselingInformation{Context: ctx, Handler: handler}
}

/*
	UpdateMTOPostCounselingInformation swagger:route PATCH /move-task-orders/{moveTaskOrderID}/post-counseling-info moveTaskOrder updateMTOPostCounselingInformation

updateMTOPostCounselingInformation

### Functionality
This endpoint **updates** the MoveTaskOrder to indicate that the Prime has completed Counseling.
This update uses the moveTaskOrderID provided in the path, updates the move status and marks child elements of the move to indicate the update.
No body object is expected for this request.

**For Full/Partial PPMs**: This action is required so that the customer can start uploading their proof of service docs.

**For other move types**: This action is required for auditing reasons so that we have a record of when the Prime counseled the customer.
*/
type UpdateMTOPostCounselingInformation struct {
	Context *middleware.Context
	Handler UpdateMTOPostCounselingInformationHandler
}

func (o *UpdateMTOPostCounselingInformation) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateMTOPostCounselingInformationParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}