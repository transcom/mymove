package publicapi

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/services"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/rateengine"
	paperworkservice "github.com/transcom/mymove/pkg/services/paperwork"
	shipmentservice "github.com/transcom/mymove/pkg/services/shipment"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

func payloadForShipmentModel(s models.Shipment) *apimessages.Shipment {
	shipmentpayload := &apimessages.Shipment{
		ID:               *handlers.FmtUUID(s.ID),
		Status:           apimessages.ShipmentStatus(s.Status),
		SourceGbloc:      payloadForGBLOC(s.SourceGBLOC),
		DestinationGbloc: payloadForGBLOC(s.DestinationGBLOC),
		GblNumber:        s.GBLNumber,
		Market:           payloadForMarkets(s.Market),
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		UpdatedAt:        strfmt.DateTime(s.UpdatedAt),

		// associations
		TrafficDistributionListID: handlers.FmtUUIDPtr(s.TrafficDistributionListID),
		TrafficDistributionList:   payloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMemberID:           strfmt.UUID(s.ServiceMemberID.String()),
		ServiceMember:             payloadForServiceMemberModel(&s.ServiceMember),
		MoveID:                    strfmt.UUID(s.MoveID.String()),
		Move:                      payloadForMoveModel(&s.Move),

		// dates
		ActualPickupDate:     handlers.FmtDatePtr(s.ActualPickupDate),
		ActualPackDate:       handlers.FmtDatePtr(s.ActualPackDate),
		ActualDeliveryDate:   handlers.FmtDatePtr(s.ActualDeliveryDate),
		BookDate:             handlers.FmtDatePtr(s.BookDate),
		RequestedPickupDate:  handlers.FmtDatePtr(s.RequestedPickupDate),
		OriginalDeliveryDate: handlers.FmtDatePtr(s.OriginalDeliveryDate),
		OriginalPackDate:     handlers.FmtDatePtr(s.OriginalPackDate),

		// calculated durations
		EstimatedPackDays:    s.EstimatedPackDays,
		EstimatedTransitDays: s.EstimatedTransitDays,

		// addresses
		PickupAddress:                payloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:    handlers.FmtBool(s.HasSecondaryPickupAddress),
		SecondaryPickupAddress:       payloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:           handlers.FmtBool(s.HasDeliveryAddress),
		DeliveryAddress:              payloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress: handlers.FmtBool(s.HasPartialSITDeliveryAddress),
		PartialSitDeliveryAddress:    payloadForAddressModel(s.PartialSITDeliveryAddress),

		// weights
		WeightEstimate:              handlers.FmtPoundPtr(s.WeightEstimate),
		ProgearWeightEstimate:       handlers.FmtPoundPtr(s.ProgearWeightEstimate),
		SpouseProgearWeightEstimate: handlers.FmtPoundPtr(s.SpouseProgearWeightEstimate),
		NetWeight:                   handlers.FmtPoundPtr(s.NetWeight),
		GrossWeight:                 handlers.FmtPoundPtr(s.GrossWeight),
		TareWeight:                  handlers.FmtPoundPtr(s.TareWeight),

		// pre-move survey
		PmSurveyConductedDate:               handlers.FmtDatePtr(s.PmSurveyConductedDate),
		PmSurveyCompletedAt:                 handlers.FmtDateTimePtr(s.PmSurveyCompletedAt),
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(s.PmSurveyPlannedPackDate),
		PmSurveyPlannedPickupDate:           handlers.FmtDatePtr(s.PmSurveyPlannedPickupDate),
		PmSurveyPlannedDeliveryDate:         handlers.FmtDatePtr(s.PmSurveyPlannedDeliveryDate),
		PmSurveyWeightEstimate:              handlers.FmtPoundPtr(s.PmSurveyWeightEstimate),
		PmSurveyProgearWeightEstimate:       handlers.FmtPoundPtr(s.PmSurveyProgearWeightEstimate),
		PmSurveySpouseProgearWeightEstimate: handlers.FmtPoundPtr(s.PmSurveySpouseProgearWeightEstimate),
		PmSurveyNotes:                       s.PmSurveyNotes,
		PmSurveyMethod:                      s.PmSurveyMethod,
	}
	tspID := s.CurrentTransportationServiceProviderID()
	if tspID != uuid.Nil {
		shipmentpayload.TransportationServiceProviderID = *handlers.FmtUUID(tspID)
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
	var shipment *models.Shipment
	var err error
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
		if err != nil {
			h.Logger().Error("Error retrieving authenticated TSP user", zap.Error(err))
			return shipmentop.NewGetShipmentForbidden()
		}
		shipment, err = models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return shipmentop.NewGetShipmentForbidden()
		}
	} else if session.IsOfficeUser() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return shipmentop.NewGetShipmentForbidden()
		}
	} else {
		return shipmentop.NewGetShipmentForbidden()
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewGetShipmentOK().WithPayload(sp)
}

