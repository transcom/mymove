package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	// "go.uber.org/zap"

	// shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForShipmentModel(s models.Shipment) *internalmessages.Shipment {
	shipmentPayload := &internalmessages.Shipment{
		ID:     strfmt.UUID(s.ID.String()),
		MoveID: strfmt.UUID(s.MoveID.String()),
		TrafficDistributionListID: fmtUUIDPtr(s.TrafficDistributionListID),
		SourceGbloc:               s.SourceGBLOC,
		Market:                    s.Market,
		Status:                    s.Status,
		BookDate:                  fmtDatePtr(s.BookDate),
		RequestedPickupDate:       fmtDatePtr(s.RequestedPickupDate),
		PickupDate:                fmtDatePtr(s.PickupDate),
		DeliveryDate:              fmtDatePtr(s.DeliveryDate),
		CreatedAt:                 strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                 strfmt.DateTime(s.UpdatedAt),
	}
	return shipmentPayload
}

/* NOTE - The code above is for the INTERNAL API. The code below is for the public API. These will, obviously,
need to be reconciled. This will be done when the NotImplemented code below is Implemented
*/

// ShipmentIndexHandler returns a list of shipments
type ShipmentIndexHandler HandlerContext

// Handle retrieves a list of all shipments
func (h ShipmentIndexHandler) Handle(p apioperations.IndexShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexShipments has not yet been implemented")
}

// GetShipmentHandler returns a particular shipment
type GetShipmentHandler HandlerContext

// Handle returns a specified shipment
func (h GetShipmentHandler) Handle(p apioperations.GetShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .getShipment has not yet been implemented")
}

// AcceptShipmentHandler allows a TSP to accept a particular shipment
type AcceptShipmentHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h AcceptShipmentHandler) Handle(p apioperations.AcceptShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .acceptShipment has not yet been implemented")
}

// RefuseShipmentHandler allows a TSP to refuse a particular shipment
type RefuseShipmentHandler HandlerContext

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h RefuseShipmentHandler) Handle(p apioperations.RefuseShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

// UpdateShipmentHandler allows a TSP to refuse a particular shipment
type UpdateShipmentHandler HandlerContext

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h UpdateShipmentHandler) Handle(p apioperations.UpdateShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

// ShipmentContactDetailsHandler allows a TSP to accept a particular shipment
type ShipmentContactDetailsHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h ShipmentContactDetailsHandler) Handle(p apioperations.ShipmentContactDetailsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentClaimsHandler allows a TSP to accept a particular shipment
type GetShipmentClaimsHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentClaimsHandler) Handle(p apioperations.GetShipmentClaimsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentDocumentsHandler allows a TSP to accept a particular shipment
type GetShipmentDocumentsHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentDocumentsHandler) Handle(p apioperations.GetShipmentDocumentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// CreateShipmentDocumentHandler allows a TSP to accept a particular shipment
type CreateShipmentDocumentHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentDocumentHandler) Handle(p apioperations.CreateShipmentDocumentParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}
