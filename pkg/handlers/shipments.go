package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
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
	deliveryAddress := addressModelFromPayload(payload.DeliveryAddress)
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

// PatchShipmentHandler Patchs a PPM
type PatchShipmentHandler HandlerContext

// Handle is the handler
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	// #nosec UUID is pattern matched by swagger and will be ok
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.db, session, shipmentID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	if shipment.MoveID != moveID {
		h.logger.Info("Move ID for Shipment does not match requested Shipment Move ID", zap.String("requested move_id", moveID.String()), zap.String("actual move_id", shipment.MoveID.String()))
		return shipmentop.NewPatchShipmentBadRequest()
	}
	patchShipmentWithPayload(shipment, params.Shipment)

	verrs, err := models.SaveShipmentAndAddresses(h.db, shipment)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	shipmentPayload := payloadForShipmentModel(*shipment)
	return shipmentop.NewPatchShipmentCreated().WithPayload(shipmentPayload)
}

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

// PublicIndexShipmentHandler returns a list of shipments
type PublicIndexShipmentHandler HandlerContext

// Handle retrieves a list of all shipments
func (h PublicIndexShipmentHandler) Handle(p publicshipmentop.IndexShipmentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .indexShipments has not yet been implemented")
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
