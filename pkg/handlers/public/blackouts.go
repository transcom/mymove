package public

import (
	"github.com/go-openapi/runtime/middleware"
	blackoutsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/blackouts"
	"github.com/transcom/mymove/pkg/handlers/utils"
)

// BlackoutIndexHandler returns a list of all the Blackouts
type BlackoutIndexHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h BlackoutIndexHandler) Handle(params blackoutsop.IndexBlackoutsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexBlackouts has not yet been implemented")
}

// CreateBlackoutHandler returns a list of all the Blackouts
type CreateBlackoutHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h CreateBlackoutHandler) Handle(params blackoutsop.CreateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .createBlackout has not yet been implemented")
}

// DeleteBlackoutHandler returns a list of all the Blackouts
type DeleteBlackoutHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h DeleteBlackoutHandler) Handle(params blackoutsop.DeleteBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .deleteBlackout has not yet been implemented")
}

// GetBlackoutHandler returns a list of all the Blackouts
type GetBlackoutHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h GetBlackoutHandler) Handle(params blackoutsop.GetBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .getBlackout has not yet been implemented")
}

// UpdateBlackoutHandler returns a list of all the Blackouts
type UpdateBlackoutHandler utils.HandlerContext

// Handle simply returns a NotImplementedError
func (h UpdateBlackoutHandler) Handle(params blackoutsop.UpdateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .updateBlackout has not yet been implemented")
}
