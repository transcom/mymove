package internalapi

import (
	"time"

	"github.com/facebookgo/clock"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	invoiceop "github.com/transcom/mymove/pkg/services/invoice"
)

func payloadForInvoiceModel(a *models.Invoice) *internalmessages.Invoice {
	if a == nil {
		return nil
	}

	return &internalmessages.Invoice{
		ID:                *handlers.FmtUUID(a.ID),
		ShipmentID:        *handlers.FmtUUID(a.ShipmentID),
		InvoiceNumber:     a.InvoiceNumber,
		ApproverFirstName: a.Approver.FirstName,
		ApproverLastName:  a.Approver.LastName,
		Status:            internalmessages.InvoiceStatus(a.Status),
		InvoicedDate:      *handlers.FmtDateTime(a.InvoicedDate),
	}
}

func payloadForShipmentModel(s models.Shipment) (*internalmessages.Shipment, error) {
	// TODO: For now, we keep the Shipment structure the same but change where the CodeOfService
	// TODO: is coming from.  Ultimately we should probably rework the structure below to more
	// TODO: closely match the database structure.
	var codeOfService *string
	if s.TrafficDistributionList != nil {
		codeOfService = &s.TrafficDistributionList.CodeOfService
	}

	var moveDatesSummary internalmessages.ShipmentMoveDatesSummary
	if s.RequestedPickupDate != nil && s.EstimatedPackDays != nil && s.EstimatedTransitDays != nil {
		summary, err := calculateMoveDatesFromShipment(&s)
		if err != nil {
			return nil, err
		}
		moveDatesSummary = internalmessages.ShipmentMoveDatesSummary{
			Pack:     handlers.FmtDateSlice(summary.PackDays),
			Pickup:   handlers.FmtDateSlice(summary.PickupDays),
			Transit:  handlers.FmtDateSlice(summary.TransitDays),
			Delivery: handlers.FmtDateSlice(summary.DeliveryDays),
		}
	}

	shipmentPayload := &internalmessages.Shipment{
		ID:               strfmt.UUID(s.ID.String()),
		Status:           internalmessages.ShipmentStatus(s.Status),
		SourceGbloc:      payloadForGBLOC(s.SourceGBLOC),
		DestinationGbloc: payloadForGBLOC(s.DestinationGBLOC),
		Market:           payloadForMarkets(s.Market),
		CodeOfService:    codeOfService,
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		UpdatedAt:        strfmt.DateTime(s.UpdatedAt),

		// associations
		TrafficDistributionListID: handlers.FmtUUIDPtr(s.TrafficDistributionListID),
		TrafficDistributionList:   payloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMemberID:           strfmt.UUID(s.ServiceMemberID.String()),
		MoveID:                    strfmt.UUID(s.MoveID.String()),

		// dates
		ActualPickupDate:     handlers.FmtDatePtr(s.ActualPickupDate),
		ActualPackDate:       handlers.FmtDatePtr(s.ActualPackDate),
		ActualDeliveryDate:   handlers.FmtDatePtr(s.ActualDeliveryDate),
		BookDate:             handlers.FmtDatePtr(s.BookDate),
		RequestedPickupDate:  handlers.FmtDatePtr(s.RequestedPickupDate),
		OriginalDeliveryDate: handlers.FmtDatePtr(s.OriginalDeliveryDate),
		OriginalPackDate:     handlers.FmtDatePtr(s.OriginalPackDate),
		MoveDatesSummary:     &moveDatesSummary,

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
		shipmentPayload.TransportationServiceProviderID = *handlers.FmtUUID(tspID)
	}
	return shipmentPayload, nil
}

