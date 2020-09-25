package primeapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/services"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
)

// UpdateMTOAgentHandler is the handler to update an address
type UpdateMTOAgentHandler struct {
	handlers.HandlerContext
	MTOAgentUpdater services.MTOAgentUpdater
}

// Handle updates an address on a shipment
func (h UpdateMTOAgentHandler) Handle(params mtoshipmentops.UpdateMTOAgentParams) middleware.Responder {
	return mtoshipmentops.NewUpdateMTOAgentNotImplemented().WithPayload(
		payloads.NotImplementedError(nil, h.GetTraceID()))
}
