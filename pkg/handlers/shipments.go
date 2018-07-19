package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	// "go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForShipmentModel(s models.Shipment) *internalmessages.Shipment {
	shipmentPayload := &internalmessages.Shipment{
		ID:     strfmt.UUID(s.ID.String()),
		MoveID: strfmt.UUID(s.MoveID.String()),
		TrafficDistributionListID:    fmtUUIDPtr(s.TrafficDistributionListID),
		SourceGbloc:                  s.SourceGBLOC,
		Market:                       s.Market,
		Status:                       s.Status,
		BookDate:                     fmtDatePtr(s.BookDate),
		RequestedPickupDate:          fmtDatePtr(s.RequestedPickupDate),
		PickupDate:                   fmtDatePtr(s.PickupDate),
		DeliveryDate:                 fmtDatePtr(s.DeliveryDate),
		CreatedAt:                    strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                    strfmt.DateTime(s.UpdatedAt),
		EstimatedPackDays:            s.EstimatedPackDays,
		EstimatedTransitDays:         s.EstimatedTransitDays,
		PickupAddress:                payloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:    s.HasSecondaryPickupAddress,
		SecondaryPickupAddress:       payloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:           s.HasDeliveryAddress,
		DeliveryAddress:              payloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress: s.HasPartialSITDeliveryAddress,
		PartialSitDeliveryAddress:    payloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:               fmtPoundPtr(s.WeightEstimate),
		ProgearWeightEstimate:        fmtPoundPtr(s.ProgearWeightEstimate),
		SpouseProgearWeightEstimate:  fmtPoundPtr(s.SpouseProgearWeightEstimate),
	}
	return shipmentPayload
}

// CreateShipmentHandler creates a Shipment
type CreateShipmentHandler HandlerContext

// Handle is the handler
func (h CreateShipmentHandler) Handle(params shipmentop.CreateShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.Shipment

	pickupAddress := addressModelFromPayload(payload.PickupAddress)
	secondaryPickupAddress := addressModelFromPayload(payload.SecondaryPickupAddress)
	deliveryAddress := addressModelFromPayload(payload.PickupAddress)
	partialSITDeliveryAddress := addressModelFromPayload(payload.PartialSitDeliveryAddress)

	newShipment := models.Shipment{
		MoveID:                       move.ID,
		Status:                       "DRAFT",
		EstimatedPackDays:            payload.EstimatedPackDays,
		EstimatedTransitDays:         payload.EstimatedTransitDays,
		WeightEstimate:               (*unit.Pound)(payload.WeightEstimate),
		ProgearWeightEstimate:        (*unit.Pound)(payload.ProgearWeightEstimate),
		SpouseProgearWeightEstimate:  (*unit.Pound)(payload.SpouseProgearWeightEstimate),
		PickupAddress:                pickupAddress,
		HasSecondaryPickupAddress:    payload.HasSecondaryPickupAddress,
		SecondaryPickupAddress:       secondaryPickupAddress,
		HasDeliveryAddress:           payload.HasDeliveryAddress,
		DeliveryAddress:              deliveryAddress,
		HasPartialSITDeliveryAddress: payload.HasPartialSitDeliveryAddress,
		PartialSITDeliveryAddress:    partialSITDeliveryAddress,
	}

	verrs, err := models.SaveShipmentAndAddresses(h.db, &newShipment)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	shipmentPayload := payloadForShipmentModel(newShipment)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return shipmentop.NewCreateShipmentCreated().WithPayload(shipmentPayload)
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
