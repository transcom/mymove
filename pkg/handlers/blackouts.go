package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	publicblackoutsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/blackouts"
)

/*
 * ------------------------------------------
 * The code below is for the INTERNAL REST API.
 * ------------------------------------------
 */

// NO CODE YET!

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

// PublicBlackoutIndexHandler returns a list of all the Blackouts
type PublicBlackoutIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicBlackoutIndexHandler) Handle(params publicblackoutsop.IndexBlackoutsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexBlackouts has not yet been implemented")
}

// PublicCreateBlackoutHandler returns a list of all the Blackouts
type PublicCreateBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicCreateBlackoutHandler) Handle(params publicblackoutsop.CreateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .createBlackout has not yet been implemented")
}

// PublicDeleteBlackoutHandler returns a list of all the Blackouts
type PublicDeleteBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicDeleteBlackoutHandler) Handle(params publicblackoutsop.DeleteBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .deleteBlackout has not yet been implemented")
}

// PublicGetBlackoutHandler returns a list of all the Blackouts
type PublicGetBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicGetBlackoutHandler) Handle(params publicblackoutsop.GetBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .getBlackout has not yet been implemented")
}

// PublicUpdateBlackoutHandler returns a list of all the Blackouts
type PublicUpdateBlackoutHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h PublicUpdateBlackoutHandler) Handle(params publicblackoutsop.UpdateBlackoutParams) middleware.Responder {
	return middleware.NotImplemented("operation .updateBlackout has not yet been implemented")
}
