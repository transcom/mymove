package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
)

// BlackoutIndexHandler returns a list of all the Blackouts
type BlackoutIndexHandler HandlerContext

// Handle simply returns a NotImplementedError
func (h BlackoutIndexHandler) Handle(params apioperations.IndexBlackoutsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexTSPs has not yet been implemented")
}