// CreateShipmentHandler creates a Shipment
type CreateShipmentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CreateShipmentHandler) Handle(params shipmentop.CreateShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := params.Shipment

	pickupAddress := addressModelFromPayload(payload.PickupAddress)
	secondaryPickupAddress := addressModelFromPayload(payload.SecondaryPickupAddress)
	deliveryAddress := addressModelFromPayload(payload.DeliveryAddress)
	partialSITDeliveryAddress := addressModelFromPayload(payload.PartialSitDeliveryAddress)
	market := "dHHG"

	var requestedPickupDate *time.Time
	if payload.RequestedPickupDate != nil {
		date := time.Time(*payload.RequestedPickupDate)
		requestedPickupDate = &date
	}

	hasSecondaryPickupAddress := false
	if payload.HasSecondaryPickupAddress != nil {
		hasSecondaryPickupAddress = *payload.HasSecondaryPickupAddress
	}

	hasDeliveryAddress := false
	if payload.HasDeliveryAddress != nil {
		hasDeliveryAddress = *payload.HasDeliveryAddress
	}

	hasPartialSitDeliveryAddress := false
	if payload.HasPartialSitDeliveryAddress != nil {
		hasPartialSitDeliveryAddress = *payload.HasPartialSitDeliveryAddress
	}

	newShipment := models.Shipment{
		MoveID:                       move.ID,
		ServiceMemberID:              session.ServiceMemberID,
		Status:                       models.ShipmentStatusDRAFT,
		RequestedPickupDate:          requestedPickupDate,
		EstimatedPackDays:            payload.EstimatedPackDays,
		EstimatedTransitDays:         payload.EstimatedTransitDays,
		WeightEstimate:               handlers.PoundPtrFromInt64Ptr(payload.WeightEstimate),
		ProgearWeightEstimate:        handlers.PoundPtrFromInt64Ptr(payload.ProgearWeightEstimate),
		SpouseProgearWeightEstimate:  handlers.PoundPtrFromInt64Ptr(payload.SpouseProgearWeightEstimate),
		PickupAddress:                pickupAddress,
		HasSecondaryPickupAddress:    hasSecondaryPickupAddress,
		SecondaryPickupAddress:       secondaryPickupAddress,
		HasDeliveryAddress:           hasDeliveryAddress,
		DeliveryAddress:              deliveryAddress,
		HasPartialSITDeliveryAddress: hasPartialSitDeliveryAddress,
		PartialSITDeliveryAddress:    partialSITDeliveryAddress,
		Market:                       &market,
	}
	if err = updateShipmentDatesWithPayload(h, &newShipment, params.Shipment); err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	verrs, err := models.SaveShipmentAndAddresses(h.DB(), &newShipment)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	shipmentPayload, err := payloadForShipmentModel(newShipment)
	if err != nil {
		h.Logger().Error("Error in shipment payload: ", zap.Error(err))
	}

	return shipmentop.NewCreateShipmentCreated().WithPayload(shipmentPayload)
}

func patchShipmentWithPremoveSurveyFields(shipment *models.Shipment, payload *internalmessages.Shipment) {
	// Premove Survey values entered by TSP agent
	requiredValue := payload.PmSurveyPlannedPackDate

	// If any PmSurvey data was sent, update all fields
	// This takes advantage of the fact that all PmSurvey data is updated at once and allows us to null out optional fields
	if requiredValue != nil {
		shipment.PmSurveyPlannedPackDate = (*time.Time)(payload.PmSurveyPlannedPackDate)
		shipment.PmSurveyConductedDate = (*time.Time)(payload.PmSurveyConductedDate)
		shipment.PmSurveyPlannedPickupDate = (*time.Time)(payload.PmSurveyPlannedPickupDate)
		shipment.PmSurveyPlannedDeliveryDate = (*time.Time)(payload.PmSurveyPlannedDeliveryDate)
		shipment.PmSurveyNotes = payload.PmSurveyNotes
		shipment.PmSurveyMethod = payload.PmSurveyMethod
		shipment.PmSurveyProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyProgearWeightEstimate)
		shipment.PmSurveySpouseProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveySpouseProgearWeightEstimate)
		shipment.PmSurveyWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyWeightEstimate)
	}
}

func patchShipmentWithPayload(shipment *models.Shipment, payload *internalmessages.Shipment) {

	if payload.ActualPickupDate != nil {
		shipment.ActualPickupDate = (*time.Time)(payload.ActualPickupDate)
	}
	if payload.ActualPackDate != nil {
		shipment.ActualPackDate = (*time.Time)(payload.ActualPackDate)
	}
	if payload.RequestedPickupDate != nil {
		shipment.RequestedPickupDate = (*time.Time)(payload.RequestedPickupDate)
	}
	if payload.EstimatedPackDays != nil {
		shipment.EstimatedPackDays = payload.EstimatedPackDays
	}
	if payload.EstimatedTransitDays != nil {
		shipment.EstimatedTransitDays = payload.EstimatedTransitDays
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

	if payload.WeightEstimate != nil {
		shipment.WeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.WeightEstimate)
	}
	if payload.ProgearWeightEstimate != nil {
		shipment.ProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.ProgearWeightEstimate)
	}
	if payload.SpouseProgearWeightEstimate != nil {
		shipment.SpouseProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.SpouseProgearWeightEstimate)
	}
}

// PatchShipmentHandler Patchs an HHG
type PatchShipmentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.DB(), session, shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	patchShipmentWithPayload(shipment, params.Shipment)
	if err = updateShipmentDatesWithPayload(h, shipment, params.Shipment); err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Premove survey info can only be edited by office users or TSPs
	if session.IsOfficeUser() {
		patchShipmentWithPremoveSurveyFields(shipment, params.Shipment)
	}

	verrs, err := models.SaveShipmentAndAddresses(h.DB(), shipment)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	shipmentPayload, err := payloadForShipmentModel(*shipment)
	if err != nil {
		h.Logger().Error("Error in shipment payload: ", zap.Error(err))
	}

	return shipmentop.NewPatchShipmentOK().WithPayload(shipmentPayload)
}

