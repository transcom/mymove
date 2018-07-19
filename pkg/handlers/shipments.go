package handlers

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

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

func payloadForShipmentModel(s models.ShipmentWithOffer) *internalmessages.ShipmentPayload {
	shipmentPayload := &internalmessages.ShipmentPayload{
		ID:                              fmtUUID(s.ID),
		PickupDate:                      fmtDate(time.Now()),
		DeliveryDate:                    fmtDate(time.Now()),
		TrafficDistributionListID:       fmtUUID(s.TrafficDistributionListID),
		TransportationServiceProviderID: fmtUUIDPtr(s.TransportationServiceProviderID),
		AdministrativeShipment:          (s.AdministrativeShipment),
		CreatedAt:                       fmtDateTime(s.CreatedAt),
		UpdatedAt:                       fmtDateTime(s.UpdatedAt),
	}
	return shipmentPayload
}

// IndexShipmentsHandler returns a list of shipments
type IndexShipmentsHandler HandlerContext

// Handle retrieves a list of all shipments
func (h IndexShipmentsHandler) Handle(p shipmentop.IndexShipmentsParams) middleware.Responder {
	var response middleware.Responder

	shipments, err := models.FetchShipments(h.db, false)

	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = shipmentop.NewIndexShipmentsBadRequest()
	} else {
		isp := make(internalmessages.IndexShipmentsPayload, len(shipments))
		for i, s := range shipments {
			isp[i] = payloadForShipmentModel(s)
		}
		response = shipmentop.NewIndexShipmentsOK().WithPayload(isp)
	}
	return response
}

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

func publicPayloadForShipmentModel(s models.Shipment) *apimessages.ShipmentPayload {
	shipmentPayload := &apimessages.ShipmentPayload{
		ID: *fmtUUID(s.ID),
		TrafficDistributionListID:    *fmtUUID(s.TrafficDistributionListID),
		PickupDate:                   *fmtDate(s.PickupDate),
		DeliveryDate:                 *fmtDate(s.DeliveryDate),
		CreatedAt:                    *fmtDateTime(s.CreatedAt),
		UpdatedAt:                    *fmtDateTime(s.UpdatedAt),
		SourceGbloc:                  apimessages.GBLOC(s.SourceGBLOC),
		Market:                       apimessages.ShipmentMarket(*s.Market),
		BookDate:                     *fmtDate(s.BookDate),
		RequestedPickupDate:          *fmtDateTime(s.RequestedPickupDate),
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
		WeightEstimate:               *fmtInt64(s.WeightEstimate.Int()),
		ProgearWeightEstimate:        *fmtInt64(s.ProgearWeightEstimate.Int()),
		SpouseProgearWeightEstimate:  *fmtInt64(s.SpouseProgearWeightEstimate.Int()),
	}
	return shipmentPayload
}

// PublicIndexShipmentsHandler returns a list of shipments
type PublicIndexShipmentsHandler HandlerContext

// Handle retrieves a list of all shipments
func (h PublicIndexShipmentsHandler) Handle(p publicshipmentop.IndexShipmentsParams) middleware.Responder {
	var response middleware.Responder

	session := auth.SessionFromRequestContext(p.HTTPRequest)
	transportationServiceProviderID := session.EntityID

	shipments := []models.Shipment{}

	err := h.db.Eager().Where(fmt.Sprintf("transportation_service_provider = '%s'", transportationServiceProviderID)).All(&shipments)

	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = publicshipmentop.NewIndexShipmentsBadRequest()
	} else {
		isp := make(apimessages.IndexShipmentsPayload, len(shipments))
		for i, s := range shipments {
			isp[i] = publicPayloadForShipmentModel(s)
		}
		response = publicshipmentop.NewIndexShipmentsOK().WithPayload(isp)
	}
	return response
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
