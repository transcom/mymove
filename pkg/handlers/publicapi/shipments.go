package publicapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func payloadForShipmentModel(s models.Shipment) *apimessages.Shipment {
	shipmentpayload := &apimessages.Shipment{
		ID: *handlers.FmtUUID(s.ID),
		TrafficDistributionList:             payloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMember:                       payloadForServiceMemberModel(s.ServiceMember),
		ActualPickupDate:                    *handlers.FmtDateTimePtr(s.ActualPickupDate),
		ActualDeliveryDate:                  *handlers.FmtDateTimePtr(s.ActualDeliveryDate),
		CreatedAt:                           strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                           strfmt.DateTime(s.UpdatedAt),
		SourceGbloc:                         apimessages.GBLOC(*s.SourceGBLOC),
		DestinationGbloc:                    apimessages.GBLOC(*s.DestinationGBLOC),
		Market:                              apimessages.ShipmentMarket(*s.Market),
		BookDate:                            *handlers.FmtDatePtr(s.BookDate),
		RequestedPickupDate:                 *handlers.FmtDateTimePtr(s.RequestedPickupDate),
		Move:                                payloadForMoveModel(&s.Move),
		Status:                              apimessages.ShipmentStatus(s.Status),
		EstimatedPackDays:                   handlers.FmtInt64(*s.EstimatedPackDays),
		EstimatedTransitDays:                handlers.FmtInt64(*s.EstimatedTransitDays),
		PickupAddress:                       payloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:           *handlers.FmtBool(s.HasSecondaryPickupAddress),
		SecondaryPickupAddress:              payloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:                  *handlers.FmtBool(s.HasDeliveryAddress),
		DeliveryAddress:                     payloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress:        *handlers.FmtBool(s.HasPartialSITDeliveryAddress),
		PartialSitDeliveryAddress:           payloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:                      handlers.FmtInt64(s.WeightEstimate.Int64()),
		ProgearWeightEstimate:               handlers.FmtInt64(s.ProgearWeightEstimate.Int64()),
		SpouseProgearWeightEstimate:         handlers.FmtInt64(s.SpouseProgearWeightEstimate.Int64()),
		ActualWeight:                        handlers.FmtPoundPtr(s.ActualWeight),
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(s.PmSurveyPlannedPackDate),
		PmSurveyPlannedPickupDate:           handlers.FmtDatePtr(s.PmSurveyPlannedPickupDate),
		PmSurveyPlannedDeliveryDate:         handlers.FmtDatePtr(s.PmSurveyPlannedDeliveryDate),
		PmSurveyWeightEstimate:              handlers.FmtPoundPtr(s.PmSurveyWeightEstimate),
		PmSurveyProgearWeightEstimate:       handlers.FmtPoundPtr(s.PmSurveyProgearWeightEstimate),
		PmSurveySpouseProgearWeightEstimate: handlers.FmtPoundPtr(s.PmSurveySpouseProgearWeightEstimate),
		PmSurveyNotes:                       s.PmSurveyNotes,
		PmSurveyMethod:                      s.PmSurveyMethod,
	}
	return shipmentpayload
}

// IndexShipmentsHandler returns a list of shipments
type IndexShipmentsHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a list of all shipments
func (h IndexShipmentsHandler) Handle(params shipmentop.IndexShipmentsParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewIndexShipmentsForbidden()
	}

	shipments, err := models.FetchShipmentsByTSP(h.DB(), tspUser.TransportationServiceProviderID,
		params.Status, params.OrderBy, params.Limit, params.Offset)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewIndexShipmentsBadRequest()
	}

	isp := make(apimessages.IndexShipments, len(shipments))
	for i, s := range shipments {
		isp[i] = payloadForShipmentModel(s)
	}
	return shipmentop.NewIndexShipmentsOK().WithPayload(isp)
}

// GetShipmentHandler returns a particular shipment
type GetShipmentHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified shipment
func (h GetShipmentHandler) Handle(params shipmentop.GetShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewGetShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewGetShipmentBadRequest()
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewGetShipmentOK().WithPayload(sp)
}