// GetShipmentInvoicesHandler returns all invoices for a shipment
type GetShipmentInvoicesHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment ID - returns list of associated invoices
func (h GetShipmentInvoicesHandler) Handle(params shipmentop.GetShipmentInvoicesParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	if !session.IsOfficeUser() {
		// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
		// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
		// See original commits in https://github.com/transcom/mymove/pull/802
		tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return shipmentop.NewGetShipmentInvoicesForbidden()
		}

		// Make sure TSP has access to this shipment
		_, err = models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return shipmentop.NewGetShipmentInvoicesForbidden()
		}
	}

	invoices, err := models.FetchInvoicesForShipment(h.DB(), shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewGetShipmentInvoicesBadRequest()
	}

	payload := payloadForInvoiceModels(invoices)
	return shipmentop.NewGetShipmentInvoicesOK().WithPayload(payload)
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

// TransportShipmentHandler allows a TSP to start transporting a particular shipment
type TransportShipmentHandler struct {
	handlers.HandlerContext
}

// Handle updates the shipment with pack and pickup dates and weights and puts it in-transit - checks that currently logged in user is authorized to act for the TSP assigned the shipment
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

	actualPackDate := (time.Time)(*params.Payload.ActualPackDate)

	err = shipment.Pack(actualPackDate)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	actualPickupDate := (time.Time)(*params.Payload.ActualPickupDate)

	err = shipment.Transport(actualPickupDate)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	shipment.NetWeight = handlers.PoundPtrFromInt64Ptr(params.Payload.NetWeight)

	if params.Payload.GrossWeight != nil {
		shipment.GrossWeight = handlers.PoundPtrFromInt64Ptr(params.Payload.GrossWeight)
	}

	if params.Payload.TareWeight != nil {
		shipment.TareWeight = handlers.PoundPtrFromInt64Ptr(params.Payload.TareWeight)
	}

	verrs, err := h.DB().ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewTransportShipmentOK().WithPayload(sp)
}

// DeliverShipmentHandler allows a TSP to start transporting a particular shipment
type DeliverShipmentHandler struct {
	handlers.HandlerContext
}

// Handle delivers the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h DeliverShipmentHandler) Handle(params shipmentop.DeliverShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewDeliverShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewDeliverShipmentBadRequest()
	}

	actualDeliveryDate := (time.Time)(*params.Payload.ActualDeliveryDate)
	engine := rateengine.NewRateEngine(h.DB(), h.Logger(), h.Planner())

	verrs, err := shipmentservice.DeliverAndPriceShipment{
		DB:     h.DB(),
		Engine: engine,
	}.Call(actualDeliveryDate, shipment)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewDeliverShipmentOK().WithPayload(sp)
}

// CompletePmSurveyHandler completes a pre-move survey for a particular shipment
type CompletePmSurveyHandler struct {
	handlers.HandlerContext
}

// Handle completes a pre-moves survey - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CompletePmSurveyHandler) Handle(params shipmentop.CompletePmSurveyParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewCompletePmSurveyForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewCompletePmSurveyBadRequest()
	}

	pmSurveyCompletedAt := time.Now()

	shipment.PmSurveyCompletedAt = &pmSurveyCompletedAt
	verrs, err := models.SaveShipment(h.DB(), shipment)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewCompletePmSurveyOK().WithPayload(sp)
}

