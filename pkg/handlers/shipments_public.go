package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"
	"time"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	publicshipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

func publicPayloadForShipmentModel(s models.Shipment) *apimessages.Shipment {
	shipmentPayload := &apimessages.Shipment{
		ID: *fmtUUID(s.ID),
		TrafficDistributionList:             publicPayloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMember:                       publicPayloadForServiceMemberModel(s.ServiceMember),
		PickupDate:                          *fmtDateTimePtr(s.PickupDate),
		DeliveryDate:                        *fmtDateTimePtr(s.DeliveryDate),
		CreatedAt:                           strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                           strfmt.DateTime(s.UpdatedAt),
		SourceGbloc:                         apimessages.GBLOC(*s.SourceGBLOC),
		DestinationGbloc:                    apimessages.GBLOC(*s.DestinationGBLOC),
		Market:                              apimessages.ShipmentMarket(*s.Market),
		BookDate:                            *fmtDatePtr(s.BookDate),
		RequestedPickupDate:                 *fmtDateTimePtr(s.RequestedPickupDate),
		Move:                                publicPayloadForMoveModel(s.Move),
		Status:                              apimessages.ShipmentStatus(s.Status),
		EstimatedPackDays:                   fmtInt64(*s.EstimatedPackDays),
		EstimatedTransitDays:                fmtInt64(*s.EstimatedTransitDays),
		PickupAddress:                       publicPayloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:           *fmtBool(s.HasSecondaryPickupAddress),
		SecondaryPickupAddress:              publicPayloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:                  *fmtBool(s.HasDeliveryAddress),
		DeliveryAddress:                     publicPayloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress:        *fmtBool(s.HasPartialSITDeliveryAddress),
		PartialSitDeliveryAddress:           publicPayloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:                      fmtInt64(s.WeightEstimate.Int64()),
		ProgearWeightEstimate:               fmtInt64(s.ProgearWeightEstimate.Int64()),
		SpouseProgearWeightEstimate:         fmtInt64(s.SpouseProgearWeightEstimate.Int64()),
		PmSurveyPackDate:                    fmtDatePtr(s.PmSurveyPackDate),
		PmSurveyPickupDate:                  fmtDatePtr(s.PmSurveyPickupDate),
		PmSurveyLatestPickupDate:            fmtDatePtr(s.PmSurveyLatestPickupDate),
		PmSurveyEarliestDeliveryDate:        fmtDatePtr(s.PmSurveyEarliestDeliveryDate),
		PmSurveyLatestDeliveryDate:          fmtDatePtr(s.PmSurveyLatestDeliveryDate),
		PmSurveyWeightEstimate:              fmtPoundPtr(s.PmSurveyWeightEstimate),
		PmSurveyProgearWeightEstimate:       fmtPoundPtr(s.PmSurveyProgearWeightEstimate),
		PmSurveySpouseProgearWeightEstimate: fmtPoundPtr(s.PmSurveySpouseProgearWeightEstimate),
		PmSurveyNotes:                       s.PmSurveyNotes,
	}
	return shipmentPayload
}

// PublicIndexShipmentsHandler returns a list of shipments
type PublicIndexShipmentsHandler HandlerContext

// Handle retrieves a list of all shipments
func (h PublicIndexShipmentsHandler) Handle(params publicshipmentop.IndexShipmentsParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.logger.Error("Missing TSP User ID")
		return publicshipmentop.NewIndexShipmentsForbidden()
	}

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.db, session.TspUserID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewIndexShipmentsForbidden()
	}

	shipments, err := models.FetchShipmentsByTSP(h.db, tspUser.TransportationServiceProviderID,
		params.Status, params.OrderBy, params.Limit, params.Offset)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewIndexShipmentsBadRequest()
	}

	isp := make(apimessages.IndexShipments, len(shipments))
	for i, s := range shipments {
		isp[i] = publicPayloadForShipmentModel(s)
	}
	return publicshipmentop.NewIndexShipmentsOK().WithPayload(isp)
}

// PublicGetShipmentHandler returns a particular shipment
type PublicGetShipmentHandler HandlerContext

// Handle returns a specified shipment
func (h PublicGetShipmentHandler) Handle(params publicshipmentop.GetShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentUUID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.logger.Error("Missing TSP User ID")
		return publicshipmentop.NewGetShipmentForbidden()
	}

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.db, session.TspUserID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentBadRequest()
	}

	sp := publicPayloadForShipmentModel(*shipment)
	return publicshipmentop.NewGetShipmentOK().WithPayload(sp)
}

