package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"

	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	publicshipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

/*
 * ------------------------------------------
 * The code below is for the INTERNAL REST API.
 * ------------------------------------------
 */
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
		WeightEstimate:               poundPtrFromInt64Ptr(payload.WeightEstimate),
		ProgearWeightEstimate:        poundPtrFromInt64Ptr(payload.ProgearWeightEstimate),
		SpouseProgearWeightEstimate:  poundPtrFromInt64Ptr(payload.SpouseProgearWeightEstimate),
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
	return shipmentop.NewCreateShipmentCreated().WithPayload(shipmentPayload)
}

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

func publicPayloadForShipmentModel(s models.Shipment) *apimessages.ShipmentPayload {
	shipmentPayload := &apimessages.ShipmentPayload{
		ID: *fmtUUID(s.ID),
		TrafficDistributionListID:    *fmtUUID(*s.TrafficDistributionListID),
		PickupDate:                   *fmtDate(*s.PickupDate),
		DeliveryDate:                 *fmtDate(*s.DeliveryDate),
		CreatedAt:                    *fmtDateTime(s.CreatedAt),
		UpdatedAt:                    *fmtDateTime(s.UpdatedAt),
		SourceGbloc:                  apimessages.GBLOC(*s.SourceGBLOC),
		Market:                       apimessages.ShipmentMarket(*s.Market),
		BookDate:                     *fmtDate(*s.BookDate),
		RequestedPickupDate:          *fmtDateTime(*s.RequestedPickupDate),
		MoveID:                       *fmtUUID(s.MoveID),
		Status:                       apimessages.ShipmentStatus(s.Status),
		EstimatedPackDays:            *fmtInt64(*s.EstimatedPackDays),
		EstimatedTransitDays:         *fmtInt64(*s.EstimatedTransitDays),
		PickupAddress:                publicPayloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:    fmtBool(s.HasSecondaryPickupAddress),
		SecondaryPickupAddress:       publicPayloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:           fmtBool(s.HasDeliveryAddress),
		DeliveryAddress:              publicPayloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress: fmtBool(s.HasPartialSITDeliveryAddress),
		PartialSitDeliveryAddress:    publicPayloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:               *fmtInt64(s.WeightEstimate.Int64()),
		ProgearWeightEstimate:        *fmtInt64(s.ProgearWeightEstimate.Int64()),
		SpouseProgearWeightEstimate:  *fmtInt64(s.SpouseProgearWeightEstimate.Int64()),
	}
	return shipmentPayload
}

// PublicIndexShipmentsHandler returns a list of shipments
type PublicIndexShipmentsHandler HandlerContext

// Handle retrieves a list of all shipments
func (h PublicIndexShipmentsHandler) Handle(p publicshipmentop.IndexShipmentsParams) middleware.Responder {

	session := auth.SessionFromRequestContext(p.HTTPRequest)
	// Possible they are coming from the wrong endpoint and thus the session is missing the EntityID
	if session.EntityID == uuid.Nil {
		return publicshipmentop.NewIndexShipmentsForbidden()
	}

	shipments, err := models.FetchShipmentsByTSP(h.db, session.EntityID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewIndexShipmentsBadRequest()
	}

	isp := make(apimessages.IndexShipmentsPayload, len(shipments))
	for i, s := range shipments {
		isp[i] = publicPayloadForShipmentModel(s)
	}
	return publicshipmentop.NewIndexShipmentsOK().WithPayload(isp)
}

// PublicGetShipmentHandler returns a particular shipment
type PublicGetShipmentHandler HandlerContext

// Handle returns a specified shipment
func (h PublicGetShipmentHandler) Handle(p publicshipmentop.GetShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .getShipment has not yet been implemented")
}

// PublicCreateShipmentAcceptHandler allows a TSP to accept a particular shipment
type PublicCreateShipmentAcceptHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicCreateShipmentAcceptHandler) Handle(p publicshipmentop.CreateShipmentAcceptParams) middleware.Responder {
	return middleware.NotImplemented("operation .acceptShipment has not yet been implemented")
}

// PublicCreateShipmentRefuseHandler allows a TSP to refuse a particular shipment
type PublicCreateShipmentRefuseHandler HandlerContext

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicCreateShipmentRefuseHandler) Handle(p publicshipmentop.CreateShipmentRefuseParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

// PublicUpdateShipmentHandler allows a TSP to refuse a particular shipment
type PublicUpdateShipmentHandler HandlerContext

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicUpdateShipmentHandler) Handle(p publicshipmentop.UpdateShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

// PublicGetShipmentContactDetailsHandler allows a TSP to accept a particular shipment
type PublicGetShipmentContactDetailsHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicGetShipmentContactDetailsHandler) Handle(p publicshipmentop.GetShipmentContactDetailsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// PublicGetShipmentClaimsHandler allows a TSP to accept a particular shipment
type PublicGetShipmentClaimsHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicGetShipmentClaimsHandler) Handle(p publicshipmentop.GetShipmentClaimsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// PublicGetShipmentDocumentsHandler allows a TSP to accept a particular shipment
type PublicGetShipmentDocumentsHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicGetShipmentDocumentsHandler) Handle(p publicshipmentop.GetShipmentDocumentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// PublicCreateShipmentDocumentHandler allows a TSP to accept a particular shipment
type PublicCreateShipmentDocumentHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicCreateShipmentDocumentHandler) Handle(p publicshipmentop.CreateShipmentDocumentParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}