func updateShipmentDatesWithPayload(h handlers.HandlerContext, shipment *models.Shipment, payload *internalmessages.Shipment) error {
	if payload.RequestedPickupDate == nil {
		return nil
	}

	moveDate := time.Time(*payload.RequestedPickupDate)

	summary, err := calculateMoveDatesFromMove(h.DB(), h.Planner(), shipment.MoveID, moveDate)
	if err != nil {
		return nil
	}

	packDays := int64(len(summary.PackDays))
	shipment.EstimatedPackDays = &packDays

	transitDays := int64(len(summary.TransitDays))
	shipment.EstimatedTransitDays = &transitDays

	deliveryDate := summary.DeliveryDays[0]
	shipment.OriginalDeliveryDate = &deliveryDate
	packDate := summary.PackDays[0]
	shipment.OriginalPackDate = &packDate

	return nil
}

// GetShipmentHandler Returns an HHG
type GetShipmentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h GetShipmentHandler) Handle(params shipmentop.GetShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.DB(), session, shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	shipmentPayload, err := payloadForShipmentModel(*shipment)
	if err != nil {
		h.Logger().Error("Error in shipment payload: ", zap.Error(err))
	}

	return shipmentop.NewGetShipmentOK().WithPayload(shipmentPayload)
}

// ApproveHHGHandler approves an HHG
type ApproveHHGHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h ApproveHHGHandler) Handle(params shipmentop.ApproveHHGParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return shipmentop.NewApproveHHGForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.DB(), session, shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	err = shipment.Approve()
	if err != nil {
		h.Logger().Error("Attempted to approve HHG, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
		return handlers.ResponseForError(h.Logger(), err)
	}
	verrs, err := h.DB().ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	shipmentPayload, err := payloadForShipmentModel(*shipment)
	if err != nil {
		h.Logger().Error("Error in shipment payload: ", zap.Error(err))
	}

	return shipmentop.NewApproveHHGOK().WithPayload(shipmentPayload)
}

// CompleteHHGHandler completes an HHG
type CompleteHHGHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CompleteHHGHandler) Handle(params shipmentop.CompleteHHGParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return shipmentop.NewCompleteHHGForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	shipment, err := models.FetchShipment(h.DB(), session, shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	err = shipment.Complete()
	if err != nil {
		h.Logger().Error("Attempted to complete HHG, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
		return handlers.ResponseForError(h.Logger(), err)
	}
	verrs, err := h.DB().ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	shipmentPayload, err := payloadForShipmentModel(*shipment)
	if err != nil {
		h.Logger().Error("Error in shipment payload: ", zap.Error(err))
	}

	return shipmentop.NewCompleteHHGOK().WithPayload(shipmentPayload)
}

// ShipmentInvoiceHandler sends an invoice through GEX to Syncada
type ShipmentInvoiceHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h ShipmentInvoiceHandler) Handle(params shipmentop.CreateAndSendHHGInvoiceParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return shipmentop.NewCreateAndSendHHGInvoiceForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	shipment, err := invoiceop.FetchShipmentForInvoice{DB: h.DB()}.Call(shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	if shipment.Status != models.ShipmentStatusDELIVERED && shipment.Status != models.ShipmentStatusCOMPLETED {
		h.Logger().Error("Shipment status not in delivered state.")
		return shipmentop.NewCreateAndSendHHGInvoicePreconditionFailed()
	}

	//for now we limit a shipment to 1 invoice
	//if invoices exists and at least one is either in process or has succeeded then return 409
	existingInvoices, err := models.FetchInvoicesForShipment(h.DB(), shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	for _, invoice := range existingInvoices {
		//if an invoice has started, is in process or has been submitted successfully then throw err
		if invoice.Status != models.InvoiceStatusSUBMISSIONFAILURE {
			payload := payloadForInvoiceModel(&invoice)
			return shipmentop.NewCreateAndSendHHGInvoiceConflict().WithPayload(payload)
		}
	}

	approver, err := models.FetchOfficeUserByID(h.DB(), session.OfficeUserID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	// before processing the invoice, save it in an in process state
	var invoice models.Invoice
	verrs, err := invoiceop.CreateInvoice{DB: h.DB(), Clock: clock.New()}.Call(*approver, &invoice, shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	invoice858CString, verrs, err := invoiceop.ProcessInvoice{
		DB:                    h.DB(),
		GexSender:             h.GexSender(),
		SendProductionInvoice: h.SendProductionInvoice(),
		ICNSequencer:          h.ICNSequencer(),
	}.Call(&invoice, shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	// Send invoice to S3 for storage if response from GEX is successful
	fs := h.FileStorer()
	verrs, err = invoiceop.StoreInvoice858C{
		DB:     h.DB(),
		Logger: h.Logger(),
		Storer: &fs,
	}.Call(*invoice858CString, &invoice, session.UserID)
	if verrs.HasAny() {
		h.Logger().Error("Failed to store invoice record to s3, with validation errors", zap.Error(verrs))
	}
	if err != nil {
		h.Logger().Error("Failed to store invoice record to s3, with error", zap.Error(err))
	}

	payload := payloadForInvoiceModel(&invoice)

	return shipmentop.NewCreateAndSendHHGInvoiceOK().WithPayload(payload)
}
