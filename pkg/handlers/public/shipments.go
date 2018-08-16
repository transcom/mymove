package public

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	publicshipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
)

func publicPayloadForShipmentModel(s models.Shipment) *apimessages.Shipment {
	shipmentPayload := &apimessages.Shipment{
		ID: *utils.FmtUUID(s.ID),
		TrafficDistributionList:      publicPayloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMember:                publicPayloadForServiceMemberModel(s.ServiceMember),
		PickupDate:                   *utils.FmtDateTimePtr(s.PickupDate),
		DeliveryDate:                 *utils.FmtDateTimePtr(s.DeliveryDate),
		CreatedAt:                    strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                    strfmt.DateTime(s.UpdatedAt),
		SourceGbloc:                  apimessages.GBLOC(*s.SourceGBLOC),
		DestinationGbloc:             apimessages.GBLOC(*s.DestinationGBLOC),
		Market:                       apimessages.ShipmentMarket(*s.Market),
		BookDate:                     *utils.FmtDatePtr(s.BookDate),
		RequestedPickupDate:          *utils.FmtDateTimePtr(s.RequestedPickupDate),
		Move:                         publicPayloadForMoveModel(s.Move),
		Status:                       apimessages.ShipmentStatus(s.Status),
		EstimatedPackDays:            utils.FmtInt64(*s.EstimatedPackDays),
		EstimatedTransitDays:         utils.FmtInt64(*s.EstimatedTransitDays),
		PickupAddress:                publicPayloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:    *utils.FmtBool(s.HasSecondaryPickupAddress),
		SecondaryPickupAddress:       publicPayloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:           *utils.FmtBool(s.HasDeliveryAddress),
		DeliveryAddress:              publicPayloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress: *utils.FmtBool(s.HasPartialSITDeliveryAddress),
		PartialSitDeliveryAddress:    publicPayloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:               utils.FmtInt64(s.WeightEstimate.Int64()),
		ProgearWeightEstimate:        utils.FmtInt64(s.ProgearWeightEstimate.Int64()),
		SpouseProgearWeightEstimate:  utils.FmtInt64(s.SpouseProgearWeightEstimate.Int64()),
	}
	return shipmentPayload
}

// IndexShipmentsHandler returns a list of shipments
type IndexShipmentsHandler utils.HandlerContext

// Handle retrieves a list of all shipments
func (h IndexShipmentsHandler) Handle(params publicshipmentop.IndexShipmentsParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger.Error("Missing TSP User ID")
		return publicshipmentop.NewIndexShipmentsForbidden()
	}

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.Db, session.TspUserID)
	if err != nil {
		h.Logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewIndexShipmentsForbidden()
	}

	shipments, err := models.FetchShipmentsByTSP(h.Db, tspUser.TransportationServiceProviderID,
		params.Status, params.OrderBy, params.Limit, params.Offset)
	if err != nil {
		h.Logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewIndexShipmentsBadRequest()
	}

	isp := make(apimessages.IndexShipments, len(shipments))
	for i, s := range shipments {
		isp[i] = publicPayloadForShipmentModel(s)
	}
	return publicshipmentop.NewIndexShipmentsOK().WithPayload(isp)
}

// GetShipmentHandler returns a particular shipment
type GetShipmentHandler utils.HandlerContext

// Handle returns a specified shipment
func (h GetShipmentHandler) Handle(params publicshipmentop.GetShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentUUID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger.Error("Missing TSP User ID")
		return publicshipmentop.NewGetShipmentForbidden()
	}

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.Db, session.TspUserID)
	if err != nil {
		h.Logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.Db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentBadRequest()
	}

	sp := publicPayloadForShipmentModel(*shipment)
	return publicshipmentop.NewGetShipmentOK().WithPayload(sp)
}

// CreateShipmentAcceptHandler allows a TSP to accept a particular shipment
type CreateShipmentAcceptHandler utils.HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentAcceptHandler) Handle(params publicshipmentop.CreateShipmentAcceptParams) middleware.Responder {
	return middleware.NotImplemented("operation .acceptShipment has not yet been implemented")
}

// CreateShipmentRejectHandler allows a TSP to refuse a particular shipment
type CreateShipmentRejectHandler utils.HandlerContext

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentRejectHandler) Handle(params publicshipmentop.CreateShipmentRejectParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

// UpdateShipmentHandler allows a TSP to refuse a particular shipment
type UpdateShipmentHandler utils.HandlerContext

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h UpdateShipmentHandler) Handle(p publicshipmentop.UpdateShipmentParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

// GetShipmentContactDetailsHandler allows a TSP to accept a particular shipment
type GetShipmentContactDetailsHandler utils.HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentContactDetailsHandler) Handle(p publicshipmentop.GetShipmentContactDetailsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentClaimsHandler allows a TSP to accept a particular shipment
type GetShipmentClaimsHandler utils.HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentClaimsHandler) Handle(p publicshipmentop.GetShipmentClaimsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentDocumentsHandler allows a TSP to accept a particular shipment
type GetShipmentDocumentsHandler utils.HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentDocumentsHandler) Handle(p publicshipmentop.GetShipmentDocumentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// CreateShipmentDocumentHandler allows a TSP to accept a particular shipment
type CreateShipmentDocumentHandler utils.HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentDocumentHandler) Handle(p publicshipmentop.CreateShipmentDocumentParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}