// PublicCreateShipmentAcceptHandler allows a TSP to accept a particular shipment
type PublicCreateShipmentAcceptHandler HandlerContext

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicCreateShipmentAcceptHandler) Handle(params publicshipmentop.CreateShipmentAcceptParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentUUID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.logger.Error("Missing TSP User ID")
		return publicshipmentop.NewGetShipmentForbidden()
	}

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.db, session.TspUserID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentBadRequest()
	}

	// Accept the shipment
	err = shipment.Accept()
	if err != nil {
		h.logger.Info("Attempted to accept shipment, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
		return responseForError(h.logger, err)
	}

	verrs, err := h.db.ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	// Get the Shipment Offer
	shipmentOffer, err := models.FetchShipmentOfferByTSP(h.db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewGetShipmentBadRequest()
	}

	// Accept the Shipment Offer
	err = shipmentOffer.Accept()
	if err != nil {
		h.logger.Info("Attempted to accept shipment offer, got invalid transition", zap.Error(err), zap.Bool("shipmentOffer_accepted", *shipmentOffer.Accepted))
		return responseForError(h.logger, err)
	}

	verrs, err = h.db.ValidateAndUpdate(shipmentOffer)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	// Do we need to update Move status???

	sp := publicPayloadForShipmentModel(*shipment)
	return publicshipmentop.NewCreateShipmentAcceptOK().WithPayload(sp)
}

// PublicCreateShipmentRejectHandler allows a TSP to refuse a particular shipment
type PublicCreateShipmentRejectHandler HandlerContext

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicCreateShipmentRejectHandler) Handle(params publicshipmentop.CreateShipmentRejectParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

func publicPatchShipmentWithPayload(shipment *models.Shipment, payload *apimessages.Shipment) {
	// Premove Survey values entered by TSP agent
	if payload.PmSurveyEarliestDeliveryDate != nil {
		shipment.PmSurveyEarliestDeliveryDate = (*time.Time)(payload.PmSurveyEarliestDeliveryDate)
	}
	if payload.PmSurveyLatestDeliveryDate != nil {
		shipment.PmSurveyLatestDeliveryDate = (*time.Time)(payload.PmSurveyLatestDeliveryDate)
	}
	if payload.PmSurveyLatestPickupDate != nil {
		shipment.PmSurveyLatestPickupDate = (*time.Time)(payload.PmSurveyLatestPickupDate)
	}
	if payload.PmSurveyNotes != nil {
		shipment.PmSurveyNotes = payload.PmSurveyNotes
	}
	if payload.PmSurveyPackDate != nil {
		shipment.PmSurveyPackDate = (*time.Time)(payload.PmSurveyPackDate)
	}
	if payload.PmSurveyPickupDate != nil {
		shipment.PmSurveyPickupDate = (*time.Time)(payload.PmSurveyPickupDate)
	}

	if payload.PmSurveyProgearWeightEstimate != nil {
		shipment.PmSurveyProgearWeightEstimate = poundPtrFromInt64Ptr(payload.PmSurveyProgearWeightEstimate)
	}
	if payload.PmSurveySpouseProgearWeightEstimate != nil {
		shipment.PmSurveySpouseProgearWeightEstimate = poundPtrFromInt64Ptr(payload.PmSurveySpouseProgearWeightEstimate)
	}
	if payload.PmSurveyWeightEstimate != nil {
		shipment.PmSurveyWeightEstimate = poundPtrFromInt64Ptr(payload.PmSurveyWeightEstimate)
	}

}

// PublicPatchShipmentHandler allows a TSP to patch a particular shipment
type PublicPatchShipmentHandler HandlerContext

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicPatchShipmentHandler) Handle(params publicshipmentop.PatchShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentUUID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.logger.Error("Missing TSP User ID")
		return publicshipmentop.NewGetShipmentForbidden()
	}

	tspUser, err := models.FetchTspUserByID(h.db, session.TspUserID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewPatchShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		return publicshipmentop.NewPatchShipmentBadRequest()
	}

	publicPatchShipmentWithPayload(shipment, params.Update)
	verrs, err := models.SaveShipmentAndAddresses(h.db, shipment)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	shipmentPayload := publicPayloadForShipmentModel(*shipment)
	return publicshipmentop.NewPatchShipmentOK().WithPayload(shipmentPayload)
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
