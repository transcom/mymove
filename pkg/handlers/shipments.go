package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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
		TrafficDistributionListID:           fmtUUIDPtr(s.TrafficDistributionListID),
		ServiceMemberID:                     strfmt.UUID(s.ServiceMemberID.String()),
		SourceGbloc:                         s.SourceGBLOC,
		DestinationGbloc:                    s.DestinationGBLOC,
		Market:                              s.Market,
		CodeOfService:                       s.CodeOfService,
		Status:                              s.Status,
		BookDate:                            fmtDatePtr(s.BookDate),
		RequestedPickupDate:                 fmtDatePtr(s.RequestedPickupDate),
		PickupDate:                          fmtDatePtr(s.PickupDate),
		DeliveryDate:                        fmtDatePtr(s.DeliveryDate),
		CreatedAt:                           strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                           strfmt.DateTime(s.UpdatedAt),
		EstimatedPackDays:                   s.EstimatedPackDays,
		EstimatedTransitDays:                s.EstimatedTransitDays,
		PickupAddress:                       payloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:           s.HasSecondaryPickupAddress,
		SecondaryPickupAddress:              payloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:                  s.HasDeliveryAddress,
		DeliveryAddress:                     payloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress:        s.HasPartialSITDeliveryAddress,
		PartialSitDeliveryAddress:           payloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:                      fmtPoundPtr(s.WeightEstimate),
		ProgearWeightEstimate:               fmtPoundPtr(s.ProgearWeightEstimate),
		SpouseProgearWeightEstimate:         fmtPoundPtr(s.SpouseProgearWeightEstimate),
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
	deliveryAddress := addressModelFromPayload(payload.DeliveryAddress)
	partialSITDeliveryAddress := addressModelFromPayload(payload.PartialSitDeliveryAddress)
	market := "dHHG"
	codeOfService := "D"

	var requestedPickupDate *time.Time
	if payload.RequestedPickupDate != nil {
		date := time.Time(*payload.RequestedPickupDate)
		requestedPickupDate = &date
	}

	newShipment := models.Shipment{
		MoveID:                       move.ID,
		ServiceMemberID:              session.ServiceMemberID,
		Status:                       "DRAFT",
		RequestedPickupDate:          requestedPickupDate,
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
		Market:        &market,
		CodeOfService: &codeOfService,
	}

	verrs, err := models.SaveShipmentAndAddresses(h.db, &newShipment)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	shipmentPayload := payloadForShipmentModel(newShipment)
	return shipmentop.NewCreateShipmentCreated().WithPayload(shipmentPayload)
}

func patchShipmentWithPremoveSurveyFields(shipment *models.Shipment, payload *internalmessages.Shipment) {
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

func patchShipmentWithPayload(shipment *models.Shipment, payload *internalmessages.Shipment) {

	if payload.PickupDate != nil {
		shipment.PickupDate = (*time.Time)(payload.PickupDate)
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
	if payload.HasSecondaryPickupAddress == false {
		shipment.SecondaryPickupAddress = nil
	} else if payload.HasSecondaryPickupAddress == true {
		if payload.SecondaryPickupAddress != nil {
			if shipment.SecondaryPickupAddress == nil {
				shipment.SecondaryPickupAddress = addressModelFromPayload(payload.SecondaryPickupAddress)
			} else {
				updateAddressWithPayload(shipment.SecondaryPickupAddress, payload.SecondaryPickupAddress)
			}
		}
	}
	shipment.HasSecondaryPickupAddress = payload.HasSecondaryPickupAddress
	if payload.HasDeliveryAddress == false {
		shipment.DeliveryAddress = nil
	} else if payload.HasDeliveryAddress == true {
		if payload.DeliveryAddress != nil {
			if shipment.DeliveryAddress == nil {
				shipment.DeliveryAddress = addressModelFromPayload(payload.DeliveryAddress)
			} else {
				updateAddressWithPayload(shipment.DeliveryAddress, payload.DeliveryAddress)
			}
		}
	}
	shipment.HasDeliveryAddress = payload.HasDeliveryAddress

	if payload.HasPartialSitDeliveryAddress == false {
		shipment.PartialSITDeliveryAddress = nil
	} else if payload.HasPartialSitDeliveryAddress == true {
		if payload.PartialSitDeliveryAddress != nil {
			if shipment.PartialSITDeliveryAddress == nil {
				shipment.PartialSITDeliveryAddress = addressModelFromPayload(payload.PartialSitDeliveryAddress)
			} else {
				updateAddressWithPayload(shipment.PartialSITDeliveryAddress, payload.PartialSitDeliveryAddress)
			}
		}
	}
	shipment.HasPartialSITDeliveryAddress = payload.HasPartialSitDeliveryAddress

	if payload.WeightEstimate != nil {
		shipment.WeightEstimate = poundPtrFromInt64Ptr(payload.WeightEstimate)
	}
	if payload.ProgearWeightEstimate != nil {
		shipment.ProgearWeightEstimate = poundPtrFromInt64Ptr(payload.ProgearWeightEstimate)
	}
	if payload.SpouseProgearWeightEstimate != nil {
		shipment.SpouseProgearWeightEstimate = poundPtrFromInt64Ptr(payload.SpouseProgearWeightEstimate)
	}
}

// PatchShipmentHandler Patchs an HHG
type PatchShipmentHandler HandlerContext

// Handle is the handler
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.db, session, shipmentID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	patchShipmentWithPayload(shipment, params.Shipment)

	if session.IsOfficeUser() {
		patchShipmentWithPremoveSurveyFields(shipment, params.Shipment)
	}

	verrs, err := models.SaveShipmentAndAddresses(h.db, shipment)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	shipmentPayload := payloadForShipmentModel(*shipment)
	return shipmentop.NewPatchShipmentOK().WithPayload(shipmentPayload)
}

// GetShipmentHandler Returns an HHG
type GetShipmentHandler HandlerContext

// Handle is the handler
func (h GetShipmentHandler) Handle(params shipmentop.GetShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.db, session, shipmentID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	shipmentPayload := payloadForShipmentModel(*shipment)
	return shipmentop.NewGetShipmentOK().WithPayload(shipmentPayload)
}
