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
		PickupDate:                          *handlers.FmtDateTimePtr(s.PickupDate),
		DeliveryDate:                        *handlers.FmtDateTimePtr(s.DeliveryDate),
		CreatedAt:                           strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                           strfmt.DateTime(s.UpdatedAt),
		SourceGbloc:                         apimessages.GBLOC(*s.SourceGBLOC),
		DestinationGbloc:                    apimessages.GBLOC(*s.DestinationGBLOC),
		Market:                              apimessages.ShipmentMarket(*s.Market),
		BookDate:                            *handlers.FmtDatePtr(s.BookDate),
		RequestedPickupDate:                 *handlers.FmtDateTimePtr(s.RequestedPickupDate),
		Move:                                payloadForMoveModel(s.Move),
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
		PmSurveyPackDate:                    handlers.FmtDatePtr(s.PmSurveyPackDate),
		PmSurveyPickupDate:                  handlers.FmtDatePtr(s.PmSurveyPickupDate),
		PmSurveyLatestPickupDate:            handlers.FmtDatePtr(s.PmSurveyLatestPickupDate),
		PmSurveyEarliestDeliveryDate:        handlers.FmtDatePtr(s.PmSurveyEarliestDeliveryDate),
		PmSurveyLatestDeliveryDate:          handlers.FmtDatePtr(s.PmSurveyLatestDeliveryDate),
		PmSurveyWeightEstimate:              handlers.FmtPoundPtr(s.PmSurveyWeightEstimate),
		PmSurveyProgearWeightEstimate:       handlers.FmtPoundPtr(s.PmSurveyProgearWeightEstimate),
		PmSurveySpouseProgearWeightEstimate: handlers.FmtPoundPtr(s.PmSurveySpouseProgearWeightEstimate),
		PmSurveyNotes:                       s.PmSurveyNotes,
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

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger().Error("Missing TSP User ID")
		return shipmentop.NewIndexShipmentsForbidden()
	}

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

	shipmentID, _ := uuid.FromString(params.ShipmentUUID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger().Error("Missing TSP User ID")
		return shipmentop.NewGetShipmentForbidden()
	}

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

// CreateShipmentAcceptHandler allows a TSP to accept a particular shipment
type CreateShipmentAcceptHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentAcceptHandler) Handle(params shipmentop.CreateShipmentAcceptParams) middleware.Responder {
	return middleware.NotImplemented("operation .acceptShipment has not yet been implemented")
}

// CreateShipmentRejectHandler allows a TSP to refuse a particular shipment
type CreateShipmentRejectHandler struct {
	handlers.HandlerContext
}

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentRejectHandler) Handle(params shipmentop.CreateShipmentRejectParams) middleware.Responder {
	return middleware.NotImplemented("operation .refuseShipment has not yet been implemented")
}

func patchShipmentWithPayload(shipment *models.Shipment, payload *apimessages.Shipment) {
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
		shipment.PmSurveyProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyProgearWeightEstimate)
	}
	if payload.PmSurveySpouseProgearWeightEstimate != nil {
		shipment.PmSurveySpouseProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveySpouseProgearWeightEstimate)
	}
	if payload.PmSurveyWeightEstimate != nil {
		shipment.PmSurveyWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyWeightEstimate)
	}

}

// PatchShipmentHandler allows a TSP to refuse a particular shipment
type PatchShipmentHandler struct {
	handlers.HandlerContext
}

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentUUID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger().Error("Missing TSP User ID")
		return shipmentop.NewGetShipmentForbidden()
	}

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