// AcceptShipmentHandler allows a TSP to accept a particular shipment
type AcceptShipmentHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h AcceptShipmentHandler) Handle(params shipmentop.AcceptShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_08_22) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewAcceptShipmentForbidden()
	}

	// Accept the shipment
	shipment, shipmentOffer, verrs, err := models.AcceptShipmentForTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil || verrs.HasAny() {
		if err == models.ErrFetchNotFound {
			h.Logger().Error("DB Query", zap.Error(err))
			return shipmentop.NewAcceptShipmentBadRequest()
		} else if err == models.ErrInvalidTransition {
			h.Logger().Info("Attempted to accept shipment, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
			h.Logger().Info("Attempted to accept shipment offer, got invalid transition", zap.Error(err), zap.Bool("shipment_offer_accepted", *shipmentOffer.Accepted))
			return shipmentop.NewAcceptShipmentConflict()
		} else {
			h.Logger().Error("Unknown Error", zap.Error(err))
			return handlers.ResponseForVErrors(h.Logger(), verrs, err)
		}
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewAcceptShipmentOK().WithPayload(sp)
}

// RejectShipmentHandler allows a TSP to refuse a particular shipment
type RejectShipmentHandler struct {
	handlers.HandlerContext
}

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h RejectShipmentHandler) Handle(params shipmentop.RejectShipmentParams) middleware.Responder {
	// set reason, set thing
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_08_22) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewRejectShipmentForbidden()
	}

	// Reject the shipment
	shipment, shipmentOffer, verrs, err := models.RejectShipmentForTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID, *params.Payload.Reason)
	if err != nil || verrs.HasAny() {
		if err == models.ErrFetchNotFound {
			h.Logger().Error("DB Query", zap.Error(err))
			return shipmentop.NewRejectShipmentBadRequest()
		} else if err == models.ErrInvalidTransition {
			h.Logger().Info("Attempted to reject shipment, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
			h.Logger().Info("Attempted to reject shipment offer, got invalid transition", zap.Error(err), zap.Bool("shipment_offer_accepted", *shipmentOffer.Accepted))
			return shipmentop.NewRejectShipmentConflict()
		} else {
			h.Logger().Error("Unknown Error", zap.Error(err))
			return handlers.ResponseForVErrors(h.Logger(), verrs, err)
		}
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewRejectShipmentOK().WithPayload(sp)
}

// TransportShipmentHandler allows a TSP to start transporting a particular shipment
type TransportShipmentHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h TransportShipmentHandler) Handle(params shipmentop.TransportShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewTransportShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewTransportShipmentBadRequest()
	}

	actualPickupDate := (time.Time)(*params.ActualPickupDate.ActualPickupDate)

	err = shipment.Transport(actualPickupDate)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	verrs, err := h.DB().ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewTransportShipmentOK().WithPayload(sp)
}

func patchShipmentWithPayload(shipment *models.Shipment, payload *apimessages.Shipment) {
	// Premove Survey values entered by TSP agent
	requiredValue := payload.PmSurveyPlannedPackDate

	// If any PmSurvey data was sent, update all fields
	// This takes advantage of the fact that all PmSurvey data is updated at once and allows us to null out optional fields
	if requiredValue != nil {
		shipment.PmSurveyPlannedDeliveryDate = (*time.Time)(payload.PmSurveyPlannedDeliveryDate)
		shipment.PmSurveyNotes = payload.PmSurveyNotes
		shipment.PmSurveyMethod = payload.PmSurveyMethod
		shipment.PmSurveyPlannedPackDate = (*time.Time)(payload.PmSurveyPlannedPackDate)
		shipment.PmSurveyPlannedPickupDate = (*time.Time)(payload.PmSurveyPlannedPickupDate)
		shipment.PmSurveyProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyProgearWeightEstimate)
		shipment.PmSurveySpouseProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveySpouseProgearWeightEstimate)
		shipment.PmSurveyWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyWeightEstimate)
	}

	if payload.ActualWeight != nil {
		shipment.ActualWeight = handlers.PoundPtrFromInt64Ptr(payload.ActualWeight)
	}

	if &payload.ActualPickupDate != nil {
		shipment.ActualPickupDate = (*time.Time)(&payload.ActualPickupDate)
	}
}

// PatchShipmentHandler allows a TSP to refuse a particular shipment
type PatchShipmentHandler struct {
	handlers.HandlerContext
}

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewPatchShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewPatchShipmentBadRequest()
	}

	patchShipmentWithPayload(shipment, params.Update)
	verrs, err := models.SaveShipmentAndAddresses(h.DB(), shipment)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	shipmentPayload := payloadForShipmentModel(*shipment)
	return shipmentop.NewPatchShipmentOK().WithPayload(shipmentPayload)
}

// GetShipmentContactDetailsHandler allows a TSP to accept a particular shipment
type GetShipmentContactDetailsHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentContactDetailsHandler) Handle(p shipmentop.GetShipmentContactDetailsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentClaimsHandler allows a TSP to accept a particular shipment
type GetShipmentClaimsHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentClaimsHandler) Handle(p shipmentop.GetShipmentClaimsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentDocumentsHandler allows a TSP to accept a particular shipment
type GetShipmentDocumentsHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentDocumentsHandler) Handle(p shipmentop.GetShipmentDocumentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// CreateShipmentDocumentHandler allows a TSP to accept a particular shipment
type CreateShipmentDocumentHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentDocumentHandler) Handle(p shipmentop.CreateShipmentDocumentParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}
