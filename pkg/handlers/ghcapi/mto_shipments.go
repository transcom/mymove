package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
)

type ListMTOShipmentsHandler struct {
	handlers.HandlerContext
}

// Handle handler that lists mto service items for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	return mtoshipmentops.NewListMTOShipmentsOK()
}