func patchShipmentWithPayload(shipment *models.Shipment, payload *apimessages.Shipment) {

	// Comparing against the zero time allows users to set dates to nil via PATCH
	zeroTime := time.Time{}

	// PM Survey fields may be updated individually in the Dates panel and so cannot be lumped into one update
	if payload.PmSurveyConductedDate != nil {
		if zeroTime == (time.Time)(*payload.PmSurveyConductedDate) {
			shipment.PmSurveyConductedDate = nil
		} else {
			shipment.PmSurveyConductedDate = (*time.Time)(payload.PmSurveyConductedDate)
		}
	}

	if payload.PmSurveyPlannedDeliveryDate != nil {
		if zeroTime == (time.Time)(*payload.PmSurveyPlannedDeliveryDate) {
			shipment.PmSurveyPlannedDeliveryDate = nil
		} else {
			shipment.PmSurveyPlannedDeliveryDate = (*time.Time)(payload.PmSurveyPlannedDeliveryDate)
		}
	}

	if payload.PmSurveyMethod != "" {
		shipment.PmSurveyMethod = payload.PmSurveyMethod
	}

	if payload.PmSurveyPlannedPackDate != nil {
		if zeroTime == (time.Time)(*payload.PmSurveyPlannedPackDate) {
			shipment.PmSurveyPlannedPackDate = nil
		} else {
			shipment.PmSurveyPlannedPackDate = (*time.Time)(payload.PmSurveyPlannedPackDate)
		}
	}

	if payload.PmSurveyPlannedPickupDate != nil {
		if zeroTime == (time.Time)(*payload.PmSurveyPlannedPickupDate) {
			shipment.PmSurveyPlannedPickupDate = nil
		} else {
			shipment.PmSurveyPlannedPickupDate = (*time.Time)(payload.PmSurveyPlannedPickupDate)
		}
	}

	if payload.PmSurveyNotes != nil {
		shipment.PmSurveyNotes = payload.PmSurveyNotes
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

	if payload.NetWeight != nil {
		shipment.NetWeight = handlers.PoundPtrFromInt64Ptr(payload.NetWeight)
	}

	if payload.GrossWeight != nil {
		shipment.GrossWeight = handlers.PoundPtrFromInt64Ptr(payload.GrossWeight)
	}

	if payload.TareWeight != nil {
		shipment.TareWeight = handlers.PoundPtrFromInt64Ptr(payload.TareWeight)
	}

	if payload.ActualPickupDate != nil {
		if zeroTime == (time.Time)(*payload.ActualPickupDate) {
			shipment.ActualPickupDate = nil
		} else {
			shipment.ActualPickupDate = (*time.Time)(payload.ActualPickupDate)
		}
	}

	if payload.ActualPackDate != nil {
		if zeroTime == (time.Time)(*payload.ActualPackDate) {
			shipment.ActualPackDate = nil
		} else {
			shipment.ActualPackDate = (*time.Time)(payload.ActualPackDate)
		}
	}

	if payload.ActualDeliveryDate != nil {
		if zeroTime == (time.Time)(*payload.ActualDeliveryDate) {
			shipment.ActualDeliveryDate = nil
		} else {
			shipment.ActualDeliveryDate = (*time.Time)(payload.ActualDeliveryDate)
		}
	}

	if payload.PickupAddress != nil {
		if shipment.PickupAddress == nil {
			shipment.PickupAddress = addressModelFromPayload(payload.PickupAddress)
		} else {
			updateAddressWithPayload(shipment.PickupAddress, payload.PickupAddress)
		}
	}
	if payload.HasSecondaryPickupAddress != nil {
		if *payload.HasSecondaryPickupAddress == false {
			shipment.SecondaryPickupAddress = nil
		} else if *payload.HasSecondaryPickupAddress == true {
			if payload.SecondaryPickupAddress != nil {
				if shipment.SecondaryPickupAddress == nil {
					shipment.SecondaryPickupAddress = addressModelFromPayload(payload.SecondaryPickupAddress)
				} else {
					updateAddressWithPayload(shipment.SecondaryPickupAddress, payload.SecondaryPickupAddress)
				}
			}
		}
		shipment.HasSecondaryPickupAddress = *payload.HasSecondaryPickupAddress
	}

	if payload.HasDeliveryAddress != nil {
		if *payload.HasDeliveryAddress == false {
			shipment.DeliveryAddress = nil
		} else if *payload.HasDeliveryAddress == true {
			if payload.DeliveryAddress != nil {
				if shipment.DeliveryAddress == nil {
					shipment.DeliveryAddress = addressModelFromPayload(payload.DeliveryAddress)
				} else {
					updateAddressWithPayload(shipment.DeliveryAddress, payload.DeliveryAddress)
				}
			}
		}
		shipment.HasDeliveryAddress = *payload.HasDeliveryAddress
	}

	if payload.HasPartialSitDeliveryAddress != nil {
		if *payload.HasPartialSitDeliveryAddress == false {
			shipment.PartialSITDeliveryAddress = nil
		} else if *payload.HasPartialSitDeliveryAddress == true {
			if payload.PartialSitDeliveryAddress != nil {
				if shipment.PartialSITDeliveryAddress == nil {
					shipment.PartialSITDeliveryAddress = addressModelFromPayload(payload.PartialSitDeliveryAddress)
				} else {
					updateAddressWithPayload(shipment.PartialSITDeliveryAddress, payload.PartialSitDeliveryAddress)
				}
			}
		}
		shipment.HasPartialSITDeliveryAddress = *payload.HasPartialSitDeliveryAddress
	}
}

// PatchShipmentHandler allows a TSP to refuse a particular shipment
type PatchShipmentHandler struct {
	handlers.HandlerContext
}

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {
	var shipment *models.Shipment
	var err error
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// authorization
	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
		if err != nil {
			h.Logger().Error("Error retrieving authenticated TSP user", zap.Error(err))
			return shipmentop.NewGetShipmentForbidden()
		}
		shipment, err = models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return shipmentop.NewPatchShipmentBadRequest()
		}
	} else if session.IsOfficeUser() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return shipmentop.NewPatchShipmentBadRequest()
		}
	} else {
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

// CreateGovBillOfLadingHandler creates a GBL PDF & uploads it as a document associated to a move doc, shipment and move
type CreateGovBillOfLadingHandler struct {
	handlers.HandlerContext
	createForm services.FormCreator
}

// Handle generates the GBL PDF & uploads it as a document associated to a move doc, shipment and move
func (h CreateGovBillOfLadingHandler) Handle(params shipmentop.CreateGovBillOfLadingParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the TSP user is authorized to update move doc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	tspUser, shipment, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
	if err != nil {
		if err.Error() == "USER_UNAUTHORIZED" {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		if err.Error() == "FETCH_FORBIDDEN" {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		return handlers.ResponseForError(h.Logger(), err)
	}
	// Don't allow GBL generation for shipments that already have a GBL move document
	extantGBLS, _ := models.FetchMoveDocumentsByTypeForShipment(h.DB(), session, models.MoveDocumentTypeGOVBILLOFLADING, shipmentID)
	if len(extantGBLS) > 0 {
		return handlers.ResponseForCustomErrors(h.Logger(), fmt.Errorf("there is already a Bill of Lading for this shipment"), http.StatusBadRequest)
	}

	// Don't allow GBL generation for incomplete orders
	orders, ordersErr := models.FetchOrder(h.DB(), shipment.Move.OrdersID)
	if ordersErr != nil {
		return handlers.ResponseForError(h.Logger(), ordersErr)
	}
	if orders.IsCompleteForGBL() != true {
		return handlers.ResponseForCustomErrors(h.Logger(), fmt.Errorf("the move is missing some information from the JPPSO. Please contact the JPPSO"), http.StatusExpectationFailed)
	}

	// Create PDF for GBL
	gbl, err := models.FetchGovBillOfLadingFormValues(h.DB(), shipmentID)
	if err != nil {
		// TODO: (andrea) Pass info of exactly what is missing in custom error message
		h.Logger().Error("Failed retrieving the GBL data.", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingExpectationFailed()
	}
	formLayout := paperwork.Form1203Layout

	template, err := paperworkservice.MakeFormTemplate(gbl, gbl.GBLNumber1, formLayout, services.GBL)
	if err != nil {
		h.Logger().Error(errors.Cause(err).Error(), zap.Error(errors.Cause(err)))
	}

	gblFile, err := h.createForm.CreateForm(template)
	if err != nil {
		h.Logger().Error(errors.Cause(err).Error(), zap.Error(errors.Cause(err)))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	uploader := uploaderpkg.NewUploader(h.DB(), h.Logger(), h.FileStorer())
	upload, verrs, err := uploader.CreateUpload(nil, *tspUser.UserID, gblFile)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	uploads := []models.Upload{*upload}

	// Create GBL move document associated to the shipment
	doc, verrs, err := shipment.Move.CreateMoveDocument(h.DB(),
		uploads,
		&shipmentID,
		models.MoveDocumentTypeGOVBILLOFLADING,
		string("Government Bill Of Lading"),
		swag.String(""),
		models.SelectedMoveType(apimessages.SelectedMoveTypeHHG),
	)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	documentPayload, err := payloadForDocumentModel(h.FileStorer(), doc.Document)
	if err != nil {
		h.Logger().Error("Error fetching document for gbl doc", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	moveDocumentPayload := &apimessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(doc.ID),
		ShipmentID:       handlers.FmtUUIDPtr(doc.ShipmentID),
		Document:         documentPayload,
		Title:            handlers.FmtStringPtr(&doc.Title),
		MoveDocumentType: apimessages.MoveDocumentType(doc.MoveDocumentType),
		Status:           apimessages.MoveDocumentStatus(doc.Status),
		Notes:            handlers.FmtStringPtr(doc.Notes),
	}

	return shipmentop.NewCreateGovBillOfLadingCreated().WithPayload(moveDocumentPayload)
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
