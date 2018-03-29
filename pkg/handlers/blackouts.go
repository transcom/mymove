package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
)

// CreateBlackoutHandler returns a list of all the Blackouts
type CreateBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h CreateBlackoutHandler) Handle(params apioperations.CreateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .createBlackout has not yet been implemented")
}

// BlackoutIndexHandler returns a list of all the Blackouts
type BlackoutIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h BlackoutIndexHandler) Handle(params apioperations.IndexBlackoutsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexBlackouts has not yet been implemented")
}

// DeleteBlackoutHandler returns a list of all the Blackouts
type DeleteBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h DeleteBlackoutHandler) Handle(params apioperations.DeleteBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .deleteBlackout has not yet been implemented")
}

// GetBlackoutHandler returns a list of all the Blackouts
type GetBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h GetBlackoutHandler) Handle(params apioperations.GetBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .getBlackout has not yet been implemented")
}

// UpdateBlackoutHandler returns a list of all the Blackouts
type UpdateBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h UpdateBlackoutHandler) Handle(params apioperations.UpdateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .updateBlackout has not yet been implemented")
}
