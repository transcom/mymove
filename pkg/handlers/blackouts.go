package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	publicblackoutsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/blackouts"
)

// BlackoutIndexHandler returns a list of all the Blackouts
type BlackoutIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h BlackoutIndexHandler) Handle(params publicblackoutsop.IndexBlackoutsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexBlackouts has not yet been implemented")
}

// CreateBlackoutHandler returns a list of all the Blackouts
type CreateBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h CreateBlackoutHandler) Handle(params publicblackoutsop.CreateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .createBlackout has not yet been implemented")
}

// DeleteBlackoutHandler returns a list of all the Blackouts
type DeleteBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h DeleteBlackoutHandler) Handle(params publicblackoutsop.DeleteBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .deleteBlackout has not yet been implemented")
}

// GetBlackoutHandler returns a list of all the Blackouts
type GetBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h GetBlackoutHandler) Handle(params publicblackoutsop.GetBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .getBlackout has not yet been implemented")
}

// UpdateBlackoutHandler returns a list of all the Blackouts
type UpdateBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h UpdateBlackoutHandler) Handle(params publicblackoutsop.UpdateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .updateBlackout has not yet been implemented")
}
