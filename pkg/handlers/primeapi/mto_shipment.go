package primeapi

import (
	"github.com/go-openapi/runtime/middleware"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
)

type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	// fetch shipment
	// convert if-unmodified-since to a time that can be compared to updated_at
	// check if shipment's updated_at is before the if-unmodified-since
	// TRUE - do the updates
	// FALSE - return 412
	return mtoshipmentops.NewUpdateMTOShipmentOK()
}
